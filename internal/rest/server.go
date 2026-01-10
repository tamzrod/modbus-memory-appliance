package rest

import (
	"net/http"
	"time"
)

// NewServer wires the REST HTTP server.
// No business logic here.
func NewServer(
	addr string,
	handlers *Handlers,
	authMiddleware func(http.Handler) http.Handler,
) *http.Server {

	mux := http.NewServeMux()

	// ---- routes ----

	mux.HandleFunc("/api/v1/health",
		handlers.HandleHealth)

	mux.HandleFunc("/api/v1/diagnostics/memory",
		handlers.HandleDiagnosticsMemory)

	mux.HandleFunc("/api/v1/diagnostics/stats",
		handlers.HandleDiagnosticsStats)

	mux.HandleFunc("/api/v1/memory/read",
		handlers.HandleMemoryRead)

	mux.HandleFunc("/api/v1/ingest",
		handlers.HandleIngest)

	mux.HandleFunc(	"/api/v1/diagnostics/mqtt",
	handlers.HandleDiagnosticsMQTT,
)
	

	// ---- middleware ----

	var h http.Handler = mux
	if authMiddleware != nil {
		h = authMiddleware(h)
	}

	// ---- server ----

	return &http.Server{
		Addr:         addr,
		Handler:      h,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
}
