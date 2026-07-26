package services

import (
	"fmt"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	fiscal_models "github.com/nutrixpos/pos/modules/fiscal_hr/models"
	core_models "github.com/nutrixpos/pos/modules/core/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const settingsCollection = "settings"

// FiscalSettingsServiceHR handles Croatian fiscal settings persistence.
type FiscalSettingsServiceHR struct {
	Config config.Config
	Logger logger.ILogger
}

// Get retrieves the Croatian fiscal settings from the database.
func (s *FiscalSettingsServiceHR) Get() (*fiscal_models.FiscalSettingsHR, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, fmt.Errorf("GetDatabaseClient: %w", err)
	}

	coll := client.Database(s.Config.Databases[0].Database).Collection(settingsCollection)

	var settings core_models.Settings
	err = coll.FindOne(nil, bson.M{}).Decode(&settings)
	if err != nil {
		return nil, fmt.Errorf("find settings: %w", err)
	}

	return &settings.FiscalHR, nil
}

// Update saves the Croatian fiscal settings to the database.
func (s *FiscalSettingsServiceHR) Update(settings *fiscal_models.FiscalSettingsHR) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return fmt.Errorf("GetDatabaseClient: %w", err)
	}

	coll := client.Database(s.Config.Databases[0].Database).Collection(settingsCollection)

	update := bson.M{
		"$set": bson.M{
			"fiscal_hr": settings,
		},
	}

	_, err = coll.UpdateOne(nil, bson.M{}, update)
	if err != nil {
		return fmt.Errorf("update fiscal settings HR: %w", err)
	}

	return nil
}
