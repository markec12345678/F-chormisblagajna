package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/helpers"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/fiscal_hr/models"
	"github.com/nutrixpos/pos/modules/fiscal_hr/services"
)

// FiscalizeOrderHandlerHR fiscalizes a paid order by sending it to Croatian CIS.
// POST /api/fiscal_hr/invoice
func FiscalizeOrderHandlerHR(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	type itemRequest struct {
		Name          string  `json:"name"`
		Quantity      float64 `json:"quantity"`
		UnitPrice     float64 `json:"unit_price"`
		TaxRate       float64 `json:"tax_rate"`
		TaxableAmount float64 `json:"taxable_amount"`
		TaxAmount     float64 `json:"tax_amount"`
	}

	type request struct {
		OrderID       string        `json:"order_id"`
		Items         []itemRequest `json:"items"`
		TotalAmount   float64       `json:"total_amount"`
		PaymentMethod string        `json:"payment_method"`
		OperatorOIB   string        `json:"operator_oib"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		settingsSvc := services.FiscalSettingsServiceHR{
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

		client, err := services.NewCISClient(settings)
		if err != nil {
			log.Error("CIS client error: " + err.Error())
			http.Error(w, "failed to initialize fiscal client", http.StatusInternalServerError)
			return
		}

		items := make([]models.InvoiceItemHR, len(req.Items))
		for i, item := range req.Items {
			items[i] = models.InvoiceItemHR{
				Name:          item.Name,
				Quantity:      item.Quantity,
				UnitPrice:     item.UnitPrice,
				TaxRate:       item.TaxRate,
				TaxableAmount: item.TaxableAmount,
				TaxAmount:     item.TaxAmount,
			}
		}

		fiscalReq := &models.InvoiceRequestHR{
			OrderID:       req.OrderID,
			Items:         items,
			TotalAmount:   req.TotalAmount,
			PaymentMethod: req.PaymentMethod,
			OperatorOIB:   req.OperatorOIB,
		}

		ctx := r.Context()
		resp, err := client.FiscalizeInvoice(ctx, fiscalReq)
		if err != nil {
			log.Error("Croatian fiscalization error: " + err.Error())
			http.Error(w, "fiscalization failed", http.StatusInternalServerError)
			return
		}

		settings.InvoiceNumber++
		if err := settingsSvc.Update(settings); err != nil {
			log.Error("update settings error: " + err.Error())
		}

		receipt := models.FiscalReceiptHR{
			OrderID:       req.OrderID,
			JIR:           resp.JIR,
			InvoiceAmount: req.TotalAmount,
			IssuedAt:      time.Now(),
		}

		if err := helpers.MarkOrderFiscalized(cfg, log, req.OrderID, resp.JIR); err != nil {
			log.Error("update order fiscal status HR: " + err.Error())
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(receipt)
	}
}

// FiscalEchoHandlerHR tests the connection to CIS.
// POST /api/fiscal_hr/echo
func FiscalEchoHandlerHR(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		settingsSvc := services.FiscalSettingsServiceHR{
			Config: cfg,
			Logger: log,
		}

		settings, err := settingsSvc.Get()
		if err != nil {
			http.Error(w, "fiscal settings not configured", http.StatusBadRequest)
			return
		}

		client, err := services.NewCISClient(settings)
		if err != nil {
			log.Error("CIS client error: " + err.Error())
			http.Error(w, "failed to initialize fiscal client", http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		if err := client.Echo(ctx); err != nil {
			log.Error("CIS echo error: " + err.Error())
			http.Error(w, "CIS echo failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}

// GetFiscalSettingsHandlerHR returns the current Croatian fiscal settings.
// GET /api/fiscal_hr/settings
func GetFiscalSettingsHandlerHR(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		settingsSvc := services.FiscalSettingsServiceHR{
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

// UpdateFiscalSettingsHandlerHR updates the Croatian fiscal settings.
// PATCH /api/fiscal_hr/settings
func UpdateFiscalSettingsHandlerHR(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var settings models.FiscalSettingsHR
		if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		settingsSvc := services.FiscalSettingsServiceHR{
			Config: cfg,
			Logger: log,
		}

		if err := settingsSvc.Update(&settings); err != nil {
			log.Error("update fiscal settings HR error: " + err.Error())
			http.Error(w, "failed to update settings", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(settings)
	}
}
