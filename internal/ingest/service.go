package ingest

import (
	"modbus-memory-appliance/internal/core"
)

// Service handles REST/MQTT ingestion into memories.
// It NEVER affects Modbus behavior.
type Service struct {
	memories map[string]*core.Memory
}

// New creates a new ingest service.
func New(memories map[string]*core.Memory) *Service {
	return &Service{
		memories: memories,
	}
}

// Ingest applies a validated command to memory.
func (s *Service) Ingest(cmd Command) error {
	// 1. Resolve memory
	mem, ok := s.memories[cmd.Memory]
	if !ok {
		return ErrUnknownMemory
	}

	// 2. Validate area
	if !isValidArea(cmd.Area) {
		return ErrInvalidArea
	}

	// 3. Validate payload presence (exactly one)
	hasBools := len(cmd.Bools) > 0
	hasValues := len(cmd.Values) > 0
	if hasBools == hasValues {
		return ErrInvalidPayload
	}

	// =====================================================
	// 🔒 STATE RULES (REST / MQTT ONLY)
	// =====================================================
	// PRE-RUN  → allow ALL areas
	// RUN      → block coils + holding registers
	if !mem.IsPreRun() {
		switch cmd.Area {
		case Coils, HoldingRegs:
			return ErrIngestDenied
		}
	}

	// =====================================================
	// AREA → PAYLOAD COMPATIBILITY
	// =====================================================
	switch cmd.Area {

	case Coils:
		if !hasBools {
			return ErrPayloadMismatch
		}
		return s.writeCoils(mem, cmd)

	case DiscreteInputs:
		if !hasBools {
			return ErrPayloadMismatch
		}
		return s.writeDiscreteInputs(mem, cmd)

	case HoldingRegs:
		if !hasValues {
			return ErrPayloadMismatch
		}
		return s.writeHoldingRegisters(mem, cmd)

	case InputRegisters:
		if !hasValues {
			return ErrPayloadMismatch
		}
		return s.writeInputRegisters(mem, cmd)
	}

	// unreachable
	return ErrInvalidArea
}

// -----------------------------------------------------
// Helpers
// -----------------------------------------------------

func isValidArea(a Area) bool {
	switch a {
	case Coils, DiscreteInputs, HoldingRegs, InputRegisters:
		return true
	default:
		return false
	}
}
