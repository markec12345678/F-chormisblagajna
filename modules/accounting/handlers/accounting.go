package handlers

import (
	"net/http"
	"time"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/accounting/services"
)

func ExportQuickBooksCSV(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startDate := r.URL.Query().Get("start_date")
		endDate := r.URL.Query().Get("end_date")

		start, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			start = time.Now().AddDate(0, 0, -30)
		}

		end, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			end = time.Now()
		}

		svc := services.AccountingService{Logger: log, Config: cfg}
		csvData, err := svc.ExportQuickBooks(start, end)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "export failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=quickbooks_export.csv")
		w.Write(csvData)
	}
}

func ExportXeroCSV(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startDate := r.URL.Query().Get("start_date")
		endDate := r.URL.Query().Get("end_date")

		start, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			start = time.Now().AddDate(0, 0, -30)
		}

		end, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			end = time.Now()
		}

		svc := services.AccountingService{Logger: log, Config: cfg}
		csvData, err := svc.ExportXero(start, end)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "export failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=xero_export.csv")
		w.Write(csvData)
	}
}
