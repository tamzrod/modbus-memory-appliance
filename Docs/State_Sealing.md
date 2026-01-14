# State Sealing

## Purpose

**State Sealing** is a lifecycle safety mechanism of the **Modbus Memory Appliance (MMA)**.

Its purpose is to prevent Modbus clients from interacting with **uninitialized or unsafe memory** during startup, while still allowing external systems (REST / MQTT) to initialize state.

This avoids unintended plant behavior such as:
- Controllers acting on default values
- Zero setpoints being interpreted as valid commands
- Control loops starting before state restoration is complete

State Sealing exists for **industrial safety**, not convenience.

---

## Default Behavior (IMPORTANT)

**State Sealing is DISABLED by default.**

This is intentional.

When State Sealing is disabled:

- MMA starts immediately in **RUN state**
- Modbus reads and writes are fully enabled
- REST and MQTT ingestion behave normally
- Behavior is identical to previous MMA versions
- No additional configuration is required

Users must **explicitly enable State Sealing** and understand its impact.

---

## Scope

State Sealing is scoped **per memory instance**.

- Each memory has its own lifecycle state
- Modbus access is enforced per memory
- REST and MQTT do **not** own lifecycle state
- Lifecycle state is not global

State Sealing is a **memory-level lifecycle mechanism**, not a Modbus feature.

---

## Lifecycle States

### 1. Pre-Run (Unsealed)

- Modbus access: ❌ blocked  
- REST / MQTT ingestion: ✅ allowed  

Used for:
- Restoring retained values
- Initializing control parameters
- Synchronizing external systems

### 2. RUN (Sealed)

- Modbus access: ✅ enabled  
- Normal operational state  

---

## State Transition (Sealing)

Transition from **Pre-Run → RUN**:

- One-time only
- Explicit
- Memory-scoped
- Irreversible until restart

---

## Gate Mechanism

- Single discrete input bit
- Writing `1` seals the memory
- Evaluated via REST / MQTT
- Modbus cannot open the gate

---

## Configuration Example

```yaml
memory:
  memories:
    plant_a:
      default: true
      coils:
        start: 0
        size: 1024
      discrete_inputs:
        start: 0
        size: 1024
      holding_registers:
        start: 0
        size: 4096
      input_registers:
        start: 0
        size: 4096

  state_sealing:
    enable: true
    gate:
      area: discrete_inputs
      address: 127
```

---

## Opening the Gate (REST)

```http
POST /api/v1/ingest
Authorization: Bearer INGEST_ONLY_TOKEN
Content-Type: application/json

{
  "memory": "plant_a",
  "area": "discrete_inputs",
  "address": 127,
  "bools": [1]
}
```

Expected log:

```
[STATE] memory transitioned to RUN via gate @ discrete_inputs[127]
```

---

## Summary

- Disabled by default
- Per-memory lifecycle
- REST / MQTT initialize state
- Discrete-input gate seals memory
- Modbus enabled only after sealing

Deterministic, explicit, and safe startup behavior.
