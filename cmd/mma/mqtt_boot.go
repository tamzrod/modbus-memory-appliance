package main

import (
	"log"
	"time"

	"modbus-memory-appliance/internal/config"
	"modbus-memory-appliance/internal/ingest"
	"modbus-memory-appliance/internal/mqtt"
)

func startMQTT(cfg *config.AppConfig, ingestSvc *ingest.Service) {
	if !cfg.MQTT.Enabled {
		return
	}

	mqttCfg := mqtt.Config{
		Enabled:  cfg.MQTT.Enabled,
		Broker:   cfg.MQTT.Broker,
		ClientID: cfg.MQTT.ClientID,
		Topic:    cfg.MQTT.Topic,
		Username: cfg.MQTT.Username,
		Password: cfg.MQTT.Password,
	}

	go func() {
		for {
			log.Printf("Starting MQTT subscriber (broker=%s)", mqttCfg.Broker)

			_, err := mqtt.NewSubscriber(mqttCfg, ingestSvc)
			if err != nil {
				log.Printf("MQTT unavailable: %v", err)
				time.Sleep(10 * time.Second)
				continue
			}

			log.Printf("MQTT subscriber exited, restarting...")
			time.Sleep(5 * time.Second)
		}
	}()
}
