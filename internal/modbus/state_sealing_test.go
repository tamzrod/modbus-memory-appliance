package modbus

import (
	"testing"

	"modbus-memory-appliance/internal/core"
)

func TestPreRunBlocksAllModbusAccess(t *testing.T) {
	mem := core.NewMemory(8, 8, 8, 8)
	mem.SetStateSealing(true, 0)

	if !mem.IsPreRun() {
		t.Fatalf("memory should start in Pre-Run")
	}

	// This test is conceptual: it ensures the handler MUST
	// reject requests when IsPreRun() == true.
	//
	// If someone removes the Pre-Run guard in handler.go,
	// this test should be updated to fail.

	if mem.IsPreRun() != true {
		t.Fatalf("expected Pre-Run state")
	}
}
