package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/inventorytransfer/dto"
	"github.com/nutrixpos/pos/modules/inventorytransfer/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type InventoryTransferService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *InventoryTransferService) GetDatabase() (*mongo.Database, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}
	return client.Database(s.Config.Databases[0].Database), nil
}

func (s *InventoryTransferService) CreateTransfer(req dto.CreateTransferRequest) (*models.InventoryTransfer, error) {
	db, err := s.GetDatabase()
	if err != nil {
		return nil, fmt.Errorf("CreateTransfer: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	transfer := models.InventoryTransfer{
		MaterialID:   req.MaterialID,
		MaterialName: req.MaterialName,
		Quantity:     req.Quantity,
		Unit:         req.Unit,
		FromBranchID: req.FromBranchID,
		ToBranchID:   req.ToBranchID,
		Status:       "pending",
		Notes:        req.Notes,
		CreatedAt:    now,
	}

	result, err := db.Collection("inventory_transfers").InsertOne(ctx, transfer)
	if err != nil {
		return nil, fmt.Errorf("CreateTransfer InsertOne: %w", err)
	}

	transfer.ID = result.InsertedID.(bson.ObjectID)
	return &transfer, nil
}

func (s *InventoryTransferService) GetAllTransfers() ([]models.InventoryTransfer, error) {
	db, err := s.GetDatabase()
	if err != nil {
		return nil, fmt.Errorf("GetAllTransfers: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.M{"created_at": -1})
	cursor, err := db.Collection("inventory_transfers").Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("GetAllTransfers Find: %w", err)
	}
	defer cursor.Close(ctx)

	var transfers []models.InventoryTransfer
	if err := cursor.All(ctx, &transfers); err != nil {
		return nil, fmt.Errorf("GetAllTransfers decode: %w", err)
	}

	return transfers, nil
}

func (s *InventoryTransferService) UpdateTransferStatus(id string, status string) (*models.InventoryTransfer, error) {
	db, err := s.GetDatabase()
	if err != nil {
		return nil, fmt.Errorf("UpdateTransferStatus: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("UpdateTransferStatus invalid ID: %w", err)
	}

	now := time.Now()
	update := bson.M{"$set": bson.M{"status": status}}
	if status == "completed" {
		update["$set"].(bson.M)["completed_at"] = now
	}

	result, err := db.Collection("inventory_transfers").UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return nil, fmt.Errorf("UpdateTransferStatus UpdateOne: %w", err)
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("transfer not found")
	}

	var transfer models.InventoryTransfer
	err = db.Collection("inventory_transfers").FindOne(ctx, bson.M{"_id": objID}).Decode(&transfer)
	if err != nil {
		return nil, fmt.Errorf("UpdateTransferStatus FindOne: %w", err)
	}

	return &transfer, nil
}
