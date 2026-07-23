package models

import (
	"encoding/json"
	"math"
	"testing"
	"time"
)

func TestJSONFloat_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    JSONFloat
		expected string
	}{
		{"positive", JSONFloat(1.5), "1.5"},
		{"negative", JSONFloat(-2.5), "-2.5"},
		{"zero", JSONFloat(0), "0"},
		{"infinity", JSONFloat(math.Inf(1)), "-1"},
		{"neg_infinity", JSONFloat(math.Inf(-1)), "-1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("MarshalJSON failed: %v", err)
			}
			if string(result) != tt.expected {
				t.Errorf("MarshalJSON(%v) = %s, want %s", tt.input, string(result), tt.expected)
			}
		})
	}
}

func TestProduct_JSON(t *testing.T) {
	product := Product{
		Id:       "prod-123",
		Name:     "Test Product",
		Price:    9.99,
		Quantity: 10,
		Ready:    5,
		Unit:     "kg",
		Materials: []Material{
			{
				Id:       "mat-1",
				Name:     "Flour",
				Quantity: 2,
			},
		},
	}

	data, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded Product
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Name != product.Name {
		t.Errorf("Name = %v, want %v", decoded.Name, product.Name)
	}
	if decoded.Price != product.Price {
		t.Errorf("Price = %v, want %v", decoded.Price, product.Price)
	}
	if len(decoded.Materials) != 1 {
		t.Errorf("Materials length = %v, want 1", len(decoded.Materials))
	}
}

func TestOrder_JSON(t *testing.T) {
	order := Order{
		Id:          "order-123",
		DisplayId:   "A-1",
		State:       "submitted",
		SubmittedAt: time.Now(),
		IsPaid:      false,
	}

	data, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded Order
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Id != order.Id {
		t.Errorf("Id = %v, want %v", decoded.Id, order.Id)
	}
	if decoded.State != order.State {
		t.Errorf("State = %v, want %v", decoded.State, order.State)
	}
}

func TestCategory_JSON(t *testing.T) {
	category := Category{
		Id:   "cat-1",
		Name: "Drinks",
	}

	data, err := json.Marshal(category)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded Category
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Name != category.Name {
		t.Errorf("Name = %v, want %v", decoded.Name, category.Name)
	}
}

func TestSettings_JSON(t *testing.T) {
	settings := Settings{
		Id: "settings-1",
	}

	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded Settings
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Id != settings.Id {
		t.Errorf("Id = %v, want %v", decoded.Id, settings.Id)
	}
}

func TestJSONFloat_Unmarshal_Number(t *testing.T) {
	input := `{"sale_price": 12.5}`
	var result struct {
		SalePrice JSONFloat `json:"sale_price"`
	}
	if err := json.Unmarshal([]byte(input), &result); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if float64(result.SalePrice) != 12.5 {
		t.Errorf("SalePrice = %v, want 12.5", result.SalePrice)
	}
}

func TestJSONFloat_MarshalZero(t *testing.T) {
	result, err := json.Marshal(JSONFloat(0))
	if err != nil {
		t.Fatalf("MarshalJSON(0) failed: %v", err)
	}
	if string(result) != "0" {
		t.Errorf("MarshalJSON(0) = %s, want \"0\"", string(result))
	}
}
