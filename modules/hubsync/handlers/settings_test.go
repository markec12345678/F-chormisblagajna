package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/modules/hubsync/models"
)

func TestGetSettings_ReturnsJSON(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/hubsync", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		hub := models.Hubsync{
			Id: "test-id",
			Settings: models.Settings{
				Enabled:      true,
				ServerHost:   "https://hub.example.com",
				SyncInterval: 60,
				BufferSize:   100,
			},
			LastSynced:   1000,
			SyncProgress: 50,
		}
		data := struct {
			Data models.Hubsync `json:"data"`
		}{Data: hub}
		json.NewEncoder(w).Encode(data)
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/hubsync", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var resp struct {
		Data models.Hubsync `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp.Data.Settings.ServerHost != "https://hub.example.com" {
		t.Errorf("expected hub host, got %s", resp.Data.Settings.ServerHost)
	}
	if !resp.Data.Settings.Enabled {
		t.Errorf("expected enabled=true")
	}
}

func TestPatchSettings_InvalidJSON(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/hubsync", func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Data models.Hubsync `json:"data"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}).Methods("PATCH")

	req := httptest.NewRequest("PATCH", "/api/hubsync", strings.NewReader("invalid json"))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid JSON, got %d", rec.Code)
	}
}

func TestHubsyncSettings_Serialization(t *testing.T) {
	settings := models.Settings{
		Enabled:       true,
		ServerHost:    "https://hub.example.com",
		Token:         "secret-token",
		SyncInterval:  120,
		BufferSize:    200,
		SyncSales:     true,
		SyncInventory: false,
	}

	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded models.Settings
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.ServerHost != "https://hub.example.com" {
		t.Errorf("expected server host, got %s", decoded.ServerHost)
	}
	if decoded.SyncInterval != 120 {
		t.Errorf("expected 120, got %d", decoded.SyncInterval)
	}
	if !decoded.SyncSales {
		t.Errorf("expected sync_sales=true")
	}
	if decoded.SyncInventory {
		t.Errorf("expected sync_inventory=false")
	}
}

func TestHubsyncModel_Serialization(t *testing.T) {
	hub := models.Hubsync{
		Id: "507f1f77bcf86cd799439011",
		Settings: models.Settings{
			Enabled:    true,
			BufferSize: 100,
		},
		LastSynced:   1700000000,
		SyncProgress: 75.5,
	}

	data, err := json.Marshal(hub)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded models.Hubsync
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Id != "507f1f77bcf86cd799439011" {
		t.Errorf("expected id, got %s", decoded.Id)
	}
	if decoded.SyncProgress != 75.5 {
		t.Errorf("expected 75.5, got %f", decoded.SyncProgress)
	}
	if !decoded.Settings.Enabled {
		t.Errorf("expected enabled=true")
	}
}

func TestHubsyncModel_DefaultValues(t *testing.T) {
	hub := models.Hubsync{}
	if hub.SyncProgress != 0 {
		t.Errorf("expected 0, got %f", hub.SyncProgress)
	}
	if hub.Settings.BufferSize != 0 {
		t.Errorf("expected 0, got %d", hub.Settings.BufferSize)
	}
}

func TestSettings_Tags(t *testing.T) {
	settings := models.Settings{}
	data, _ := json.Marshal(settings)
	var raw map[string]interface{}
	json.Unmarshal(data, &raw)

	expectedFields := []string{"enabled", "server_host", "token", "sync_interval", "buffer_size", "sync_sales", "sync_inventory"}
	for _, field := range expectedFields {
		if _, ok := raw[field]; !ok {
			t.Errorf("expected field %s in JSON output", field)
		}
	}
}
