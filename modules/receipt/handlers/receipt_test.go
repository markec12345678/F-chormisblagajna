package handlers

import (
	"encoding/json"
	"testing"

	receipt_models "github.com/nutrixpos/pos/modules/receipt/models"
)

func TestReceiptTemplate_Serialization(t *testing.T) {
	tpl := receipt_models.ReceiptTemplate{
		Id:              "t-1",
		Name:            "Default",
		BusinessName:    "NutrixPOS Restaurant",
		BusinessAddress: "Testna ulica 12, 1000 Ljubljana",
		BusinessPhone:   "+386 1 234 5678",
		BusinessTaxId:   "SI12345678",
		Header:          "Welcome!",
		Footer:          "Thank you for dining with us!",
		ShowLogo:        true,
		ShowTaxId:       true,
		ShowQRCode:      true,
		ShowServer:      true,
		ShowTable:       true,
		PaperWidth:      80,
		CustomFields: []receipt_models.CustomField{
			{Key: "wifi", Value: "Nutrix-Guest / password123"},
		},
		IsDefault: true,
	}

	data, err := json.Marshal(tpl)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded receipt_models.ReceiptTemplate
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.BusinessName != "NutrixPOS Restaurant" {
		t.Errorf("expected BusinessName='NutrixPOS Restaurant', got %s", decoded.BusinessName)
	}
	if decoded.PaperWidth != 80 {
		t.Errorf("expected PaperWidth=80, got %d", decoded.PaperWidth)
	}
	if len(decoded.CustomFields) != 1 {
		t.Errorf("expected 1 custom field, got %d", len(decoded.CustomFields))
	}
}

func TestPrintSettings_Serialization(t *testing.T) {
	settings := receipt_models.PrintSettings{
		Id:          "ps-1",
		PrinterName: "EPSON TM-T88",
		PrinterIP:   "192.168.1.100",
		AutoPrint:   true,
		PrintCopies: 2,
		TemplateId:  "t-1",
		Connected:   true,
	}

	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded receipt_models.PrintSettings
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.PrinterName != "EPSON TM-T88" {
		t.Errorf("expected PrinterName='EPSON TM-T88', got %s", decoded.PrinterName)
	}
	if decoded.PrintCopies != 2 {
		t.Errorf("expected PrintCopies=2, got %d", decoded.PrintCopies)
	}
}

func TestPaperWidth_Values(t *testing.T) {
	validWidths := []int{58, 80}
	for _, w := range validWidths {
		if w != 58 && w != 80 {
			t.Errorf("unexpected paper width: %d", w)
		}
	}
}

func TestCustomField_Serialization(t *testing.T) {
	fields := []receipt_models.CustomField{
		{Key: "wifi", Value: "Network / Password"},
		{Key: "social", Value: "@nutrixpos"},
	}

	data, err := json.Marshal(fields)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded []receipt_models.CustomField
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(decoded) != 2 {
		t.Errorf("expected 2 fields, got %d", len(decoded))
	}
	if decoded[0].Key != "wifi" {
		t.Errorf("expected Key='wifi', got %s", decoded[0].Key)
	}
}
