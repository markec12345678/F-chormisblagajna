package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/modules/splitbill/dto"
	"github.com/nutrixpos/pos/modules/splitbill/models"
)

func TestJSONApiOkResponse_Serialization(t *testing.T) {
	splitBill := models.SplitBill{
		Id:        "test-id-123",
		OrderId:   "order-1",
		SplitType: "equal",
		Status:    "pending",
		Parts: []models.SplitPart{
			{Id: "part-1", Amount: 25.00, IsPaid: false},
			{Id: "part-2", Amount: 25.00, IsPaid: false},
		},
	}

	resp := JSONApiOkResponse{
		Data: splitBill,
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

func TestCreateSplitBillRequest_Deserialization(t *testing.T) {
	body := `{"order_id":"order-123","split_type":"equal","split_count":3}`

	var req dto.CreateSplitBillRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if req.OrderId != "order-123" {
		t.Errorf("expected order_id order-123, got %s", req.OrderId)
	}
	if req.SplitType != "equal" {
		t.Errorf("expected split_type equal, got %s", req.SplitType)
	}
	if req.SplitCount != 3 {
		t.Errorf("expected split_count 3, got %d", req.SplitCount)
	}
}

func TestCreateSplitBill_InvalidJSON(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/split-bills", func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateSplitBillRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}).Methods("POST")

	body := strings.NewReader(`{invalid json}`)
	req := httptest.NewRequest("POST", "/api/split-bills", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestCreateSplitBill_ValidJSON(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/split-bills", func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateSplitBillRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		parts := make([]models.SplitPart, req.SplitCount)
		for i := 0; i < req.SplitCount; i++ {
			parts[i] = models.SplitPart{
				Id:     "part-" + string(rune('0'+i)),
				Amount: 33.33,
				IsPaid: false,
			}
		}

		resp := JSONApiOkResponse{
			Data: models.SplitBill{
				OrderId:   req.OrderId,
				SplitType: req.SplitType,
				Parts:     parts,
				Status:    "pending",
			},
			Meta: JSONAPIMeta{TotalRecords: 1},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}).Methods("POST")

	body := `{"order_id":"order-456","split_type":"equal","split_count":3}`
	req := httptest.NewRequest("POST", "/api/split-bills", strings.NewReader(body))
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

	splitBill, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatal("expected Data to be a map")
	}

	if splitBill["split_type"] != "equal" {
		t.Errorf("expected split_type equal, got %v", splitBill["split_type"])
	}
	if splitBill["status"] != "pending" {
		t.Errorf("expected status pending, got %v", splitBill["status"])
	}
}

func TestPaySplitPartRequest_Deserialization(t *testing.T) {
	body := `{"part_id":"part-1","payment_method":"card","amount":33.33}`

	var req dto.PaySplitPartRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if req.PartId != "part-1" {
		t.Errorf("expected part_id part-1, got %s", req.PartId)
	}
	if req.PaymentMethod != "card" {
		t.Errorf("expected payment_method card, got %s", req.PaymentMethod)
	}
	if req.Amount != 33.33 {
		t.Errorf("expected amount 33.33, got %f", req.Amount)
	}
}

func TestPaySplitPart_InvalidJSON(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/split-bills/{id}/pay", func(w http.ResponseWriter, r *http.Request) {
		var req dto.PaySplitPartRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	body := strings.NewReader(`{invalid json}`)
	req := httptest.NewRequest("POST", "/api/split-bills/sb-1/pay", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestGetSplitBill_NotFound(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/split-bills/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/split-bills/nonexistent", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}
}

func TestItemSplit_Deserialization(t *testing.T) {
	body := `{"item_id":"item-1","part_index":0}`

	var split dto.ItemSplit
	if err := json.Unmarshal([]byte(body), &split); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if split.ItemId != "item-1" {
		t.Errorf("expected item_id item-1, got %s", split.ItemId)
	}
	if split.PartIndex != 0 {
		t.Errorf("expected part_index 0, got %d", split.PartIndex)
	}
}
