package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/fiscal/models"
	"github.com/nutrixpos/pos/modules/fiscal/services"
)

// StornoHandler issues a corrective (storno) invoice for a previously fiscalized receipt.
// POST /api/fiscal/storno
func StornoHandler(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	type request struct {
		OriginalOrderID string  `json:"original_order_id"`
		Reason          string  `json:"reason"`
		Amount          float64 `json:"amount"`
		TaxNumber       int     `json:"tax_number"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if req.OriginalOrderID == "" {
			http.Error(w, "original_order_id is required", http.StatusBadRequest)
			return
		}

		if req.Amount <= 0 {
			http.Error(w, "amount must be positive", http.StatusBadRequest)
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

		// Storno is a negative-amount invoice referencing the original
		items := []models.InvoiceItem{
			{
				Name:          "Storno: " + req.Reason,
				Quantity:      1,
				UnitPrice:     -req.Amount,
				TaxRate:       22.0,
				TaxableAmount: -req.Amount / 1.22,
				TaxAmount:     -req.Amount + req.Amount/1.22,
			},
		}

		ctx := r.Context()
		resp, err := client.FiscalizeInvoice(ctx, req.OriginalOrderID, items, -req.Amount, time.Now())
		if err != nil {
			log.Error("storno fiscalization error: " + err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		settingsSvc.Update(settings)

		result := models.FiscalReceipt{
			OrderID:       req.OriginalOrderID,
			EOR:           resp.UniqueInvoiceID,
			InvoiceAmount: -req.Amount,
			IssuedAt:      time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}
