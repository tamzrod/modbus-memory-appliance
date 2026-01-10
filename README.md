# Modbus Memory Appliance (MMA)

A **deterministic, minimal, and opinionated Modbus TCP memory core** written in Go.

This project provides a **raw Modbus memory appliance** designed for reliability, correctness, and observability in industrial and control-system environments.

---

## What This Is

MMA is a **memory appliance**, not a framework.

Its job is simple:

> **Store and serve raw Modbus memory correctly, under all conditions.**

It supports:
- Modbus TCP
- Atomic JSON ingest via MQTT and REST
- Strict validation
- Non-fatal adapters
- High-throughput, low-overhead operation

---

## What This Is NOT

This project intentionally does **not**:
- parse float32 / float64
- scale values
- interpret register meaning
- provide a GUI
- embed control logic
- replace PLCs
- implement OPC UA / IEC 61850

If you need those features, build them **upstream**.

---

## Core Principles

- **Dumb core, smart edges**
- **Deterministic behavior**
- **Atomic writes**
- **Adapter isolation**
- **Raw memory only**

In control systems:

> **Boring is reliable.**

---

## Supported Modbus Areas

| Area | Type |
|----|----|
| Holding Registers | uint16 |
| Input Registers | uint16 |
| Coils | bool (0 / 1) |
| Discrete Inputs | bool (0 / 1) |

All addressing is **zero-based**.

---

## Ingest (MQTT / REST)

MMA supports **atomic JSON ingest** for all Modbus areas.

- Same JSON format for MQTT and REST
- Entire batch is validated first
- Any invalid value rejects the entire batch
- No partial writes
- No silent truncation

### Example (Holding Registers)

```json
{
  "memory": "plant_a",
  "area": "holding_registers",
  "address": 0,
  "values": [100, 200, 300]
}
