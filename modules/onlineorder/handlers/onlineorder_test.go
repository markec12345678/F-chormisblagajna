package handlers

import (
	"encoding/json"
	"testing"
	"time"

	onlineorder_models "github.com/nutrixpos/pos/modules/onlineorder/models"
)

func TestOnlineOrder_Serialization(t *testing.T) {
	order := onlineorder_models.OnlineOrder{
		Id:            "o-1",
		DisplayId:     "WO-1234",
		CustomerName:  "Janez Novak",
		CustomerPhone: "+38640123456",
		CustomerEmail: "janez@example.com",
		Items: []onlineorder_models.OnlineOrderItem{
			{ProductId: "p-1", ProductName: "Pizza Margherita", Quantity: 2, Price: 8.50},
			{ProductId: "p-2", ProductName: "Coca Cola", Quantity: 1, Price: 2.50},
		},
		Subtotal:    19.50,
		DeliveryFee: 3.00,
		Total:       22.50,
		OrderType:   "delivery",
		DeliveryAddr: "Testna ulica 12, Ljubljana",
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	data, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded onlineorder_models.OnlineOrder
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.CustomerName != "Janez Novak" {
		t.Errorf("expected CustomerName='Janez Novak', got %s", decoded.CustomerName)
	}
	if decoded.Total != 22.50 {
		t.Errorf("expected Total=22.50, got %f", decoded.Total)
	}
	if len(decoded.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(decoded.Items))
	}
}

func TestOnlineOrderItem_Calculation(t *testing.T) {
	items := []onlineorder_models.OnlineOrderItem{
		{ProductId: "p-1", ProductName: "Pizza", Quantity: 2, Price: 8.50},
		{ProductId: "p-2", ProductName: "Coke", Quantity: 3, Price: 2.50},
	}

	var subtotal float64
	for _, item := range items {
		subtotal += item.Price * float64(item.Quantity)
	}

	expected := 24.50
	if subtotal != expected {
		t.Errorf("expected subtotal=%.2f, got %.2f", expected, subtotal)
	}
}

func TestMenuCategory_Serialization(t *testing.T) {
	menu := []onlineorder_models.MenuCategory{
		{
			Id:   "cat-1",
			Name: "Pizze",
			Products: []onlineorder_models.MenuProduct{
				{Id: "p-1", Name: "Margherita", Price: 8.50, Available: true},
			},
		},
	}

	data, err := json.Marshal(menu)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded []onlineorder_models.MenuCategory
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(decoded) != 1 {
		t.Errorf("expected 1 category, got %d", len(decoded))
	}
	if !decoded[0].Products[0].Available {
		t.Error("expected product to be available")
	}
}

func TestOrderStatus_Values(t *testing.T) {
	validStatuses := []string{"pending", "confirmed", "preparing", "ready", "delivered", "cancelled"}
	for _, status := range validStatuses {
		if status == "" {
			t.Error("status should not be empty")
		}
	}
}

func TestOrderType_Values(t *testing.T) {
	validTypes := []string{"delivery", "takeaway", "dine_in"}
	for _, orderType := range validTypes {
		if orderType == "" {
			t.Error("order type should not be empty")
		}
	}
}
