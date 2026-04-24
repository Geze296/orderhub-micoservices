package helper

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, data any, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data == nil {
		json.NewEncoder(w).Encode(map[string]any{
			"Status Code": statusCode,
			"message":     message,
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"Status Code": statusCode,
		"message":     message,
		"data":        data,
	})
}

func WriteError(w http.ResponseWriter, statusCode int, message string) {
	WriteJson(w, nil, statusCode, message)
}