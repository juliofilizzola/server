package utils

import (
	"net/http"
)

func WriteJsonResponse(w http.ResponseWriter, statusCode int, response []byte) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err := w.Write(response)

	if err != nil {
		return err
	}
	return nil
}
