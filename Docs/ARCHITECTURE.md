# Modbus Memory Appliance – Architecture

## Purpose

The Modbus Memory Appliance (MMA) is designed to be a **deterministic, minimal, and opinionated Modbus TCP memory core**.

It is not a framework, not a PLC, and not a data modeler.

Its job is simple:
> **store raw Modbus memory correctly and predictably, under all conditions.**

---

## Core Design Philosophy

### 1. Dumb Core, Smart Edges

The core memory:
- stores raw values only
- enforces bounds and atomicity
- knows nothing about meaning

All intelligence lives **outside**:
- Node-RED
- PLCs
- SCADA
- Control algorithms

This prevents:
- hidden logic
- implicit behavior
- vendor-style magic

---

### 2. Determinism Above All

Every operation has a deterministic outcome:

- valid write → fully applied
- invalid write → fully rejected
- no partial state
- no silent truncation

This applies equally to:
- MQTT
- REST
- Modbus writes

---

## Memory Model

- Zero-based addressing
- Fixed-width storage
- No dynamic resizing at runtime

### Supported Areas

| Area | Storage Type |
|-----|--------------|
| Holding Registers | uint16 |
| Input Registers | uint16 |
| Coils | bool (0 / 1) |
| Discrete Inputs | bool (0 / 1) |

The memory layer does **not** distinguish semantics beyond size and bounds.

---

## Atomic Ingest

All ingest operations are **transactional**:

- Entire batch validated first
- Any invalid value aborts the batch
- Memory remains unchanged on failure

This is intentional.

Partial writes are more dangerous than rejected writes in control systems.

---

## Adapter Model

Adapters are **untrusted**.

Examples:
- MQTT
- REST
- future protocols

Rules:
- adapters may fail
- adapters may disconnect
- adapters may misbehave

**The core must continue running regardless.**

An adapter failure must never:
- crash Modbus
- corrupt memory
- block other adapters

---

## Modbus TCP Role

Modbus TCP is the **authoritative interface** for external devices.

- It reads from core memory
- It writes only where allowed by policy
- It never bypasses validation

The appliance behaves like a well-defined Modbus slave:
- predictable
- boring
- reliable

---

## Configuration Model

- YAML-driven
- Explicit policies
- No implicit defaults beyond safety

Configuration controls:
- memory layout
- unit ID routing
- port access
- function code permissions

Configuration is **runtime-only** and intentionally excluded from Git.

---

## Observability

The appliance emits:
- structured logs to stdout
- no direct dependency on logging backends

Designed for:
- Promtail
- Loki
- Grafana

Errors are logged **only when necessary**.

No log spam on success paths.

---

## Non-Goals (Very Important)

This project intentionally does NOT:

- parse float32 / float64
- scale values
- interpret register meaning
- provide a GUI
- replace PLC logic
- auto-discover devices
- implement OPC UA or IEC 61850

If you need those features:
- build them upstream
- keep the core clean

---

## Why Modbus Still Matters

Modbus remains widely deployed because it is:
- simple
- predictable
- universally supported

This project embraces Modbus’s strengths instead of hiding them.

---

## Stability Contract

The following are considered **stable guarantees**:

- raw memory only
- atomic writes
- adapter isolation
- deterministic behavior

Breaking these would constitute a **major version change**.

---

## Intended Audience

This project is for:
- control engineers
- SCADA integrators
- systems programmers
- anyone who values correctness over convenience

---

## Summary

The Modbus Memory Appliance is designed to be:

- boring
- predictable
- strict
- fast

Because in control systems:

> **boring is reliable.**
