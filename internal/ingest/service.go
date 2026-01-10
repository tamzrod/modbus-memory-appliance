// internal/ingest/service.go
package ingest

import (
	"modbus-memory-appliance/internal/core"
)

func isValidArea(a Area) bool {
	switch a {
	case Coils, DiscreteInputs, HoldingRegs, InputRegisters:
		return true
	default:
		return false
	}
}

type Service struct {
	memories map[string]*core.Memory
}

func New(memories map[string]*core.Memory) *Service {
	return &Service{memories: memories}
}

func (s *Service) Ingest(cmd Command) error {
	// 1. Memory exists
	mem, ok := s.memories[cmd.Memory]
	if !ok {
		return ErrUnknownMemory
	}

	// 2. Area must be valid FIRST
	if !isValidArea(cmd.Area) {
		return ErrInvalidArea
	}

	// 3. Exactly one payload
	hasBools := len(cmd.Bools) > 0
	hasValues := len(cmd.Values) > 0
	if hasBools == hasValues {
		return ErrInvalidPayload
	}

	// 4. Area ↔ payload compatibility
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
