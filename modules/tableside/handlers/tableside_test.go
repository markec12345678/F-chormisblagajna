package handlers

import (
	"encoding/json"
	"testing"
	"time"

	ts_models "github.com/nutrixpos/pos/modules/tableside/models"
)

func TestTableSession_Serialization(t *testing.T) {
	session := ts_models.TableSession{
		Id:         "ts-1",
		TableLabel: "Table 5",
		Zone:       "Main Hall",
		QrToken:    "abc123",
		Active:     true,
		GuestCount: 4,
		OpenedAt:   time.Now(),
	}

	data, err := json.Marshal(session)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded ts_models.TableSession
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.TableLabel != "Table 5" {
		t.Errorf("expected TableLabel='Table 5', got %s", decoded.TableLabel)
	}
	if !decoded.Active {
		t.Error("expected Active=true")
	}
}

func TestTableOrder_Serialization(t *testing.T) {
	order := ts_models.TableOrder{
		Id:        "to-1",
		SessionId: "ts-1",
		Items: []ts_models.TableOrderItem{
			{ProductId: "p1", ProductName: "Burger", Quantity: 2, UnitPrice: 8.50},
			{ProductId: "p2", ProductName: "Fries", Quantity: 1, UnitPrice: 3.50},
		},
		Status:   "pending",
		Subtotal: 20.50,
		PlacedAt: time.Now(),
	}

	data, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded ts_models.TableOrder
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(decoded.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(decoded.Items))
	}
	if decoded.Subtotal != 20.50 {
		t.Errorf("expected Subtotal=20.50, got %f", decoded.Subtotal)
	}
}

func TestTableOrderItem_Serialization(t *testing.T) {
	item := ts_models.TableOrderItem{
		ProductId:   "p1",
		ProductName: "Pizza",
		Quantity:    1,
		UnitPrice:   12.00,
		Notes:       "Extra cheese",
	}

	data, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded ts_models.TableOrderItem
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.ProductName != "Pizza" {
		t.Errorf("expected ProductName='Pizza', got %s", decoded.ProductName)
	}
}

func TestQrInfo_Serialization(t *testing.T) {
	qr := ts_models.QrInfo{
		TableLabel: "Table 3",
		Token:      "token123",
		Url:        "http://localhost/tableside/menu/token123",
		Host:       "http://localhost",
	}

	data, err := json.Marshal(qr)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded ts_models.QrInfo
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Url != "http://localhost/tableside/menu/token123" {
		t.Errorf("unexpected URL: %s", decoded.Url)
	}
}
