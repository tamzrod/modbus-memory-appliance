package rest

import "net/http"

func (h *Handlers) HandleDiagnosticsMemory(w http.ResponseWriter, r *http.Request) {
	if !h.EnableDiagnostics {
		writeJSON(w, http.StatusForbidden, reject("diagnostics disabled"))
		return
	}

	m := h.Memory

	writeJSON(w, http.StatusOK, map[string]any{
		"memories": map[string]any{
			"default": map[string]any{
				"coils":             len(m.Coils),
				"discrete_inputs":   len(m.DiscreteInputs),
				"holding_registers": len(m.HoldingRegs),
				"input_registers":   len(m.InputRegs),
			},
		},
	})
}
