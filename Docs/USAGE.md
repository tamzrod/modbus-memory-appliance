# Modbus Memory Appliance – Usage Guide

## Purpose

The **Modbus Memory Appliance (MMA)** is a deterministic Modbus TCP memory server with optional MQTT and REST ingestion.

It provides:
- raw **uint16 register storage**
- **atomic batch writes**
- strict validation
- non-fatal adapters (MQTT / REST cannot crash Modbus)

This appliance is intentionally **dumb and fast**.
All interpretation, scaling, and control logic belongs **outside** the appliance.

---

## High-Level Flow

```
[ MQTT / REST ]  ──▶  [ INGEST + VALIDATION ]  ──▶  [ MEMORY ]  ──▶  [ Modbus TCP ]
                                            (uint16 only)
```

---

## 1. Configuration (`config.yaml`)

`config.yaml` is **runtime-only** and is not committed to Git.

### Minimal Example

```yaml
memory:
  plant_a:
    size: 1000

routing:
  unit_id_map:
    1: plant_a

ports:
  502:
    allow_unit_ids: [1]
    allow_memories: [plant_a]
    allow_function_codes: [3, 4, 6, 16]

mqtt:
  enabled: true
  broker: tcp://mosquitto:1883
  topic: modbus/ingest
  client_id: mma-plant-a
```

---

## 2. Memory Model

- Memory is **zero-based**
- Each register is **uint16**
- Valid range: `0 – 65535`
- Out-of-range values are rejected

---

## 3. Modbus TCP Usage

### Supported Function Codes

- FC 3 – Read Holding Registers
- FC 4 – Read Input Registers
- FC 6 – Write Single Register
- FC 16 – Write Multiple Registers

### Example (modpoll)

```bash
modpoll -m tcp -p 502 -a 1 -r 0 -c 10 localhost
```

---

## 4. MQTT Ingest (Primary Write Path)

### Topic

Defined in `config.yaml`:

```yaml
mqtt:
  topic: modbus/ingest
```

---

## 5. MQTT JSON Payload Format (Canonical)

### Single Register Write

```json
{
  "memory": "plant_a",
  "area": "holding_registers",
  "address": 0,
  "values": [123]
}
```

### Batch Write (Atomic)

```json
{
  "memory": "plant_a",
  "area": "holding_registers",
  "address": 0,
  "values": [10, 20, 30, 40]
}
```

---

## 6. Validation Rules

- Values must be integers
- Range must be `0–65535`
- Any invalid value rejects the entire batch
- No partial writes

---

## 7. Error Logging

Malformed ingest is logged in a Loki-friendly format:

```
ingest_reject src=192.168.1.45 memory=plant_a area=holding_registers addr=0 index=2 value=999999 reason=uint16_out_of_range
```

---

## 8. Docker Usage

```bash
docker run -d \
  --name modbus-memory \
  --restart unless-stopped \
  -p 502:502 \
  -v $(pwd)/config.yaml:/app/config.yaml:ro \
  rodtamin/modbus-memory-appliance:0.1.1
```

---

## Design Guarantees

- Atomic writes
- Deterministic behavior
- No silent truncation
- Adapter failures never crash Modbus
