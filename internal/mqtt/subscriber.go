package mqtt

import (
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"modbus-memory-appliance/internal/ingest"
)

type Subscriber struct {
	client mqtt.Client
	ingest  *ingest.Service
	topic   string
}

func NewSubscriber(cfg Config, ingestSvc *ingest.Service) (*Subscriber, error) {
	opts := mqtt.NewClientOptions().
		AddBroker(cfg.Broker).
		SetClientID(cfg.ClientID)

	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
		opts.SetPassword(cfg.Password)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	s := &Subscriber{
		client: client,
		ingest:  ingestSvc,
		topic:   cfg.Topic,
	}

	if token := client.Subscribe(cfg.Topic, 0, s.handleMessage); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	log.Printf("mqtt ingest subscribed: %s", cfg.Topic)

	return s, nil
}

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
