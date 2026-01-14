package modbus

import (
	"encoding/binary"
	"net"

	"modbus-memory-appliance/internal/core"
)

// handleConn handles a single Modbus TCP connection.
func handleConn(conn net.Conn, resolve MemoryResolver) {
	defer conn.Close()

	for {
		mbap, err := readMBAP(conn)
		if err != nil {
			return
		}

		pdu, err := readPDU(conn, mbap.Length)
		if err != nil {
			return
		}

		// Resolve memory (routing only)
		mem := resolve(mbap.UnitID, pdu.Function)

		if mem == nil {
			resp := exception(pdu.Function, 0x02) // Illegal Data Address
			mbap.Length = uint16(len(resp) + 1)
			writeMBAP(conn, mbap)
			conn.Write(resp)
			continue
		}

		// 🔒 STATE SEALING HARD GATE
		// If memory is Pre-Run, ALL Modbus access is denied
		if mem.IsPreRun() {
			resp := exception(pdu.Function, 0x01) // Illegal Function
			mbap.Length = uint16(len(resp) + 1)
			writeMBAP(conn, mbap)
			conn.Write(resp)
			continue
		}

		resp := handlePDU(pdu, mem)

		mbap.Length = uint16(len(resp) + 1)
		writeMBAP(conn, mbap)
		conn.Write(resp)
	}
}

// handlePDU executes a Modbus PDU against a RUN-state memory.
func handlePDU(pdu PDU, mem *core.Memory) []byte {
	switch pdu.Function {

	case 0x03: // Read Holding Registers
		addr := binary.BigEndian.Uint16(pdu.Data[0:2])
		count := binary.BigEndian.Uint16(pdu.Data[2:4])

		values, err := mem.ReadHoldingRegs(
			extToInternal(addr),
			int(count),
		)
		if err != nil {
			return exception(pdu.Function, 0x02)
		}

		out := make([]byte, 2+len(values)*2)
		out[0] = pdu.Function
		out[1] = uint8(len(values) * 2)

		for i, v := range values {
			binary.BigEndian.PutUint16(out[2+i*2:], v)
		}
		return out

	case 0x06: // Write Single Register
		addr := binary.BigEndian.Uint16(pdu.Data[0:2])
		val := binary.BigEndian.Uint16(pdu.Data[2:4])

		if err := mem.WriteHoldingRegs(
			extToInternal(addr),
			[]uint16{val},
		); err != nil {
			return exception(pdu.Function, 0x02)
		}

		return append([]byte{pdu.Function}, pdu.Data...)

	case 0x04: // Read Input Registers
		addr := binary.BigEndian.Uint16(pdu.Data[0:2])
		count := binary.BigEndian.Uint16(pdu.Data[2:4])

		values, err := mem.ReadInputRegs(
			extToInternal(addr),
			int(count),
		)
		if err != nil {
			return exception(pdu.Function, 0x02)
		}

		out := make([]byte, 2+len(values)*2)
		out[0] = pdu.Function
		out[1] = uint8(len(values) * 2)

		for i, v := range values {
			binary.BigEndian.PutUint16(out[2+i*2:], v)
		}
		return out

	case 0x01: // Read Coils
		addr := binary.BigEndian.Uint16(pdu.Data[0:2])
		count := binary.BigEndian.Uint16(pdu.Data[2:4])

		values, err := mem.ReadCoils(
			extToInternal(addr),
			int(count),
		)
		if err != nil {
			return exception(pdu.Function, 0x02)
		}

		data := packBools(values)

		out := make([]byte, 2+len(data))
		out[0] = pdu.Function
		out[1] = uint8(len(data))
		copy(out[2:], data)
		return out

	case 0x02: // Read Discrete Inputs
		addr := binary.BigEndian.Uint16(pdu.Data[0:2])
		count := binary.BigEndian.Uint16(pdu.Data[2:4])

		values, err := mem.ReadDiscreteInputs(
			extToInternal(addr),
			int(count),
		)
		if err != nil {
			return exception(pdu.Function, 0x02)
		}

		data := packBools(values)

		out := make([]byte, 2+len(data))
		out[0] = pdu.Function
		out[1] = uint8(len(data))
		copy(out[2:], data)
		return out

	case 0x05: // Write Single Coil
		addr := binary.BigEndian.Uint16(pdu.Data[0:2])
		val := binary.BigEndian.Uint16(pdu.Data[2:4])

		var b bool
		switch val {
		case 0xFF00:
			b = true
		case 0x0000:
			b = false
		default:
			return exception(pdu.Function, 0x03)
		}

		if err := mem.WriteCoils(
			extToInternal(addr),
			[]bool{b},
		); err != nil {
			return exception(pdu.Function, 0x02)
		}

		return append([]byte{pdu.Function}, pdu.Data...)

	case 0x0F: // Write Multiple Coils
		addr := binary.BigEndian.Uint16(pdu.Data[0:2])
		count := binary.BigEndian.Uint16(pdu.Data[2:4])
		byteCount := int(pdu.Data[4])

		if byteCount != (int(count)+7)/8 {
			return exception(pdu.Function, 0x03)
		}

		values := unpackBools(pdu.Data[5:], int(count))

		if err := mem.WriteCoils(
			extToInternal(addr),
			values,
		); err != nil {
			return exception(pdu.Function, 0x02)
		}

		out := make([]byte, 5)
		out[0] = pdu.Function
		binary.BigEndian.PutUint16(out[1:3], addr)
		binary.BigEndian.PutUint16(out[3:5], count)
		return out

	case 0x10: // Write Multiple Holding Registers
		addr := binary.BigEndian.Uint16(pdu.Data[0:2])
		count := binary.BigEndian.Uint16(pdu.Data[2:4])
		byteCount := int(pdu.Data[4])

		if byteCount != int(count)*2 {
			return exception(pdu.Function, 0x03)
		}

		values := make([]uint16, count)
		for i := 0; i < int(count); i++ {
			values[i] = binary.BigEndian.Uint16(pdu.Data[5+i*2:])
		}

		if err := mem.WriteHoldingRegs(
			extToInternal(addr),
			values,
		); err != nil {
			return exception(pdu.Function, 0x02)
		}

		out := make([]byte, 5)
		out[0] = pdu.Function
		binary.BigEndian.PutUint16(out[1:3], addr)
		binary.BigEndian.PutUint16(out[3:5], count)
		return out

	default:
		return exception(pdu.Function, 0x01)
	}
}

// exception builds a Modbus exception response.
func exception(fc uint8, code uint8) []byte {
	return []byte{fc | 0x80, code}
}

func packBools(values []bool) []byte {
	byteCount := (len(values) + 7) / 8
	out := make([]byte, byteCount)

	for i, v := range values {
		if v {
			out[i/8] |= 1 << (i % 8)
		}
	}
	return out
}

func unpackBools(data []byte, count int) []bool {
	out := make([]bool, count)
	for i := 0; i < count; i++ {
		if data[i/8]&(1<<(i%8)) != 0 {
			out[i] = true
		}
	}
	return out
}
