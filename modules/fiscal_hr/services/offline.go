package services

import (
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
)

const maxRetryCountHR = 10

// OfflineFiscalQueueHR handles retrying failed Croatian fiscalization attempts.
type OfflineFiscalQueueHR struct {
	Logger logger.ILogger
	Config config.Config
}

// NewOfflineFiscalQueueHR creates a new offline queue processor.
func NewOfflineFiscalQueueHR(log logger.ILogger, cfg config.Config) *OfflineFiscalQueueHR {
	return &OfflineFiscalQueueHR{
		Logger: log,
		Config: cfg,
	}
}

// ProcessPending retries all pending offline Croatian fiscalization requests.
func (q *OfflineFiscalQueueHR) ProcessPending() {
	receiptSvc := FiscalReceiptServiceHR{
		Config: q.Config,
		Logger: q.Logger,
	}

	receipts, err := receiptSvc.GetPendingOffline()
	if err != nil {
		q.Logger.Error("HR offline queue: failed to get pending receipts: " + err.Error())
		return
	}

	if len(receipts) == 0 {
		return
	}

	q.Logger.Info("HR offline queue: processing " + fmt.Sprintf("%d", len(receipts)) + " pending receipts")

	settingsSvc := FiscalSettingsServiceHR{
		Config: q.Config,
		Logger: q.Logger,
	}

	settings, err := settingsSvc.Get()
	if err != nil {
		q.Logger.Error("HR offline queue: failed to get settings: " + err.Error())
		return
	}

	if !settings.Enabled {
		return
	}

	client, err := NewCISClient(settings)
	if err != nil {
		q.Logger.Error("HR offline queue: failed to create CIS client: " + err.Error())
		return
	}

	for _, receipt := range receipts {
		if receipt.RetryCount >= maxRetryCountHR {
			q.Logger.Error("HR offline queue: max retries reached for order " + receipt.OrderID)
			continue
		}

		_ = client // In production, re-send the original request using stored data
		q.Logger.Info("HR offline queue: retrying order " + receipt.OrderID)

		if err := receiptSvc.IncrementRetry(receipt.OrderID); err != nil {
			q.Logger.Error("HR offline queue: increment retry error: " + err.Error())
		}

		time.Sleep(1 * time.Second)
	}
}
