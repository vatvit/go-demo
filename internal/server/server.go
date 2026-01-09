package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/vatvit/go-demo/internal/config"
	"github.com/vatvit/go-demo/internal/handler"
)

type Server struct {
	config  *config.Config
	handler *handler.Handler
	server  *http.Server
}

func New(cfg *config.Config) *Server {
	h := handler.New()

	return &Server{
		config:  cfg,
		handler: h,
		server: &http.Server{
			Addr:    ":" + cfg.Port,
			Handler: h.Routes(),
		},
	}
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

	fmt.Println("Server stopped gracefully")
	return nil
}
