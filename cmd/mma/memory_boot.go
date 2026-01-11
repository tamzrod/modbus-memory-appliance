package main

import (
	"log"

	"modbus-memory-appliance/internal/config"
	"modbus-memory-appliance/internal/core"
)

func buildMemories(cfg *config.AppConfig) map[string]*core.Memory {
	memories, err := config.BuildMemories(&cfg.Memory)
	if err != nil {
		log.Fatal(err)
	}
	return memories
}
