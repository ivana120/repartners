package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ivana120/repartners/internal/calculator"
)

func TestHandler_Calculate(t *testing.T) {
	calc := calculator.New([]int{250, 500, 1000, 2000, 5000})
	handler := NewHandler(calc)

	tests := []struct {
		name           string
		method         string
		body           interface{}
		wantStatus     int
		wantSuccess    bool
		wantTotalItems int
	}{
		{
			name:           "valid order",
			method:         http.MethodPost,
			body:           CalculateRequest{OrderQty: 251},
			wantStatus:     http.StatusOK,
			wantSuccess:    true,
			wantTotalItems: 500,
		},
		{
			name:        "invalid method",
			method:      http.MethodGet,
			body:        CalculateRequest{OrderQty: 100},
			wantStatus:  http.StatusMethodNotAllowed,
			wantSuccess: false,
		},
		{
			name:        "invalid order quantity",
			method:      http.MethodPost,
			body:        CalculateRequest{OrderQty: 0},
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
		},
		{
			name:        "negative order quantity",
			method:      http.MethodPost,
			body:        CalculateRequest{OrderQty: -10},
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(tt.method, "/api/calculate", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.Calculate(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status: want %d, got %d", tt.wantStatus, rec.Code)
			}

			var resp CalculateResponse
			json.NewDecoder(rec.Body).Decode(&resp)

			if resp.Success != tt.wantSuccess {
				t.Errorf("success: want %v, got %v", tt.wantSuccess, resp.Success)
			}

			if tt.wantSuccess && resp.TotalItems != tt.wantTotalItems {
				t.Errorf("total_items: want %d, got %d", tt.wantTotalItems, resp.TotalItems)
			}
		})
	}
}

func TestHandler_GetPackSizes(t *testing.T) {
	sizes := []int{250, 500, 1000}
	calc := calculator.New(sizes)
	handler := NewHandler(calc)

	req := httptest.NewRequest(http.MethodGet, "/api/pack-sizes", nil)
	rec := httptest.NewRecorder()

	handler.GetPackSizes(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status: want %d, got %d", http.StatusOK, rec.Code)
	}

	var resp PackSizesResponse
	json.NewDecoder(rec.Body).Decode(&resp)

	if !resp.Success {
		t.Error("expected success to be true")
	}

	if len(resp.PackSizes) != len(sizes) {
		t.Errorf("pack sizes count: want %d, got %d", len(sizes), len(resp.PackSizes))
	}
}

func TestHandler_UpdatePackSizes(t *testing.T) {
	calc := calculator.New([]int{100, 200})
	handler := NewHandler(calc)

	tests := []struct {
		name        string
		body        interface{}
		wantStatus  int
		wantSuccess bool
	}{
		{
			name:        "valid update",
			body:        PackSizesRequest{PackSizes: []int{50, 100, 150}},
			wantStatus:  http.StatusOK,
			wantSuccess: true,
		},
		{
			name:        "empty pack sizes",
			body:        PackSizesRequest{PackSizes: []int{}},
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
		},
		{
			name:        "negative pack size",
			body:        PackSizesRequest{PackSizes: []int{50, -10}},
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
		},
		{
			name:        "zero pack size",
			body:        PackSizesRequest{PackSizes: []int{50, 0}},
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPut, "/api/pack-sizes", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.UpdatePackSizes(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status: want %d, got %d", tt.wantStatus, rec.Code)
			}

			var resp PackSizesResponse
			json.NewDecoder(rec.Body).Decode(&resp)

			if resp.Success != tt.wantSuccess {
				t.Errorf("success: want %v, got %v", tt.wantSuccess, resp.Success)
			}
		})
	}
}

func TestHandler_HealthCheck(t *testing.T) {
	calc := calculator.New([]int{100})
	handler := NewHandler(calc)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler.HealthCheck(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status: want %d, got %d", http.StatusOK, rec.Code)
	}

	var resp map[string]string
	json.NewDecoder(rec.Body).Decode(&resp)

	if resp["status"] != "healthy" {
		t.Errorf("status: want 'healthy', got '%s'", resp["status"])
	}
}

func TestEnableCORS(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	corsHandler := EnableCORS(handler)

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	rec := httptest.NewRecorder()
	corsHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("OPTIONS status: want %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("expected CORS header to be set")
	}
}
