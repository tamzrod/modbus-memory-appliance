package rest

import (
	"modbus-memory-appliance/internal/config"
	"modbus-memory-appliance/internal/core"
	"modbus-memory-appliance/internal/ingest"
)

type Handlers struct {
	// Config truth (diagnostics)
	MemoryConfig *config.MemoryConfig

	// Runtime truth (operations)
	Memories map[string]*core.Memory

	Ingest *ingest.Service
	Stats  *Stats

	MQTTStatus func() any

	EnableIngest      bool
	EnableRead        bool
	EnableDiagnostics bool
}
