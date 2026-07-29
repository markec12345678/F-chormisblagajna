package handlers

import (
	"encoding/json"
	"testing"
	"time"

	auditlog_models "github.com/nutrixpos/pos/modules/auditlog/models"
)

func TestAuditLogEntry_Serialization(t *testing.T) {
	entry := auditlog_models.AuditLogEntry{
		Id:         "a-1",
		Action:     "create",
		Resource:   "product",
		ResourceId: "p-1",
		UserId:     "u-1",
		Username:   "admin",
		Details:    map[string]string{"name": "Pizza Margherita", "price": "8.50"},
		IpAddress:  "127.0.0.1",
		CreatedAt:  time.Now(),
	}

	data, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded auditlog_models.AuditLogEntry
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Action != "create" {
		t.Errorf("expected Action='create', got %s", decoded.Action)
	}
	if decoded.Resource != "product" {
		t.Errorf("expected Resource='product', got %s", decoded.Resource)
	}
	if decoded.Details["name"] != "Pizza Margherita" {
		t.Errorf("expected detail name='Pizza Margherita', got %s", decoded.Details["name"])
	}
}

func TestAuditLogSummary_Serialization(t *testing.T) {
	summary := auditlog_models.AuditLogSummary{
		TotalEntries: 150,
		ByAction: []auditlog_models.ActionSummary{
			{Action: "create", Count: 50},
			{Action: "update", Count: 70},
			{Action: "delete", Count: 30},
		},
		ByResource: []auditlog_models.ResourceSummary{
			{Resource: "order", Count: 80},
			{Resource: "product", Count: 70},
		},
	}

	data, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded auditlog_models.AuditLogSummary
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.TotalEntries != 150 {
		t.Errorf("expected TotalEntries=150, got %d", decoded.TotalEntries)
	}
	if len(decoded.ByAction) != 3 {
		t.Errorf("expected 3 actions, got %d", len(decoded.ByAction))
	}
}

func TestAuditActions(t *testing.T) {
	validActions := []string{"create", "update", "delete", "login", "logout", "fiscalize"}
	for _, action := range validActions {
		if action == "" {
			t.Error("action should not be empty")
		}
	}
}

func TestAuditResources(t *testing.T) {
	validResources := []string{"order", "product", "material", "customer", "settings", "user", "category"}
	for _, resource := range validResources {
		if resource == "" {
			t.Error("resource should not be empty")
		}
	}
}
