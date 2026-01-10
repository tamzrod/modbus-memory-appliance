package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

func (s *UnitIDSelector) UnmarshalYAML(n *yaml.Node) error {
	switch n.Kind {
	case yaml.ScalarNode:
		if n.Value == "all" {
			s.All = true
			s.List = nil
			return nil
		}
		return fmt.Errorf("unit_ids must be 'all' or a list")

	case yaml.SequenceNode:
		for _, c := range n.Content {
			var v int
			if err := c.Decode(&v); err != nil {
				return err
			}
			if v < 0 || v > 255 {
				return fmt.Errorf("unit_id %d out of range (0..255)", v)
			}
			s.List = append(s.List, uint8(v))
		}
		return nil
	}
	return fmt.Errorf("invalid unit_ids format")
}

func (s *MemorySelector) UnmarshalYAML(n *yaml.Node) error {
	switch n.Kind {
	case yaml.ScalarNode:
		if n.Value == "all" {
			s.All = true
			s.List = nil
			return nil
		}
		return fmt.Errorf("memories must be 'all' or a list")

	case yaml.SequenceNode:
		for _, c := range n.Content {
			var v string
			if err := c.Decode(&v); err != nil {
				return err
			}
			if v == "" {
				return fmt.Errorf("memory name cannot be empty")
			}
			s.List = append(s.List, v)
		}
		return nil
	}
	return fmt.Errorf("invalid memories format")
}
