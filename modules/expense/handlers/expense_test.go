package handlers

import (
	"encoding/json"
	"testing"
	"time"

	expense_models "github.com/nutrixpos/pos/modules/expense/models"
)

func TestExpense_Serialization(t *testing.T) {
	expense := expense_models.Expense{
		Id:          "e-1",
		Description: "Office supplies",
		Amount:      150.00,
		Category:    "supplies",
		Date:        time.Now(),
		Notes:       "Monthly supplies order",
		CreatedBy:   "admin",
		CreatedAt:   time.Now(),
	}

	data, err := json.Marshal(expense)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded expense_models.Expense
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Amount != 150.00 {
		t.Errorf("expected Amount=150, got %f", decoded.Amount)
	}
	if decoded.Category != "supplies" {
		t.Errorf("expected Category='supplies', got %s", decoded.Category)
	}
}

func TestExpenseSummary_Serialization(t *testing.T) {
	summary := expense_models.ExpenseSummary{
		TotalExpenses: 5000.00,
		MonthlyTotal:  1200.00,
		ByCategory: []expense_models.CategorySummary{
			{Category: "rent", Total: 2000.00, Count: 1},
			{Category: "supplies", Total: 1500.00, Count: 5},
			{Category: "utilities", Total: 1500.00, Count: 3},
		},
	}

	data, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded expense_models.ExpenseSummary
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.TotalExpenses != 5000.00 {
		t.Errorf("expected TotalExpenses=5000, got %f", decoded.TotalExpenses)
	}
	if len(decoded.ByCategory) != 3 {
		t.Errorf("expected 3 categories, got %d", len(decoded.ByCategory))
	}
}

func TestCategorySummary_Calculation(t *testing.T) {
	cat := expense_models.CategorySummary{
		Category: "supplies",
	}

	cat.Total = 150.00 + 200.00 + 75.00
	cat.Count = 3

	if cat.Total != 425.00 {
		t.Errorf("expected Total=425, got %f", cat.Total)
	}
	if cat.Count != 3 {
		t.Errorf("expected Count=3, got %d", cat.Count)
	}
}

func TestDateParsing_ExpenseSummary(t *testing.T) {
	startDate := "2024-01-01"
	endDate := "2024-01-31"

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		t.Fatalf("failed to parse start date: %v", err)
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		t.Fatalf("failed to parse end date: %v", err)
	}

	if start.After(end) {
		t.Error("start date should be before end date")
	}
}
