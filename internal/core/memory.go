package core

// ----------- COILS -----------

func (m *Memory) ReadCoils(addr, count int) ([]bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if err := checkBoolRange(m.Coils, addr, count); err != nil {
		return nil, err
	}

	out := make([]bool, count)
	copy(out, m.Coils[addr:addr+count])
	return out, nil
}

func (m *Memory) WriteCoils(addr int, values []bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := checkBoolRange(m.Coils, addr, len(values)); err != nil {
		return err
	}

	copy(m.Coils[addr:], values)
	return nil
}

// ----------- DISCRETE INPUTS -----------

func (m *Memory) ReadDiscreteInputs(addr, count int) ([]bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if err := checkBoolRange(m.DiscreteInputs, addr, count); err != nil {
		return nil, err
	}

	out := make([]bool, count)
	copy(out, m.DiscreteInputs[addr:addr+count])
	return out, nil
}
// write digital

// ----------- DISCRETE INPUTS -----------

func (m *Memory) WriteDiscreteInputs(addr int, values []bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := checkBoolRange(m.DiscreteInputs, addr, len(values)); err != nil {
		return err
	}

	copy(m.DiscreteInputs[addr:], values)
	return nil
}


// ----------- HOLDING REGISTERS -----------

func (m *Memory) ReadHoldingRegs(addr, count int) ([]uint16, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if err := checkUint16Range(m.HoldingRegs, addr, count); err != nil {
		return nil, err
	}

	out := make([]uint16, count)
	copy(out, m.HoldingRegs[addr:addr+count])
	return out, nil
}

func (m *Memory) WriteHoldingRegs(addr int, values []uint16) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := checkUint16Range(m.HoldingRegs, addr, len(values)); err != nil {
		return err
	}

	copy(m.HoldingRegs[addr:], values)
	return nil
}

// ----------- INPUT REGISTERS -----------

func (m *Memory) ReadInputRegs(addr, count int) ([]uint16, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if err := checkUint16Range(m.InputRegs, addr, count); err != nil {
		return nil, err
	}

	out := make([]uint16, count)
	copy(out, m.InputRegs[addr:addr+count])
	return out, nil
}

func NewMemory(
	coilCount,
	discreteCount,
	holdingCount,
	inputCount int,
) *Memory {
	return &Memory{
		Coils:          make([]bool, coilCount),
		DiscreteInputs: make([]bool, discreteCount),
		HoldingRegs:    make([]uint16, holdingCount),
		InputRegs:      make([]uint16, inputCount),
	}
}

// ----------- INPUT REGISTERS -----------

func (m *Memory) WriteInputRegs(addr int, values []uint16) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := checkUint16Range(m.InputRegs, addr, len(values)); err != nil {
		return err
	}

	copy(m.InputRegs[addr:], values)
	return nil
}
