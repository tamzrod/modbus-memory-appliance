package config

// Ports is a map keyed by TCP port number.
type Ports map[uint16]PortPolicy

// PortPolicy defines access policy for a TCP port.
// Policy only — never routing.
type PortPolicy struct {
	UnitIDs       UnitIDSelector    `yaml:"unit_ids"`
	Memories      MemorySelector    `yaml:"memories"`
	Access        AccessMode        `yaml:"access"`
	FunctionCodes *FunctionCodeACL  `yaml:"function_codes,omitempty"`
}

// AccessMode defines read/write capability.
type AccessMode string

const (
	AccessReadOnly  AccessMode = "read-only"
	AccessReadWrite AccessMode = "read-write"
)

type UnitIDSelector struct {
	All  bool
	List []uint8
}

type MemorySelector struct {
	All  bool
	List []string
}

type FunctionCodeACL struct {
	Allow []uint8 `yaml:"allow,omitempty"`
	Deny  []uint8 `yaml:"deny,omitempty"`
}
