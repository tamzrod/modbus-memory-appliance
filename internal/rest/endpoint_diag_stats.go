package rest

import "net/http"

func (h *Handlers) HandleDiagnosticsStats(w http.ResponseWriter, r *http.Request) {
	if !h.EnableDiagnostics {
		writeJSON(w, http.StatusForbidden, reject("diagnostics disabled"))
		return
	}

	writeJSON(w, http.StatusOK, h.Stats.Snapshot())
}
