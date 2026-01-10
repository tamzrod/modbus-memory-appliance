package config

import "fmt"

type RoutingConfig struct {
	UnitIDMap map[uint8]string `yaml:"unit_id_map"`
}

func (r *RoutingConfig) Validate(mem MemoryConfig) error {
	if len(r.UnitIDMap) == 0 {
		return fmt.Errorf("unit_id_map must not be empty (strict mapping enabled)")
	}

	for unitID, memID := range r.UnitIDMap {
		if _, ok := mem.Memories[memID]; !ok {
			return fmt.Errorf(
				"unit_id_map[%d] references unknown memory '%s'",
				unitID, memID,
			)
		}
	}

	return nil
}
