package config

// Root of the memory configuration
type MemoryConfig struct {
	Memories map[string]MemoryBlock `yaml:"memories"`
}

// One named memory (e.g. "memory1")
type MemoryBlock struct {
	Default bool `yaml:"default"`
	Coils            AreaConfig `yaml:"coils"`
	DiscreteInputs   AreaConfig `yaml:"discrete_inputs"`
	HoldingRegisters AreaConfig `yaml:"holding_registers"`
	InputRegisters   AreaConfig `yaml:"input_registers"`
}

// One Modbus address space definition
type AreaConfig struct {
	Start int `yaml:"start"`
	Size  int `yaml:"size"`
}
type UnitIDMap map[uint8]string

