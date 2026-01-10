package rest

import (
	"errors"
	"net/http"

	"modbus-memory-appliance/internal/ingest"
)

func writeIngestError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ingest.ErrUnknownMemory):
		http.Error(w, err.Error(), http.StatusNotFound)

	case errors.Is(err, ingest.ErrInvalidArea),
		errors.Is(err, ingest.ErrInvalidPayload),
		errors.Is(err, ingest.ErrInvalidBoolean),
		errors.Is(err, ingest.ErrPayloadMismatch):
		http.Error(w, err.Error(), http.StatusBadRequest)

	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
