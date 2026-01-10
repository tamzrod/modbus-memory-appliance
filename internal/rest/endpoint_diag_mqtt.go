package rest

import "net/http"

// HandleDiagnosticsMQTT exposes MQTT runtime status via REST.
// It is READ-ONLY diagnostics.
// REST does NOT know MQTT internals.
func (h *Handlers) HandleDiagnosticsMQTT(w http.ResponseWriter, r *http.Request) {
	if !h.EnableDiagnostics {
		writeJSON(w, http.StatusForbidden, reject("diagnostics disabled"))
		return
	}

	// If MQTT is not wired (nil), report disabled
	if h.MQTTStatus == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"enabled": false,
		})
		return
	}

	// Delegate status generation to injected function
	writeJSON(w, http.StatusOK, h.MQTTStatus())
}
