# Modbus TCP Stress Test Report  
**Go Modbus Memory Appliance**

**Audience:** Public / Technical  
**Date:** 2026-01  
**Scope:** Transport-level performance and stability  
**Non-Goal:** PLC, SCADA, or control logic comparison

---

## 1. Overview

This document presents stress-test results for a **Go-based Modbus TCP server** (Modbus Memory Appliance) and compares its behavior against **typical commercial Modbus servers** and **Node-RED Modbus implementations**.

The goal is to evaluate:
- Throughput
- Latency (average and tail)
- Stability under load
- Saturation behavior

The server under test is intentionally designed as a **dumb, deterministic transport appliance**:
- Raw memory only (`uint16`, `bool`)
- No parsing, scaling, or tag model
- No scan cycle
- No UI, alarms, or historian logic
- All intelligence lives outside the Modbus core

---

## 2. Test Environment

- **Target:** `127.0.0.1:502`
- **Protocol:** Modbus TCP
- **Function Code:** FC3 (Read Holding Registers)
- **Register Count:** 10
- **Unit ID:** 1
- **Workers:** 20 concurrent clients
- **Ramp Steps:** 100 → 500 → 1000 → 2000 req/s
- **Step Duration:** 5 seconds
- **Strict Mode:** Disabled
- **Transport:** Plain TCP (no TLS)

> Note: This is a localhost test intended to measure **server-side cost and contention**, not network latency.

---

## 3. Results Summary

### Throughput & Latency

| Target Rate | Achieved Throughput | Avg Latency | p99 Latency | Errors |
|------------|--------------------|------------|-------------|--------|
| 100 req/s  | 100.0 req/s        | 0.507 ms   | 0.524 ms    | 0 |
| 500 req/s  | 500.6 req/s        | 0.534 ms   | 0.524 ms    | 0 |
| 1000 req/s | 999.4 req/s        | 0.493 ms   | 0.524 ms    | 0 |
| 2000 req/s | 1723.2 req/s       | 0.635 ms   | 1.049 ms    | 0 |

### Key Observations
- Linear scaling up to **1000 req/s**
- Sub-millisecond latency maintained under normal load
- Graceful saturation at 2000 req/s
- No protocol errors, exceptions, or timeouts
- Latency distribution remains tight (no long tail explosion)

---

## 4. What Is p99?

**p99 latency** means:

> 99% of requests complete at or below this latency.

Example:
- **p99 = 1.0 ms** → 99 out of 100 requests finished in ≤ 1.0 ms

Why it matters:
- Average latency can hide jitter
- p99 reflects **worst-case experience for almost all clients**
- Industrial systems care about predictability, not just speed

In this test:
- p99 stays ~**0.524 ms** up to 1000 req/s
- p99 rises to ~**1.049 ms** at saturation
- Maximum observed spike was **3.5 ms**, with no error cascade

---

## 5. Comparative Context (Typical Industry Behavior)

The following comparisons reflect **commonly observed field behavior**, integrator experience, and vendor documentation — **not marketing peak claims**.

### Throughput & Latency Comparison

| System Type | Typical Throughput | Avg Latency | p99 Latency |
|------------|-------------------|-------------|-------------|
| PLC / Gateway | 50–500 req/s | 5–20 ms | 10–50 ms |
| PC-based Commercial Server | 300–1500 req/s | 2–10 ms | 5–25 ms |
| Node-RED Modbus | 100–400 req/s | 15–40 ms | 30–80 ms |
| **Go Modbus Memory Appliance** | **1700+ req/s** | **~0.5 ms** | **~1 ms** |

---

## 6. Why the Difference Exists

### Commercial / PLC Systems
- Scan-cycle driven
- Tag parsing and scaling
- Shared execution with UI, alarms, logging
- Designed for configurability and safety, not raw throughput

### Node-RED Modbus
- JavaScript event loop
- Message passing between nodes
- Flow orchestration overhead
- Competes with user logic and dashboards

### Go Modbus Memory Appliance
- Raw memory access only
- Minimal lock scope (RWMutex)
- No background tasks
- No semantic interpretation
- Each request handled independently

The performance difference is **architectural**, not tuning-based.

---

## 7. Stability Under Load

At saturation:
- Throughput degrades smoothly
- Latency increases gradually
- No error storms
- No runaway queues
- No client timeouts observed

This behavior is characteristic of a **well-designed transport appliance**, not an application server.

---

## 8. Limitations of This Test

To keep this report honest:

- Localhost only (no NIC or switch latency)
- Single function code (FC3)
- Single unit ID
- No TLS
- No concurrent REST / MQTT traffic

Even with these constraints, observed performance exceeds typical industrial requirements by a wide margin.

---

## 9. Intended Use Case

This server is **not**:
- A PLC
- A SCADA
- A control engine
- A tag database

It **is** intended to be:
- A deterministic Modbus backplane
- A shared memory appliance
- A stable foundation for Node-RED, PLCs, gateways, and higher-level logic

---

## 10. Conclusion

The Go Modbus Memory Appliance demonstrates:
- Sub-millisecond tail latency
- High throughput headroom
- Predictable saturation behavior
- Strong suitability as a transport-level Modbus server

The results validate the design goal:
> **Keep the Modbus core dumb, fast, and deterministic.**

All intelligence belongs above it.

---

_End of report._
