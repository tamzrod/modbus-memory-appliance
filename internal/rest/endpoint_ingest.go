package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"modbus-memory-appliance/internal/ingest"
)

func (h *Handlers) HandleIngest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, reject("method not allowed"))
		return
	}

	if !h.EnableIngest {
		writeJSON(w, http.StatusForbidden, reject("ingest disabled"))
		return
	}

	var req ingestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Stats.IncRejected()
		writeJSON(w, http.StatusBadRequest, reject("invalid json"))
		return
	}

	req.Area = strings.TrimSpace(req.Area)

	// REST ingest rule: only discrete_inputs and input_registers
	if req.Area != "discrete_inputs" && req.Area != "input_registers" {
		h.Stats.IncRejected()
		writeJSON(
			w,
			http.StatusForbidden,
			rejectWith("area is not writable via ingest", "area", req.Area),
		)
		return
	}

	// Payload validation
	if req.Area == "discrete_inputs" {
		if len(req.Bools) == 0 || len(req.Values) != 0 {
			h.Stats.IncRejected()
			writeJSON(
				w,
				http.StatusBadRequest,
				reject("bools required for discrete_inputs"),
			)
			return
		}
	}

	if req.Area == "input_registers" {
		if len(req.Values) == 0 || len(req.Bools) != 0 {
			h.Stats.IncRejected()
			writeJSON(
				w,
				http.StatusBadRequest,
				reject("values required for input_registers"),
			)
			return
		}
	}

	cmd := ingest.Command{
		Memory:  req.Memory,
		Area:    ingest.Area(req.Area), // simple cast, no enums
		Address: req.Address,
		Bools:   req.Bools,
		Values:  req.Values,
	}

	h.Stats.IncIngest()
	h.Stats.IncIngestBatch()

	if err := h.Ingest.Ingest(cmd); err != nil {
		h.Stats.IncRejected()
		h.Stats.IncIngestRejected()
		writeIngestError(w, err)
		return
	}

	written := len(req.Bools) + len(req.Values)
	h.Stats.AddWrittenRegs(uint32(written))

	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "accepted",
		"memory":  "default",
		"written": written,
	})
}
