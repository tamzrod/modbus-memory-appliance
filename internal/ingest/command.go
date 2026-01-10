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
	Memory  string
	Area    Area
	Address uint16

	// Exactly ONE of these must be set
	Bools  []int    // numeric booleans: 0 or 1
	Values []uint16 // input registers
}
