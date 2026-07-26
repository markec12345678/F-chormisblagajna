package helpers

import (
	"fmt"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// MarkOrderFiscalized updates an order's fiscal status after successful fiscalization.
func MarkOrderFiscalized(cfg config.Config, log logger.ILogger, orderID, fiscalID string) error {
	client, err := common.GetDatabaseClient(log, &cfg)
	if err != nil {
		return fmt.Errorf("GetDatabaseClient: %w", err)
	}

	coll := client.Database(cfg.Databases[0].Database).Collection("orders")

	update := bson.M{
		"$set": bson.M{
			"is_fiscalized": true,
			"fiscal_id":     fiscalID,
		},
	}

	_, err = coll.UpdateOne(nil, bson.M{"id": orderID}, update)
	if err != nil {
		return fmt.Errorf("update order fiscal status: %w", err)
	}

	return nil
}
