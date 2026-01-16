// internal/rawingest/header_v1.go
// PURPOSE: Parse and validate Raw Ingest V1 fixed header.
// ALLOWED: byte parsing, basic validation
// FORBIDDEN: sockets, CRC, memory access, side effects

package rawingest

import "encoding/binary"

// v1Header represents the decoded fixed header.
type v1Header struct {
	Version uint8
	Flags   uint8
	Area    MemoryArea
	MemID   uint16
	Address uint16
	Count   uint16
}

// parseV1Header parses and validates the fixed-size V1 header.
// The caller MUST ensure buf length >= v1HeaderSize.
func parseV1Header(buf []byte) (v1Header, error) {
	if len(buf) < v1HeaderSize {
		return v1Header{}, ErrRejected
	}

	// Magic check
	if buf[0] != magic0 || buf[1] != magic1 {
		return v1Header{}, ErrRejected
	}

	// Version check
	if buf[2] != protocolVersionV1 {
		return v1Header{}, ErrRejected
	}

	h := v1Header{
		Version: buf[2],
		Flags:   buf[3],
		Area:    MemoryArea(buf[4]),
		MemID:   binary.BigEndian.Uint16(buf[6:8]),
		Address: binary.BigEndian.Uint16(buf[8:10]),
		Count:   binary.BigEndian.Uint16(buf[10:12]),
	}

	// Area validation
	switch h.Area {
	case AreaCoils, AreaDiscreteInputs, AreaHoldingRegs, AreaInputRegs:
		// valid
	default:
		return v1Header{}, ErrRejected
	}

	// Count must be non-zero
	if h.Count == 0 {
		return v1Header{}, ErrRejected
	}

	return h, nil
}
