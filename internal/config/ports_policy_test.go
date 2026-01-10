package config

import "testing"

func TestAllowsUnitID(t *testing.T) {
	p := PortPolicy{
		UnitIDs: UnitIDSelector{All: false, List: []uint8{1, 2}},
	}

	if !p.AllowsUnitID(1) {
		t.Fatal("expected unit 1 allowed")
	}
	if p.AllowsUnitID(3) {
		t.Fatal("expected unit 3 denied")
	}
}

func TestAllowsMemory(t *testing.T) {
	p := PortPolicy{
		Memories: MemorySelector{All: false, List: []string{"plant_a"}},
	}

	if !p.AllowsMemory("plant_a") {
		t.Fatal("expected plant_a allowed")
	}
	if p.AllowsMemory("plant_b") {
		t.Fatal("expected plant_b denied")
	}
}

func TestAllowsFunctionCode_ReadOnly(t *testing.T) {
	p := PortPolicy{
		Access: AccessReadOnly,
	}

	if !p.AllowsFunctionCode(3) { // Read Holding
		t.Fatal("expected FC03 allowed")
	}
	if p.AllowsFunctionCode(6) { // Write Single Register
		t.Fatal("expected FC06 denied")
	}
}

func TestAllowsFunctionCode_AllowList(t *testing.T) {
	p := PortPolicy{
		Access: AccessReadWrite,
		FunctionCodes: &FunctionCodeACL{
			Allow: []uint8{3},
		},
	}

	if !p.AllowsFunctionCode(3) {
		t.Fatal("expected FC03 allowed")
	}
	if p.AllowsFunctionCode(4) {
		t.Fatal("expected FC04 denied")
	}
}
