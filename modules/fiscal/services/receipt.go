package services

import (
	"fmt"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/fiscal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const fiscalReceiptsCollection = "fiscal_receipts"

// FiscalReceiptService handles fiscal receipt storage and queries.
type FiscalReceiptService struct {
	Config config.Config
	Logger logger.ILogger
}

// Save stores a fiscal receipt in the database.
func (s *FiscalReceiptService) Save(receipt *models.FiscalReceipt) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return fmt.Errorf("GetDatabaseClient: %w", err)
	}

	coll := client.Database(s.Config.Databases[0].Database).Collection(fiscalReceiptsCollection)
	_, err = coll.InsertOne(nil, receipt)
	if err != nil {
		return fmt.Errorf("insert fiscal receipt: %w", err)
	}

	return nil
}

// GetByOrderID retrieves the fiscal receipt for a given order.
func (s *FiscalReceiptService) GetByOrderID(orderID string) (*models.FiscalReceipt, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, fmt.Errorf("GetDatabaseClient: %w", err)
	}

	coll := client.Database(s.Config.Databases[0].Database).Collection(fiscalReceiptsCollection)

	var receipt models.FiscalReceipt
	err = coll.FindOne(nil, bson.M{"order_id": orderID}).Decode(&receipt)
	if err != nil {
		return nil, fmt.Errorf("find fiscal receipt: %w", err)
	}

	return &receipt, nil
}

// GetPendingOffline retrieves all receipts that failed to fiscalize and need retry.
func (s *FiscalReceiptService) GetPendingOffline() ([]models.FiscalReceipt, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, fmt.Errorf("GetDatabaseClient: %w", err)
	}

	coll := client.Database(s.Config.Databases[0].Database).Collection(fiscalReceiptsCollection)

	cursor, err := coll.Find(nil, bson.M{"pending_offline": true})
	if err != nil {
		return nil, fmt.Errorf("find pending receipts: %w", err)
	}

	var receipts []models.FiscalReceipt
	if err := cursor.All(nil, &receipts); err != nil {
		return nil, fmt.Errorf("decode pending receipts: %w", err)
	}

	return receipts, nil
}

// MarkFiscalized marks a pending receipt as successfully fiscalized with the EOR.
func (s *FiscalReceiptService) MarkFiscalized(orderID, eor string) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return fmt.Errorf("GetDatabaseClient: %w", err)
	}

	coll := client.Database(s.Config.Databases[0].Database).Collection(fiscalReceiptsCollection)

	update := bson.M{
		"$set": bson.M{
			"pending_offline": false,
			"eor":             eor,
		},
	}

	_, err = coll.UpdateOne(nil, bson.M{"order_id": orderID}, update)
	return err
}

// IncrementRetry increments the retry counter for a pending receipt.
func (s *FiscalReceiptService) IncrementRetry(orderID string) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return fmt.Errorf("GetDatabaseClient: %w", err)
	}

	coll := client.Database(s.Config.Databases[0].Database).Collection(fiscalReceiptsCollection)

	update := bson.M{
		"$inc": bson.M{"retry_count": 1},
	}

	_, err = coll.UpdateOne(nil, bson.M{"order_id": orderID}, update)
	return err
}
