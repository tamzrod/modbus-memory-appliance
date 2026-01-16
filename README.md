# Modbus Memory Appliance (MMA)

A deterministic, minimal, and opinionated Modbus TCP memory core written in Go.

---

## Purpose

MMA is a **raw Modbus memory appliance**, not a framework and not a PLC.

Its single responsibility is:

> **Store and serve raw Modbus memory correctly, deterministically, and safely.**

Everything else — scaling, semantics, control logic, visualization, AI — lives **outside** the core.

---

## What This Is

MMA provides:

* Modbus TCP server (control plane)
* Deterministic in-memory Modbus data model
* Atomic memory updates
* Strict validation
* **Multiple ingest paths (Raw / REST / MQTT)**
* Transport adapters (Modbus / REST / MQTT / Raw TCP)
* Config-driven access control and safety boundaries
* High-throughput, low-overhead operation

MMA can run as:

* a standalone binary
* a containerized service
* an embedded library

---

## What This Is NOT

MMA intentionally does **not**:

* Parse float32 / float64
* Scale values or apply engineering units
* Interpret register meaning
* Embed control logic
* Provide a GUI
* Replace PLCs
* Act as a historian
* Implement OPC UA, IEC 61850, or DNP3

If you need these features, build them **upstream**.

---

## Core Design Principles

### 1. Dumb Core, Smart Edges

The core knows nothing about semantics:

* No scaling
* No engineering units
* No business rules
* No hidden meaning

Adapters and upstream systems decide meaning.

---

### 2. Deterministic Behavior

* Same input → same output
* No background mutation
* No timers altering state
* No hidden heuristics

Determinism is a feature, not an optimization.

---

### 3. Atomic Writes (All Paths)

All writes are **all-or-nothing**:

* Entire payload validated first
* Any invalid value rejects the whole batch
* No partial writes
* No silent truncation

This applies equally to:

* Modbus writes
* REST ingest
* MQTT ingest
* **Raw Ingest**

---

### 4. Adapter Isolation

Protocols are adapters, not dependencies:

* Modbus TCP (control plane)
* REST (device ingest plane)
* MQTT (device ingest plane)
* **Raw Ingest (alignment-only replication path)**

Adapters:

* cannot bypass validation
* cannot corrupt memory
* can fail independently

---

### 5. Raw Memory Only

MMA stores exactly what Modbus defines:

* Bits and registers
* Zero-based addressing
* Bounds-checked access

No reinterpretation. No promotion. No assumptions.

---

## Supported Modbus Areas

| Area              | Type         |
| ----------------- | ------------ |
| Holding Registers | uint16       |
| Input Registers   | uint16       |
| Coils             | bool (0 / 1) |
| Discrete Inputs   | bool (0 / 1) |

* All addressing is **zero-based**
* Out-of-range access is rejected

---

## Memory Model

* Memory is pre-allocated at startup
* Fixed size per area
* No dynamic resizing
* Thread-safe read/write locking

This guarantees:

* Predictable latency
* No fragmentation
* No runtime surprises

---

## Ingest Paths (Explicit Separation)

MMA supports **three distinct write paths**, each with a fixed role.

### 1. Modbus TCP — Control Plane

* Client-driven
* Subject to unit-ID, function-code, and port policy
* Used for **intent and control**
* Never bypasses safety rules

---

### 2. REST / MQTT — Device Ingest Plane

* Used by gateways, simulators, edge applications
* Canonical JSON schema (shared)
* Full validation before write
* Atomic batch semantics

Used when **meaning exists upstream**.

Example:

```json
{
  "memory": "plant_a",
  "area": "holding_registers",
  "address": 0,
  "values": [100, 200, 300]
}
```

Failure behavior:

* Any invalid value rejects the entire request
* No partial writes
* Memory remains unchanged

---

### 3. Raw Ingest — Alignment-Only Replication Path

> **Raw Ingest is a direct, alignment-only TCP write path into MMA memory.**

Raw Ingest is **not semantic ingest**.

It exists solely for **data replication** from weak or high-volume devices.

#### Raw Ingest Characteristics

* Stateless
* Write-only
* Panic-free
* Bounds-checked only
* One packet = one atomic write
* No retries
* No timers
* No freshness logic
* No control intent

Raw Ingest performs **alignment, not decode**:

* bytes → uint16 registers
* bytes → bit arrays
* sequential write starting at address

If Raw Ingest understands *meaning*, it is a bug.

---

## Raw Ingest — Packet Sample (Wire-Level)

This section exists to remove ambiguity for implementers.

### Packet Layout (Version 1)

All multi-byte fields are **Big-Endian**.

```
[ Magic(2) ][ Ver(1) ][ Flags(1) ]
[ Area(1) ][ Rsv(1) ][ MemoryID(2) ]
[ Address(2) ][ Count(2) ]
[ Payload(N) ]
[ CRC32(4) ]
```

### Area Values

| Value | Area              |
| ----: | ----------------- |
|  0x01 | Coils             |
|  0x02 | Discrete Inputs   |
|  0x03 | Holding Registers |
|  0x04 | Input Registers   |

---

### Example — Write Holding Registers

**Intent:** write `[100, 200, 300]` starting at address `0`

**Payload:**

```
00 64  00 C8  01 2C
```

**Packet (hex, CRC omitted):**

```
52 49 01 00 03 00 00 01 00 00 00 03
00 64 00 C8 01 2C
```

One packet → one atomic write.

---

### Example — Write Discrete Inputs

**Intent:** write `[1,0,1,1,0]` starting at address `0`

**Payload (bit-packed):**

```
0001101b → 0x0D
```

**Packet (hex, CRC omitted):**

```
52 49 01 00 02 00 00 01 00 00 00 05
0D
```

---

### Rejection Rule

If **any** of the following fail:

* bounds check
* payload length
* CRC
* structure

Result:

* response = `0x01`
* **memory remains unchanged**

---

## Security & Safety Model

### IP Filtering (Per Port)

* Allow-list based
* Config-driven
* IPv4 and IPv6 aware
* Enforced at TCP accept layer

---

### Access Control

Per-port policy supports:

* Unit ID filtering
* Memory selection
* Function code allow/deny
* Read-only or read-write modes

Policy is enforced **before memory access**.

Raw Ingest is **socket-isolated** and never shares control paths.

---

## Configuration Philosophy

* Code defines **capability**
* Config defines **behavior**

No runtime mutation.
No hidden defaults.
Restart to change behavior.

---

## Status

* Stable deterministic core
* Config-driven safety model
* **Raw Ingest formally supported**
* Production-safe architecture

Future protocols may be added **only as adapters**.
