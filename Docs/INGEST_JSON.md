# Modbus Memory Appliance – Ingest JSON Format (MQTT / REST)

This document defines the **complete, canonical JSON ingest format** used by **all ingest adapters** (MQTT and REST).

The format is **transport-agnostic**, **strict**, and **atomic**.

---

## Core Rules (Non‑Negotiable)

- Same JSON format for MQTT and REST
- All writes are **atomic** (all‑or‑nothing)
- No parsing, scaling, or semantics
- Memory stores **raw values only**
- Validation happens before any write
- Any malformed value rejects the entire batch

---

## Common JSON Fields

| Field     | Type    | Required | Description |
|----------|---------|----------|-------------|
| memory   | string  | yes      | Target memory name |
| area     | string  | yes      | Register / bit area |
| address  | number  | yes      | Zero‑based start address |
| values   | array   | yes      | Values to write |

---

## Supported Areas

| Area Name             | Type     | Value Encoding |
|----------------------|----------|----------------|
| holding_registers    | uint16   | 0–65535 |
| input_registers      | uint16   | 0–65535 |
| coils                | bool     | 0 / 1 |
| discrete_inputs      | bool     | 0 / 1 |

> ⚠️ Even bit‑based areas use integers.  
> `0 = false`, `1 = true`.

---

# HOLDING REGISTERS

## Single Holding Register Write

```json
{
  "memory": "plant_a",
  "area": "holding_registers",
  "address": 0,
  "values": [123]
}
```

## Batch Holding Register Write

```json
{
  "memory": "plant_a",
  "area": "holding_registers",
  "address": 10,
  "values": [100, 200, 300, 400]
}
```

---

# INPUT REGISTERS

> Input registers are writable **only via ingest adapters**.  
> Modbus clients cannot write input registers.

## Single Input Register Write

```json
{
  "memory": "plant_a",
  "area": "input_registers",
  "address": 0,
  "values": [555]
}
```

## Batch Input Register Write

```json
{
  "memory": "plant_a",
  "area": "input_registers",
  "address": 5,
  "values": [10, 20, 30]
}
```

---

# COILS

> Coils represent boolean values.

## Single Coil Write

```json
{
  "memory": "plant_a",
  "area": "coils",
  "address": 0,
  "values": [1]
}
```

## Batch Coil Write

```json
{
  "memory": "plant_a",
  "area": "coils",
  "address": 10,
  "values": [1, 0, 1, 1]
}
```

---

# DISCRETE INPUTS

> Discrete inputs are writable **only via ingest adapters**.

## Single Discrete Input Write

```json
{
  "memory": "plant_a",
  "area": "discrete_inputs",
  "address": 0,
  "values": [0]
}
```

## Batch Discrete Input Write

```json
{
  "memory": "plant_a",
  "area": "discrete_inputs",
  "address": 20,
  "values": [1, 1, 0, 0, 1]
}
```

---

## Validation Rules (All Areas)

### Register Areas
- Values must be integers
- Range: `0–65535`

### Bit Areas
- Values must be `0` or `1`
- Any other value is rejected

### Atomicity
- If **any value is invalid**, the entire batch is rejected
- No partial writes
- Memory remains unchanged

---

## Malformed Example (Rejected)

```json
{
  "memory": "plant_a",
  "area": "coils",
  "address": 0,
  "values": [1, 0, 2]
}
```

Result:
- ❌ Rejected (invalid bit value `2`)
- ❌ No state change
- ✅ Error logged

---

## Error Logging (Loki‑Friendly)

Example log:

```
ingest_reject src=192.168.1.45 memory=plant_a area=coils addr=0 index=2 value=2 reason=invalid_bit
```

Logged fields:
- `src` – source IP or client ID
- `area` – target area
- `index` – index in values array
- `value` – offending value
- `reason` – validation failure

---

## MQTT Example

```bash
mosquitto_pub \
  -t modbus/ingest \
  -m '{
    "memory":"plant_a",
    "area":"holding_registers",
    "address":0,
    "values":[10,20,30]
  }'
```

---

## REST Example

```http
POST /api/ingest
Content-Type: application/json

{
  "memory":"plant_a",
  "area":"input_registers",
  "address":5,
  "values":[111,222]
}
```

---

## Compatibility Guarantee

This ingest format is **stable**.

Any future extensions will be:
- additive
- backward compatible
- explicitly documented

---

## Recommended Location

```
Docs/
  └── INGEST_JSON.md
```
