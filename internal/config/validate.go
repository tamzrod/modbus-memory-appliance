package config

import (
	"errors"
	"fmt"
	"math"
)

func (c *MemoryConfig) Validate() error {
	if len(c.Memories) == 0 {
		return errors.New("no memories defined")
	}

	defaultCount := 0

	for memID, mem := range c.Memories {
		if memID == "" {
			return errors.New("memory ID cannot be empty")
		}

		if mem.Default {
			defaultCount++
		}

		if err := validateArea(mem.Coils, memID, "coils"); err != nil {
			return err
		}
		if err := validateArea(mem.DiscreteInputs, memID, "discrete_inputs"); err != nil {
			return err
		}
		if err := validateArea(mem.HoldingRegisters, memID, "holding_registers"); err != nil {
			return err
		}
		if err := validateArea(mem.InputRegisters, memID, "input_registers"); err != nil {
			return err
		}
	}

	if defaultCount == 0 {
		return errors.New("no default memory defined")
	}
	if defaultCount > 1 {
		return errors.New("multiple default memories defined")
	}

	return nil
}

func validateArea(a AreaConfig, memID, areaName string) error {
	if a.Start < 0 {
		return fmt.Errorf(
			"memory '%s': %s.start must be >= 0",
			memID, areaName,
		)
	}

	if a.Size < 0 {
		return fmt.Errorf(
			"memory '%s': %s.size must be >= 0",
			memID, areaName,
		)
	}

	if a.Size > 0 && a.Start > math.MaxInt-a.Size {
		return fmt.Errorf(
			"memory '%s': %s range overflows integer limits",
			memID, areaName,
		)
	}

	return nil
}
