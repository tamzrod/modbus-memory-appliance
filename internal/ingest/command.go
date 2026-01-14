// internal/ingest/command.go
package ingest

type Area string

const (
	Coils          Area = "coils"
	DiscreteInputs Area = "discrete_inputs"
	HoldingRegs    Area = "holding_registers"
	InputRegisters Area = "input_registers"
)

// Command is the single canonical ingestion command.
// Used by REST, MQTT, and any future transport.
type Command struct {
    Memory  string   `json:"memory"`
    Area    Area     `json:"area"`
    Address uint16   `json:"address"`
    Bools   []int    `json:"bools,omitempty"`
    Values  []uint16 `json:"values,omitempty"`
}

