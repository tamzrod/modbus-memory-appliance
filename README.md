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

* Modbus TCP server
* Deterministic in-memory Modbus data model
* Atomic memory updates
* Strict validation
* Transport adapters (Modbus / REST / MQTT)
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

### 3. Atomic Writes

All writes are **all-or-nothing**:

* Entire payload validated first
* Any invalid value rejects the whole batch
* No partial writes
* No silent truncation

This applies equally to:

* Modbus writes
* REST ingest
* MQTT ingest

---

### 4. Adapter Isolation

Protocols are adapters, not dependencies:

* Modbus TCP
* REST
* MQTT

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

## Why Zero-Based Addressing

Internally, MMA uses zero-based addressing **exclusively**.

Reasons:

* Eliminates off-by-one ambiguity
* Matches Modbus PDU semantics
* Simplifies validation and bounds checking
* Avoids vendor-specific address offsets

If a client uses 1-based notation, it must translate upstream.

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

## Ingest (REST / MQTT)

MMA supports **atomic JSON ingest** using the same schema for REST and MQTT.

### Properties

* Same JSON format for REST and MQTT
* Full validation before write
* Batch atomicity
* Deterministic failure behavior

---

## REST Ingest Examples

### HTTP Endpoint

```
POST /ingest
Authorization: Bearer <TOKEN>
Content-Type: application/json
```

### Example: Write Holding Registers

```json
{
  "memory": "plant_a",
  "area": "holding_registers",
  "address": 0,
  "values": [100, 200, 300]
}
```

### Example: Write Coils

```json
{
  "memory": "plant_a",
  "area": "coils",
  "address": 10,
  "values": [1, 0, 1, 1]
}
```

### REST Failure Behavior

* Any invalid value rejects the entire request
* No partial writes
* HTTP error returned
* Memory remains unchanged

---

## MQTT Ingest Examples

### Topic

```
mqtt/ingest
```

### Example: Write Input Registers

```json
{
  "memory": "plant_a",
  "area": "input_registers",
  "address": 50,
  "values": [1234, 5678]
}
```

### Example: Write Discrete Inputs

```json
{
  "memory": "plant_a",
  "area": "discrete_inputs",
  "address": 0,
  "values": [0, 1, 1, 0]
}
```

### MQTT Failure Behavior

* Payload is fully validated before write
* Invalid payload is rejected
* No partial writes
* Memory remains unchanged

---

### Example: Holding Registers Ingest

```json
{
  "memory": "plant_a",
  "area": "holding_registers",
  "address": 0,
  "values": [100, 200, 300]
}
```

If **any** value is invalid:

* nothing is written
* memory remains unchanged

---

## Security & Safety Model

### IP Filtering (Per Port)

* Allow-list based
* Config-driven
* IPv4 and IPv6 aware
* Enforced at TCP accept layer

If a client is not allowed:

* connection is closed immediately
* no Modbus response is sent

---

### Access Control

Per-port policy supports:

* Unit ID filtering
* Memory selection
* Function code allow/deny
* Read-only or read-write modes

Policy is enforced **before memory access**.

---

## Failure Behavior (Explicit Guarantees)

MMA guarantees:

* Invalid input never mutates memory
* Adapter failures do not crash the core
* Memory corruption is structurally impossible
* Illegal access returns protocol-correct errors

Failure is **visible and deterministic**.

---

## Deployment Models

### Standalone Binary

* Single static executable
* No runtime dependencies
* Suitable for edge devices

### Docker / Container

* One process per container
* Stateless except memory
* Easy horizontal replication

### Embedded Library

* Core can be compiled into another system
* External system controls lifecycle

---

## Configuration Philosophy

* Code defines **capability**
* Config defines **behavior**

No runtime mutation.
No hidden defaults.
Restart to change behavior.

---

## Where MMA Fits Best

* Power plant controllers
* Microgrid controllers
* Industrial gateways
* Protocol translation hubs
* Deterministic ingest buffers
* Simulation and test harnesses

---

## What MMA Is Not

MMA is not:

* a PLC
* a SCADA
* a rules engine
* an AI controller

It is a **foundational primitive**, not a finished product.

---

## Philosophy

Industrial systems fail when:

* semantics leak into transport
* logic and memory are coupled
* safety depends on convention
* systems attempt to be clever

MMA exists to do one thing correctly:

> **Be a boring, deterministic, trustworthy Modbus memory appliance.**

Smart systems belong upstream.
The core should never be smart.

---

## Status

* Stable deterministic core
* Config-driven safety model
* Production-safe architecture

Future protocols may be added **only as adapters**.
