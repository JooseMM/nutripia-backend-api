package response

import (
	"encoding/json"
	"net/http"
)

type ApiResponse[T any] struct {
	Success bool `json:"success"`
	Data    *T   `json:"data"`
}

func WriteJSON(w http.ResponseWriter, status uint16, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(status))
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
