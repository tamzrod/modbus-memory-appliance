package config

type AppConfig struct {
	Memory   MemoryConfig
	Ports    map[uint16]PortPolicy
	Routing  RoutingConfig

	MQTT     MQTTConfig   // 👈 ADD THIS LINE
}
