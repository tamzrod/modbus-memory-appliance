# Modbus Memory Appliance — Architecture

## Purpose

The Modbus Memory Appliance (MMA) is designed to be a **deterministic, minimal, and opinionated Modbus TCP memory core**.

It is not a framework, not a PLC, and not a data modeler.

Its single responsibility is:

> **Store raw Modbus memory correctly and predictably, under all conditions.**

---

## Core Design Philosophy

### 1. Dumb Core, Smart Edges

The core memory:

* stores raw values only
* enforces bounds and atomicity
* knows nothing about meaning, time, or control logic

All intelligence lives **outside** the core:

* Node-RED
* PLCs
* SCADA systems
* Control algorithms
* External services

This prevents:

* hidden logic
* implicit behavior
* vendor-style magic

---

### 2. Determinism Above All

Every write operation has a deterministic outcome:

* valid write → fully applied
* invalid write → fully rejected
* no partial state
* no silent truncation

This applies equally to **all write paths**:

* Modbus TCP writes
* REST ingest
* MQTT ingest
* Raw Ingest

---

## Memory Model

* Zero-based addressing
* Fixed-width storage
* No dynamic resizing at runtime

### Supported Memory Areas

| Area              | Storage Type |
| ----------------- | ------------ |
| Holding Registers | uint16       |
| Input Registers   | uint16       |
| Coils             | bool         |
| Discrete Inputs   | bool         |

The memory layer does **not** interpret meaning beyond size and bounds.

---

## Adapter-Based Architecture

MMA is built around **adapter isolation**.

Adapters are:

* independent
* untrusted
* replaceable

They converge **only at the memory layer**.

No adapter may:

* depend on another adapter
* call another adapter
* assume ordering guarantees

---

## Write Paths (Authoritative)

There are **two distinct write modes**, representing **different operational scenarios**.

They must be treated as **mutually exclusive at the data-source level**.

```
FIELD DEVICE (RAW)                     LOGICAL DATA SOURCE
(meter / logger / combiner)            (controller / app / gateway)
        |                                      |
        v                                      v
┌────────────────────┐              ┌────────────────────────┐
│     RAW INGEST     │              │     REST / MQTT        │
│      (TCP)         │              │       INGEST           │
│ • blind writes     │              │ • semantic writes      │
│ • no decoding      │              │ • rules & validation   │
│ • no state         │              │ • control-aware        │
└─────────┬──────────┘              └─────────┬──────────────┘
          |                                      |
          └───────────────┬──────────────────────┘
                          v
               ┌──────────────────────────┐
               │        MMA MEMORY        │
               │  (raw registers only)   │
               └─────────┬────────────────┘
                         v
                      READERS
```

---

## Raw Ingest (Formal Role)

Raw Ingest is a **transport-level write adapter**.

It exists to support:

* weak or resource-limited field devices
* high fan-out read scenarios
* mirroring or buffering use cases
* environments where decoding is wasteful or impossible

Raw Ingest characteristics:

* TCP-based binary protocol
* stateless operation
* blind, alignment-only writes
* bounds checking only
* atomic memory writes

Raw Ingest **must not**:

* decode values
* scale values
* apply rules or policies
* track freshness or timestamps
* infer semantics

The authoritative specification for Raw Ingest is defined in:

* `Docs/RAW_INGEST.md`

---

## Logical Ingest (REST / MQTT)

Logical Ingest is a **semantic write adapter**.

It exists to support:

* validated control inputs
* rule-based writes
* state-aware behavior
* higher-level applications

Logical Ingest:

* validates payloads
* applies policy and permissions
* may reject writes for semantic reasons

---

## Mandatory Separation Rule

A single data source **must not** use Raw Ingest and Logical Ingest simultaneously.

These paths represent **different operational modes** and must only converge at memory.

Violating this rule is a configuration error.

---

## Modbus TCP Role

Modbus TCP is the **authoritative external interface**.

* Reads from core memory
* Writes only where allowed by policy
* Never bypasses validation

The appliance behaves as a well-defined Modbus slave:

* predictable
* boring
* reliable

---

## Configuration Model

* YAML-driven
* Explicit policies
* No implicit defaults beyond safety

Configuration controls:

* memory layout
* unit ID routing
* port access
* function code permissions
* adapter enablement (Modbus / REST / MQTT / Raw)

Configuration is **runtime-only** and intentionally excluded from Git.

---

## Observability

The appliance emits:

* structured logs to stdout
* no dependency on logging backends

Designed for:

* Promtail
* Loki
* Grafana

Errors are logged **only when necessary**.

---

## Non-Goals (Strict)

This project intentionally does NOT:

* parse float32 / float64
* scale values
* interpret register meaning
* provide a GUI
* replace PLC logic
* auto-discover devices
* implement OPC UA or IEC-61850

Raw Ingest does **not** change these non-goals.

---

## Stability Contract

The following are **stable guarantees**:

* raw memory only
* atomic writes
* adapter isolation
* deterministic behavior
* no semantic leakage into the core

Breaking these requires a **major version bump**.

---

## Summary

The Modbus Memory Appliance is designed to be:

* boring
* predictable
* strict
* fast

Because in control systems:

> **boring is reliable.**
