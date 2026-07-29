package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/reservation/dto"
	"github.com/nutrixpos/pos/modules/reservation/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ReservationService struct {
	Logger logger.ILogger
	Config config.Config
}

type GetReservationsParams struct {
	PageNumber int
	PageSize   int
	BranchId   string
	Date       string
	Status     string
}

func (s *ReservationService) GetReservations(params GetReservationsParams) ([]models.Reservation, int, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, 0, fmt.Errorf("GetReservations: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("reservations")

	filter := bson.M{}
	if params.BranchId != "" {
		filter["branch_id"] = params.BranchId
	}
	if params.Date != "" {
		filter["date"] = params.Date
	}
	if params.Status != "" {
		filter["status"] = params.Status
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64((params.PageNumber - 1) * params.PageSize))
	findOptions.SetLimit(int64(params.PageSize))
	findOptions.SetSort(bson.D{{Key: "date", Value: 1}, {Key: "time", Value: 1}})

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("GetReservations: %w", err)
	}
	defer cursor.Close(ctx)

	reservations := make([]models.Reservation, 0)
	for cursor.Next(ctx) {
		var reservation models.Reservation
		if err := cursor.Decode(&reservation); err != nil {
			return nil, 0, fmt.Errorf("GetReservations: %w", err)
		}
		reservations = append(reservations, reservation)
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("GetReservations: %w", err)
	}

	return reservations, int(count), nil
}

func (s *ReservationService) GetReservation(id string) (models.Reservation, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return models.Reservation{}, fmt.Errorf("GetReservation: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("reservations")

	var reservation models.Reservation
	err = collection.FindOne(ctx, bson.M{"id": id}).Decode(&reservation)
	if err != nil {
		return models.Reservation{}, fmt.Errorf("GetReservation: %w", err)
	}

	return reservation, nil
}

func (s *ReservationService) CreateReservation(req dto.CreateReservationRequest) (models.Reservation, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return models.Reservation{}, fmt.Errorf("CreateReservation: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("reservations")

	if req.Status == "" {
		req.Status = "confirmed"
	}

	now := time.Now()
	reservation := models.Reservation{
		Id:            bson.NewObjectID().Hex(),
		CustomerName:  req.CustomerName,
		CustomerPhone: req.CustomerPhone,
		CustomerEmail: req.CustomerEmail,
		TableId:       req.TableId,
		BranchId:      req.BranchId,
		Date:          req.Date,
		Time:          req.Time,
		GuestCount:    req.GuestCount,
		Status:        req.Status,
		Notes:         req.Notes,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	_, err = collection.InsertOne(ctx, reservation)
	if err != nil {
		return models.Reservation{}, fmt.Errorf("CreateReservation: %w", err)
	}

	return reservation, nil
}

func (s *ReservationService) UpdateReservation(id string, req dto.UpdateReservationRequest) (models.Reservation, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return models.Reservation{}, fmt.Errorf("UpdateReservation: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("reservations")

	update := bson.M{"updated_at": time.Now()}
	if req.CustomerName != "" {
		update["customer_name"] = req.CustomerName
	}
	if req.CustomerPhone != "" {
		update["customer_phone"] = req.CustomerPhone
	}
	if req.CustomerEmail != "" {
		update["customer_email"] = req.CustomerEmail
	}
	if req.TableId != "" {
		update["table_id"] = req.TableId
	}
	if req.BranchId != "" {
		update["branch_id"] = req.BranchId
	}
	if req.Date != "" {
		update["date"] = req.Date
	}
	if req.Time != "" {
		update["time"] = req.Time
	}
	if req.GuestCount != nil {
		update["guest_count"] = *req.GuestCount
	}
	if req.Status != "" {
		update["status"] = req.Status
	}
	if req.Notes != "" {
		update["notes"] = req.Notes
	}

	_, err = collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": update})
	if err != nil {
		return models.Reservation{}, fmt.Errorf("UpdateReservation: %w", err)
	}

	var updatedReservation models.Reservation
	err = collection.FindOne(ctx, bson.M{"id": id}).Decode(&updatedReservation)
	if err != nil {
		return models.Reservation{}, fmt.Errorf("UpdateReservation: %w", err)
	}

	return updatedReservation, nil
}

func (s *ReservationService) DeleteReservation(id string) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return fmt.Errorf("DeleteReservation: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("reservations")

	_, err = collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return fmt.Errorf("DeleteReservation: %w", err)
	}

	return nil
}
