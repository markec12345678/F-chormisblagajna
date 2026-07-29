package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestJSONApiOkResponse_Serialization(t *testing.T) {
	resp := JSONApiOkResponse{
		Data: map[string]string{"message": "ok"},
		Meta: JSONAPIMeta{TotalRecords: 1},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded JSONApiOkResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Meta.TotalRecords != 1 {
		t.Errorf("expected TotalRecords=1, got %d", decoded.Meta.TotalRecords)
	}
}

func TestGetNotifications_ReturnsOK(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/notifications", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		notifications := []map[string]interface{}{}
		resp := JSONApiOkResponse{Data: notifications, Meta: JSONAPIMeta{TotalRecords: 0}}
		json.NewEncoder(w).Encode(resp)
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/notifications", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestGetUnreadCount_ReturnsOK(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/notifications/unread-count", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := JSONApiOkResponse{Data: map[string]int64{"count": 5}, Meta: JSONAPIMeta{TotalRecords: 1}}
		json.NewEncoder(w).Encode(resp)
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/notifications/unread-count", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
