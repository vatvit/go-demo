package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

const (
	statusOK        = "ok"
	statusDegraded  = "degraded"
	statusUnhealthy = "unhealthy"

	serviceConnected = "connected"
	healthTimeout    = 2 * time.Second
)

type healthResponse struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
}

func (h *Handler) Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), healthTimeout)
		defer cancel()

		services := make(map[string]string)
		var mu sync.Mutex
		var wg sync.WaitGroup

		// Check MongoDB concurrently
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := serviceConnected
			if err := h.mongo.Ping(ctx); err != nil {
				result = "error: " + err.Error()
			}
			mu.Lock()
			services["mongodb"] = result
			mu.Unlock()
		}()

		// Check Redis concurrently
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := serviceConnected
			if err := h.redis.Ping(ctx); err != nil {
				result = "error: " + err.Error()
			}
			mu.Lock()
			services["redis"] = result
			mu.Unlock()
		}()

		wg.Wait()

		status := determineStatus(services)
		httpCode := http.StatusOK
		if status == statusUnhealthy {
			httpCode = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpCode)
		if err := json.NewEncoder(w).Encode(healthResponse{
			Status:   status,
			Services: services,
		}); err != nil {
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		}
	}
}

func determineStatus(services map[string]string) string {
	errorCount := 0
	for _, status := range services {
		if status != serviceConnected {
			errorCount++
		}
	}

	switch errorCount {
	case 0:
		return statusOK
	case len(services):
		return statusUnhealthy
	default:
		return statusDegraded
	}
}
