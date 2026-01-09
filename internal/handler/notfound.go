package handler

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) NotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(errorResponse{Error: "not found"}); err != nil {
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		}
	}
}
