package utils

import (
	"log/slog"
	"net/http"
)

func WriteJsonResponse(w http.ResponseWriter, statusCode int, response []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if _, err := w.Write(response); err != nil {
		slog.Error("error writing response", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
