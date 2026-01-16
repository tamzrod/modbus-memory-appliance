# RAW INGEST -- Binary Memory Injection Protocol

## Purpose

**Raw Ingest** is a privileged, binary transport designed to write **raw
values directly into MMA memory** with minimal overhead.

It is intentionally: - semantic-free - policy-free - format-agnostic
(within strict structural rules)

Raw ingest exists to support: - high-speed shadow memory mirroring -
deterministic bulk writes - gateway / bridge replication - offline
preload and restore - DMA-style memory injection

It is **not** a user API and **not** a control protocol.

------------------------------------------------------------------------

## Position in MMA Architecture

Raw ingest is a **transport**, not an ingest service.

    REST / MQTT (semantic ingest)
            |
    Modbus TCP (control plane)
            |
    RAW INGEST (binary / DMA)
            |
    CORE MEMORY (dumb, locked)

Raw ingest sits **closest to memory**.

------------------------------------------------------------------------

## Design Principles

Raw ingest obeys the following non-negotiable rules:

-   No semantics
-   No scaling
-   No validation of meaning
-   No state sealing
-   No events
-   No side effects
-   No memory allocation based on content

**One packet = one atomic memory write.**

------------------------------------------------------------------------

## Memory Access Rules

Raw ingest **MAY write to all memory areas**:

    Area Code Memory Area         Format
  ----------- ------------------- --------
           00 Coils               bit
           01 Discrete Inputs     bit
           02 Holding Registers   u16
           03 Input Registers     u16

Rules: - Addressing is **zero-based** - Bounds are strictly enforced -
Writes are atomic - Last writer wins

------------------------------------------------------------------------

## Wire Protocol

Raw ingest uses a **fixed binary protocol** over TCP.

### Byte Order

-   Multi-byte fields: **big-endian**
-   Bit payloads: **Modbus LSB-first**

------------------------------------------------------------------------

## Packet Format

    +------------+----------+------+--------+--------+----------+
    | memory_id  | area     | fmt  | start  | count  | payload  |
    | 1 byte     | 1 byte   | 1    | 2      | 2      | N bytes |
    +------------+----------+------+--------+--------+----------+

Header size: **8 bytes**

------------------------------------------------------------------------

## Validation Rules

Packets are rejected if **any** rule fails:

1.  memory_id exists
2.  area is valid
3.  format is valid
4.  format matches area
5.  start + count within bounds
6.  payload length matches expected size
7.  packet size ≤ configured maximum

On failure: - no memory is modified - server returns `0x01`

------------------------------------------------------------------------

## Response

Server responds with **1 byte**:

-   `0x00` success
-   `0x01` failure

Connection remains open.

------------------------------------------------------------------------

## Configuration Example

``` yaml
raw_ingest:
  enabled: true
  listen: ":9000"
  max_packet_bytes: 2048
```

------------------------------------------------------------------------

## Final Note

Raw ingest exists to keep the **core memory dumb and fast**.

All intelligence belongs **above** it.
