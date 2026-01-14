// internal/config/factory.go
package config

import (
	"fmt"

	"modbus-memory-appliance/internal/core"
)

// BuildMemories creates runtime core.Memory instances
// from a validated MemoryConfig.
func BuildMemories(cfg *MemoryConfig) (map[string]*core.Memory, error) {
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

		// =========================
		// Apply State Sealing (optional, per memory)
		// =========================
		if block.StateSealing != nil && block.StateSealing.Enable {
			mem.SetStateSealing(
				true,
				block.StateSealing.Gate.Address,
			)

			fmt.Printf(
				"[BOOT] memory=%s state_sealing=enabled prerun=%v gate=%d\n",
				memID,
				mem.IsPreRun(),
				mem.GateAddress(),
			)
		} else {
			fmt.Printf(
				"[BOOT] memory=%s state_sealing=disabled\n",
				memID,
			)
		}

		memories[memID] = mem
	}

	return memories, nil
}
