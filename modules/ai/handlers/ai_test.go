package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestJSONApiOkResponse_Serialization(t *testing.T) {
	resp := JSONApiOkResponse{
		Data: map[string]string{"test": "value"},
		Meta: JSONAPIMeta{
			TotalRecords: 1,
			PageNumber:   1,
			PageSize:     20,
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

func TestAISearch_MissingQuery(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/ai/search", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Query string `json:"query"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if req.Query == "" {
			http.Error(w, "query is required", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	body := `{"query":"","branch_id":"b1","language":"sl","limit":10}`
	req := httptest.NewRequest("POST", "/api/ai/search", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty query, got %d", rec.Code)
	}
}

func TestAISearch_InvalidJSON(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/ai/search", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Query string `json:"query"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	req := httptest.NewRequest("POST", "/api/ai/search", strings.NewReader(`{invalid}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid JSON, got %d", rec.Code)
	}
}

func TestVoiceOrder_MissingAudio(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/ai/voice", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			AudioBase64 string `json:"audio_base64"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if req.AudioBase64 == "" {
			http.Error(w, "audio_base64 is required", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	body := `{"audio_base64":"","language":"sl","branch_id":"b1"}`
	req := httptest.NewRequest("POST", "/api/ai/voice", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty audio, got %d", rec.Code)
	}
}

func TestVoiceOrder_ValidRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/ai/voice", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			AudioBase64 string `json:"audio_base64"`
			Language    string `json:"language"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if req.AudioBase64 == "" {
			http.Error(w, "audio_base64 is required", http.StatusBadRequest)
			return
		}

		resp := JSONApiOkResponse{
			Data: map[string]interface{}{
				"transcript": "test",
				"items":      []interface{}{},
				"confidence": 0.0,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}).Methods("POST")

	body := `{"audio_base64":"dGVzdA==","language":"sl","branch_id":"b1"}`
	req := httptest.NewRequest("POST", "/api/ai/voice", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 for valid request, got %d", rec.Code)
	}

	var resp JSONApiOkResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
}

func TestSmartSuggestions_DefaultLimit(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/ai/suggestions", func(w http.ResponseWriter, r *http.Request) {
		resp := JSONApiOkResponse{
			Data: []interface{}{},
			Meta: JSONAPIMeta{TotalRecords: 0},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/ai/suggestions", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var resp JSONApiOkResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
}

func TestSmartSuggestions_WithBranchId(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/ai/suggestions", func(w http.ResponseWriter, r *http.Request) {
		branchId := r.URL.Query().Get("branch_id")
		if branchId != "branch-123" {
			t.Errorf("expected branch_id 'branch-123', got '%s'", branchId)
		}

		resp := JSONApiOkResponse{
			Data: []interface{}{},
			Meta: JSONAPIMeta{TotalRecords: 0},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/ai/suggestions?branch_id=branch-123", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
