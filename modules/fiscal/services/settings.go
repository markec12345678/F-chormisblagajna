package services

import (
	"fmt"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/fiscal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const fiscalSettingsCollection = "fiscal_settings"

// FiscalSettingsService handles reading and writing FURS fiscal settings.
type FiscalSettingsService struct {
	Config config.Config
	Logger logger.ILogger
}

// Get retrieves the fiscal settings from the database.
func (s *FiscalSettingsService) Get() (*models.FiscalSettings, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, fmt.Errorf("GetDatabaseClient: %w", err)
	}

	coll := client.Database(s.Config.Databases[0].Database).Collection(fiscalSettingsCollection)

	var settings models.FiscalSettings
	err = coll.FindOne(nil, bson.M{}).Decode(&settings)
	if err != nil {
		return nil, fmt.Errorf("find fiscal settings: %w", err)
	}

	return &settings, nil
}

// Update saves the fiscal settings to the database.
func (s *FiscalSettingsService) Update(settings *models.FiscalSettings) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return fmt.Errorf("GetDatabaseClient: %w", err)
	}

	coll := client.Database(s.Config.Databases[0].Database).Collection(fiscalSettingsCollection)

	// Upsert: create if not exists, update if exists
	filter := bson.M{}
	update := bson.M{"$set": settings}

	_, err = coll.UpdateOne(nil, filter, update)
	if err != nil {
		return fmt.Errorf("upsert fiscal settings: %w", err)
	}

	return nil
}
