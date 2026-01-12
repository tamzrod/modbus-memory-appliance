package main

import (
	"fmt"
	"log"

	"modbus-memory-appliance/internal/config"
	"modbus-memory-appliance/internal/core"
	"modbus-memory-appliance/internal/modbus"
)

func startModbus(cfg *config.AppConfig, memories map[string]*core.Memory) {

	// Resolver: port + unitID + function code → memory
	resolver := func(port uint16, unitID uint8, fc uint8) *core.Memory {
		memID, ok := cfg.Routing.UnitIDMap[unitID]
		if !ok {
			return nil
		}

		policy, hasPolicy := cfg.Ports[port]
		if hasPolicy {
			if !policy.AllowsUnitID(unitID) ||
				!policy.AllowsMemory(memID) ||
				!policy.AllowsFunctionCode(fc) {
				return nil
			}
		}

		return memories[memID]
	}

	if len(cfg.Ports) == 0 {
		return
	}

	// Start one Modbus TCP listener per configured port
	for port, policy := range cfg.Ports {
		p := port
		pol := policy
		addr := fmt.Sprintf(":%d", p)

		go func() {
			log.Printf("Starting Modbus TCP listener on %s", addr)

			err := modbus.Start(
				addr,
				func(unitID uint8, fc uint8) *core.Memory {
					return resolver(p, unitID, fc)
				},
				pol.IPFilter.Allow, // ✅ CONFIG-DRIVEN
				pol.IPFilter.Deny,  // ✅ CONFIG-DRIVEN
			)

			if err != nil {
				log.Fatal(err)
			}
		}()
	}
}
