package config

// AppConfig is the root runtime configuration.
// There must be exactly ONE definition of this type.
type AppConfig struct {
	// Core
	Memory  MemoryConfig
	Routing RoutingConfig
	Ports   Ports

	// Ingest / control plane
	REST RESTConfig
	MQTT MQTTConfig
}
