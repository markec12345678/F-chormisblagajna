package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type KitchenService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *KitchenService) GetDatabase() (*mongo.Database, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}
	return client.Database(s.Config.Databases[0].Database), nil
}

func (s *KitchenService) GetStations() ([]bson.M, error) {
	db, err := s.GetDatabase()
	if err != nil {
		return nil, fmt.Errorf("GetStations: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.Collection("kitchen_stations").Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("GetStations Find: %w", err)
	}
	defer cursor.Close(ctx)

	var stations []bson.M
	if err := cursor.All(ctx, &stations); err != nil {
		return nil, fmt.Errorf("GetStations decode: %w", err)
	}

	return stations, nil
}

func (s *KitchenService) CreateStation(name, branchID string) (*bson.M, error) {
	db, err := s.GetDatabase()
	if err != nil {
		return nil, fmt.Errorf("CreateStation: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	station := bson.M{
		"name":       name,
		"branch_id":  branchID,
		"active":     true,
		"created_at": now,
	}

	result, err := db.Collection("kitchen_stations").InsertOne(ctx, station)
	if err != nil {
		return nil, fmt.Errorf("CreateStation InsertOne: %w", err)
	}

	station["_id"] = result.InsertedID
	return &station, nil
}

func (s *KitchenService) UpdateItemStatus(orderID string, itemIndex int, status string) error {
	db, err := s.GetDatabase()
	if err != nil {
		return fmt.Errorf("UpdateItemStatus: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			fmt.Sprintf("items.%d.status", itemIndex): status,
			"updated_at": now,
		},
	}

	if status == "preparing" {
		update["$set"].(bson.M)[fmt.Sprintf("items.%d.started_at", itemIndex)] = now
	} else if status == "ready" {
		update["$set"].(bson.M)[fmt.Sprintf("items.%d.ready_at", itemIndex)] = now
	}

	result, err := db.Collection("orders").UpdateOne(ctx, bson.M{"_id": orderID}, update)
	if err != nil {
		return fmt.Errorf("UpdateItemStatus UpdateOne: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

func (s *KitchenService) GetOrdersByStation(station string) ([]bson.M, error) {
	db, err := s.GetDatabase()
	if err != nil {
		return nil, fmt.Errorf("GetOrdersByStation: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"state": bson.M{"$in": []string{"pending", "in_progress"}},
	}
	if station != "" {
		filter["station"] = station
	}

	opts := options.Find().SetSort(bson.M{"submitted_at": 1})
	cursor, err := db.Collection("orders").Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("GetOrdersByStation Find: %w", err)
	}
	defer cursor.Close(ctx)

	var orders []bson.M
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, fmt.Errorf("GetOrdersByStation decode: %w", err)
	}

	return orders, nil
}
