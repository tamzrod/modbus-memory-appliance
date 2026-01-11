package main

import (
	"log"

	"modbus-memory-appliance/internal/config"
	"modbus-memory-appliance/internal/core"
	"modbus-memory-appliance/internal/ingest"
	"modbus-memory-appliance/internal/rest"
)

func startREST(cfg *config.AppConfig, memories map[string]*core.Memory, ingestSvc *ingest.Service) {
	if !cfg.REST.Enabled {
		return
	}

	handlers := &rest.Handlers{
		MemoryConfig:     &cfg.Memory,
		Ingest:           ingestSvc,
		Stats:            rest.NewStats(),
		EnableIngest:     true,
		EnableRead:       true,
		EnableDiagnostics:true,
	}

	go func() {
		log.Printf("Starting REST server on %s", cfg.REST.Address)
		srv := rest.NewServer(cfg.REST.Address, handlers, nil)
		log.Fatal(srv.ListenAndServe())
	}()
}
