package handler

import (
	"encoding/json"
	"net/http"
)

type healthResponse struct {
	Status string `json:"status"`
}

func (h *Handler) Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(healthResponse{Status: "ok"}); err != nil {
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		}
	}
}
