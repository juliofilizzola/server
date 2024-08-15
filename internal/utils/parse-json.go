package utils

import (
	"encoding/json"
	"net/http"
)

func ParseJson[T any](data *http.Request, body *T) error {
	if err := json.NewDecoder(data.Body).Decode(&body); err != nil {
		return err
	}
	return nil
}
