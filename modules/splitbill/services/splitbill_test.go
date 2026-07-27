package services

import (
	"testing"

	"github.com/nutrixpos/pos/modules/splitbill/dto"
	"github.com/nutrixpos/pos/modules/splitbill/models"
)

func TestSplitPart_Model(t *testing.T) {
	part := models.SplitPart{
		Id:            "part-1",
		Amount:        25.50,
		PaymentMethod: "cash",
		IsPaid:        true,
	}

	if part.Id != "part-1" {
		t.Errorf("expected id part-1, got %s", part.Id)
	}
	if part.Amount != 25.50 {
		t.Errorf("expected amount 25.50, got %f", part.Amount)
	}
	if part.PaymentMethod != "cash" {
		t.Errorf("expected payment_method cash, got %s", part.PaymentMethod)
	}
	if !part.IsPaid {
		t.Error("expected is_paid true")
	}
}

func TestSplitBill_Model(t *testing.T) {
	bill := models.SplitBill{
		Id:        "sb-1",
		OrderId:   "order-1",
		SplitType: "equal",
		Status:    "pending",
		Parts: []models.SplitPart{
			{Id: "p1", Amount: 30.00, IsPaid: false},
			{Id: "p2", Amount: 30.00, IsPaid: false},
		},
	}

	if bill.OrderId != "order-1" {
		t.Errorf("expected order_id order-1, got %s", bill.OrderId)
	}
	if bill.SplitType != "equal" {
		t.Errorf("expected split_type equal, got %s", bill.SplitType)
	}
	if len(bill.Parts) != 2 {
		t.Errorf("expected 2 parts, got %d", len(bill.Parts))
	}
	if bill.Status != "pending" {
		t.Errorf("expected status pending, got %s", bill.Status)
	}
}

func TestCreateSplitBillRequest_Equal(t *testing.T) {
	req := dto.CreateSplitBillRequest{
		OrderId:    "order-123",
		SplitType:  "equal",
		SplitCount: 4,
	}

	if req.SplitCount != 4 {
		t.Errorf("expected split_count 4, got %d", req.SplitCount)
	}
}

func TestCreateSplitBillRequest_Custom(t *testing.T) {
	req := dto.CreateSplitBillRequest{
		OrderId:       "order-456",
		SplitType:     "custom",
		CustomAmounts: []float64{10.00, 20.00, 30.00},
	}

	if len(req.CustomAmounts) != 3 {
		t.Errorf("expected 3 amounts, got %d", len(req.CustomAmounts))
	}
	total := 0.0
	for _, a := range req.CustomAmounts {
		total += a
	}
	if total != 60.00 {
		t.Errorf("expected total 60.00, got %f", total)
	}
}

func TestCreateSplitBillRequest_ByItem(t *testing.T) {
	req := dto.CreateSplitBillRequest{
		OrderId:   "order-789",
		SplitType: "by_item",
		ItemSplits: []dto.ItemSplit{
			{ItemId: "item-1", PartIndex: 0},
			{ItemId: "item-2", PartIndex: 0},
			{ItemId: "item-3", PartIndex: 1},
		},
	}

	if len(req.ItemSplits) != 3 {
		t.Errorf("expected 3 item splits, got %d", len(req.ItemSplits))
	}
	if req.ItemSplits[0].PartIndex != 0 {
		t.Errorf("expected part_index 0, got %d", req.ItemSplits[0].PartIndex)
	}
	if req.ItemSplits[2].PartIndex != 1 {
		t.Errorf("expected part_index 1, got %d", req.ItemSplits[2].PartIndex)
	}
}

func TestPaySplitPartRequest(t *testing.T) {
	req := dto.PaySplitPartRequest{
		PartId:        "part-2",
		PaymentMethod: "card",
		Amount:        45.99,
	}

	if req.PartId != "part-2" {
		t.Errorf("expected part_id part-2, got %s", req.PartId)
	}
	if req.PaymentMethod != "card" {
		t.Errorf("expected payment_method card, got %s", req.PaymentMethod)
	}
	if req.Amount != 45.99 {
		t.Errorf("expected amount 45.99, got %f", req.Amount)
	}
}

func TestSplitBill_StatusValues(t *testing.T) {
	validStatuses := []string{"pending", "partial", "paid"}
	for _, status := range validStatuses {
		bill := models.SplitBill{Status: status}
		if bill.Status != status {
			t.Errorf("expected status %s, got %s", status, bill.Status)
		}
	}
}

func TestSplitBill_EmptyParts(t *testing.T) {
	bill := models.SplitBill{
		Id:        "sb-empty",
		OrderId:   "order-empty",
		SplitType: "equal",
		Parts:     []models.SplitPart{},
		Status:    "pending",
	}

	if len(bill.Parts) != 0 {
		t.Errorf("expected 0 parts, got %d", len(bill.Parts))
	}
}

func TestSplitBill_AllPartsPaid(t *testing.T) {
	bill := models.SplitBill{
		Id:      "sb-allpaid",
		OrderId: "order-allpaid",
		Parts: []models.SplitPart{
			{Id: "p1", Amount: 25.00, IsPaid: true},
			{Id: "p2", Amount: 25.00, IsPaid: true},
		},
	}

	allPaid := true
	for _, p := range bill.Parts {
		if !p.IsPaid {
			allPaid = false
			break
		}
	}

	if !allPaid {
		t.Error("expected all parts to be paid")
	}
}

func TestSplitBill_PartialPayment(t *testing.T) {
	bill := models.SplitBill{
		Id:      "sb-partial",
		OrderId: "order-partial",
		Parts: []models.SplitPart{
			{Id: "p1", Amount: 25.00, IsPaid: true},
			{Id: "p2", Amount: 25.00, IsPaid: false},
		},
	}

	allPaid := true
	for _, p := range bill.Parts {
		if !p.IsPaid {
			allPaid = false
			break
		}
	}

	if allPaid {
		t.Error("expected not all parts to be paid")
	}

	paidCount := 0
	for _, p := range bill.Parts {
		if p.IsPaid {
			paidCount++
		}
	}
	if paidCount != 1 {
		t.Errorf("expected 1 paid part, got %d", paidCount)
	}
}

func TestSplitBill_TotalAmount(t *testing.T) {
	parts := []models.SplitPart{
		{Amount: 10.50},
		{Amount: 20.25},
		{Amount: 15.75},
	}

	total := 0.0
	for _, p := range parts {
		total += p.Amount
	}

	if total != 46.50 {
		t.Errorf("expected total 46.50, got %f", total)
	}
}
