// internal/core/memory.go
package core

import (
	"log"
	"sync"
)

// ===========================
// Run State
// ===========================
type RunState uint8

const (
	StateRun RunState = iota
	StatePreRun
)

// ===========================
// State Sealing Gate
// ===========================
type StateSealingGate struct {
	Enabled  bool
	GateAddr int
}

// ===========================
// Memory (RAW, THREAD-SAFE)
// ===========================
type Memory struct {
	Coils          []bool
	DiscreteInputs []bool
	HoldingRegs    []uint16
	InputRegs      []uint16

	state RunState
	seal  *StateSealingGate

	mu sync.RWMutex
}

// ===========================
// Constructor
// ===========================
func NewMemory(
	coilCount,
	discreteCount,
	holdingCount,
	inputCount int,
) *Memory {
	return &Memory{
		Coils:          make([]bool, coilCount),
		DiscreteInputs: make([]bool, discreteCount),
		HoldingRegs:    make([]uint16, holdingCount),
		InputRegs:      make([]uint16, inputCount),
		state:          StateRun,
	}
}

// ===========================
// State Sealing
// ===========================
func (m *Memory) SetStateSealing(enable bool, gateAddr int) {
	if !enable {
		return
	}

	m.seal = &StateSealingGate{
		Enabled:  true,
		GateAddr: gateAddr,
	}
	m.state = StatePreRun
}

func (m *Memory) HasStateSealing() bool {
	return m.seal != nil && m.seal.Enabled
}

func (m *Memory) IsPreRun() bool {
	return m.HasStateSealing() && m.state == StatePreRun
}

func (m *Memory) GateAddress() int {
	if !m.HasStateSealing() {
		return -1
	}
	return m.seal.GateAddr
}

// ===========================
// Internal sealing transition
// ===========================
func (m *Memory) transitionToRunIfGateHit(writeAddr int, values []bool) {
	if !m.HasStateSealing() || m.state != StatePreRun {
		return
	}

	g := m.seal.GateAddr
	if writeAddr <= g && writeAddr+len(values) > g {
		if values[g-writeAddr] {
			m.state = StateRun
			log.Printf("[STATE] memory transitioned to RUN via gate @ discrete_inputs[%d]", g)
		}
	}
}

// ===========================
// COILS
// ===========================
func (m *Memory) ReadCoils(addr, count int) ([]bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if err := checkBoolRange(m.Coils, addr, count); err != nil {
		return nil, err
	}

	out := make([]bool, count)
	copy(out, m.Coils[addr:addr+count])
	return out, nil
}

func (m *Memory) WriteCoils(addr int, values []bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := checkBoolRange(m.Coils, addr, len(values)); err != nil {
		return err
	}

	copy(m.Coils[addr:], values)
	return nil
}

// ===========================
// DISCRETE INPUTS
// ===========================
func (m *Memory) ReadDiscreteInputs(addr, count int) ([]bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if err := checkBoolRange(m.DiscreteInputs, addr, count); err != nil {
		return nil, err
	}

	out := make([]bool, count)
	copy(out, m.DiscreteInputs[addr:addr+count])
	return out, nil
}

func (m *Memory) WriteDiscreteInputs(addr int, values []bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := checkBoolRange(m.DiscreteInputs, addr, len(values)); err != nil {
		return err
	}

	copy(m.DiscreteInputs[addr:], values)

	// 🔒 State Sealing gate check
	m.transitionToRunIfGateHit(addr, values)

	return nil
}

// ===========================
// HOLDING REGISTERS
// ===========================
func (m *Memory) ReadHoldingRegs(addr, count int) ([]uint16, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if err := checkUint16Range(m.HoldingRegs, addr, count); err != nil {
		return nil, err
	}

	out := make([]uint16, count)
	copy(out, m.HoldingRegs[addr:addr+count])
	return out, nil
}

func (m *Memory) WriteHoldingRegs(addr int, values []uint16) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := checkUint16Range(m.HoldingRegs, addr, len(values)); err != nil {
		return err
	}

	copy(m.HoldingRegs[addr:], values)
	return nil
}

// ===========================
// INPUT REGISTERS
// ===========================
func (m *Memory) ReadInputRegs(addr, count int) ([]uint16, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if err := checkUint16Range(m.InputRegs, addr, count); err != nil {
		return nil, err
	}

	out := make([]uint16, count)
	copy(out, m.InputRegs[addr:addr+count])
	return out, nil
}

func (m *Memory) WriteInputRegs(addr int, values []uint16) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := checkUint16Range(m.InputRegs, addr, len(values)); err != nil {
		return err
	}

	copy(m.InputRegs[addr:], values)
	return nil
}
