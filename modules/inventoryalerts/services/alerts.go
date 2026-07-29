package services

import (
	"context"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	alert_models "github.com/nutrixpos/pos/modules/inventoryalerts/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type InventoryAlertsService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *InventoryAlertsService) GetRules() ([]alert_models.InventoryAlertRule, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("inventory_alert_rules")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	rules := make([]alert_models.InventoryAlertRule, 0)
	for cursor.Next(ctx) {
		var rule alert_models.InventoryAlertRule
		if err := cursor.Decode(&rule); err != nil {
			continue
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

func (s *InventoryAlertsService) SaveRule(rule *alert_models.InventoryAlertRule) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("inventory_alert_rules")

	if rule.Id == "" {
		rule.Id = bson.NewObjectID().Hex()
		rule.CreatedAt = time.Now()
		_, err = collection.InsertOne(ctx, rule)
	} else {
		_, err = collection.UpdateOne(ctx,
			bson.M{"id": rule.Id},
			bson.M{"$set": rule},
		)
	}

	return err
}

func (s *InventoryAlertsService) DeleteRule(id string) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("inventory_alert_rules")

	_, err = collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (s *InventoryAlertsService) GetAlerts() ([]alert_models.InventoryAlert, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("inventory_alerts")

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	findOptions.SetLimit(100)

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	alerts := make([]alert_models.InventoryAlert, 0)
	for cursor.Next(ctx) {
		var alert alert_models.InventoryAlert
		if err := cursor.Decode(&alert); err != nil {
			continue
		}
		alerts = append(alerts, alert)
	}

	return alerts, nil
}

func (s *InventoryAlertsService) GetUnreadAlerts() ([]alert_models.InventoryAlert, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("inventory_alerts")

	cursor, err := collection.Find(ctx, bson.M{"is_read": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	alerts := make([]alert_models.InventoryAlert, 0)
	for cursor.Next(ctx) {
		var alert alert_models.InventoryAlert
		if err := cursor.Decode(&alert); err != nil {
			continue
		}
		alerts = append(alerts, alert)
	}

	return alerts, nil
}

func (s *InventoryAlertsService) MarkAsRead(id string) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("inventory_alerts")

	_, err = collection.UpdateOne(ctx,
		bson.M{"id": id},
		bson.M{"$set": bson.M{"is_read": true}},
	)
	return err
}

func (s *InventoryAlertsService) GetSummary() (*alert_models.AlertSummary, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rulesCol := client.Database(s.Config.Databases[0].Database).Collection("inventory_alert_rules")
	alertsCol := client.Database(s.Config.Databases[0].Database).Collection("inventory_alerts")

	totalActive, _ := rulesCol.CountDocuments(ctx, bson.M{"is_active": true})
	unreadCount, _ := alertsCol.CountDocuments(ctx, bson.M{"is_read": false})
	criticalCount, _ := alertsCol.CountDocuments(ctx, bson.M{"severity": "critical", "is_read": false})
	lowCount, _ := alertsCol.CountDocuments(ctx, bson.M{"severity": "low", "is_read": false})

	return &alert_models.AlertSummary{
		TotalActive:   int(totalActive),
		UnreadCount:   int(unreadCount),
		CriticalCount: int(criticalCount),
		LowCount:      int(lowCount),
	}, nil
}
