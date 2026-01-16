// cmd/mma/raw_ingest_boot.go
// PURPOSE: Start Raw Ingest transport (boot wiring only).
// ALLOWED: config gating, adapter wiring, goroutine start
// FORBIDDEN: protocol logic, memory logic, semantics

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"modbus-memory-appliance/internal/config"
	"modbus-memory-appliance/internal/core"
	"modbus-memory-appliance/internal/rawingest"
)

// ---- Adapter: core.Memory -> rawingest.RawWritableMemory ----

type rawIngestMemoryAdapter struct {
	mem *core.Memory
}

func (a *rawIngestMemoryAdapter) WriteCoils(addr uint16, v []bool) error {
	return a.mem.WriteCoils(int(addr), v)
}

func (a *rawIngestMemoryAdapter) WriteDiscreteInputs(addr uint16, v []bool) error {
	return a.mem.WriteDiscreteInputs(int(addr), v)
}

func (a *rawIngestMemoryAdapter) WriteHoldingRegisters(addr uint16, v []uint16) error {
	return a.mem.WriteHoldingRegs(int(addr), v)
}

func (a *rawIngestMemoryAdapter) WriteInputRegisters(addr uint16, v []uint16) error {
	return a.mem.WriteInputRegs(int(addr), v)
}

// ---- Boot wiring ----

func startRawIngest(cfg *config.AppConfig, memories map[string]*core.Memory) {

	//start debug
	log.Printf("[DEBUG] RawIngest enabled=%v listen=%s",
	cfg.RawIngest.Enabled,
	cfg.RawIngest.Listen,
)

    //end debug

	// Self-gate (same pattern as other transports)
	if !cfg.RawIngest.Enabled {
		return
	}

	// Resolver: raw ingest memory_id (uint16) -> memories map key.
	// Raw ingest writes directly to memory; no RoutingConfig is used here.
	resolver := rawingest.MemoryResolverFunc(func(id uint16) (rawingest.RawWritableMemory, bool) {
		key := fmt.Sprintf("%d", id)

		mem, ok := memories[key]
		if !ok {
			return nil, false
		}

		return &rawIngestMemoryAdapter{mem: mem}, true
	})

	readTimeout := time.Duration(cfg.RawIngest.ReadTimeoutMS) * time.Millisecond

	srv, err := rawingest.NewServer(
		cfg.RawIngest.Listen,
		resolver,
		cfg.RawIngest.MaxPacketBytes,
		readTimeout,
	)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Printf("Starting Raw Ingest listener on %s", cfg.RawIngest.Listen)
		if err := srv.Start(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
}
