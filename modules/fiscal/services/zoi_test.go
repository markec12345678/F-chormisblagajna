package services

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/nutrixpos/pos/modules/fiscal/models"
	"math/big"
	"strings"
	"testing"
	"time"
)

func testPrivateKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate RSA key: %v", err)
	}
	return key
}

func TestCalculateZOI_Returns32CharHex(t *testing.T) {
	key := testPrivateKey(t)
	issueTime := time.Date(2026, 7, 25, 14, 30, 0, 0, time.UTC)

	zoi, err := CalculateZOI(key, 10115609, issueTime, "00001-00001-00001", "BP105", "0001", 19.15)
	if err != nil {
		t.Fatalf("CalculateZOI error: %v", err)
	}

	if len(zoi) != 32 {
		t.Errorf("expected 32 char ZOI, got %d chars: %s", len(zoi), zoi)
	}

	// Must be valid hex
	for _, ch := range zoi {
		if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f')) {
			t.Errorf("invalid hex char in ZOI: %c", ch)
		}
	}
}

func TestCalculateZOI_Deterministic(t *testing.T) {
	key := testPrivateKey(t)
	issueTime := time.Date(2026, 7, 25, 14, 30, 0, 0, time.UTC)

	zoi1, err := CalculateZOI(key, 10115609, issueTime, "00001-00001-00001", "BP105", "0001", 19.15)
	if err != nil {
		t.Fatalf("CalculateZOI error: %v", err)
	}

	zoi2, err := CalculateZOI(key, 10115609, issueTime, "00001-00001-00001", "BP105", "0001", 19.15)
	if err != nil {
		t.Fatalf("CalculateZOI error: %v", err)
	}

	if zoi1 != zoi2 {
		t.Errorf("ZOI not deterministic: %s != %s", zoi1, zoi2)
	}
}

func TestCalculateZOI_DifferentAmounts(t *testing.T) {
	key := testPrivateKey(t)
	issueTime := time.Date(2026, 7, 25, 14, 30, 0, 0, time.UTC)

	zoi1, _ := CalculateZOI(key, 10115609, issueTime, "00001-00001-00001", "BP105", "0001", 19.15)
	zoi2, _ := CalculateZOI(key, 10115609, issueTime, "00001-00001-00001", "BP105", "0001", 25.00)

	if zoi1 == zoi2 {
		t.Error("expected different ZOI for different amounts")
	}
}

func TestCalculateZOI_DifferentTaxNumbers(t *testing.T) {
	key := testPrivateKey(t)
	issueTime := time.Date(2026, 7, 25, 14, 30, 0, 0, time.UTC)

	zoi1, _ := CalculateZOI(key, 10115609, issueTime, "00001-00001-00001", "BP105", "0001", 19.15)
	zoi2, _ := CalculateZOI(key, 12345678, issueTime, "00001-00001-00001", "BP105", "0001", 19.15)

	if zoi1 == zoi2 {
		t.Error("expected different ZOI for different tax numbers")
	}
}

func TestCalculateQRData_Length60(t *testing.T) {
	zoiHex := "51a40dcabb147d1c76d843ee98f951aa"
	issueTime := time.Date(2026, 7, 25, 14, 30, 0, 0, time.UTC)

	qr, err := CalculateQRData(zoiHex, 10115609, issueTime)
	if err != nil {
		t.Fatalf("CalculateQRData error: %v", err)
	}

	if len(qr) != 60 {
		t.Errorf("expected 60 char QR data, got %d chars: %s", len(qr), qr)
	}

	// Must be all digits
	for _, ch := range qr {
		if ch < '0' || ch > '9' {
			t.Errorf("invalid digit in QR data: %c", ch)
		}
	}
}

func TestCalculateQRData_CheckDigit(t *testing.T) {
	zoiHex := "51a40dcabb147d1c76d843ee98f951aa"
	issueTime := time.Date(2026, 7, 25, 14, 30, 0, 0, time.UTC)

	qr, err := CalculateQRData(zoiHex, 10115609, issueTime)
	if err != nil {
		t.Fatalf("CalculateQRData error: %v", err)
	}

	// Verify check digit: sum of first 59 digits mod 10 should equal last digit
	data := qr[:59]
	sum := 0
	for _, ch := range data {
		sum += int(ch - '0')
	}
	expectedCheck := sum % 10
	actualCheck := int(qr[59] - '0')

	if expectedCheck != actualCheck {
		t.Errorf("check digit mismatch: expected %d, got %d (data: %s)", expectedCheck, actualCheck, qr)
	}
}

func TestCalculateQRData_ZOIDecimalPadding(t *testing.T) {
	zoiHex := "51a40dcabb147d1c76d843ee98f951aa"
	issueTime := time.Date(2026, 7, 25, 14, 30, 0, 0, time.UTC)

	qr, err := CalculateQRData(zoiHex, 10115609, issueTime)
	if err != nil {
		t.Fatalf("CalculateQRData error: %v", err)
	}

	// ZOI decimal should occupy first 39 digits
	zoiDecimal := qr[:39]
	zoiBig := new(big.Int)
	_, ok := zoiBig.SetString(zoiHex, 16)
	if !ok {
		t.Fatal("invalid ZOI hex")
	}

	expected := strings.Repeat("0", 39-len(zoiBig.String())) + zoiBig.String()
	if zoiDecimal != expected {
		t.Errorf("ZOI decimal mismatch: got %s, expected %s", zoiDecimal, expected)
	}
}

func TestCalculateQRData_TaxNumberPadding(t *testing.T) {
	zoiHex := "51a40dcabb147d1c76d843ee98f951aa"
	issueTime := time.Date(2026, 7, 25, 14, 30, 0, 0, time.UTC)

	qr, err := CalculateQRData(zoiHex, 10115609, issueTime)
	if err != nil {
		t.Fatalf("CalculateQRData error: %v", err)
	}

	taxStr := qr[39:47]
	if taxStr != "10115609" {
		t.Errorf("tax number in QR data wrong: got %s, expected 10115609", taxStr)
	}
}

func TestCalculateQRData_DatetimeFormat(t *testing.T) {
	zoiHex := "51a40dcabb147d1c76d843ee98f951aa"
	issueTime := time.Date(2026, 7, 25, 14, 30, 0, 0, time.UTC)

	qr, err := CalculateQRData(zoiHex, 10115609, issueTime)
	if err != nil {
		t.Fatalf("CalculateQRData error: %v", err)
	}

	// DateTime should be at positions 47-58: YYMMDDHHmmss
	dateStr := qr[47:59]
	expected := "260725143000"
	if dateStr != expected {
		t.Errorf("datetime in QR data wrong: got %s, expected %s", dateStr, expected)
	}
}

func TestFormatInvoiceNumber(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{1, "00000-00000-00001"},
		{42, "00000-00000-00042"},
		{99999, "00000-00000-99999"},
		{100000, "00000-00001-00000"},
		{123456789, "00000-01234-56789"},
	}

	for _, tt := range tests {
		result := FormatInvoiceNumber(tt.input)
		if result != tt.expected {
			t.Errorf("FormatInvoiceNumber(%d) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestFormatInvoiceNumber_SeparatorCount(t *testing.T) {
	result := FormatInvoiceNumber(1)
	parts := strings.Split(result, "-")
	if len(parts) != 3 {
		t.Errorf("expected 3 parts separated by '-', got %d: %s", len(parts), result)
	}
}

func TestFormatIssueDateTime(t *testing.T) {
	utcTime := time.Date(2026, 7, 25, 14, 30, 0, 0, time.UTC)
	result := FormatIssueDateTime(utcTime)
	expected := "2026-07-25T14:30:00Z"
	if result != expected {
		t.Errorf("FormatIssueDateTime = %s, expected %s", result, expected)
	}
}

func TestFormatIssueDateTime_ConvertsToUTC(t *testing.T) {
	localTime := time.Date(2026, 7, 25, 16, 30, 0, 0, time.FixedZone("CET", 2*3600))
	result := FormatIssueDateTime(localTime)
	expected := "2026-07-25T14:30:00Z"
	if result != expected {
		t.Errorf("FormatIssueDateTime (CET) = %s, expected %s", result, expected)
	}
}

func TestFormatAmount_TwoDecimals(t *testing.T) {
	amount := formatAmount(19.15)
	if amount != "19.15" {
		t.Errorf("formatAmount(19.15) = %s, expected 19.15", amount)
	}
}

func TestFormatAmount_ZeroPad(t *testing.T) {
	amount := formatAmount(5.0)
	if amount != "5.00" {
		t.Errorf("formatAmount(5.0) = %s, expected 5.00", amount)
	}
}

func TestBuildTaxesPerSeller_SingleRate(t *testing.T) {
	items := []models.InvoiceItem{
		{Name: "Item 1", Quantity: 2, UnitPrice: 10.0, TaxRate: 22.0, TaxableAmount: 20.0, TaxAmount: 4.40},
	}

	result := buildTaxesPerSeller(items)
	if len(result) != 1 {
		t.Fatalf("expected 1 tax group, got %d", len(result))
	}

	if len(result[0].VAT) != 1 {
		t.Fatalf("expected 1 VAT entry, got %d", len(result[0].VAT))
	}

	vat := result[0].VAT[0]
	if vat.TaxRate != 22.0 {
		t.Errorf("expected tax rate 22.0, got %.1f", vat.TaxRate)
	}
	if vat.TaxableAmount != 20.0 {
		t.Errorf("expected taxable amount 20.0, got %.2f", vat.TaxableAmount)
	}
	if vat.TaxAmount != 4.40 {
		t.Errorf("expected tax amount 4.40, got %.2f", vat.TaxAmount)
	}
}

func TestBuildTaxesPerSeller_MultipleRates(t *testing.T) {
	items := []models.InvoiceItem{
		{Name: "A", Quantity: 1, UnitPrice: 10.0, TaxRate: 22.0, TaxableAmount: 10.0, TaxAmount: 2.20},
		{Name: "B", Quantity: 1, UnitPrice: 10.0, TaxRate: 9.5, TaxableAmount: 10.0, TaxAmount: 0.95},
		{Name: "C", Quantity: 1, UnitPrice: 10.0, TaxRate: 22.0, TaxableAmount: 10.0, TaxAmount: 2.20},
	}

	result := buildTaxesPerSeller(items)
	if len(result) != 1 {
		t.Fatalf("expected 1 tax group, got %d", len(result))
	}

	if len(result[0].VAT) != 2 {
		t.Fatalf("expected 2 VAT entries, got %d", len(result[0].VAT))
	}

	// Should be grouped: 22% → 20.00 taxable, 4.40 tax; 9.5% → 10.00 taxable, 0.95 tax
	for _, vat := range result[0].VAT {
		if vat.TaxRate == 22.0 {
			if vat.TaxableAmount != 20.0 {
				t.Errorf("22%% taxable: expected 20.0, got %.2f", vat.TaxableAmount)
			}
			if vat.TaxAmount != 4.40 {
				t.Errorf("22%% tax: expected 4.40, got %.2f", vat.TaxAmount)
			}
		} else if vat.TaxRate == 9.5 {
			if vat.TaxableAmount != 10.0 {
				t.Errorf("9.5%% taxable: expected 10.0, got %.2f", vat.TaxableAmount)
			}
		} else {
			t.Errorf("unexpected tax rate: %.1f", vat.TaxRate)
		}
	}
}

func TestBuildTaxesPerSeller_EmptyItems(t *testing.T) {
	items := []models.InvoiceItem{}
	result := buildTaxesPerSeller(items)
	if len(result) != 1 {
		t.Fatalf("expected 1 tax group even for empty, got %d", len(result))
	}
	if len(result[0].VAT) != 0 {
		t.Errorf("expected 0 VAT entries for empty items, got %d", len(result[0].VAT))
	}
}

func TestRoundTo2(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{19.155, 19.16},
		{19.154, 19.15},
		{10.0, 10.0},
		{0.125, 0.13},
	}

	for _, tt := range tests {
		result := roundTo2(tt.input)
		if result != tt.expected {
			t.Errorf("roundTo2(%.3f) = %.2f, expected %.2f", tt.input, result, tt.expected)
		}
	}
}
