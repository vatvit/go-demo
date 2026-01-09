package handler

import "net/http"

type Handler struct{}

func New() *Handler {
	return &Handler{}
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
