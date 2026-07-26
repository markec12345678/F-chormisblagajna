package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/helpers"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/fiscal/models"
	"github.com/nutrixpos/pos/modules/fiscal/services"
)

// FiscalizeOrderHandler fiscalizes a paid order by sending it to FURS.
// POST /api/fiscal/invoice
func FiscalizeOrderHandler(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	type itemRequest struct {
		Name          string  `json:"name"`
		Quantity      float64 `json:"quantity"`
		UnitPrice     float64 `json:"unit_price"`
		TaxRate       float64 `json:"tax_rate"`
		TaxableAmount float64 `json:"taxable_amount"`
		TaxAmount     float64 `json:"tax_amount"`
	}

	type request struct {
		OrderID     string        `json:"order_id"`
		Items       []itemRequest `json:"items"`
		TotalAmount float64       `json:"total_amount"`
		TaxNumber   int           `json:"tax_number"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		settingsSvc := services.FiscalSettingsService{
			Config: cfg,
			Logger: log,
		}

		settings, err := settingsSvc.Get()
		if err != nil {
			http.Error(w, "fiscal settings not configured", http.StatusBadRequest)
			return
		}

		if !settings.Enabled {
			http.Error(w, "fiscalization is disabled", http.StatusBadRequest)
			return
		}

		if req.TaxNumber > 0 {
			settings.TaxNumber = req.TaxNumber
		}

		client, err := services.NewFURSClient(settings)
		if err != nil {
			log.Error("furs client error: " + err.Error())
			http.Error(w, "failed to initialize fiscal client", http.StatusInternalServerError)
			return
		}

		items := make([]models.InvoiceItem, len(req.Items))
		for i, item := range req.Items {
			items[i] = models.InvoiceItem{
				Name:          item.Name,
				Quantity:      item.Quantity,
				UnitPrice:     item.UnitPrice,
				TaxRate:       item.TaxRate,
				TaxableAmount: item.TaxableAmount,
				TaxAmount:     item.TaxAmount,
			}
		}

		ctx := r.Context()
		resp, err := client.FiscalizeInvoice(ctx, req.OrderID, items, req.TotalAmount, time.Now())
		if err != nil {
			log.Error("fiscalization error: " + err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		settingsSvc.Update(settings)

		result := models.FiscalReceipt{
			OrderID:       req.OrderID,
			EOR:           resp.UniqueInvoiceID,
			InvoiceAmount: req.TotalAmount,
			IssuedAt:      time.Now(),
		}

		if err := helpers.MarkOrderFiscalized(cfg, log, req.OrderID, resp.UniqueInvoiceID); err != nil {
			log.Error("update order fiscal status: " + err.Error())
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

// FiscalEchoHandler tests the connection to FURS.
// POST /api/fiscal/echo
func FiscalEchoHandler(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		settingsSvc := services.FiscalSettingsService{
			Config: cfg,
			Logger: log,
		}

		settings, err := settingsSvc.Get()
		if err != nil {
			http.Error(w, "fiscal settings not configured", http.StatusBadRequest)
			return
		}

		client, err := services.NewFURSClient(settings)
		if err != nil {
			log.Error("furs client error: " + err.Error())
			http.Error(w, "failed to initialize fiscal client", http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		if err := client.Echo(ctx); err != nil {
			log.Error("furs echo error: " + err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}

// GetFiscalSettingsHandler returns the current fiscal settings.
// GET /api/fiscal/settings
func GetFiscalSettingsHandler(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		settingsSvc := services.FiscalSettingsService{
			Config: cfg,
			Logger: log,
		}

		settings, err := settingsSvc.Get()
		if err != nil {
			http.Error(w, "fiscal settings not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(settings)
	}
}

// UpdateFiscalSettingsHandler updates the fiscal settings.
// PATCH /api/fiscal/settings
func UpdateFiscalSettingsHandler(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var settings models.FiscalSettings
		if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		settingsSvc := services.FiscalSettingsService{
			Config: cfg,
			Logger: log,
		}

		if err := settingsSvc.Update(&settings); err != nil {
			log.Error("update fiscal settings error: " + err.Error())
			http.Error(w, "failed to update settings", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(settings)
	}
}
