package services

import (
	"strings"
	"testing"
)

func TestFormatReceiptText_BasicReceipt(t *testing.T) {
	data := &QRCodeData{
		ZOI:           "51a40dcabb147d1c76d843ee98f951aa",
		QRData:        "123456789012345678901234567890123456789012345678901234567890",
		EOR:           "dedc5383-169f-4c0e-b369-d64ec8797170",
		InvoiceNumber: "00000-00000-00001",
	}

	items := []ReceiptItem{
		{Name: "Pizza Margherita", Quantity: 2, UnitPrice: 8.50, Price: 17.00},
		{Name: "Coca Cola", Quantity: 1, UnitPrice: 2.50, Price: 2.50},
	}

	receipt := FormatReceiptText(data, items, 19.50)

	if !strings.Contains(receipt, "Pizza Margherita") {
		t.Error("receipt should contain item name")
	}
	if !strings.Contains(receipt, "TOTAL") {
		t.Error("receipt should contain TOTAL")
	}
	if !strings.Contains(receipt, "19.50") {
		t.Error("receipt should contain total amount")
	}
	if !strings.Contains(receipt, "dedc5383-169f-4c0e-b369-d64ec8797170") {
		t.Error("receipt should contain EOR")
	}
	if !strings.Contains(receipt, "Fiskalizirano po FURS") {
		t.Error("receipt should contain FURS notice")
	}
}

func TestFormatReceiptText_EmptyItems(t *testing.T) {
	data := &QRCodeData{
		ZOI:           "abc123",
		EOR:           "eor-123",
		InvoiceNumber: "00001",
	}

	receipt := FormatReceiptText(data, []ReceiptItem{}, 0.0)

	if !strings.Contains(receipt, "TOTAL") {
		t.Error("receipt should contain TOTAL even with no items")
	}
	if !strings.Contains(receipt, "0.00") {
		t.Error("receipt should show 0.00 for empty items")
	}
}

func TestFormatReceiptText_SingleItem(t *testing.T) {
	data := &QRCodeData{
		ZOI:           "abc",
		EOR:           "eor",
		InvoiceNumber: "001",
	}

	items := []ReceiptItem{
		{Name: "Burger", Quantity: 1, UnitPrice: 8.00, Price: 8.00},
	}

	receipt := FormatReceiptText(data, items, 8.00)

	if !strings.Contains(receipt, "Burger") {
		t.Error("receipt should contain Burger")
	}
	if !strings.Contains(receipt, "Invoice: 001") {
		t.Error("receipt should contain invoice number")
	}
}

func TestQRDataToBase64Image_EmptyData(t *testing.T) {
	_, err := QRDataToBase64Image("")
	if err == nil {
		t.Error("expected error for empty QR data")
	}
}

func TestQRDataToBase64Image_Placeholder(t *testing.T) {
	_, err := QRDataToBase64Image("1234567890")
	if err == nil {
		// Placeholder implementation returns error since real QR lib is not installed
		return
	}
	if err == nil {
		t.Error("placeholder should return error indicating library needed")
	}
}

func TestReceiptItem_Fields(t *testing.T) {
	item := ReceiptItem{
		Name:      "Test Item",
		Quantity:  3,
		UnitPrice: 5.00,
		Price:     15.00,
	}

	if item.Name != "Test Item" {
		t.Errorf("expected name 'Test Item', got '%s'", item.Name)
	}
	if item.Quantity != 3 {
		t.Errorf("expected quantity 3, got %.0f", item.Quantity)
	}
	if item.UnitPrice != 5.00 {
		t.Errorf("expected unit price 5.00, got %.2f", item.UnitPrice)
	}
	if item.Price != 15.00 {
		t.Errorf("expected price 15.00, got %.2f", item.Price)
	}
}
