package services

import (
	"context"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	cd_models "github.com/nutrixpos/pos/modules/customerdisplay/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CustomerDisplayService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *CustomerDisplayService) getAllConfigsCollection() (*mongo.Collection, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}
	return client.Database(s.Config.Databases[0].Database).Collection("customer_displays"), nil
}

func (s *CustomerDisplayService) GetAllConfigs() ([]cd_models.DisplayConfig, error) {
	coll, err := s.getAllConfigsCollection()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	configs := make([]cd_models.DisplayConfig, 0)
	for cursor.Next(ctx) {
		var c cd_models.DisplayConfig
		if err := cursor.Decode(&c); err != nil {
			continue
		}
		configs = append(configs, c)
	}

	return configs, nil
}

func (s *CustomerDisplayService) GetConfig(id string) (*cd_models.DisplayConfig, error) {
	coll, err := s.getAllConfigsCollection()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var config cd_models.DisplayConfig
	err = coll.FindOne(ctx, bson.M{"id": id}).Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (s *CustomerDisplayService) SaveConfig(config *cd_models.DisplayConfig) error {
	coll, err := s.getAllConfigsCollection()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if config.Id == "" {
		config.Id = bson.NewObjectID().Hex()
	}

	_, err = coll.ReplaceOne(
		ctx,
		bson.M{"id": config.Id},
		config,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *CustomerDisplayService) DeleteConfig(id string) error {
	coll, err := s.getAllConfigsCollection()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = coll.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (s *CustomerDisplayService) GetDisplayContent(displayId string) (*cd_models.DisplayContent, error) {
	config, err := s.GetConfig(displayId)
	if err != nil {
		return nil, err
	}

	if !config.Active {
		return &cd_models.DisplayContent{Items: []cd_models.DisplayItem{}, Interval: config.AutoSlideInterval, Theme: config.Theme, WelcomeMsg: config.WelcomeMessage}, nil
	}

	items := make([]cd_models.DisplayItem, 0)

	if config.ShowPromotions {
		items = append(items, cd_models.DisplayItem{Type: "promotions", Content: "promotion_placeholder"})
	}

	if config.ShowMenu {
		items = append(items, cd_models.DisplayItem{Type: "menu", Content: config.MenuCategoryIds})
	}

	if config.ShowOrderStatus {
		orders, _ := s.getPendingOrdersInfo()
		items = append(items, cd_models.DisplayItem{Type: "order_status", Content: orders})
	}

	return &cd_models.DisplayContent{
		Items:      items,
		Interval:   config.AutoSlideInterval,
		Theme:      config.Theme,
		WelcomeMsg: config.WelcomeMessage,
	}, nil
}

func (s *CustomerDisplayService) getPendingOrdersInfo() ([]map[string]interface{}, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll := client.Database(s.Config.Databases[0].Database).Collection("orders")
	cursor, err := coll.Find(ctx, bson.M{
		"status": bson.M{"$in": []string{"pending", "preparing"}},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	orders := make([]map[string]interface{}, 0)
	for cursor.Next(ctx) {
		var order map[string]interface{}
		if err := cursor.Decode(&order); err != nil {
			continue
		}
		orders = append(orders, order)
	}

	return orders, nil
}
