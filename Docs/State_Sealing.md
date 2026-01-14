# State Sealing

## Purpose

**State Sealing** is a lifecycle safety feature of the Modbus Memory Appliance (MMA).

It prevents Modbus clients from accessing **uninitialized or unsafe memory** during startup.
This avoids unintended plant behavior such as controllers interpreting default values
(e.g. zero setpoints) as valid commands.

State Sealing is designed for **industrial safety**, not convenience.

---

## Default Behavior (IMPORTANT)

**State Sealing is DISABLED by default.**

This is intentional.

When disabled:
- MMA starts directly in **RUN state**
- Modbus reads and writes work immediately
- Behavior matches previous MMA versions
- No additional configuration is required

End users must **explicitly enable State Sealing** and understand its effects.
This avoids confusion and prevents the system from appearing “broken” by default.

---

## Scope

State Sealing is scoped **per memory instance**.

- Each memory has its own lifecycle state
- Modbus ports enforce the state of the memory they are bound to
- Transports do not own lifecycle state

State Sealing is a **memory-level lifecycle mechanism**, not a Modbus feature.

---

## Lifecycle States

Each memory exists in exactly one of the following states.

### 1. Pre-Run (Unsealed)

The memory exists but is **not operationally valid**.

This state is used to:
- restore critical state
- initialize control values
- synchronize external systems
- prevent unsafe plant execution

> This state only exists when State Sealing is enabled.

---

### 2. RUN (Sealed)

The memory is operational and may safely participate in control loops.

- This is the default state
- This is the only state when State Sealing is disabled

---

## State Transition (State Sealing)

**State Sealing** is the act of transitioning a memory from Pre-Run to RUN.

Properties:
- One-time only
- Explicit
- Memory-scoped
- Irreversible until restart

Once sealed, restore authority is permanently revoked.

---

## Configuration

State Sealing is **opt-in per memory**.

### Disabled (default)

```yaml
memories:
  main:
    coils: 128
    holding_registers: 256
    discrete_inputs: 128
    input_registers: 128
