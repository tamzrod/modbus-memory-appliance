# Modbus Memory Appliance – Design Philosophy

## Purpose

The Modbus Memory Appliance (MMA) is a **deterministic, safety-first Modbus TCP memory server**.

Its sole responsibility is to:
- store raw Modbus data
- enforce strict bounds and validation
- serve data predictably

Everything else belongs **outside** the appliance.

---

## Core Principles

### 1. Determinism Over Convenience

If MMA is running, its behavior is:
- predictable
- repeatable
- explainable

There are no:
- hidden defaults
- time-based changes
- implicit behavior

If behavior changes, the process must restart.

---

### 2. Configuration Is Immutable at Runtime

Configuration is:
- loaded once at startup
- validated strictly
- stored in memory
- never reloaded automatically

Any configuration change:
- requires validation
- requires an explicit restart
- never mutates live runtime state

This prevents undefined behavior and partial application.

---

### 3. Restart Is a Feature, Not a Failure

Restarting is:
- explicit
- safe
- intentional

MMA is designed so that:
- restart restores a known-good state
- no runtime drift exists
- operators always know when changes apply

---

### 4. Memory Is Dumb and Honest

Memory:
- stores only raw values (`uint16`, `bool`)
- performs no scaling
- performs no interpretation
- performs no control logic

All semantics belong to:
- PLCs
- controllers
- Node-RED
- SCADA
- external systems

MMA never guesses user intent.

---

### 5. Strict Validation, Loud Failure

MMA prefers:
- rejection over silent correction
- explicit errors over assumptions
- failing fast over failing later

Invalid operations:
- reject entire batches
- never partially apply
- are logged clearly

---

### 6. Adapters Must Never Affect the Core

REST, MQTT, and other adapters:
- are optional
- are non-fatal
- cannot crash Modbus
- cannot corrupt memory

The Modbus core must remain operational even if adapters fail.

---

### 7. Security Is Layered, Not Embedded

Security is enforced through:
- middleware
- configuration
- network controls
- external gateways

MMA does not:
- manage users
- manage sessions
- implement OAuth
- mutate secrets at runtime

Security can be hardened progressively without redesign.

---

### 8. Recovery Is Mandatory

Invalid configuration must never brick the system.

If configuration is malformed:
- MMA enters recovery mode
- Modbus and MQTT are disabled
- REST remains available
- configuration can be inspected and repaired
- restart applies the fix

---

### 9. The Core Must Remain Boring

The core of MMA should:
- feel boring
- change rarely
- be easy to audit
- be easy to reason about

Innovation belongs at the edges.

---

## Non-Goals

MMA is NOT:
- a PLC
- a SCADA system
- a protocol translator
- a control engine
- a configuration UI
- a database

Attempting to make it so violates its purpose.

---

## Design Contract

Any change to MMA must answer:

1. Does this preserve determinism?
2. Does this require runtime mutation?
3. Does this introduce ambiguity?
4. Does this increase blast radius?
5. Does this belong outside the core?

If the answer is unclear, the change does not belong in MMA.

---

## Final Statement

MMA exists to be trusted.

Trust comes from:
- predictability
- restraint
- clarity
- correctness

This project chooses boring correctness over clever features.
