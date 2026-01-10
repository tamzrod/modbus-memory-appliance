package mqtt

import (
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"modbus-memory-appliance/internal/ingest"
)

/*
Subscriber is a thin MQTT → ingest adapter.
It owns:
- MQTT connection lifecycle
- topic subscription
- message decoding

It does NOT:
- parse semantics
- touch memory directly
- expose control endpoints
*/

type Subscriber struct {
	client mqtt.Client
	ingest *ingest.Service
	topic  string
}

// NewSubscriber creates and maintains an MQTT ingest subscriber.
// It blocks only during connect & subscribe.
// Runtime reconnects are handled internally by the client.
func NewSubscriber(cfg Config, ingestSvc *ingest.Service) (*Subscriber, error) {
	// pessimistic default
	setConnected(false)

	opts := mqtt.NewClientOptions().
		AddBroker(cfg.Broker).
		SetClientID(cfg.ClientID)

	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
		opts.SetPassword(cfg.Password)
	}

	// Connection lifecycle hooks
	opts.SetOnConnectHandler(func(_ mqtt.Client) {
		log.Printf("mqtt connected")
		setConnected(true)
	})

	opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		log.Printf("mqtt connection lost: %v", err)
		setConnected(false)
	})

	client := mqtt.NewClient(opts)

	// Initial connect (blocking)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		setConnected(false)
		return nil, token.Error()
	}

	s := &Subscriber{
		client: client,
		ingest: ingestSvc,
		topic:  cfg.Topic,
	}

	// Subscribe to ingest topic
	if token := client.Subscribe(cfg.Topic, 0, s.handleMessage); token.Wait() && token.Error() != nil {
		setConnected(false)
		return nil, token.Error()
	}

	log.Printf("mqtt ingest subscribed: %s", cfg.Topic)

	return s, nil
}

// handleMessage converts MQTT payload → ingest.Command
// No retries, no buffering, no semantics.
func (s *Subscriber) handleMessage(_ mqtt.Client, msg mqtt.Message) {
	var cmd ingest.Command

	if err := json.Unmarshal(msg.Payload(), &cmd); err != nil {
		log.Printf("mqtt ingest: invalid json: %v", err)
		return
	}

	if err := s.ingest.Ingest(cmd); err != nil {
		log.Printf("mqtt ingest rejected: %v", err)
		return
	}
}
