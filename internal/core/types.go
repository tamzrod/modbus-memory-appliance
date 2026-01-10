package core

import "sync"

// Memory is the locked, raw Modbus memory model.
// ZERO-BASED. RAW ONLY. THREAD-SAFE.
type Memory struct {
	Coils          []bool
	DiscreteInputs []bool
	HoldingRegs    []uint16
	InputRegs      []uint16

	mu sync.RWMutex
}
