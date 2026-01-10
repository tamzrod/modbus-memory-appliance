package mqtt

type Config struct {
	Enabled  bool
	Broker   string
	ClientID string
	Username string
	Password string
	Topic    string // e.g. "mqtt/ingest"
}
