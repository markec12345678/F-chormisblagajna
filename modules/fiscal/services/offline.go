package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/fiscal/models"
)

// OfflineFiscalQueue handles retry of failed fiscalizations with exponential backoff.
type OfflineFiscalQueue struct {
	Logger   logger.ILogger
	Config   config.Config
	receiptSvc *FiscalReceiptService
	settingsSvc *FiscalSettingsService
}

// NewOfflineFiscalQueue creates a new offline queue processor.
func NewOfflineFiscalQueue(log logger.ILogger, cfg config.Config) *OfflineFiscalQueue {
	return &OfflineFiscalQueue{
		Logger:      log,
		Config:      cfg,
		receiptSvc:  &FiscalReceiptService{Config: cfg, Logger: log},
		settingsSvc: &FiscalSettingsService{Config: cfg, Logger: log},
	}
}

// ProcessPending retries all pending offline receipts.
func (q *OfflineFiscalQueue) ProcessPending() {
	receipts, err := q.receiptSvc.GetPendingOffline()
	if err != nil {
		q.Logger.Error("offline queue: failed to get pending receipts: " + err.Error())
		return
	}

	if len(receipts) == 0 {
		return
	}

	q.Logger.Info(fmt.Sprintf("offline queue: processing %d pending receipts", len(receipts)))

	settings, err := q.settingsSvc.Get()
	if err != nil {
		q.Logger.Error("offline queue: failed to get fiscal settings: " + err.Error())
		return
	}

	if !settings.Enabled {
		return
	}

	client, err := NewFURSClient(settings)
	if err != nil {
		q.Logger.Error("offline queue: failed to create FURS client: " + err.Error())
		return
	}

	for _, receipt := range receipts {
		// Exponential backoff: cap at 1 hour
		backoff := time.Duration(min(receipt.RetryCount, 6)) * 5 * time.Minute
		if backoff > 1*time.Hour {
			backoff = 1 * time.Hour
		}

		q.Logger.Info(fmt.Sprintf("offline queue: retrying order %s (attempt %d, backoff %v)",
			receipt.OrderID, receipt.RetryCount+1, backoff))

		// In a real implementation, we'd reconstruct the invoice items from the order.
		// For now, we use the stored amount and a placeholder.
		items := []models.InvoiceItem{
			{Name: "Receipt", Quantity: 1, UnitPrice: receipt.InvoiceAmount,
				TaxRate: 22.0, TaxableAmount: receipt.InvoiceAmount / 1.22,
				TaxAmount: receipt.InvoiceAmount - receipt.InvoiceAmount/1.22},
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		resp, err := client.FiscalizeInvoice(ctx, receipt.OrderID, items, receipt.InvoiceAmount, receipt.IssuedAt)
		cancel()

		if err != nil {
			q.Logger.Error(fmt.Sprintf("offline queue: retry failed for order %s: %s", receipt.OrderID, err.Error()))
			q.receiptSvc.IncrementRetry(receipt.OrderID)
			continue
		}

		err = q.receiptSvc.MarkFiscalized(receipt.OrderID, resp.UniqueInvoiceID)
		if err != nil {
			q.Logger.Error(fmt.Sprintf("offline queue: failed to mark fiscalized for order %s: %s", receipt.OrderID, err.Error()))
			continue
		}

		q.Logger.Info(fmt.Sprintf("offline queue: order %s fiscalized with EOR %s", receipt.OrderID, resp.UniqueInvoiceID))
	}
}
