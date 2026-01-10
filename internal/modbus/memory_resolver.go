package modbus

import "modbus-memory-appliance/internal/core"

// MemoryResolver resolves memory based on Unit ID and Function Code.
// Policy decisions live outside Modbus.
type MemoryResolver func(unitID uint8, functionCode uint8) *core.Memory
