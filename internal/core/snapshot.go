package core

type Snapshot struct {
	Coils          []bool
	DiscreteInputs []bool
	HoldingRegs    []uint16
	InputRegs      []uint16
}

func (m *Memory) Snapshot() Snapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return Snapshot{
		Coils:          append([]bool(nil), m.Coils...),
		DiscreteInputs: append([]bool(nil), m.DiscreteInputs...),
		HoldingRegs:    append([]uint16(nil), m.HoldingRegs...),
		InputRegs:      append([]uint16(nil), m.InputRegs...),
	}
}
