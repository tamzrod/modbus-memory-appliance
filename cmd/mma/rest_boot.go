package main

import (
	"log"
	"net/http"

	"modbus-memory-appliance/internal/config"
	"modbus-memory-appliance/internal/core"
	"modbus-memory-appliance/internal/ingest"
	"modbus-memory-appliance/internal/rest"
)

func startREST(
	cfg *config.AppConfig,
	memories map[string]*core.Memory,
	ingestSvc *ingest.Service,
) {
	// ---- config gate ----
	if !cfg.REST.Enabled {
		log.Println("REST disabled by config")
		return
	}

	if cfg.REST.Address == "" {
		log.Println("REST enabled but no address configured")
		return
	}

	// ---- handlers ----
	handlers := &rest.Handlers{
		MemoryConfig:      &cfg.Memory,
		Ingest:            ingestSvc,
		Stats:             rest.NewStats(),
		EnableIngest:      true,
		EnableRead:        true,
		EnableDiagnostics: true,
	}

	// ---- AUTH (HARD-WIRED FOR NOW) ----
	// This uses your EXISTING TokenSet implementation
	// Change the token value if you want
	tokenSet := rest.NewTokenSet(
		true, // enable auth
		[]string{
			"INGEST_ONLY_TOKEN",
		},
	)

	authMiddleware := func(next http.Handler) http.Handler {
		return tokenSet.Require(next, handlers.Stats)
	}

	// ---- start server ----
	go func() {
		log.Printf("Starting REST server on %s", cfg.REST.Address)

		srv := rest.NewServer(
			cfg.REST.Address,
			handlers,
			authMiddleware, // 🔐 AUTH IS NOW ACTIVE
		)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("REST server error: %v", err)
		}
	}()
}
