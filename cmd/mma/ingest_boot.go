// cmd/mma/ingest_boot.go
package main

import (
	"modbus-memory-appliance/internal/core"
	"modbus-memory-appliance/internal/ingest"
)

func buildIngest(memories map[string]*core.Memory) *ingest.Service {
	return ingest.New(memories)
}
