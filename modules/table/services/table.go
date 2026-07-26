package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/table/dto"
	"github.com/nutrixpos/pos/modules/table/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type TableService struct {
	Logger logger.ILogger
	Config config.Config
}

type GetTablesParams struct {
	PageNumber int
	PageSize   int
}

func (s *TableService) GetTables(params GetTablesParams) ([]models.Table, int, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, 0, fmt.Errorf("GetTables: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("tables")

	findOptions := options.Find()
	findOptions.SetSkip(int64((params.PageNumber - 1) * params.PageSize))
	findOptions.SetLimit(int64(params.PageSize))

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("GetTables: %w", err)
	}
	defer cursor.Close(ctx)

	tables := make([]models.Table, 0)
	for cursor.Next(ctx) {
		var table models.Table
		if err := cursor.Decode(&table); err != nil {
			return nil, 0, fmt.Errorf("GetTables: %w", err)
		}
		tables = append(tables, table)
	}

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("GetTables: %w", err)
	}

	return tables, int(count), nil
}

func (s *TableService) GetTable(id string) (models.Table, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return models.Table{}, fmt.Errorf("GetTable: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("tables")

	var table models.Table
	err = collection.FindOne(ctx, bson.M{"id": id}).Decode(&table)
	if err != nil {
		return models.Table{}, fmt.Errorf("GetTable: %w", err)
	}

	return table, nil
}

func (s *TableService) CreateTable(req dto.CreateTableRequest) (models.Table, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return models.Table{}, fmt.Errorf("CreateTable: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("tables")

	now := time.Now()
	table := models.Table{
		Id:        bson.NewObjectID().Hex(),
		Number:    req.Number,
		Name:      req.Name,
		Capacity:  req.Capacity,
		Zone:      req.Zone,
		Status:    "available",
		BranchId:  req.BranchId,
		CreatedAt: now,
		UpdatedAt: now,
	}

	table.QRCode = s.GenerateQRCode(table.Id, table.Number)

	_, err = collection.InsertOne(ctx, table)
	if err != nil {
		return models.Table{}, fmt.Errorf("CreateTable: %w", err)
	}

	return table, nil
}

func (s *TableService) UpdateTable(id string, req dto.UpdateTableRequest) (models.Table, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return models.Table{}, fmt.Errorf("UpdateTable: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("tables")

	update := bson.M{"updated_at": time.Now()}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Capacity != 0 {
		update["capacity"] = req.Capacity
	}
	if req.Zone != "" {
		update["zone"] = req.Zone
	}
	if req.Status != "" {
		update["status"] = req.Status
	}

	_, err = collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": update})
	if err != nil {
		return models.Table{}, fmt.Errorf("UpdateTable: %w", err)
	}

	var updatedTable models.Table
	err = collection.FindOne(ctx, bson.M{"id": id}).Decode(&updatedTable)
	if err != nil {
		return models.Table{}, fmt.Errorf("UpdateTable: %w", err)
	}

	return updatedTable, nil
}

func (s *TableService) DeleteTable(id string) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return fmt.Errorf("DeleteTable: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("tables")

	_, err = collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return fmt.Errorf("DeleteTable: %w", err)
	}

	return nil
}

func (s *TableService) UpdateTableStatus(id string, status string) (models.Table, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return models.Table{}, fmt.Errorf("UpdateTableStatus: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("tables")

	_, err = collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.M{
		"status":     status,
		"updated_at": time.Now(),
	}})
	if err != nil {
		return models.Table{}, fmt.Errorf("UpdateTableStatus: %w", err)
	}

	var table models.Table
	err = collection.FindOne(ctx, bson.M{"id": id}).Decode(&table)
	if err != nil {
		return models.Table{}, fmt.Errorf("UpdateTableStatus: %w", err)
	}

	return table, nil
}

func (s *TableService) GenerateQRCode(tableId string, tableNumber int) string {
	frontendUrl := "http://localhost:8080"
	return fmt.Sprintf("%s/order?table=%d", frontendUrl, tableNumber)
}
