# raw_ingest Transport Specification

## Overview

**raw_ingest** is a private, binary transport used between **Mirror processes** and the **Modbus Memory Appliance (MMA)**.

Its sole responsibility is to perform **atomic, contiguous memory writes** into MMA memory using **raw byte buffers**.

This transport is **internal, trusted, and point-to-point**.

---

## Design Principles

- Binary, not text
- Fixed header, fixed order
- No interpretation of values
- No retries
- No semantics
- One packet = one atomic write
- Optimized for memcpy
- Optimized for Kubernetes / horizontal scaling

---

## Architecture Position

```
┌──────────────────────────────────────────────┐
│            External Field Domain             │
│                                              │
│   [ Field Device ]                            │
│          ↓                                   │
│   [ Mirror (polling, retries, policy) ]      │
│          ↓                                   │
│   ( raw_ingest binary packets )              │
│                                              │
└──────────────────────────────────────────────┘
                    │
                    │ TCP
                    ▼
┌──────────────────────────────────────────────┐
│        Modbus Memory Appliance (MMA)          │
│                                              │
│   [ raw_ingest transport ]                   │
│          ↓                                   │
│   [ Memory Core (authoritative state) ]      │
│          ↓                                   │
│   [ Modbus TCP servers ]                     │
│                                              │
└──────────────────────────────────────────────┘
```

**Mirrors and field devices are NOT part of MMA.**

---

## Non-Goals

raw_ingest does **NOT**:

- poll devices
- retry writes
- buffer data
- scale values
- interpret bits or registers
- apply quality or timestamps
- generate defaults
- create or resize memory

---

## Transport

- **Protocol:** TCP
- **Direction:** Inbound only
- **Connection:** Server (MMA listens)
- **Security:** Internal network only
- **Exposure:** Never exposed to SCADA

---

## Packet Model

```
[ FIXED HEADER (8 bytes) ]
[ RAW DATA BUFFER (N bytes) ]
```

---

## Fixed Header (8 bytes)

| Offset | Size | Field | Type | Description |
|------|------|------|------|------------|
| 0 | 2 | memory_id | uint16 | Target memory identifier |
| 2 | 1 | area | uint8 | Memory area selector |
| 3 | 1 | format | uint8 | Data encoding |
| 4 | 2 | start | uint16 | Start offset (element index) |
| 6 | 2 | count | uint16 | Number of elements |

**All integers are big-endian.**

---

## Area Codes

```
00 = coils
01 = discrete_inputs
02 = holding_registers
03 = input_registers
```

---

## Format Codes

```
00 = bit   (packed bits, Modbus LSB-first)
01 = u16   (unsigned 16-bit registers, big-endian)
```

---

## Data Buffer Rules

### For `bit` format
- Payload size = `ceil(count / 8)` bytes
- Bit order: **LSB-first**, Modbus compatible
- Unused bits in final byte are ignored

### For `u16` format
- Payload size = `count × 2` bytes
- Big-endian encoding
- No scaling or conversion

---

## Atomicity Rule

> **One packet = one atomic memory write**

- Either all bytes are written
- Or nothing is written
- No partial commits
- One lock per packet

---

## Validation Rules (MMA Side)

Before writing, raw_ingest **MUST** validate:

1. `memory_id` exists
2. `area` exists
3. `format` matches area
4. `start + count` is within bounds
5. Data length matches expected size
6. Packet size ≤ configured maximum

Failure → **reject packet, write nothing**

---

## ACK / Response

### Success
```
0x00
```

### Error
```
0x01
```

---

## Packet Size Limits

- **Optimal:** ≤ 1024 bytes
- **Hard limit:** 2048 bytes
- **Typical usage:** 200–500 bytes

---

## Memory ID Mapping

Memory IDs are mapped via configuration on **both sides**.

```yaml
memory_map:
  1: mvps01_scb01
  2: mvps01_scb02
  3: mvps18_scb03
```

---

## Mental Model

> **raw_ingest is DMA over TCP.**

- Header = address
- Payload = bytes
- MMA copies memory
- No opinions
- No drama

---

## Summary

- Binary, fixed, minimal protocol
- Modbus-style efficiency without Modbus semantics
- One packet = one contiguous write
- Mirror owns policy
- MMA owns memory
- SCADA sees calm, stable data
