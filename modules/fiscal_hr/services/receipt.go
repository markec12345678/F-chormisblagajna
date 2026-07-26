package services

import (
	"fmt"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/fiscal_hr/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const fiscalReceiptsHRCollection = "fiscal_receipts_hr"

// FiscalReceiptServiceHR handles Croatian fiscal receipt storage and queries.
type FiscalReceiptServiceHR struct {
	Config config.Config
	Logger logger.ILogger
}

// Save stores a Croatian fiscal receipt in the database.
func (s *FiscalReceiptServiceHR) Save(receipt *models.FiscalReceiptHR) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return fmt.Errorf("GetDatabaseClient: %w", err)
	}

	coll := client.Database(s.Config.Databases[0].Database).Collection(fiscalReceiptsHRCollection)
	_, err = coll.InsertOne(nil, receipt)
	if err != nil {
		return fmt.Errorf("insert fiscal receipt HR: %w", err)
	}

	return nil
}

// GetByOrderID retrieves the Croatian fiscal receipt for a given order.
func (s *FiscalReceiptServiceHR) GetByOrderID(orderID string) (*models.FiscalReceiptHR, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, fmt.Errorf("GetDatabaseClient: %w", err)
	}

	coll := client.Database(s.Config.Databases[0].Database).Collection(fiscalReceiptsHRCollection)

	var receipt models.FiscalReceiptHR
	err = coll.FindOne(nil, bson.M{"order_id": orderID}).Decode(&receipt)
	if err != nil {
		return nil, fmt.Errorf("find fiscal receipt HR: %w", err)
	}

	return &receipt, nil
}

// GetPendingOffline retrieves all Croatian receipts that failed to fiscalize.
func (s *FiscalReceiptServiceHR) GetPendingOffline() ([]models.FiscalReceiptHR, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, fmt.Errorf("GetDatabaseClient: %w", err)
	}

	coll := client.Database(s.Config.Databases[0].Database).Collection(fiscalReceiptsHRCollection)

	cursor, err := coll.Find(nil, bson.M{"pending_offline": true})
	if err != nil {
		return nil, fmt.Errorf("find pending receipts HR: %w", err)
	}

	var receipts []models.FiscalReceiptHR
	if err := cursor.All(nil, &receipts); err != nil {
		return nil, fmt.Errorf("decode pending receipts HR: %w", err)
	}

	return receipts, nil
}

// MarkFiscalized marks a pending receipt as successfully fiscalized with the JIR.
func (s *FiscalReceiptServiceHR) MarkFiscalized(orderID, jir string) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return fmt.Errorf("GetDatabaseClient: %w", err)
	}

	coll := client.Database(s.Config.Databases[0].Database).Collection(fiscalReceiptsHRCollection)

	update := bson.M{
		"$set": bson.M{
			"pending_offline": false,
			"jir":             jir,
		},
	}

	_, err = coll.UpdateOne(nil, bson.M{"order_id": orderID}, update)
	return err
}

// IncrementRetry increments the retry counter for a pending Croatian receipt.
func (s *FiscalReceiptServiceHR) IncrementRetry(orderID string) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return fmt.Errorf("GetDatabaseClient: %w", err)
	}

	coll := client.Database(s.Config.Databases[0].Database).Collection(fiscalReceiptsHRCollection)

	update := bson.M{
		"$inc": bson.M{"retry_count": 1},
	}

	_, err = coll.UpdateOne(nil, bson.M{"order_id": orderID}, update)
	return err
}
