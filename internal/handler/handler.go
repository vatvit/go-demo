package handler

import (
	"net/http"

	"github.com/vatvit/go-demo/internal/mongodb"
	"github.com/vatvit/go-demo/internal/redis"
)

type Handler struct {
	mongo *mongodb.Client
	redis *redis.Client
}

func New(mongo *mongodb.Client, redis *redis.Client) *Handler {
	return &Handler{
		mongo: mongo,
		redis: redis,
	}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health())

	// Cache the 404 handler to avoid allocation per request
	notFoundHandler := h.NotFound()

	// Wrap with custom 404 handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the route exists
		_, pattern := mux.Handler(r)
		if pattern == "" {
			notFoundHandler(w, r)
			return
		}
		mux.ServeHTTP(w, r)
	})
}
