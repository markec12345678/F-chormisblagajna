package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/report/services"
)

func ExportSalesReportCSV(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startDate := r.URL.Query().Get("start_date")
		endDate := r.URL.Query().Get("end_date")

		svc := services.ReportService{Logger: logger, Config: config}
		report, err := svc.GetSalesReport(startDate, endDate)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment;filename=sales_report.csv")

		writer := csv.NewWriter(w)
		defer writer.Flush()

		writer.Write([]string{"Metric", "Value"})
		writer.Write([]string{"Period", report.Period})
		writer.Write([]string{"Total Revenue", fmt.Sprintf("%.2f", report.TotalRevenue)})
		writer.Write([]string{"Total Orders", strconv.Itoa(report.TotalOrders)})
		writer.Write([]string{"Average Order", fmt.Sprintf("%.2f", report.AverageOrder)})
		writer.Write([]string{"Total Items", strconv.Itoa(report.TotalItems)})
		writer.Write([]string{"Refund Amount", fmt.Sprintf("%.2f", report.RefundAmount)})
		writer.Write([]string{"Net Revenue", fmt.Sprintf("%.2f", report.NetRevenue)})
		writer.Write([]string{""})
		writer.Write([]string{"Top Products"})
		writer.Write([]string{"Name", "Quantity", "Revenue"})
		for _, p := range report.TopProducts {
			writer.Write([]string{p.Name, strconv.Itoa(p.Quantity), fmt.Sprintf("%.2f", p.Revenue)})
		}
	}
}

func ExportInventoryReportCSV(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc := services.ReportService{Logger: logger, Config: config}
		report, err := svc.GetInventoryReport()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment;filename=inventory_report.csv")

		writer := csv.NewWriter(w)
		defer writer.Flush()

		writer.Write([]string{"Metric", "Value"})
		writer.Write([]string{"Total Materials", strconv.Itoa(report.TotalMaterials)})
		writer.Write([]string{"Low Stock Count", strconv.Itoa(report.LowStockCount)})
		writer.Write([]string{"Out of Stock Count", strconv.Itoa(report.OutOfStockCount)})
		writer.Write([]string{"Total Value", fmt.Sprintf("%.2f", report.TotalValue)})
		writer.Write([]string{""})
		writer.Write([]string{"Low Stock Items"})
		writer.Write([]string{"Name", "Quantity", "Unit", "Value"})
		for _, item := range report.LowStockItems {
			writer.Write([]string{item.Name, fmt.Sprintf("%.2f", item.Quantity), item.Unit, fmt.Sprintf("%.2f", item.Value)})
		}
	}
}
