package ingest

import (
	"testing"

	"modbus-memory-appliance/internal/core"
)

// newTestService creates a Service with a memory in PRE-RUN state
// so ingest restore (including coils & holding registers) is allowed.
func newTestService() *Service {
	mem := core.NewMemory(
		10, // coils
		10, // discrete inputs
		10, // holding registers
		10, // input registers
	)

	// 🔑 Force PRE-RUN for tests that expect ingest to succeed
	mem.SetStateSealing(true, 999) // gate address unused in tests

	return New(map[string]*core.Memory{
		"test": mem,
	})
}

func TestIngest_DiscreteInputs(t *testing.T) {
	tests := []struct {
		name    string
		cmd     Command
		wantErr error
	}{
		{
			name: "valid numeric booleans",
			cmd: Command{
				Memory:  "test",
				Area:    DiscreteInputs,
				Address: 0,
				Bools:   []int{1, 0, 1},
			},
			wantErr: nil,
		},
		{
			name: "invalid boolean value",
			cmd: Command{
				Memory:  "test",
				Area:    DiscreteInputs,
				Address: 0,
				Bools:   []int{1, 2},
			},
			wantErr: ErrInvalidBoolean,
		},
		{
			name: "missing bools payload",
			cmd: Command{
				Memory:  "test",
				Area:    DiscreteInputs,
				Address: 0,
			},
			wantErr: ErrInvalidPayload,
		},
		{
			name: "payload mismatch",
			cmd: Command{
				Memory:  "test",
				Area:    DiscreteInputs,
				Address: 0,
				Values:  []uint16{1},
			},
			wantErr: ErrPayloadMismatch,
		},
	}

	svc := newTestService()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.Ingest(tt.cmd)

			if tt.wantErr == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.wantErr != nil && err != tt.wantErr {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestIngest_InputRegisters(t *testing.T) {
	tests := []struct {
		name    string
		cmd     Command
		wantErr error
	}{
		{
			name: "valid input registers",
			cmd: Command{
				Memory:  "test",
				Area:    InputRegisters,
				Address: 2,
				Values:  []uint16{100, 200},
			},
			wantErr: nil,
		},
		{
			name: "missing values payload",
			cmd: Command{
				Memory:  "test",
				Area:    InputRegisters,
				Address: 0,
			},
			wantErr: ErrInvalidPayload,
		},
		{
			name: "payload mismatch",
			cmd: Command{
				Memory:  "test",
				Area:    InputRegisters,
				Address: 0,
				Bools:   []int{1},
			},
			wantErr: ErrPayloadMismatch,
		},
	}

	svc := newTestService()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.Ingest(tt.cmd)

			if tt.wantErr == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.wantErr != nil && err != tt.wantErr {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestIngest_Guards(t *testing.T) {
	svc := newTestService()

	tests := []struct {
		name string
		cmd  Command
		err  error
	}{
		{
			name: "unknown memory",
			cmd: Command{
				Memory: "nope",
				Area:   DiscreteInputs,
				Bools:  []int{1},
			},
			err: ErrUnknownMemory,
		},
		{
			name: "both payloads provided",
			cmd: Command{
				Memory: "test",
				Area:   DiscreteInputs,
				Bools:  []int{1},
				Values: []uint16{1},
			},
			err: ErrInvalidPayload,
		},
		{
			name: "invalid area",
			cmd: Command{
				Memory: "test",
				Area:   Area("invalid_area"),
				Bools:  []int{1},
			},
			err: ErrInvalidArea,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := svc.Ingest(tt.cmd); err != tt.err {
				t.Fatalf("expected %v, got %v", tt.err, err)
			}
		})
	}
}

func TestIngest_CoilsAndHoldingRegisters(t *testing.T) {
	svc := newTestService()

	tests := []struct {
		name string
		cmd  Command
		err  error
	}{
		{
			name: "write coils",
			cmd: Command{
				Memory:  "test",
				Area:    Coils,
				Address: 0,
				Bools:   []int{1, 0, 1},
			},
			err: nil,
		},
		{
			name: "write holding registers",
			cmd: Command{
				Memory:  "test",
				Area:    HoldingRegs,
				Address: 1,
				Values:  []uint16{100, 200},
			},
			err: nil,
		},
		{
			name: "coil invalid boolean",
			cmd: Command{
				Memory:  "test",
				Area:    Coils,
				Address: 0,
				Bools:   []int{2},
			},
			err: ErrInvalidBoolean,
		},
		{
			name: "holding regs wrong payload",
			cmd: Command{
				Memory:  "test",
				Area:    HoldingRegs,
				Address: 0,
				Bools:   []int{1},
			},
			err: ErrPayloadMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := svc.Ingest(tt.cmd); err != tt.err {
				t.Fatalf("expected %v, got %v", tt.err, err)
			}
		})
	}
}
