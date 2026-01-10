package config

import (
	"fmt"

	"modbus-memory-appliance/internal/core"
)

// BuildMemories creates runtime core.Memory instances
// from a validated MemoryConfig.
func BuildMemories(cfg *MemoryConfig) (map[string]*core.Memory, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	memories := make(map[string]*core.Memory)

	for memID, block := range cfg.Memories {
		mem := core.NewMemory(
			block.Coils.Size,
			block.DiscreteInputs.Size,
			block.HoldingRegisters.Size,
			block.InputRegisters.Size,
		)

		if mem == nil {
			return nil, fmt.Errorf(
				"failed to create memory '%s'",
				memID,
			)
		}

		memories[memID] = mem
	}

	return memories, nil
}
