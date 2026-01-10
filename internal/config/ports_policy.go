package config

// AllowsUnitID returns true if the given unit ID is allowed by the port policy.
func (p PortPolicy) AllowsUnitID(unitID uint8) bool {
	if p.UnitIDs.All {
		return true
	}

	for _, uid := range p.UnitIDs.List {
		if uid == unitID {
			return true
		}
	}
	return false
}

// AllowsMemory returns true if the given memory ID is allowed by the port policy.
func (p PortPolicy) AllowsMemory(memoryID string) bool {
	if p.Memories.All {
		return true
	}

	for _, mem := range p.Memories.List {
		if mem == memoryID {
			return true
		}
	}
	return false
}

// AllowsFunctionCode returns true if the function code is allowed by policy.
// NOTE: This does NOT enforce Modbus semantics yet.
// That happens later in the resolver.
func (p PortPolicy) AllowsFunctionCode(fc uint8) bool {
	// Step 1: access mode gate
	if p.Access == AccessReadOnly {
		switch fc {
		case 1, 2, 3, 4:
			// allowed reads
		default:
			return false
		}
	}

	// Step 2: optional function code overrides
	if p.FunctionCodes == nil {
		return true
	}

	// Deny list wins
	for _, deny := range p.FunctionCodes.Deny {
		if deny == fc {
			return false
		}
	}

	// Allow list (if present)
	if len(p.FunctionCodes.Allow) > 0 {
		for _, allow := range p.FunctionCodes.Allow {
			if allow == fc {
				return true
			}
		}
		return false
	}

	return true
}
