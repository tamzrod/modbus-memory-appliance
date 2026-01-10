package mqtt

import "sync/atomic"

type Status struct {
	Enabled   bool   `json:"enabled"`
	Connected bool   `json:"connected"`
	Broker    string `json:"broker"`
	Topic     string `json:"topic"`
}

var connected atomic.Bool

func setConnected(v bool) {
	connected.Store(v)
}

func GetStatus(cfg Config) Status {
	return Status{
		Enabled:   cfg.Enabled,
		Connected: connected.Load(),
		Broker:    cfg.Broker,
		Topic:     cfg.Topic,
	}
}
