# REST API Documentation

## Overview

The Modbus Memory Appliance (MMA) provides a comprehensive REST API for:
- **Health checks** (liveness)
- **Memory operations** (read/ingest)
- **Diagnostics** (memory layout, stats, MQTT status)

**Base URL:** `http://localhost:8080/api/v1`

---

## Authentication

Most endpoints require Bearer token authentication. The `/health` endpoint is always public.

### Authenticated Request Format
```
Authorization: Bearer <TOKEN>
```

**Default Token:** `INGEST_ONLY_TOKEN` (configured in `cmd/mma/rest_boot.go`)

### Response on Unauthorized (401)
```
HTTP/1.1 401 Unauthorized
(no body)
```

---

## Endpoints

### 1. Health Check
**Always public (no authentication required)**

```
GET /api/v1/health
```

#### Response (200 OK)
```json
{
  "status": "ok"
}
```

#### Use Case
Liveness probe for monitoring and container orchestration.

---

### 2. Read Memory
**Requires authentication**

```
GET /api/v1/memory/read?memory=<NAME>&area=<AREA>&address=<ADDR>&count=<COUNT>
Authorization: Bearer <TOKEN>
```

#### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `memory` | string | ✅ Yes | Memory instance name (e.g., `default`, `plant_a`) |
| `area` | string | ✅ Yes | Memory area: `coils`, `discrete_inputs`, `holding_registers`, `input_registers` |
| `address` | int | ✅ Yes | Starting address (0 or higher) |
| `count` | int | ✅ Yes | Number of values to read (1 or higher) |

#### Success Response (200 OK)
```json
{
  "values": [100, 200, 300]
}
```

#### Error Responses
| Status | Scenario |
|--------|----------|
| **403 Forbidden** | Read disabled in config |
| **400 Bad Request** | Missing/invalid query parameters |
| **404 Not Found** | Memory instance not found |

#### Examples

##### Read 3 holding registers starting at address 0
```bash
curl -X GET \
  "http://localhost:8080/api/v1/memory/read?memory=default&area=holding_registers&address=0&count=3" \
  -H "Authorization: Bearer INGEST_ONLY_TOKEN"
```

Response:
```json
{
  "values": [100, 200, 300]
}
```

##### Read coils (boolean values)
```bash
curl -X GET \
  "http://localhost:8080/api/v1/memory/read?memory=default&area=coils&address=10&count=4" \
  -H "Authorization: Bearer INGEST_ONLY_TOKEN"
```

Response:
```json
{
  "values": [1, 0, 1, 1]
}
```

---

### 3. Ingest (Write via REST)
**Requires authentication (Bearer token)**

```
POST /api/v1/ingest
Content-Type: application/json
Authorization: Bearer <TOKEN>
```

#### Request Body
```json
{
  "memory": "default",
  "area": "discrete_inputs|input_registers",
  "address": 10,
  "bools": [1, 0, 1],
  "values": [100, 200]
}
```

#### Payload Rules
- **For `discrete_inputs`:** Use `bools` array only, leave `values` empty/omitted
- **For `input_registers`:** Use `values` array only, leave `bools` empty/omitted
- **Writable areas via REST/MQTT ingest:**
  - ✅ `discrete_inputs` - sensor data (read-only from Modbus perspective)
  - ✅ `input_registers` - sensor data (read-only from Modbus perspective)
- **Not writable via REST/MQTT (use Modbus TCP instead):**
  - ❌ `coils` - output/control (written by Modbus master)
  - ❌ `holding_registers` - configuration/output (written by Modbus master)

#### Success Response (200 OK)
```json
{
  "status": "accepted",
  "memory": "default",
  "written": 3
}
```

#### Error Responses
| Status | Scenario |
|--------|----------|
| **403 Forbidden** | Ingest disabled OR area is read-only (coils, holding_registers) |
| **400 Bad Request** | Invalid JSON, wrong payload format, or value mismatch |
| **405 Method Not Allowed** | Wrong HTTP method (must be POST) |

#### Examples

##### Write discrete inputs
```bash
curl -X POST http://localhost:8080/api/v1/ingest \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer INGEST_ONLY_TOKEN" \
  -d '{
    "memory": "default",
    "area": "discrete_inputs",
    "address": 0,
    "bools": [1, 0, 1, 1]
  }'
```

Response:
```json
{
  "status": "accepted",
  "memory": "default",
  "written": 4
}
```

##### Write input registers
```bash
curl -X POST http://localhost:8080/api/v1/ingest \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer INGEST_ONLY_TOKEN" \
  -d '{
    "memory": "default",
    "area": "input_registers",
    "address": 5,
    "values": [111, 222, 333]
  }'
```

Response:
```json
{
  "status": "accepted",
  "memory": "default",
  "written": 3
}
```

---

### 4. Diagnostics: Memory Layout
**Requires authentication**

```
GET /api/v1/diagnostics/memory
Authorization: Bearer <TOKEN>
```

#### Response (200 OK)
```json
{
  "memories": {
    "default": {
      "coils": { "offset": 0, "count": 100 },
      "discrete_inputs": { "offset": 100, "count": 100 },
      "holding_registers": { "offset": 200, "count": 100 },
      "input_registers": { "offset": 300, "count": 100 }
    }
  }
}
```

#### Use Case
- Discover configured memory layout
- Understand address ranges per area
- Configuration truth (not runtime state)

#### Error Response
| Status | Scenario |
|--------|----------|
| **403 Forbidden** | Diagnostics disabled in config |

---

### 5. Diagnostics: Statistics
**Requires authentication**

```
GET /api/v1/diagnostics/stats
Authorization: Bearer <TOKEN>
```

#### Response (200 OK)
```json
{
  "rest": {
    "requests": 150,
    "reads": 45,
    "ingest": 30,
    "rejected": 5,
    "unauthorized": 2
  },
  "ingest": {
    "batches": 30,
    "written": 120,
    "rejected": 8
  }
}
```

#### Field Definitions
| Field | Description |
|-------|-------------|
| `rest.requests` | Total HTTP requests |
| `rest.reads` | Successful memory reads |
| `rest.ingest` | Successful ingest requests |
| `rest.rejected` | Rejected ingest requests (validation errors) |
| `rest.unauthorized` | Missing/invalid Bearer tokens |
| `ingest.batches` | Total ingest batches processed |
| `ingest.written` | Total registers written |
| `ingest.rejected` | Rejected ingest operations |

#### Use Case
Monitoring and observability.

#### Error Response
| Status | Scenario |
|--------|----------|
| **403 Forbidden** | Diagnostics disabled in config |

---

### 6. Diagnostics: MQTT Status
**Requires authentication**

```
GET /api/v1/diagnostics/mqtt
Authorization: Bearer <TOKEN>
```

#### Response (200 OK) — MQTT Enabled
```json
{
  "enabled": true,
  "connected": true,
  "broker": "mqtt.example.com:1883",
  "topic": "modbus/ingest",
  "client_id": "mma-client-001"
}
```

#### Response (200 OK) — MQTT Disabled
```json
{
  "enabled": false
}
```

#### Use Case
Health check for MQTT integration.

#### Error Response
| Status | Scenario |
|--------|----------|
| **403 Forbidden** | Diagnostics disabled in config |

---

## Error Handling

### Common Error Response Format
```json
{
  "status": "rejected",
  "error": "descriptive error message"
}
```

### HTTP Status Codes
| Code | Meaning |
|------|---------|
| **200 OK** | Success |
| **400 Bad Request** | Invalid request format, missing parameters, validation failure |
| **401 Unauthorized** | Missing or invalid Bearer token |
| **403 Forbidden** | Endpoint disabled, read-only area, or insufficient permissions |
| **404 Not Found** | Memory instance not found |
| **405 Method Not Allowed** | Wrong HTTP method (e.g., GET on POST endpoint) |
| **500 Internal Server Error** | Unexpected server error |

---

## Configuration

### Enable/Disable Features (in `config.yaml`)

```yaml
rest:
  enabled: true
  address: "0.0.0.0:8080"

# Endpoint control (hard-wired in rest_boot.go for now)
# To disable: set handler flags in Handlers struct
```

### Authentication (in `cmd/mma/rest_boot.go`)

```go
tokenSet := rest.NewTokenSet(
  true, // enable auth
  []string{
    "INGEST_ONLY_TOKEN",
    "YOUR_TOKEN_HERE",
  },
)
```

---

## Unified Payload Format

### REST Ingest ↔ MQTT Ingest
Both use the same JSON schema:

```json
{
  "memory": "instance_name",
  "area": "discrete_inputs|input_registers",
  "address": 0,
  "bools": [1, 0, 1],
  "values": [100, 200, 300]
}
```

**Key Difference:**
- **REST:** `memory` field is hardcoded to `"default"` (input ignored)
- **MQTT:** `memory` field is used as-is from the payload

---

## Best Practices

### 1. Token Management
- Use environment variables for tokens in production
- Rotate tokens regularly
- Use different tokens for different services if possible

### 2. Batch Operations
- Group related writes into single ingest calls
- Atomicity guaranteed per request
- Failed requests leave memory unchanged

### 3. Monitoring
- Poll `/health` periodically for liveness checks
- Monitor `/diagnostics/stats` for operational insight
- Check `/diagnostics/mqtt` for MQTT connectivity

### 4. Error Handling
- Validate query parameters before sending requests
- Handle 400/403/404 responses explicitly
- Retry on 500 (server errors) with exponential backoff

### 5. Performance
- Read timeout: 10 seconds
- Write timeout: 10 seconds
- Idle timeout: 30 seconds

---

## Compatibility

This REST API is **stable**. Future extensions will be:
- Additive (new endpoints, not breaking changes)
- Backward compatible (existing endpoints unchanged)
- Explicitly documented
