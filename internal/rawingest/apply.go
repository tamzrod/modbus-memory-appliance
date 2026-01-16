// internal/rawingest/apply.go
// PURPOSE: Apply a validated header + payload into a resolved memory.
// ALLOWED: alignment + single atomic write into memory
// FORBIDDEN: sockets, retries, semantics, logging, side effects beyond memory write

package rawingest

// applyPayload aligns payload and writes it into mem.
// Caller guarantees: header is validated, payload length matches payloadSize(area,count).
func applyPayload(mem RawWritableMemory, h v1Header, payload []byte) error {
	switch h.Area {
	case AreaCoils:
		vals, err := alignBits(payload, h.Count)
		if err != nil {
			return ErrRejected
		}
		if err := mem.WriteCoils(h.Address, vals); err != nil {
			return ErrRejected
		}
		return nil

	case AreaDiscreteInputs:
		vals, err := alignBits(payload, h.Count)
		if err != nil {
			return ErrRejected
		}
		if err := mem.WriteDiscreteInputs(h.Address, vals); err != nil {
			return ErrRejected
		}
		return nil

	case AreaHoldingRegs:
		vals, err := alignRegs(payload, h.Count)
		if err != nil {
			return ErrRejected
		}
		if err := mem.WriteHoldingRegisters(h.Address, vals); err != nil {
			return ErrRejected
		}
		return nil

	case AreaInputRegs:
		vals, err := alignRegs(payload, h.Count)
		if err != nil {
			return ErrRejected
		}
		if err := mem.WriteInputRegisters(h.Address, vals); err != nil {
			return ErrRejected
		}
		return nil

	default:
		return ErrRejected
	}
}
