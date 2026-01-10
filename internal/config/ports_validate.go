package config

import "fmt"

func (c *AppConfig) ValidatePorts() error {
	// Ports are optional
	if c.Ports == nil {
		return nil
	}

	for port, p := range c.Ports {
		if port == 0 {
			return fmt.Errorf("invalid port 0")
		}

		if !p.UnitIDs.All && len(p.UnitIDs.List) == 0 {
			return fmt.Errorf("ports.%d.unit_ids cannot be empty", port)
		}

		if !p.Memories.All && len(p.Memories.List) == 0 {
			return fmt.Errorf("ports.%d.memories cannot be empty", port)
		}

		if p.Access != AccessReadOnly && p.Access != AccessReadWrite {
			return fmt.Errorf("ports.%d.access invalid", port)
		}

		// Validate unit IDs exist in routing
		for _, uid := range p.UnitIDs.List {
			if _, ok := c.Routing.UnitIDMap[uid]; !ok {
				return fmt.Errorf(
					"ports.%d: unit_id %d not in routing.unit_id_map",
					port, uid,
				)
			}
		}

		// Validate memories exist
		for _, mem := range p.Memories.List {
			if _, ok := c.Memory.Memories[mem]; !ok {
				return fmt.Errorf(
					"ports.%d: unknown memory %q",
					port, mem,
				)
			}
		}
	}

	return nil
}
