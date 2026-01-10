package rest

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func reject(msg string) map[string]any {
	return map[string]any{
		"status": "rejected",
		"error":  msg,
	}
}

func rejectWith(msg, key string, val any) map[string]any {
	return map[string]any{
		"status": "rejected",
		"error":  msg,
		key:      val,
	}
}
