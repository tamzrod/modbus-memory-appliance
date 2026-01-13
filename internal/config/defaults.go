package config

import "fmt"

// DefaultMaxConnections is the hardcoded fallback limit for concurrent
// TCP connections per port when not specified in config.
//
// This is a mechanical safety guard, not a security feature.
const DefaultMaxConnections = 32

// DefaultMemoryID returns the memory ID marked as default.
func (c *MemoryConfig) DefaultMemoryID() (string, error) {
	for id, mem := range c.Memories {
		if mem.Default {
			return id, nil
		}
	}
	return "", fmt.Errorf("no default memory found")
}
