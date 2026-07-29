package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/modules/report/models"
)

func TestJSONApiOkResponse_Serialization(t *testing.T) {
	report := models.SalesReport{
		Period:       "2026-07-01 to 2026-07-31",
		TotalRevenue: 15000.50,
		TotalOrders:  450,
		AverageOrder: 33.33,
		TopProducts:  []models.ProductStat{{Name: "Pizza", Quantity: 120, Revenue: 1440}},
	}

	resp := JSONApiOkResponse{
		Data: report,
		Meta: JSONAPIMeta{TotalRecords: 1},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded JSONApiOkResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
}

func TestGetSalesReport_ReturnsOK(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/reports/sales", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		report := models.SalesReport{TotalOrders: 0, TopProducts: make([]models.ProductStat, 0)}
		resp := JSONApiOkResponse{Data: report, Meta: JSONAPIMeta{TotalRecords: 1}}
		json.NewEncoder(w).Encode(resp)
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/reports/sales?start_date=2026-07-01&end_date=2026-07-31", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestGetInventoryReport_ReturnsOK(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/reports/inventory", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		report := models.InventoryReport{LowStockItems: make([]models.LowStockItem, 0)}
		resp := JSONApiOkResponse{Data: report, Meta: JSONAPIMeta{TotalRecords: 1}}
		json.NewEncoder(w).Encode(resp)
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/reports/inventory", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestGetDashboardStats_ReturnsOK(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/reports/dashboard", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := JSONApiOkResponse{Data: models.DashboardStats{}, Meta: JSONAPIMeta{TotalRecords: 1}}
		json.NewEncoder(w).Encode(resp)
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/reports/dashboard", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestSalesReport_TopProducts(t *testing.T) {
	report := models.SalesReport{
		TopProducts: []models.ProductStat{
			{Name: "Pizza", Quantity: 100, Revenue: 1200},
			{Name: "Burger", Quantity: 80, Revenue: 640},
		},
	}

	if len(report.TopProducts) != 2 {
		t.Errorf("expected 2 top products, got %d", len(report.TopProducts))
	}
	if report.TopProducts[0].Name != "Pizza" {
		t.Errorf("expected Pizza, got %s", report.TopProducts[0].Name)
	}
}
