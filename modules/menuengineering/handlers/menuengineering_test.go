package handlers

import (
	"encoding/json"
	"testing"
	"time"

	menu_models "github.com/nutrixpos/pos/modules/menuengineering/models"
)

func TestMenuAnalysis_Serialization(t *testing.T) {
	summary := menu_models.MenuEngineeringSummary{
		TotalItems:   10,
		TotalRevenue: 5000.00,
		AvgProfit:    25.50,
		TopStars: []menu_models.MenuItemAnalysis{
			{
				ProductId:     "p-1",
				ProductName:   "Pizza Margherita",
				TotalSold:     150,
				TotalRevenue:  1500.00,
				TotalCost:     450.00,
				ProfitMargin:  70.0,
				ProfitPerItem: 7.00,
				Quadrant:      "star",
			},
		},
		TopPlowhorses: []menu_models.MenuItemAnalysis{},
		TopPuzzles:    []menu_models.MenuItemAnalysis{},
		TopDogs:       []menu_models.MenuItemAnalysis{},
	}

	data, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded menu_models.MenuEngineeringSummary
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.TotalItems != 10 {
		t.Errorf("expected TotalItems=10, got %d", decoded.TotalItems)
	}
	if decoded.TotalRevenue != 5000.00 {
		t.Errorf("expected TotalRevenue=5000, got %f", decoded.TotalRevenue)
	}
	if len(decoded.TopStars) != 1 {
		t.Errorf("expected 1 star, got %d", len(decoded.TopStars))
	}
	if decoded.TopStars[0].Quadrant != "star" {
		t.Errorf("expected quadrant=star, got %s", decoded.TopStars[0].Quadrant)
	}
}

func TestMenuItemAnalysis_ProfitCalculation(t *testing.T) {
	item := menu_models.MenuItemAnalysis{
		TotalSold:    100,
		TotalRevenue: 1000.00,
		TotalCost:    300.00,
	}

	item.ProfitMargin = ((item.TotalRevenue - item.TotalCost) / item.TotalRevenue) * 100
	item.ProfitPerItem = (item.TotalRevenue - item.TotalCost) / float64(item.TotalSold)

	if item.ProfitMargin != 70.0 {
		t.Errorf("expected ProfitMargin=70, got %f", item.ProfitMargin)
	}
	if item.ProfitPerItem != 7.0 {
		t.Errorf("expected ProfitPerItem=7, got %f", item.ProfitPerItem)
	}
}

func TestQuadrantClassification(t *testing.T) {
	tests := []struct {
		name       string
		profit     float64
		popularity float64
		avgProfit  float64
		avgPop     float64
		expected   string
	}{
		{"star", 30.0, 150.0, 20.0, 100.0, "star"},
		{"plowhorse", 10.0, 150.0, 20.0, 100.0, "plowhorse"},
		{"puzzle", 30.0, 50.0, 20.0, 100.0, "puzzle"},
		{"dog", 10.0, 50.0, 20.0, 100.0, "dog"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			highProfit := tt.profit >= tt.avgProfit
			highPopularity := tt.popularity >= tt.avgPop

			var quadrant string
			switch {
			case highProfit && highPopularity:
				quadrant = "star"
			case !highProfit && highPopularity:
				quadrant = "plowhorse"
			case highProfit && !highPopularity:
				quadrant = "puzzle"
			default:
				quadrant = "dog"
			}

			if quadrant != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, quadrant)
			}
		})
	}
}

func TestDateParsing_MenuAnalysis(t *testing.T) {
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
