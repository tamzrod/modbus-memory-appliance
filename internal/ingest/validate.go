// internal/ingest/validate.go
package ingest

import "modbus-memory-appliance/internal/core"

func (s *Service) writeDiscreteInputs(mem *core.Memory, cmd Command) error {
	bools := make([]bool, len(cmd.Bools))

	for i, v := range cmd.Bools {
		switch v {
		case 0:
			bools[i] = false
		case 1:
			bools[i] = true
		default:
			return ErrInvalidBoolean
		}
	}

	return mem.WriteDiscreteInputs(
		int(cmd.Address),
		bools,
	)
}

func (s *Service) writeInputRegisters(mem *core.Memory, cmd Command) error {
	return mem.WriteInputRegs(
		int(cmd.Address),
		cmd.Values,
	)
}
func (s *Service) writeCoils(mem *core.Memory, cmd Command) error {
	bools := make([]bool, len(cmd.Bools))

	for i, v := range cmd.Bools {
		switch v {
		case 0:
			bools[i] = false
		case 1:
			bools[i] = true
		default:
			return ErrInvalidBoolean
		}
	}

	return mem.WriteCoils(
		int(cmd.Address),
		bools,
	)
}

func (s *Service) writeHoldingRegisters(mem *core.Memory, cmd Command) error {
	return mem.WriteHoldingRegs(
		int(cmd.Address),
		cmd.Values,
	)
}
