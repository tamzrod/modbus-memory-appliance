package ingest

import (
	"testing"

	"modbus-memory-appliance/internal/core"
)

func newTestMemoryWithSealing(gateAddr int) *core.Memory {
	mem := core.NewMemory(16, 16, 16, 16)
	mem.SetStateSealing(true, gateAddr)
	return mem
}

func TestPreRun_AllowsFullRestore(t *testing.T) {
	mem := newTestMemoryWithSealing(5)

	svc := New(map[string]*core.Memory{
		"plant": mem,
	})

	err := svc.Ingest(Command{
		Memory:  "plant",
		Area:    HoldingRegs,
		Address: 0,
		Values:  []uint16{123, 456},
	})
	if err != nil {
		t.Fatalf("expected holding ingest to succeed in Pre-Run, got %v", err)
	}

	err = svc.Ingest(Command{
		Memory:  "plant",
		Area:    Coils,
		Address: 0,
		Bools:   []int{1, 0, 1},
	})
	if err != nil {
		t.Fatalf("expected coil ingest to succeed in Pre-Run, got %v", err)
	}
}

func TestGateSealsMemory(t *testing.T) {
	mem := newTestMemoryWithSealing(3)

	svc := New(map[string]*core.Memory{
		"plant": mem,
	})

	// Write gate DI = true
	err := svc.Ingest(Command{
		Memory:  "plant",
		Area:    DiscreteInputs,
		Address: 3,
		Bools:   []int{1},
	})
	if err != nil {
		t.Fatalf("unexpected error writing gate DI: %v", err)
	}

	if mem.IsPreRun() {
		t.Fatalf("expected memory to be sealed after gate flip")
	}
}

func TestRun_DeniesHoldingAndCoilIngest(t *testing.T) {
	mem := newTestMemoryWithSealing(1)
	mem.Seal() // force RUN

	svc := New(map[string]*core.Memory{
		"plant": mem,
	})

	err := svc.Ingest(Command{
		Memory:  "plant",
		Area:    HoldingRegs,
		Address: 0,
		Values:  []uint16{10},
	})
	if err == nil {
		t.Fatalf("expected holding ingest to fail in RUN")
	}

	err = svc.Ingest(Command{
		Memory:  "plant",
		Area:    Coils,
		Address: 0,
		Bools:   []int{1},
	})
	if err == nil {
		t.Fatalf("expected coil ingest to fail in RUN")
	}
}

func TestRun_AllowsDIandIRIngest(t *testing.T) {
	mem := newTestMemoryWithSealing(1)
	mem.Seal()

	svc := New(map[string]*core.Memory{
		"plant": mem,
	})

	err := svc.Ingest(Command{
		Memory:  "plant",
		Area:    DiscreteInputs,
		Address: 0,
		Bools:   []int{1},
	})
	if err != nil {
		t.Fatalf("expected DI ingest to succeed in RUN, got %v", err)
	}

	err = svc.Ingest(Command{
		Memory:  "plant",
		Area:    InputRegisters,
		Address: 0,
		Values:  []uint16{99},
	})
	if err != nil {
		t.Fatalf("expected IR ingest to succeed in RUN, got %v", err)
	}
}
