package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/splitbill/dto"
	"github.com/nutrixpos/pos/modules/splitbill/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type SplitBillService struct {
	Logger logger.ILogger
	Config config.Config
}

type GetSplitBillsParams struct {
	PageNumber int
	PageSize   int
}

func (s *SplitBillService) CreateSplitBill(req dto.CreateSplitBillRequest, orderTotal float64) (models.SplitBill, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return models.SplitBill{}, fmt.Errorf("CreateSplitBill: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("split_bills")

	var parts []models.SplitPart
	now := time.Now()

	switch req.SplitType {
	case "equal":
		count := req.SplitCount
		if count < 2 {
			count = 2
		}
		eachAmount := math.Round(orderTotal/float64(count)*100) / 100
		for i := 0; i < count; i++ {
			amount := eachAmount
			if i == count-1 {
				remaining := math.Round((orderTotal-eachAmount*float64(count-1))*100) / 100
				amount = eachAmount + remaining
			}
			parts = append(parts, models.SplitPart{
				Id:     bson.NewObjectID().Hex(),
				Amount: amount,
				IsPaid: false,
			})
		}
	case "custom":
		for _, amt := range req.CustomAmounts {
			parts = append(parts, models.SplitPart{
				Id:     bson.NewObjectID().Hex(),
				Amount: amt,
				IsPaid: false,
			})
		}
	case "by_item":
		partMap := make(map[int]float64)
		for _, is := range req.ItemSplits {
			partMap[is.PartIndex] += 0
		}
		count := len(partMap)
		if count < 2 {
			count = 2
		}
		eachAmount := math.Round(orderTotal/float64(count)*100) / 100
		for i := 0; i < count; i++ {
			amount := eachAmount
			if i == count-1 {
				remaining := math.Round((orderTotal-eachAmount*float64(count-1))*100) / 100
				amount = eachAmount + remaining
			}
			parts = append(parts, models.SplitPart{
				Id:     bson.NewObjectID().Hex(),
				Amount: amount,
				IsPaid: false,
			})
		}
	default:
		return models.SplitBill{}, fmt.Errorf("CreateSplitBill: unknown split_type %s", req.SplitType)
	}

	splitBill := models.SplitBill{
		Id:        bson.NewObjectID().Hex(),
		OrderId:   req.OrderId,
		SplitType: req.SplitType,
		Parts:     parts,
		Status:    "pending",
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err = collection.InsertOne(ctx, splitBill)
	if err != nil {
		return models.SplitBill{}, fmt.Errorf("CreateSplitBill: %w", err)
	}

	return splitBill, nil
}

func (s *SplitBillService) GetSplitBill(id string) (models.SplitBill, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return models.SplitBill{}, fmt.Errorf("GetSplitBill: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("split_bills")

	var splitBill models.SplitBill
	err = collection.FindOne(ctx, bson.M{"id": id}).Decode(&splitBill)
	if err != nil {
		return models.SplitBill{}, fmt.Errorf("GetSplitBill: %w", err)
	}

	return splitBill, nil
}

func (s *SplitBillService) GetSplitBillByOrder(orderId string) (models.SplitBill, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return models.SplitBill{}, fmt.Errorf("GetSplitBillByOrder: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("split_bills")

	var splitBill models.SplitBill
	err = collection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&splitBill)
	if err != nil {
		return models.SplitBill{}, fmt.Errorf("GetSplitBillByOrder: %w", err)
	}

	return splitBill, nil
}

func (s *SplitBillService) PaySplitPart(splitBillId string, req dto.PaySplitPartRequest) (models.SplitBill, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return models.SplitBill{}, fmt.Errorf("PaySplitPart: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("split_bills")

	splitBill, err := s.GetSplitBill(splitBillId)
	if err != nil {
		return models.SplitBill{}, fmt.Errorf("PaySplitPart: %w", err)
	}

	now := time.Now()
	allPaid := true
	for i, part := range splitBill.Parts {
		if part.Id == req.PartId {
			splitBill.Parts[i].IsPaid = true
			splitBill.Parts[i].PaymentMethod = req.PaymentMethod
			splitBill.Parts[i].PaidAt = &now
		}
		if !splitBill.Parts[i].IsPaid {
			allPaid = false
		}
	}

	status := "partial"
	if allPaid {
		status = "paid"
	}

	_, err = collection.UpdateOne(ctx, bson.M{"id": splitBillId}, bson.M{"$set": bson.M{
		"parts":      splitBill.Parts,
		"status":     status,
		"updated_at": now,
	}})
	if err != nil {
		return models.SplitBill{}, fmt.Errorf("PaySplitPart: %w", err)
	}

	splitBill.Status = status
	splitBill.UpdatedAt = now
	return splitBill, nil
}

func (s *SplitBillService) GetSplitBills(params GetSplitBillsParams) ([]models.SplitBill, int, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, 0, fmt.Errorf("GetSplitBills: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("split_bills")

	findOptions := options.Find()
	findOptions.SetSkip(int64((params.PageNumber - 1) * params.PageSize))
	findOptions.SetLimit(int64(params.PageSize))

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("GetSplitBills: %w", err)
	}
	defer cursor.Close(ctx)

	bills := make([]models.SplitBill, 0)
	for cursor.Next(ctx) {
		var bill models.SplitBill
		if err := cursor.Decode(&bill); err != nil {
			return nil, 0, fmt.Errorf("GetSplitBills: %w", err)
		}
		bills = append(bills, bill)
	}

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("GetSplitBills: %w", err)
	}

	return bills, int(count), nil
}
