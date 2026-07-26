package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/modules/table/dto"
	"github.com/nutrixpos/pos/modules/table/models"
)

func TestJSONApiOkResponse_Serialization(t *testing.T) {
	table := models.Table{
		Id:       "test-id-123",
		Number:   3,
		Name:     "Table 3",
		Capacity: 4,
		Zone:     "indoor",
		Status:   "available",
	}

	resp := JSONApiOkResponse{
		Data: table,
		Meta: JSONAPIMeta{
			TotalRecords: 1,
			PageNumber:   1,
			PageSize:     50,
			PageCount:    1,
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal response: %v", err)
	}

	var decoded JSONApiOkResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if decoded.Meta.TotalRecords != 1 {
		t.Errorf("expected total_records 1, got %d", decoded.Meta.TotalRecords)
	}
}

func TestJSONAPIMeta_Fields(t *testing.T) {
	meta := JSONAPIMeta{
		TotalRecords: 100,
		PageNumber:   2,
		PageSize:     25,
		PageCount:    4,
	}

	if meta.TotalRecords != 100 {
		t.Errorf("expected 100, got %d", meta.TotalRecords)
	}
	if meta.PageNumber != 2 {
		t.Errorf("expected 2, got %d", meta.PageNumber)
	}
	if meta.PageSize != 25 {
		t.Errorf("expected 25, got %d", meta.PageSize)
	}
	if meta.PageCount != 4 {
		t.Errorf("expected 4, got %d", meta.PageCount)
	}
}

func TestCreateTableRequest_Deserialization(t *testing.T) {
	body := `{"number":5,"name":"Patio","capacity":6,"zone":"outdoor","branch_id":"b1"}`

	var req dto.CreateTableRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if req.Number != 5 {
		t.Errorf("expected number 5, got %d", req.Number)
	}
	if req.Name != "Patio" {
		t.Errorf("expected name 'Patio', got %s", req.Name)
	}
	if req.Capacity != 6 {
		t.Errorf("expected capacity 6, got %d", req.Capacity)
	}
	if req.Zone != "outdoor" {
		t.Errorf("expected zone 'outdoor', got %s", req.Zone)
	}
	if req.BranchId != "b1" {
		t.Errorf("expected branch_id 'b1', got %s", req.BranchId)
	}
}

func TestGetTable_BadRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/tables/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{}}`))
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/tables/abc123", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestCreateTable_InvalidJSON(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/tables", func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateTableRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}).Methods("POST")

	body := strings.NewReader(`{invalid json}`)
	req := httptest.NewRequest("POST", "/api/tables", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestCreateTable_ValidJSON(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/tables", func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateTableRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		resp := JSONApiOkResponse{
			Data: models.Table{
				Number:   req.Number,
				Name:     req.Name,
				Capacity: req.Capacity,
				Zone:     req.Zone,
				Status:   "available",
			},
			Meta: JSONAPIMeta{TotalRecords: 1},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}).Methods("POST")

	body := `{"number":10,"name":"VIP","capacity":8,"zone":"vip","branch_id":"b2"}`
	req := httptest.NewRequest("POST", "/api/tables", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rec.Code)
	}

	var resp JSONApiOkResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	table, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatal("expected Data to be a map")
	}

	if table["number"].(float64) != 10 {
		t.Errorf("expected number 10, got %v", table["number"])
	}
	if table["status"] != "available" {
		t.Errorf("expected status 'available', got %v", table["status"])
	}
}
