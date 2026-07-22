package dto

import (
	"encoding/json"
	"testing"
)

func TestOrderItemRefundRequest_JSON(t *testing.T) {
	request := OrderItemRefundRequest{
		OrderId:     "order-123",
		ItemId:      "item-1",
		ProductId:   "prod-1",
		Reason:      "Wrong order",
		RefundValue: 15.50,
		Destination: DTOOrderItemRefundDestination_Inventory,
		MaterialRefunds: []OrderItemRefundMaterialDTO{
			{
				MaterialId:         "mat-1",
				EntryId:            "entry-1",
				InventoryReturnQty: 2,
				DisposeQty:         0,
				WasteQty:           0,
				Comment:            "returned to stock",
			},
		},
		ProductAdd: []OrderItemRefundProductAddDTO{
			{
				ProductId: "prod-2",
				Quantity:  1,
				Comment:   "replacement",
			},
		},
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded OrderItemRefundRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.OrderId != request.OrderId {
		t.Errorf("OrderId = %v, want %v", decoded.OrderId, request.OrderId)
	}
	if decoded.RefundValue != request.RefundValue {
		t.Errorf("RefundValue = %v, want %v", decoded.RefundValue, request.RefundValue)
	}
	if len(decoded.MaterialRefunds) != 1 {
		t.Errorf("MaterialRefunds length = %v, want 1", len(decoded.MaterialRefunds))
	}
	if len(decoded.ProductAdd) != 1 {
		t.Errorf("ProductAdd length = %v, want 1", len(decoded.ProductAdd))
	}
}

func TestRefundDestination_Constants(t *testing.T) {
	if DTOOrderItemRefundDestination_Inventory != "inventory" {
		t.Errorf("Inventory constant = %v, want 'inventory'", DTOOrderItemRefundDestination_Inventory)
	}
	if DTOOrderItemRefundDestination_Disposals != "disposals" {
		t.Errorf("Disposals constant = %v, want 'disposals'", DTOOrderItemRefundDestination_Disposals)
	}
	if DTOOrderItemRefundDestination_Waste != "waste" {
		t.Errorf("Waste constant = %v, want 'waste'", DTOOrderItemRefundDestination_Waste)
	}
	if DTOOrderItemRefundDestination_Custom != "custom" {
		t.Errorf("Custom constant = %v, want 'custom'", DTOOrderItemRefundDestination_Custom)
	}
}

func TestOrderItemRefundMaterialDTO_JSON(t *testing.T) {
	dto := OrderItemRefundMaterialDTO{
		MaterialId:         "mat-1",
		EntryId:            "entry-1",
		InventoryReturnQty: 1.5,
		DisposeQty:         0.5,
		WasteQty:           0,
		Comment:            "partial return",
	}

	data, err := json.Marshal(dto)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded OrderItemRefundMaterialDTO
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.MaterialId != dto.MaterialId {
		t.Errorf("MaterialId = %v, want %v", decoded.MaterialId, dto.MaterialId)
	}
	if decoded.InventoryReturnQty != dto.InventoryReturnQty {
		t.Errorf("InventoryReturnQty = %v, want %v", decoded.InventoryReturnQty, dto.InventoryReturnQty)
	}
}
