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
If the configuration is invalid, **the server will not start**.

### Minimal Example

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

rest:
  enabled: true
  address: ":8080"
```

---

## 2. Memory Model

- Memory is **zero-based**
- Each register is **uint16**
- Valid range: `0 – 65535`
- Out-of-range values are rejected
- No partial writes (atomic batches only)

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

## 6. REST API Usage

Base URL (default):

```
http://localhost:8080/api/v1
```

### 6.1 Health Check

```
GET /health
```

Response:

```json
{ "status": "ok" }
```

### 6.2 Diagnostics – Memory

```
GET /diagnostics/memory
```

### 6.3 Diagnostics – MQTT

```
GET /diagnostics/mqtt
```

### 6.4 Memory Read

```
GET /memory/read
```

Required query parameters:

- `memory`
- `area`
- `address`
- `count`

Example:

```bash
curl "http://localhost:8080/api/v1/memory/read?memory=plant_a&area=holding_registers&address=0&count=4"
```

---

## 7. Validation Rules

- Values must be integers
- Range must be `0–65535`
- Memory must exist
- Area must be valid
- Any invalid value rejects the entire batch
- No partial writes

---

## 8. Docker Usage

```bash
docker run -d   --name modbus-memory   --restart unless-stopped   -p 502:502   -p 8080:8080   -v $(pwd)/config.yaml:/app/config.yaml:ro   rodtamin/modbus-memory-appliance:0.1.1
```

---

## Design Guarantees

- Atomic writes
- Deterministic behavior
- No silent truncation
- No implicit defaults
- Adapter failures never crash Modbus
