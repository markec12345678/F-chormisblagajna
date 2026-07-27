package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/branch/dto"
	"github.com/nutrixpos/pos/modules/branch/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type BranchService struct {
	Logger logger.ILogger
	Config config.Config
}

type GetBranchesParams struct {
	PageNumber int
	PageSize   int
}

func (s *BranchService) GetBranches(params GetBranchesParams) ([]models.Branch, int, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, 0, fmt.Errorf("GetBranches: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("branches")

	findOptions := options.Find()
	findOptions.SetSkip(int64((params.PageNumber - 1) * params.PageSize))
	findOptions.SetLimit(int64(params.PageSize))

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("GetBranches: %w", err)
	}
	defer cursor.Close(ctx)

	branches := make([]models.Branch, 0)
	for cursor.Next(ctx) {
		var branch models.Branch
		if err := cursor.Decode(&branch); err != nil {
			return nil, 0, fmt.Errorf("GetBranches: %w", err)
		}
		branches = append(branches, branch)
	}

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("GetBranches: %w", err)
	}

	return branches, int(count), nil
}

func (s *BranchService) GetBranch(id string) (models.Branch, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return models.Branch{}, fmt.Errorf("GetBranch: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("branches")

	var branch models.Branch
	err = collection.FindOne(ctx, bson.M{"id": id}).Decode(&branch)
	if err != nil {
		return models.Branch{}, fmt.Errorf("GetBranch: %w", err)
	}

	return branch, nil
}

func (s *BranchService) CreateBranch(req dto.CreateBranchRequest) (models.Branch, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return models.Branch{}, fmt.Errorf("CreateBranch: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("branches")

	now := time.Now()
	branch := models.Branch{
		Id:        bson.NewObjectID().Hex(),
		Name:      req.Name,
		Address:   req.Address,
		Phone:     req.Phone,
		Email:     req.Email,
		IsActive:  req.IsActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err = collection.InsertOne(ctx, branch)
	if err != nil {
		return models.Branch{}, fmt.Errorf("CreateBranch: %w", err)
	}

	return branch, nil
}

func (s *BranchService) UpdateBranch(id string, req dto.UpdateBranchRequest) (models.Branch, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return models.Branch{}, fmt.Errorf("UpdateBranch: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("branches")

	update := bson.M{"updated_at": time.Now()}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Address != "" {
		update["address"] = req.Address
	}
	if req.Phone != "" {
		update["phone"] = req.Phone
	}
	if req.Email != "" {
		update["email"] = req.Email
	}
	if req.IsActive != nil {
		update["is_active"] = *req.IsActive
	}

	_, err = collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": update})
	if err != nil {
		return models.Branch{}, fmt.Errorf("UpdateBranch: %w", err)
	}

	var updatedBranch models.Branch
	err = collection.FindOne(ctx, bson.M{"id": id}).Decode(&updatedBranch)
	if err != nil {
		return models.Branch{}, fmt.Errorf("UpdateBranch: %w", err)
	}

	return updatedBranch, nil
}

func (s *BranchService) DeleteBranch(id string) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return fmt.Errorf("DeleteBranch: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("branches")

	_, err = collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return fmt.Errorf("DeleteBranch: %w", err)
	}

	return nil
}

func (s *BranchService) GetBranchStats(id string) (map[string]interface{}, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, fmt.Errorf("GetBranchStats: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	branch, err := s.GetBranch(id)
	if err != nil {
		return nil, fmt.Errorf("GetBranchStats: %w", err)
	}

	ordersCol := client.Database(s.Config.Databases[0].Database).Collection("orders")
	tablesCol := client.Database(s.Config.Databases[0].Database).Collection("tables")

	orderCount, _ := ordersCol.CountDocuments(ctx, bson.M{"branch_id": id})
	tableCount, _ := tablesCol.CountDocuments(ctx, bson.M{"branch_id": id})

	return map[string]interface{}{
		"branch":      branch,
		"order_count": orderCount,
		"table_count": tableCount,
	}, nil
}
