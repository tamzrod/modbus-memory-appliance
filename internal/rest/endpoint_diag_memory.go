// File: endpoint_diag_memory.go
// Endpoint: GET /api/v1/diagnostics/memory
// Purpose: Report configured memory layouts (config truth)

package rest

import "net/http"

func (h *Handlers) HandleDiagnosticsMemory(w http.ResponseWriter, r *http.Request) {
	if !h.EnableDiagnostics {
		writeJSON(w, http.StatusForbidden, reject("diagnostics disabled"))
		return
	}

	out := make(map[string]any)

	for name, mem := range h.MemoryConfig.Memories {
		out[name] = map[string]any{
			"default": mem.Default,
			"coils": map[string]any{
				"start": mem.Coils.Start,
				"size":  mem.Coils.Size,
			},
			"discrete_inputs": map[string]any{
				"start": mem.DiscreteInputs.Start,
				"size":  mem.DiscreteInputs.Size,
			},
			"holding_registers": map[string]any{
				"start": mem.HoldingRegisters.Start,
				"size":  mem.HoldingRegisters.Size,
			},
			"input_registers": map[string]any{
				"start": mem.InputRegisters.Start,
				"size":  mem.InputRegisters.Size,
			},
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"memories": out,
	})
}
