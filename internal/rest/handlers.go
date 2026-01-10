package rest

import (
	"encoding/json"
	"net/http"

	"modbus-memory-appliance/internal/ingest"
)

type Handlers struct {
	Ingest *ingest.Service
}

func (h *Handlers) handleIngest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ingestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	cmd := ingest.Command{
		Memory:  req.Memory,
		Area:    ingest.Area(req.Area),
		Address: req.Address,
		Bools:   req.Bools,
		Values:  req.Values,
	}

	if err := h.Ingest.Ingest(cmd); err != nil {
		writeIngestError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
