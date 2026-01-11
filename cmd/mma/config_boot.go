package main

import (
	"log"

	"modbus-memory-appliance/internal/config"
)

func loadConfig() *config.AppConfig {
	cfg, err := config.LoadOrCreate()
	if err != nil {
		log.Fatal(err)
	}

	if err := cfg.ValidatePorts(); err != nil {
		log.Fatal(err)
	}

	return cfg // ✅ NOT &cfg
}
