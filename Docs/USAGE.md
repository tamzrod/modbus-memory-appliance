# Modbus Memory Appliance (MMA) — Usage Guide

**Version:** Architecture‑Aligned (2026)

---

## 1. Purpose

The **Modbus Memory Appliance (MMA)** is a deterministic, in‑memory Modbus TCP server with optional REST and MQTT ingestion adapters.

It provides:

* Raw **bool** (0/1) and **uint16** memory storage
* **Atomic reads and writes** (all‑or‑nothing)
* Strict bounds and validation
* Transport isolation (REST/MQTT cannot crash Modbus)

This appliance is intentionally **dumb and fast**.
All interpretation, scaling, units, floats, and control logic live **outside** MMA.

---

## 2. What MMA Is Not

MMA is **not**:

* A PLC
* A SCADA system
* A protocol translator
* A rules engine
* A historian or database
* A persistence layer

If logic or meaning is required, it belongs in **Node‑RED, PLCs, PPC logic, or external services**.

---

## 3. High‑Level Architecture

```
           (device ingestion)
        ┌───────────┐
        │ MQTT /    │
        │ REST      │
        └─────┬─────┘
              │
              ▼
     ┌───────────────────┐
     │  IngestService     │
     │  (DI / IR only)    │
     └─────┬─────────────┘
           │
           ▼
     ┌───────────────────┐
     │  Memory Appliance  │
     │  (raw bool/uint16) │
     └─────┬─────────────┘
           │
           ▼
     ┌───────────────────┐
     │   Modbus TCP       │
     │ (client control)  │
     └───────────────────┘
```

**Separation of concerns (locked):**

* **Modbus TCP** → client control plane
* **REST / MQTT** → device ingestion plane

---

## 4. Memory Model (Locked)

* Internal addressing is **zero‑based**
* Memory areas:

  * Coils → `bool`
  * Discrete Inputs (DI) → `bool`
  * Holding Registers (HR) → `uint16`
  * Input Registers (IR) → `uint16`
* Valid register range: `0–65535`
* No floats, no scaling, no semantics
* Memory layout defined **only at startup**

### Atomicity Guarantees

* Single write → atomic
* Batch write → atomic
* No partial or mixed state is observable
* Invalid batch → entire batch rejected

---

## 5. Configuration (`config.yaml`)

### Rules

* Loaded **once at startup**
* Stored in RAM
* Immutable at runtime
* Invalid config → **recovery mode**
* Restart required to apply changes

### Minimal Example

```yaml
memory:
  memories:
    plant_a:
      default: true
      coils: { start: 0, size: 1024 }
      discrete_inputs: { start: 0, size: 1024 }
      holding_registers: { start: 0, size: 4096 }
      input_registers: { start: 0, size: 4096 }

routing:
  unit_id_map:
    1: plant_a

ports:
  502:
    allow_unit_ids: [1]
    allow_memories: [plant_a]
    allow_function_codes: [3, 4, 6, 16]

rest:
  enable: true

mqtt:
  enable: true
```

---

## 6. Modbus TCP Usage (Client Plane)

### Supported Function Codes

* FC 3 — Read Holding Registers
* FC 4 — Read Input Registers
* FC 6 — Write Single Holding Register
* FC 16 — Write Multiple Holding Registers

### Modbus Rules

* Modbus **cannot write**:

  * Discrete Inputs
  * Input Registers
* All writes are bounds‑checked and atomic
* Port policies strictly enforced

### Example (modpoll)

```bash
modpoll -m tcp -p 502 -a 1 -r 0 -c 10 localhost
```

---

## 7. Unified Ingestion (REST & MQTT)

REST and MQTT share a **single canonical ingest model** and the same internal `IngestService`.

### Ingestion Scope (Locked)

* Allowed areas:

  * `discrete_inputs`
  * `input_registers`
* Disallowed:

  * `coils`
  * `holding_registers`

No Modbus function codes are simulated.

---

## 8. Canonical Ingest Payload (JSON)

### Boolean Ingest (DI)

Boolean values are encoded as **0 or 1** (not `true/false`).

```json
{
  "memory": "plant_a",
  "area": "discrete_inputs",
  "address": 0,
  "bools": [1, 0, 1]
}
```

### Register Ingest (IR — Atomic)

```json
{
  "memory": "plant_a",
  "area": "input_registers",
  "address": 10,
  "values": [1234, 5678, 9012]
}
```

### Validation Rules

* Entire payload validated before write
* Any invalid value → no write
* Bounds enforced
* Atomic commit

---

## 9. REST API

### Base URL

```
http://localhost:8080/api/v1
```

### Core Endpoints

| Method | Endpoint              | Purpose               |
| ------ | --------------------- | --------------------- |
| GET    | `/health`             | Liveness check        |
| GET    | `/diagnostics/memory` | Memory layout & stats |
| GET    | `/diagnostics/mqtt`   | MQTT status           |
| GET    | `/diagnostics/stats`  | Counters              |
| POST   | `/ingest`             | Canonical ingest      |
| GET    | `/memory/read`        | Direct memory read    |

---

## 10. Recovery Mode

If configuration validation fails:

* Modbus TCP → disabled
* MQTT → disabled
* REST → enabled
* Memory not exposed
* Config can be inspected and repaired
* Restart required after repair

---

## 11. Docker Usage

```bash
docker run -d \
  --name mma \
  --restart unless-stopped \
  -p 502:502 \
  -p 8080:8080 \
  -v $(pwd)/config.yaml:/app/config.yaml:ro \
  rodtamin/modbus-memory-appliance:latest
```

---

## 12. Windows — Run as a Service

MMA can run natively on **Windows** as a **Windows Service**.

### Build

```bash
GOOS=windows GOARCH=amd64 go build -o mma.exe ./cmd/mma
```

### Install Service (Administrator)

```powershell
mma.exe install
```

### Start / Stop

```powershell
mma.exe start
mma.exe stop
```

### Uninstall

```powershell
mma.exe uninstall
```

### Recommended Layout

```
C:\mma\
 ├─ mma.exe
 └─ config.yaml
```

---

## 13. Design Guarantees (Non‑Negotiable)

* Deterministic behavior
* Atomic writes
* Strict bounds enforcement
* No silent truncation
* No runtime memory mutation
* Adapter failures never affect Modbus core

---

## Final Mental Model

> **MMA is a dumb, fast, in‑memory Modbus server where all intelligence lives outside.**
