package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/modules/report/models"
)

func TestExportSalesReportCSV_ReturnsCSV(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/reports/sales/export", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment;filename=sales_report.csv")
		w.Write([]byte("Metric,Value\nTotal Revenue,15000.50\n"))
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/reports/sales/export?start_date=2026-07-01&end_date=2026-07-31", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	if rec.Header().Get("Content-Type") != "text/csv" {
		t.Errorf("expected text/csv, got %s", rec.Header().Get("Content-Type"))
	}

	if rec.Header().Get("Content-Disposition") != "attachment;filename=sales_report.csv" {
		t.Errorf("expected attachment disposition, got %s", rec.Header().Get("Content-Disposition"))
	}
}

func TestExportInventoryReportCSV_ReturnsCSV(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/reports/inventory/export", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment;filename=inventory_report.csv")
		w.Write([]byte("Metric,Value\nTotal Materials,50\n"))
	}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/reports/inventory/export", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	if rec.Header().Get("Content-Type") != "text/csv" {
		t.Errorf("expected text/csv, got %s", rec.Header().Get("Content-Type"))
	}
}

func TestSalesReport_TopProducts_JSON(t *testing.T) {
	report := models.SalesReport{
		TopProducts: []models.ProductStat{
			{Name: "Pizza", Quantity: 100, Revenue: 1200},
			{Name: "Burger", Quantity: 80, Revenue: 640},
			{Name: "Pasta", Quantity: 60, Revenue: 480},
		},
	}

	data, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded models.SalesReport
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(decoded.TopProducts) != 3 {
		t.Errorf("expected 3 products, got %d", len(decoded.TopProducts))
	}
}
