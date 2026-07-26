package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/fiscal_hr/models"
	"github.com/nutrixpos/pos/modules/fiscal_hr/services"
)

// StornoHandlerHR reverses a fiscalized invoice by sending a storno to Croatian CIS.
// POST /api/fiscal_hr/storno
func StornoHandlerHR(cfg config.Config, log logger.ILogger) http.HandlerFunc {
	type request struct {
		OrderID     string  `json:"order_id"`
		JIR         string  `json:"jir"`
		InvoiceDate string  `json:"invoice_date"`
		TotalAmount float64 `json:"total_amount"`
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

		// Build storno as a negative invoice referencing the original JIR
		oib := settings.OIB
		datVrijeme := services.FormatDateTimeHR(time.Now())
		brOznRac := fmt.Sprintf("%d", settings.InvoiceNumber)
		oznPosPr := settings.BusinessPremiseID
		oznNapUr := settings.ElectronicDeviceID
		operatorOIB := settings.OperatorOIB

		pdvEntries := []services.PDVEntry{
			{TaxRate: 25, TaxableAmount: -req.TotalAmount, TaxAmount: 0},
		}

		zki, err := services.CalculateZKI(
			client.PrivateKey(),
			oib,
			datVrijeme,
			brOznRac,
			oznPosPr,
			oznNapUr,
			services.FormatAmountHR(-req.TotalAmount),
		)
		if err != nil {
			log.Error("ZKI error: " + err.Error())
			http.Error(w, "failed to calculate ZKI", http.StatusInternalServerError)
			return
		}

		bodyXML := services.BuildRacunXML(
			oib, datVrijeme, brOznRac, oznPosPr, oznNapUr,
			pdvEntries, -req.TotalAmount, "G", operatorOIB, zki,
		)

		signedBody, err := services.SignEnvelope(
			services.WrapInSOAP(bodyXML),
			client.PrivateKey(),
			client.Certificate(),
		)
		if err != nil {
			log.Error("XML sign error: " + err.Error())
			http.Error(w, "failed to sign XML", http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		resp, err := client.SendRaw(ctx, signedBody)
		if err != nil {
			log.Error("CIS storno error: " + err.Error())
			http.Error(w, "storno failed", http.StatusInternalServerError)
			return
		}

		jir, err := services.ParseRacunOdgovor(resp)
		if err != nil {
			log.Error("parse storno response: " + err.Error())
			http.Error(w, "failed to parse storno response", http.StatusInternalServerError)
			return
		}

		settings.InvoiceNumber++
		if err := settingsSvc.Update(settings); err != nil {
			log.Error("update settings error: " + err.Error())
		}

		receipt := models.FiscalReceiptHR{
			OrderID:       req.OrderID,
			JIR:           jir,
			InvoiceAmount: -req.TotalAmount,
			IssuedAt:      time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(receipt)
	}
}
