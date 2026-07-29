package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/scheduling/dto"
	"github.com/nutrixpos/pos/modules/scheduling/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type SchedulingService struct {
	Logger logger.ILogger
	Config config.Config
}

type GetShiftsParams struct {
	PageNumber int
	PageSize   int
	BranchId   string
	StartDate  string
	EndDate    string
}

func (s *SchedulingService) GetShifts(params GetShiftsParams) ([]models.Shift, int, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, 0, fmt.Errorf("GetShifts: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("shifts")

	filter := bson.M{}
	if params.BranchId != "" {
		filter["branch_id"] = params.BranchId
	}
	if params.StartDate != "" && params.EndDate != "" {
		filter["date"] = bson.M{"$gte": params.StartDate, "$lte": params.EndDate}
	} else if params.StartDate != "" {
		filter["date"] = bson.M{"$gte": params.StartDate}
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64((params.PageNumber - 1) * params.PageSize))
	findOptions.SetLimit(int64(params.PageSize))
	findOptions.SetSort(bson.D{{Key: "date", Value: 1}, {Key: "start_time", Value: 1}})

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("GetShifts: %w", err)
	}
	defer cursor.Close(ctx)

	shifts := make([]models.Shift, 0)
	for cursor.Next(ctx) {
		var shift models.Shift
		if err := cursor.Decode(&shift); err != nil {
			return nil, 0, fmt.Errorf("GetShifts: %w", err)
		}
		shifts = append(shifts, shift)
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("GetShifts: %w", err)
	}

	return shifts, int(count), nil
}

func (s *SchedulingService) CreateShift(req dto.CreateShiftRequest) (models.Shift, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return models.Shift{}, fmt.Errorf("CreateShift: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("shifts")

	if req.Status == "" {
		req.Status = "scheduled"
	}

	now := time.Now()
	shift := models.Shift{
		Id:         bson.NewObjectID().Hex(),
		EmployeeId: req.EmployeeId,
		BranchId:   req.BranchId,
		Date:       req.Date,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Role:       req.Role,
		Status:     req.Status,
		Notes:      req.Notes,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	_, err = collection.InsertOne(ctx, shift)
	if err != nil {
		return models.Shift{}, fmt.Errorf("CreateShift: %w", err)
	}

	return shift, nil
}

func (s *SchedulingService) UpdateShift(id string, req dto.UpdateShiftRequest) (models.Shift, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return models.Shift{}, fmt.Errorf("UpdateShift: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("shifts")

	update := bson.M{"updated_at": time.Now()}
	if req.EmployeeId != "" {
		update["employee_id"] = req.EmployeeId
	}
	if req.BranchId != "" {
		update["branch_id"] = req.BranchId
	}
	if req.Date != "" {
		update["date"] = req.Date
	}
	if req.StartTime != "" {
		update["start_time"] = req.StartTime
	}
	if req.EndTime != "" {
		update["end_time"] = req.EndTime
	}
	if req.Role != "" {
		update["role"] = req.Role
	}
	if req.Status != "" {
		update["status"] = req.Status
	}
	if req.Notes != "" {
		update["notes"] = req.Notes
	}

	_, err = collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": update})
	if err != nil {
		return models.Shift{}, fmt.Errorf("UpdateShift: %w", err)
	}

	var updatedShift models.Shift
	err = collection.FindOne(ctx, bson.M{"id": id}).Decode(&updatedShift)
	if err != nil {
		return models.Shift{}, fmt.Errorf("UpdateShift: %w", err)
	}

	return updatedShift, nil
}

func (s *SchedulingService) DeleteShift(id string) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return fmt.Errorf("DeleteShift: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("shifts")

	_, err = collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return fmt.Errorf("DeleteShift: %w", err)
	}

	return nil
}
