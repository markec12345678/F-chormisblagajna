package services

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"
)

func TestCalculateZKI(t *testing.T) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)

	oib := "12345678901"
	datVrijeme := "26.07.2026T15:00:00"
	brOznRac := "1"
	oznPosPr := "1"
	oznNapUr := "1"
	iznosUkupno := "12.34"

	zki, err := CalculateZKI(key, oib, datVrijeme, brOznRac, oznPosPr, oznNapUr, iznosUkupno)
	if err != nil {
		t.Fatalf("CalculateZKI: %v", err)
	}

	// ZKI is MD5 of RSA signature → 32-char hex
	if len(zki) != 32 {
		t.Errorf("ZKI length = %d, want 32", len(zki))
	}

	// Same inputs should produce same ZKI (deterministic)
	zki2, err := CalculateZKI(key, oib, datVrijeme, brOznRac, oznPosPr, oznNapUr, iznosUkupno)
	if err != nil {
		t.Fatalf("CalculateZKI second: %v", err)
	}

	if zki != zki2 {
		t.Errorf("ZKI not deterministic: %s != %s", zki, zki2)
	}
}

func TestFormatDateTimeHR(t *testing.T) {
	ts := time.Date(2026, 7, 26, 15, 30, 45, 0, time.UTC)
	got := FormatDateTimeHR(ts)
	want := "26.07.2026T15:30:45"
	if got != want {
		t.Errorf("FormatDateTimeHR = %q, want %q", got, want)
	}
}

func TestFormatDateHR(t *testing.T) {
	ts := time.Date(2026, 1, 5, 10, 0, 0, 0, time.UTC)
	got := FormatDateHR(ts)
	want := "05.01.2026"
	if got != want {
		t.Errorf("FormatDateHR = %q, want %q", got, want)
	}
}

func TestFormatAmountHR(t *testing.T) {
	tests := []struct {
		input float64
		want  string
	}{
		{12.345, "12.35"},
		{12.3, "12.30"},
		{0, "0.00"},
		{100, "100.00"},
		{99999.99, "99999.99"},
	}

	for _, tt := range tests {
		got := FormatAmountHR(tt.input)
		if got != tt.want {
			t.Errorf("FormatAmountHR(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestValidateOIB(t *testing.T) {
	tests := []struct {
		oib  string
		want bool
	}{
		{"12345678901", false}, // invalid checksum
		{"98765432100", false}, // invalid checksum
		{"1234", false},        // too short
		{"123456789012", false}, // too long
		{"ABCDEFGHIJK", false},  // not digits
		{"", false},             // empty
	}

	for _, tt := range tests {
		got := ValidateOIB(tt.oib)
		if got != tt.want {
			t.Errorf("ValidateOIB(%q) = %v, want %v", tt.oib, got, tt.want)
		}
	}
}

func TestValidateOIBKnownGood(t *testing.T) {
	// OIB 12345678901 has valid checksum:
	// sum = (1+0)+(2+10)+(3+0)+(4+10)+(5+0)+(6+10)+(7+0)+(8+10)+(9+0)+(0+10) = 1+12+3+14+5+16+7+18+9+10 = 95
	// 11 - (95%11) = 11 - 7 = 4 → checkDigit = 4
	// So OIB 12345678904 is valid
	if !ValidateOIB("12345678904") {
		t.Error("ValidateOIB(12345678904) = false, want true")
	}
}

func TestPaymentMethodCode(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"cash", "G"},
		{"gotovina", "G"},
		{"card", "K"},
		{"kartica", "K"},
		{"transfer", "T"},
		{"transakcija", "T"},
		{"other", "O"},
		{"", "O"},
		{"unknown", "O"},
	}

	for _, tt := range tests {
		got := PaymentMethodCode(tt.input)
		if got != tt.want {
			t.Errorf("PaymentMethodCode(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
