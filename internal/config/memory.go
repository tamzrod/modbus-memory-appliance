// internal/config/memory.go
package config

import "fmt"

// =========================
// Memory Configuration Root
// =========================

type MemoryConfig struct {
	Memories map[string]MemoryBlock `yaml:"memories"`
}

// =========================
// Memory Block
// =========================

type MemoryBlock struct {
	Default bool `yaml:"default"`

	Coils            AreaConfig `yaml:"coils"`
	DiscreteInputs   AreaConfig `yaml:"discrete_inputs"`
	HoldingRegisters AreaConfig `yaml:"holding_registers"`
	InputRegisters   AreaConfig `yaml:"input_registers"`

	StateSealing *StateSealingConfig `yaml:"state_sealing,omitempty"`
}

// =========================
// Area Definition
// =========================

type AreaConfig struct {
	Start int `yaml:"start"`
	Size  int `yaml:"size"`
}

// =========================
// State Sealing Config
// =========================

type StateSealingConfig struct {
	Enable bool `yaml:"enable"`
	Gate   GateConfig `yaml:"gate"`
}

type GateConfig struct {
	Area    string `yaml:"area"`
	Address int    `yaml:"address"`
}

// =========================
// Validation
// =========================

// Validate validates the entire memory configuration.
func (c *MemoryConfig) Validate() error {
	if len(c.Memories) == 0 {
		return fmt.Errorf("no memories defined")
	}

	hasDefault := false

	for name, mem := range c.Memories {
		if name == "" {
			return fmt.Errorf("memory name cannot be empty")
		}

		if mem.Default {
			hasDefault = true
		}

		if err := validateArea(mem.Coils, "coils", name); err != nil {
			return err
		}
		if err := validateArea(mem.DiscreteInputs, "discrete_inputs", name); err != nil {
			return err
		}
		if err := validateArea(mem.HoldingRegisters, "holding_registers", name); err != nil {
			return err
		}
		if err := validateArea(mem.InputRegisters, "input_registers", name); err != nil {
			return err
		}

		if mem.StateSealing != nil && mem.StateSealing.Enable {
			if err := validateStateSealing(mem); err != nil {
				return err
			}
		}
	}

	if !hasDefault {
		return fmt.Errorf("no default memory defined")
	}

	return nil
}

func validateArea(a AreaConfig, areaName, memName string) error {
	if a.Size <= 0 {
		return fmt.Errorf(
			"memory '%s': area '%s' size must be > 0",
			memName,
			areaName,
		)
	}
	if a.Start < 0 {
		return fmt.Errorf(
			"memory '%s': area '%s' start must be >= 0",
			memName,
			areaName,
		)
	}
	return nil
}

func validateStateSealing(mem MemoryBlock) error {
	g := mem.StateSealing.Gate

	if g.Address < 0 {
		return fmt.Errorf("state_sealing gate address must be >= 0")
	}

	if g.Area != "discrete_inputs" {
		return fmt.Errorf(
			"state_sealing gate area must be 'discrete_inputs', got '%s'",
			g.Area,
		)
	}

	return nil
}
