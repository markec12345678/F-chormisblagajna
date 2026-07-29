package handlers

import (
	"encoding/json"
	"testing"
	"time"

	employee_models "github.com/nutrixpos/pos/modules/employee/models"
)

func TestPerformanceSummary_Serialization(t *testing.T) {
	summary := employee_models.PerformanceSummary{
		TotalEmployees:  5,
		TotalRevenue:    15000.00,
		AvgSalesPerHour: 187.50,
		TopPerformers: []employee_models.EmployeePerformance{
			{
				EmployeeId:    "emp-1",
				EmployeeName:  "Janez Novak",
				TotalSales:    5000.00,
				OrderCount:    120,
				AvgOrderValue: 41.67,
				SalesPerHour:  625.00,
				Rank:          1,
			},
		},
	}

	data, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded employee_models.PerformanceSummary
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.TotalEmployees != 5 {
		t.Errorf("expected TotalEmployees=5, got %d", decoded.TotalEmployees)
	}
	if decoded.TotalRevenue != 15000.00 {
		t.Errorf("expected TotalRevenue=15000, got %f", decoded.TotalRevenue)
	}
	if len(decoded.TopPerformers) != 1 {
		t.Errorf("expected 1 top performer, got %d", len(decoded.TopPerformers))
	}
	if decoded.TopPerformers[0].Rank != 1 {
		t.Errorf("expected Rank=1, got %d", decoded.TopPerformers[0].Rank)
	}
}

func TestEmployeePerformance_Calculations(t *testing.T) {
	emp := employee_models.EmployeePerformance{
		TotalSales: 1000.00,
		OrderCount: 25,
		TotalHours: 8.0,
	}

	emp.AvgOrderValue = emp.TotalSales / float64(emp.OrderCount)
	emp.SalesPerHour = emp.TotalSales / emp.TotalHours

	if emp.AvgOrderValue != 40.0 {
		t.Errorf("expected AvgOrderValue=40, got %f", emp.AvgOrderValue)
	}
	if emp.SalesPerHour != 125.0 {
		t.Errorf("expected SalesPerHour=125, got %f", emp.SalesPerHour)
	}
}

func TestProductStat_Serialization(t *testing.T) {
	stat := employee_models.ProductStat{
		ProductId:   "p-1",
		ProductName: "Pizza Margherita",
		Quantity:    50,
		Revenue:     500.00,
	}

	data, err := json.Marshal(stat)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded employee_models.ProductStat
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.ProductName != "Pizza Margherita" {
		t.Errorf("expected ProductName='Pizza Margherita', got %s", decoded.ProductName)
	}
	if decoded.Quantity != 50 {
		t.Errorf("expected Quantity=50, got %d", decoded.Quantity)
	}
}

func TestDateParsing_EmployeePerformance(t *testing.T) {
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
