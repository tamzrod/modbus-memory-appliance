package config

type MQTTConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Broker   string `yaml:"broker"`
	ClientID string `yaml:"client_id"`
	Topic    string `yaml:"topic"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type RESTConfig struct {
	Enabled bool   `yaml:"enabled"`
	Address string `yaml:"address"`
}

type AppConfig struct {
	Ports   map[uint16]PortPolicy `yaml:"ports"`
	Memory  MemoryConfig          `yaml:"memory"`
	MQTT    MQTTConfig            `yaml:"mqtt"`
	Routing RoutingConfig         `yaml:"routing"`
	REST    RESTConfig            `yaml:"rest"`
}
