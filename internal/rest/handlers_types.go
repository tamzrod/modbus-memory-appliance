package rest

import (
	"encoding/json"
	"net/http"

	"modbus-memory-appliance/internal/core"
	"modbus-memory-appliance/internal/ingest"
)

/*
Shared REST types and helpers.
NO endpoint logic here.
*/

type Handlers struct {
	// Core dependencies
	Memory *core.Memory
	Ingest *ingest.Service
	Stats  *Stats

	// Feature flags
	EnableIngest      bool
	EnableRead        bool
	EnableDiagnostics bool
}

/* -----------------------------
   Shared helpers
------------------------------*/

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

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
