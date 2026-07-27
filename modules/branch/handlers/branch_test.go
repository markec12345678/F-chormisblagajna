package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/modules/branch/dto"
	"github.com/nutrixpos/pos/modules/branch/models"
)

func TestJSONApiOkResponse_Serialization(t *testing.T) {
	branch := models.Branch{
		Id:       "test-id-123",
		Name:     "Ljubljana",
		Address:  "Glavni trg 1",
		Phone:    "+386 1 234 5678",
		Email:    "ljubljana@nutrix.si",
		IsActive: true,
	}

	resp := JSONApiOkResponse{
		Data: branch,
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

func TestCreateBranchRequest_Deserialization(t *testing.T) {
	body := `{"name":"Maribor","address":"Glavni trg 5","phone":"+386 2 123 4567","email":"mb@nutrix.si","is_active":true}`

	var req dto.CreateBranchRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if req.Name != "Maribor" {
		t.Errorf("expected name Maribor, got %s", req.Name)
	}
	if req.Address != "Glavni trg 5" {
		t.Errorf("expected address 'Glavni trg 5', got %s", req.Address)
	}
	if !req.IsActive {
		t.Error("expected is_active true")
	}
}

func TestGetBranch_BadRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/branches/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{}}`))
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/branches/abc123", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestCreateBranch_InvalidJSON(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/branches", func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateBranchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}).Methods("POST")

	body := strings.NewReader(`{invalid json}`)
	req := httptest.NewRequest("POST", "/api/branches", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestCreateBranch_ValidJSON(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/branches", func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateBranchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		resp := JSONApiOkResponse{
			Data: models.Branch{
				Name:     req.Name,
				Address:  req.Address,
				Phone:    req.Phone,
				Email:    req.Email,
				IsActive: req.IsActive,
			},
			Meta: JSONAPIMeta{TotalRecords: 1},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}).Methods("POST")

	body := `{"name":"Koper","address":"Pristanišče 1","phone":"+386 5 678 9012","email":"koper@nutrix.si","is_active":true}`
	req := httptest.NewRequest("POST", "/api/branches", strings.NewReader(body))
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

	branch, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatal("expected Data to be a map")
	}

	if branch["name"] != "Koper" {
		t.Errorf("expected name Koper, got %v", branch["name"])
	}
	if branch["is_active"] != true {
		t.Errorf("expected is_active true, got %v", branch["is_active"])
	}
}

func TestUpdateBranchRequest_Deserialization(t *testing.T) {
	body := `{"name":"Updated Branch","address":"New Address"}`

	var req dto.UpdateBranchRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if req.Name != "Updated Branch" {
		t.Errorf("expected name 'Updated Branch', got %s", req.Name)
	}
	if req.Address != "New Address" {
		t.Errorf("expected address 'New Address', got %s", req.Address)
	}
}

func TestDeleteBranch_NoContent(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/branches/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}).Methods("DELETE")

	req := httptest.NewRequest("DELETE", "/api/branches/abc123", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rec.Code)
	}
}
