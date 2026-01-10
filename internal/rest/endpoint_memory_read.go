// File: endpoint_memory_read.go
// Endpoint: GET /api/v1/memory/read
// Purpose: Read memory values (explicit memory selection)

package rest

import (
	"net/http"
	"strconv"
)

func (h *Handlers) HandleMemoryRead(w http.ResponseWriter, r *http.Request) {
	if !h.EnableRead {
		writeJSON(w, http.StatusForbidden, reject("read disabled"))
		return
	}

	memName := r.URL.Query().Get("memory")
	if memName == "" {
		writeJSON(w, http.StatusBadRequest, reject("memory is required"))
		return
	}

	mem, ok := h.Memories[memName]
	if !ok {
		writeJSON(w, http.StatusNotFound, reject("memory not found"))
		return
	}

	area := r.URL.Query().Get("area")
	addrStr := r.URL.Query().Get("address")
	countStr := r.URL.Query().Get("count")

	if area == "" || addrStr == "" || countStr == "" {
		writeJSON(w, http.StatusBadRequest, reject("area, address, and count are required"))
		return
	}

	addr, err := strconv.Atoi(addrStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, reject("invalid address"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, reject("invalid count"))
		return
	}

	switch area {
	case "coils":
		vals, err := mem.ReadCoils(addr, count)
		if err != nil {
			writeIngestError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"values": vals})

	case "discrete_inputs":
		vals, err := mem.ReadDiscreteInputs(addr, count)
		if err != nil {
			writeIngestError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"values": vals})

	case "holding_registers":
		vals, err := mem.ReadHoldingRegs(addr, count)
		if err != nil {
			writeIngestError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"values": vals})

	case "input_registers":
		vals, err := mem.ReadInputRegs(addr, count)
		if err != nil {
			writeIngestError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"values": vals})

	default:
		writeJSON(w, http.StatusBadRequest, reject("unknown area"))
	}
}
