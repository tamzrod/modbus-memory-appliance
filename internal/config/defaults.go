package config

import "fmt"

// DefaultMemoryID returns the memory ID marked as default.
func (c *MemoryConfig) DefaultMemoryID() (string, error) {
	for id, mem := range c.Memories {
		if mem.Default {
			return id, nil
		}
	}
	return "", fmt.Errorf("no default memory found")
}
