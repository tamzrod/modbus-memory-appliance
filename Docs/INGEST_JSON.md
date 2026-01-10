# Modbus Memory Appliance – MQTT / REST JSON Ingest Format

This document defines the **canonical JSON ingest format** used by both **MQTT** and **REST** adapters.

The format is **stable, strict, and shared** across all ingest transports.

---

## Design Rules (Read First)

- Same JSON schema for MQTT and REST
- No partial writes
- Atomic batch semantics
- No parsing or scaling
- Values are written as raw `uint16`
- Any malformed value rejects the entire batch

---

## Common Fields

| Field     | Type    | Required | Description |
|----------|---------|----------|-------------|
| memory   | string  | yes      | Target memory name |
| area     | string  | yes      | Register area |
| address  | number  | yes      | Zero-based start address |
| values   | array   | yes      | Array of integer values |

---

## Supported Areas

Currently supported:

```json
"area": "holding_registers"
```

Planned (future):
- `input_registers`
- `coils`
- `discrete_inputs`

---

## Single Register Write

### JSON Payload

```json
{
  "memory": "plant_a",
  "area": "holding_registers",
  "address": 0,
  "values": [123]
}
```

### Result

- Holding Register 0 = `123`
- Atomic (single value)

---

## Batch Register Write (Atomic)

### JSON Payload

```json
{
  "memory": "plant_a",
  "area": "holding_registers",
  "address": 10,
  "values": [10, 20, 30, 40]
}
```

### Result

| Address | Value |
|--------|-------|
| 10     | 10    |
| 11     | 20    |
| 12     | 30    |
| 13     | 40    |

All values are written **together** or **not at all**.

---

## Validation Rules

- Values must be integers
- Range must be `0–65535`
- Negative values are rejected
- Overflow values are rejected
- Any invalid value rejects the entire batch

---

## Example: Malformed Payload (Rejected)

### JSON Payload

```json
{
  "memory": "plant_a",
  "area": "holding_registers",
  "address": 0,
  "values": [100, 200, 999999]
}
```

### Result

- ❌ No registers written
- ❌ No partial state
- ✅ Error logged

---

## Error Logging (Loki-Friendly)

Malformed ingest is logged as a single structured log line:

```
ingest_reject src=192.168.1.45 memory=plant_a area=holding_registers addr=0 index=2 value=999999 reason=uint16_out_of_range
```

Fields:
- `src` – source IP or client identifier
- `index` – index of invalid value in `values[]`
- `value` – offending value
- `reason` – rejection reason

---

## MQTT Usage

### Topic

Defined in `config.yaml`:

```yaml
mqtt:
  topic: modbus/ingest
```

### Publish Example

```bash
mosquitto_pub \
  -t modbus/ingest \
  -m '{
    "memory":"plant_a",
    "area":"holding_registers",
    "address":0,
    "values":[1,2,3,4]
  }'
```

---

## REST Usage

### Endpoint (example)

```
POST /api/ingest
```

### Body

```json
{
  "memory": "plant_a",
  "area": "holding_registers",
  "address": 20,
  "values": [500, 600]
}
```

REST uses **the exact same validation and semantics** as MQTT.

---

## What This Format Is NOT

- ❌ Not a PLC language
- ❌ Not semantic data
- ❌ Not scaled values
- ❌ Not floats

It is **raw register memory ingestion**.

---

## Compatibility Guarantee

This JSON schema is considered **stable**.

Future extensions will be:
- additive
- backward compatible
- explicitly documented

---

## Recommended File Placement

```
Docs/
  └── INGEST_JSON.md
```

