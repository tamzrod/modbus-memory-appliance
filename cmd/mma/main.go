package main

import (
	"fmt"
	"log"
	"time"

	"modbus-memory-appliance/internal/config"
	"modbus-memory-appliance/internal/core"
	"modbus-memory-appliance/internal/ingest"
	"modbus-memory-appliance/internal/modbus"
	"modbus-memory-appliance/internal/mqtt"
	"modbus-memory-appliance/internal/rest"
)

func main() {
	// --------------------------------
	// Load or create configuration
	// --------------------------------
	cfg, err := config.LoadOrCreate()
	if err != nil {
		log.Println(err)
		return
	}

	// --------------------------------
	// Validate ports configuration
	// --------------------------------
	if err := cfg.ValidatePorts(); err != nil {
		log.Fatal(err)
	}

	// --------------------------------
	// Build memory instances
	// --------------------------------
	memories, err := config.BuildMemories(&cfg.Memory)
	if err != nil {
		log.Fatal(err)
	}

	// --------------------------------
	// Create INGEST service
	// --------------------------------
	ingestSvc := ingest.New(memories)

	// --------------------------------
	// MQTT config (SINGLE SOURCE OF TRUTH)
	// --------------------------------
	mqttCfg := mqtt.Config{
		Enabled:  cfg.MQTT.Enabled,
		Broker:   cfg.MQTT.Broker,
		ClientID: cfg.MQTT.ClientID,
		Topic:    cfg.MQTT.Topic,
		Username: cfg.MQTT.Username,
		Password: cfg.MQTT.Password,
	}

	// --------------------------------
	// REST handlers
	// --------------------------------
	handlers := &rest.Handlers{
		MemoryConfig: &cfg.Memory,
		Ingest: ingestSvc,
		Stats:  rest.NewStats(),

		// REST observes MQTT via function (no coupling)
		MQTTStatus: func() any {
			return mqtt.GetStatus(mqttCfg)
		},

		EnableIngest:      true,
		EnableRead:        true,
		EnableDiagnostics: true,
	}

	// --------------------------------
	// Policy-aware resolver (Modbus only)
	// --------------------------------
	resolver := func(port uint16, unitID uint8, fc uint8) *core.Memory {
		memID, ok := cfg.Routing.UnitIDMap[unitID]
		if !ok {
			return nil
		}

		policy, hasPolicy := cfg.Ports[port]
		if hasPolicy {
			if !policy.AllowsUnitID(unitID) {
				return nil
			}
			if !policy.AllowsMemory(memID) {
				return nil
			}
			if !policy.AllowsFunctionCode(fc) {
				return nil
			}
		}

		return memories[memID]
	}

	// ====================================
	// MODBUS TCP LISTENERS
	// ====================================
	if len(cfg.Ports) > 0 {
		for port := range cfg.Ports {
			p := port
			addr := fmt.Sprintf(":%d", p)

			log.Printf("Starting Modbus TCP listener on %s", addr)

			go func() {
				if err := modbus.Start(addr, func(unitID uint8, fc uint8) *core.Memory {
					return resolver(p, unitID, fc)
				}); err != nil {
					log.Fatalf("Modbus listener on port %d failed: %v", p, err)
				}
			}()
		}
	} else {
		const defaultPort = 5020
		addr := fmt.Sprintf(":%d", defaultPort)

		log.Printf("No ports defined, starting legacy Modbus TCP on %s", addr)

		if err := modbus.Start(addr, func(unitID uint8, fc uint8) *core.Memory {
			return resolver(defaultPort, unitID, fc)
		}); err != nil {
			log.Fatal(err)
		}
	}

	// ====================================
	// MQTT INGESTION (NON-FATAL, RETRYING)
	// ====================================
	if mqttCfg.Enabled {
		go func() {
			for {
				log.Printf("Starting MQTT subscriber (broker=%s)", mqttCfg.Broker)

				_, err := mqtt.NewSubscriber(mqttCfg, ingestSvc)
				if err != nil {
					log.Printf("MQTT unavailable: %v", err)
					log.Printf("Retrying MQTT in 10 seconds...")
					time.Sleep(10 * time.Second)
					continue
				}

				// NewSubscriber blocks while connected
				log.Printf("MQTT subscriber exited, restarting...")
				time.Sleep(5 * time.Second)
			}
		}()
	}

	// ====================================
	// REST HTTP SERVER
	// ====================================
	if cfg.REST.Enabled {
		addr := cfg.REST.Address

		log.Printf("Starting REST server on %s", addr)

		go func() {
			srv := rest.NewServer(addr, handlers, nil)
			log.Fatal(srv.ListenAndServe())
		}()
	}

	// --------------------------------
	// Block forever
	// --------------------------------
	select {}
}
