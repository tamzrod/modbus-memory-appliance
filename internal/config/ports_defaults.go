package config

func normalizePortDefaults(p *PortPolicy) {
	if p.MaxConnections <= 0 {
		p.MaxConnections = DefaultMaxConnections
	}
}
