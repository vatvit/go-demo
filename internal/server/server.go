package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/vatvit/go-demo/internal/config"
	"github.com/vatvit/go-demo/internal/handler"
	"github.com/vatvit/go-demo/internal/mongodb"
	"github.com/vatvit/go-demo/internal/redis"
)

type Server struct {
	config  *config.Config
	handler *handler.Handler
	server  *http.Server
	mongo   *mongodb.Client
	redis   *redis.Client
}

func New(cfg *config.Config) (*Server, error) {
	mongoClient, err := mongodb.New(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
	}

	redisClient, err := redis.New(cfg.RedisAddr)
	if err != nil {
		mongoClient.Close(context.Background())
		return nil, fmt.Errorf("failed to create Redis client: %w", err)
	}

	h := handler.New(mongoClient, redisClient)

	return &Server{
		config:  cfg,
		handler: h,
		server: &http.Server{
			Addr:    ":" + cfg.Port,
			Handler: h.Routes(),
		},
		mongo: mongoClient,
		redis: redisClient,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		fmt.Printf("Server starting on port %s\n", s.config.Port)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		return s.Shutdown()
	case err := <-errCh:
		return err
	}
}

func (s *Server) Shutdown() error {
	fmt.Println("Server shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	if s.mongo != nil {
		if err := s.mongo.Close(ctx); err != nil {
			fmt.Printf("MongoDB close error: %v\n", err)
		}
	}

	if s.redis != nil {
		if err := s.redis.Close(); err != nil {
			fmt.Printf("Redis close error: %v\n", err)
		}
	}

	fmt.Println("Server stopped gracefully")
	return nil
}
