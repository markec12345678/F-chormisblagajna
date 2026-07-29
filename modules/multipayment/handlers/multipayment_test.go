package handlers

import (
	"encoding/json"
	"testing"
	"time"

	mp_models "github.com/nutrixpos/pos/modules/multipayment/models"
)

func TestPaymentPart_Serialization(t *testing.T) {
	payment := mp_models.PaymentPart{
		Id:            "pp-1",
		OrderId:       "o-1",
		Amount:        10.50,
		PaymentMethod: "cash",
		Reference:     "",
		Notes:         "Cash payment",
		ReceivedBy:    "cashier-1",
		CreatedAt:     time.Now(),
	}

	data, err := json.Marshal(payment)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded mp_models.PaymentPart
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Amount != 10.50 {
		t.Errorf("expected Amount=10.50, got %f", decoded.Amount)
	}
	if decoded.PaymentMethod != "cash" {
		t.Errorf("expected PaymentMethod='cash', got %s", decoded.PaymentMethod)
	}
}

func TestPaymentSummary_Serialization(t *testing.T) {
	summary := mp_models.PaymentSummary{
		OrderId:     "o-1",
		TotalDue:    25.00,
		TotalPaid:   15.00,
		Remaining:   10.00,
		IsFullyPaid: false,
		Payments: []mp_models.PaymentPart{
			{Id: "pp-1", Amount: 15.00, PaymentMethod: "cash"},
		},
		MethodBreakdown: []mp_models.MethodAmount{
			{Method: "cash", Total: 15.00, Count: 1},
		},
	}

	data, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded mp_models.PaymentSummary
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.TotalDue != 25.00 {
		t.Errorf("expected TotalDue=25, got %f", decoded.TotalDue)
	}
	if decoded.IsFullyPaid {
		t.Error("expected IsFullyPaid=false")
	}
}

func TestDailyPayments_Serialization(t *testing.T) {
	daily := mp_models.DailyPayments{
		Date:       "2024-01-15",
		TotalCash:  150.00,
		TotalCard:  300.00,
		GrandTotal: 450.00,
		Count:      25,
	}

	data, err := json.Marshal(daily)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded mp_models.DailyPayments
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.GrandTotal != 450.00 {
		t.Errorf("expected GrandTotal=450, got %f", decoded.GrandTotal)
	}
}

func TestPaymentMethods(t *testing.T) {
	methods := []string{"cash", "card", "voucher", "mobile", "gift_card"}
	for _, m := range methods {
		if m == "" {
			t.Error("method should not be empty")
		}
	}
}
