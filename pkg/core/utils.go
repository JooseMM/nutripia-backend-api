package core

import (
	"encoding/json"
	"net/http"
)

func DecodeJSON[T any](w http.ResponseWriter, r *http.Request) (T, []string) {
	var payload T

	// Limit the request body size to prevent memory exhaustion attacks
	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB limit

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return payload, []string{"Invalid request payload: " + err.Error()}
	}

	return payload, nil
}
