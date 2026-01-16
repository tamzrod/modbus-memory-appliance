# RAW INGEST — CODING RULES (AUTHORITATIVE)

> **Purpose of this document**  
> This file defines the **exact rules required to implement Raw Ingest correctly**.  
> If something is not explicitly listed here, it **does not belong in Raw Ingest**.

---

## 1. Definition (One Sentence)

> **Raw Ingest is a direct, alignment-only TCP write path into MMA memory.**

No decode. No semantics. No rules.

---

## 2. Architectural Position

Raw Ingest is a **peer writer** to REST/MQTT ingest.

```
FIELD DEVICE (RAW)                     SOME DEVICE (LOGICAL)
(meter / logger / combiner)            (controller / gateway / app)
        |                                      |
        v                                      v
┌────────────────────┐              ┌────────────────────────┐
│     RAW INGEST     │              │     REST / MQTT        │
│      (TCP)         │              │       INGEST           │
│ • blind writes     │              │ • semantic writes      │
│ • no decode        │              │ • rules & validation   │
└─────────┬──────────┘              └─────────┬──────────────┘
          |                                      |
          └───────────────┬──────────────────────┘
                          v
               ┌──────────────────────────┐
               │        MMA MEMORY        │
               │  (raw registers only)   │
               └─────────┬────────────────┘
                         v
                      READERS
```

**Locks:**
- Devices are not shared
- Paths do not cross
- Memory is the only convergence point

---

## 3. Hard Rules (Non‑Negotiable)

Raw Ingest **MUST**:
- Write directly to `core.Memory`
- Be stateless
- Be write‑only
- Be panic‑free
- Perform bounds checking only

Raw Ingest **MUST NOT**:
- Import or reference `internal/ingest`
- Decode Modbus function codes
- Apply scaling or typing
- Track freshness or timestamps
- Zero memory on silence
- Contain control logic

If Raw Ingest understands *meaning*, it is a bug.

---

## 4. Alignment vs Decode

Raw Ingest performs **alignment only**.

Alignment:
- Bytes → uint16 registers
- Bytes → bit arrays (coils / discretes)
- Sequential write starting at address

Decode (forbidden):
- Understanding data types
- Applying Modbus semantics
- Interpreting function codes

---

## 5. Packet Processing Model

- **One packet = one atomic write**
- No partial writes
- Reject on any error

Processing order:
1. Receive full packet
2. Validate structure
3. Validate bounds
4. Align payload
5. Write to memory
6. Return status

---

## 6. Packet Format (Version 1)

### 6.1 Endianness
- Header fields: Big‑Endian
- Register payload: Big‑Endian uint16

---

### 6.2 Frame Layout

```
[ Magic(2) ][ Ver(1) ][ Flags(1) ]
[ Area(1) ][ Rsv(1) ][ MemoryID(2) ]
[ Address(2) ][ Count(2) ]
[ Payload(N) ]
[ CRC32(4) ]
```

---

### 6.3 Header Fields

- Magic: `0x52 0x49` (`RI`)
- Version: `0x01`
- Flags: ignored
- Area:
  - `0x01` Coils
  - `0x02` Discrete Inputs
  - `0x03` Holding Registers
  - `0x04` Input Registers
- MemoryID: MMA memory selector
- Address: zero‑based internal address
- Count:
  - bits for coils/discretes
  - registers for holding/input

---

### 6.4 Payload Length

- Coils / Discretes: `ceil(count / 8)` bytes
- Registers: `count * 2` bytes

Mismatch → reject.

---

### 6.5 CRC

- CRC32 (IEEE)
- Covers everything except CRC field
- CRC failure → reject

---

## 7. Write Semantics

- All‑or‑nothing writes
- No retries
- No memory modification on failure

If Raw Ingest stops sending:
- Memory remains unchanged

---

## 8. Response

Single‑byte response:
- `0x00` OK
- `0x01` REJECTED

Error reason is logged internally only.

---

## 9. Configuration Surface

Raw Ingest owns **only socket safety configuration**:

```yaml
raw_ingest:
  enabled: true
  listen: ":9000"
  max_packet_bytes: 4096
```

Notes:
- No semantic timeouts
- No freshness timers
- No retry logic

Any timeout is **purely a TCP read‑deadline for socket safety**, not a data rule.

---

## 10. Minimal Coding Checklist

- [ ] TCP listener
- [ ] Packet size guard
- [ ] CRC32 validation
- [ ] Area/address/count bounds check
- [ ] Payload alignment only
- [ ] Direct write to `core.Memory`
- [ ] Atomic reject on error

Nothing else.

---

## 11. Final Lock

> **Raw Ingest is intentionally dumb.**  
> This is the feature that makes the system scalable and safe.

---

**End of Raw Ingest Rules**

