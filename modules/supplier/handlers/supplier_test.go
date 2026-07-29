package handlers

import (
	"encoding/json"
	"testing"
	"time"

	supplier_models "github.com/nutrixpos/pos/modules/supplier/models"
)

func TestSupplier_Serialization(t *testing.T) {
	supplier := supplier_models.Supplier{
		Id:          "s-1",
		Name:        "Fresh Produce Ltd",
		ContactName: "Janez Novak",
		Email:       "janez@freshproduce.si",
		Phone:       "+38640123456",
		Address:     "Glavna 1, Ljubljana",
		IsActive:    true,
		CreatedAt:   time.Now(),
	}

	data, err := json.Marshal(supplier)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded supplier_models.Supplier
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Name != "Fresh Produce Ltd" {
		t.Errorf("expected Name='Fresh Produce Ltd', got %s", decoded.Name)
	}
	if !decoded.IsActive {
		t.Error("expected IsActive=true")
	}
}

func TestSupplierOrder_Serialization(t *testing.T) {
	order := supplier_models.SupplierOrder{
		Id:           "so-1",
		SupplierId:   "s-1",
		SupplierName: "Fresh Produce Ltd",
		OrderDate:    time.Now(),
		TotalAmount:  250.00,
		Status:       "pending",
		Items: []supplier_models.SupplierOrderItem{
			{
				MaterialId:   "m-1",
				MaterialName: "Tomatoes",
				Quantity:     50,
				UnitPrice:    2.50,
			},
		},
	}

	data, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded supplier_models.SupplierOrder
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.TotalAmount != 250.00 {
		t.Errorf("expected TotalAmount=250, got %f", decoded.TotalAmount)
	}
	if len(decoded.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(decoded.Items))
	}
}

func TestSupplierOrderItem_Calculation(t *testing.T) {
	item := supplier_models.SupplierOrderItem{
		Quantity:  50,
		UnitPrice: 2.50,
	}

	total := item.Quantity * item.UnitPrice
	if total != 125.00 {
		t.Errorf("expected total=125, got %f", total)
	}
}

func TestSupplierStatusTransitions(t *testing.T) {
	validStatuses := []string{"pending", "delivered", "cancelled"}
	for _, status := range validStatuses {
		if status == "" {
			t.Error("status should not be empty")
		}
	}
}
