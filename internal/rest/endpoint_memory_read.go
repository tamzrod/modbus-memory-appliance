package rest

import (
	"net/http"
	"strconv"
	"strings"
)

func (h *Handlers) HandleMemoryRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, reject("method not allowed"))
		return
	}

	if !h.EnableRead {
		writeJSON(w, http.StatusForbidden, reject("read disabled"))
		return
	}

	area := strings.TrimSpace(r.URL.Query().Get("area"))
	addrStr := strings.TrimSpace(r.URL.Query().Get("address"))
	countStr := strings.TrimSpace(r.URL.Query().Get("count"))

	if area == "" || addrStr == "" || countStr == "" {
		h.Stats.IncRejected()
		writeJSON(w, http.StatusBadRequest, reject("missing query parameters"))
		return
	}

	address, err := strconv.Atoi(addrStr)
	if err != nil || address < 0 {
		h.Stats.IncRejected()
		writeJSON(w, http.StatusBadRequest, reject("invalid address"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 {
		h.Stats.IncRejected()
		writeJSON(w, http.StatusBadRequest, reject("invalid count"))
		return
	}

	h.Stats.IncReads()

	switch area {

	case "coils":
		v, err := h.Memory.ReadCoils(address, count)
		if err != nil {
			h.Stats.IncRejected()
			writeJSON(w, http.StatusUnprocessableEntity, reject(err.Error()))
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"memory":  "default",
			"area":    area,
			"address": address,
			"count":   count,
			"bools":   v,
		})

	case "discrete_inputs":
		v, err := h.Memory.ReadDiscreteInputs(address, count)
		if err != nil {
			h.Stats.IncRejected()
			writeJSON(w, http.StatusUnprocessableEntity, reject(err.Error()))
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"memory":  "default",
			"area":    area,
			"address": address,
			"count":   count,
			"bools":   v,
		})

	case "holding_registers":
		v, err := h.Memory.ReadHoldingRegs(address, count)
		if err != nil {
			h.Stats.IncRejected()
			writeJSON(w, http.StatusUnprocessableEntity, reject(err.Error()))
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"memory":  "default",
			"area":    area,
			"address": address,
			"count":   count,
			"values":  v,
		})

	case "input_registers":
		v, err := h.Memory.ReadInputRegs(address, count)
		if err != nil {
			h.Stats.IncRejected()
			writeJSON(w, http.StatusUnprocessableEntity, reject(err.Error()))
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"memory":  "default",
			"area":    area,
			"address": address,
			"count":   count,
			"values":  v,
		})

	default:
		h.Stats.IncRejected()
		writeJSON(w, http.StatusBadRequest, reject("invalid area"))
	}
}
