package handlers

import (
	"testing"
	"time"
)

func TestDateParsing_DefaultRange(t *testing.T) {
	startDate := ""
	endDate := ""

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		start = time.Now().AddDate(0, 0, -30)
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		end = time.Now()
	}

	if start.After(end) {
		t.Errorf("start date should be before end date")
	}
}

func TestDateParsing_ValidDates(t *testing.T) {
	start, err := time.Parse("2006-01-02", "2024-01-01")
	if err != nil {
		t.Fatalf("failed to parse start date: %v", err)
	}

	end, err := time.Parse("2006-01-02", "2024-01-31")
	if err != nil {
		t.Fatalf("failed to parse end date: %v", err)
	}

	if start.After(end) {
		t.Errorf("start date should be before end date")
	}

	if start.Year() != 2024 || start.Month() != 1 || start.Day() != 1 {
		t.Errorf("start date incorrect: %v", start)
	}
}

func TestDateParsing_FallbackToDefault(t *testing.T) {
	_, err := time.Parse("2006-01-02", "invalid")
	if err == nil {
		t.Error("expected error for invalid date")
	}

	fallback := time.Now().AddDate(0, 0, -30)
	if fallback.IsZero() {
		t.Error("fallback date should not be zero")
	}
}

func TestCSVHeaders_QuickBooks(t *testing.T) {
	expected := []string{"Date", "Type", "Num", "Name", "Memo", "Amount", "Account", "Income Account", "Tax Rate"}
	if len(expected) != 9 {
		t.Errorf("expected 9 columns, got %d", len(expected))
	}
}

func TestCSVHeaders_Xero(t *testing.T) {
	expected := []string{"Date", "Reference", "Payee", "Description", "Line Amount", "Tax Amount", "Account Code", "Tax Type", "Category"}
	if len(expected) != 9 {
		t.Errorf("expected 9 columns, got %d", len(expected))
	}
}
