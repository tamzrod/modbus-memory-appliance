package config

// RawIngestConfig configures the Raw Ingest core transport.
// This is a memory-write transport (no semantics).
type RawIngestConfig struct {
	Enabled        bool   `yaml:"enabled"`
	Listen         string `yaml:"listen"`
	MaxPacketBytes int    `yaml:"max_packet_bytes"`
	ReadTimeoutMS  int    `yaml:"read_timeout_ms"`
}
