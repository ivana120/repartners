package api

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/ivana120/repartners/internal/calculator"
)

type Handler struct {
	calc *calculator.Calculator
	mu   sync.RWMutex
}

func NewHandler(calc *calculator.Calculator) *Handler {
	return &Handler{calc: calc}
}

type CalculateRequest struct {
	OrderQty int `json:"order_qty"`
}

type CalculateResponse struct {
	Success    bool        `json:"success"`
	Packs      map[int]int `json:"packs,omitempty"`
	TotalItems int         `json:"total_items,omitempty"`
	TotalPacks int         `json:"total_packs,omitempty"`
	Error      string      `json:"error,omitempty"`
}

type PackSizesRequest struct {
	PackSizes []int `json:"pack_sizes"`
}

type PackSizesResponse struct {
	Success   bool   `json:"success"`
	PackSizes []int  `json:"pack_sizes,omitempty"`
	Error     string `json:"error,omitempty"`
}

func (h *Handler) Calculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	h.mu.RLock()
	result, err := h.calc.Calculate(req.OrderQty)
	h.mu.RUnlock()

	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, CalculateResponse{
		Success:    true,
		Packs:      result.Packs,
		TotalItems: result.TotalItems,
		TotalPacks: result.TotalPacks,
	})
}

func (h *Handler) GetPackSizes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	h.mu.RLock()
	sizes := h.calc.PackSizes()
	h.mu.RUnlock()

	h.writeJSON(w, http.StatusOK, PackSizesResponse{
		Success:   true,
		PackSizes: sizes,
	})
}

func (h *Handler) UpdatePackSizes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req PackSizesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(req.PackSizes) == 0 {
		h.writeError(w, http.StatusBadRequest, "pack sizes cannot be empty")
		return
	}

	for _, size := range req.PackSizes {
		if size <= 0 {
			h.writeError(w, http.StatusBadRequest, "pack sizes must be positive")
			return
		}
	}

	h.mu.Lock()
	h.calc.SetPackSizes(req.PackSizes)
	sizes := h.calc.PackSizes()
	h.mu.Unlock()

	h.writeJSON(w, http.StatusOK, PackSizesResponse{
		Success:   true,
		PackSizes: sizes,
	})
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
	})
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]interface{}{
		"success": false,
		"error":   message,
	})
}

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
