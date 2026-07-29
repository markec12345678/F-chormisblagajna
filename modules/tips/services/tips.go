package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/tips/dto"
	"github.com/nutrixpos/pos/modules/tips/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type TipsService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *TipsService) GetDatabase() (*mongo.Database, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}
	return client.Database(s.Config.Databases[0].Database), nil
}

func (s *TipsService) RecordTip(req dto.RecordTipRequest) (*models.Tip, error) {
	db, err := s.GetDatabase()
	if err != nil {
		return nil, fmt.Errorf("RecordTip: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	tip := models.Tip{
		OrderID:       req.OrderID,
		EmployeeID:    req.EmployeeID,
		EmployeeName:  req.EmployeeName,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		BranchID:      req.BranchID,
		Date:          now,
		CreatedAt:     now,
	}

	result, err := db.Collection("tips").InsertOne(ctx, tip)
	if err != nil {
		return nil, fmt.Errorf("RecordTip InsertOne: %w", err)
	}

	tip.ID = result.InsertedID.(bson.ObjectID)
	return &tip, nil
}

func (s *TipsService) GetTipsByEmployee(startDate, endDate string) ([]models.TipSummary, error) {
	db, err := s.GetDatabase()
	if err != nil {
		return nil, fmt.Errorf("GetTipsByEmployee: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}
	if startDate != "" || endDate != "" {
		dateFilter := bson.M{}
		if startDate != "" {
			t, err := time.Parse("2006-01-02", startDate)
			if err == nil {
				dateFilter["$gte"] = t
			}
		}
		if endDate != "" {
			t, err := time.Parse("2006-01-02", endDate)
			if err == nil {
				t = t.Add(24 * time.Hour)
				dateFilter["$lte"] = t
			}
		}
		filter["date"] = dateFilter
	}

	pipeline := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id":          "$employee_id",
			"employee_name": bson.M{"$first": "$employee_name"},
			"total_tips":   bson.M{"$sum": "$amount"},
			"tip_count":    bson.M{"$sum": 1},
			"average_tip":  bson.M{"$avg": "$amount"},
		}},
		{"$sort": bson.M{"total_tips": -1}},
	}

	cursor, err := db.Collection("tips").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("GetTipsByEmployee Aggregate: %w", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("GetTipsByEmployee decode: %w", err)
	}

	var summaries []models.TipSummary
	for _, r := range results {
		empID, _ := r["_id"].(string)
		empName, _ := r["employee_name"].(string)
		total, _ := r["total_tips"].(float64)
		count, _ := r["tip_count"].(int32)
		avg, _ := r["average_tip"].(float64)

		summaries = append(summaries, models.TipSummary{
			EmployeeID:   empID,
			EmployeeName: empName,
			TotalTips:    total,
			TipCount:     int(count),
			AverageTip:   avg,
		})
	}

	return summaries, nil
}

func (s *TipsService) PayoutTips(req dto.PayoutTipsRequest) (*models.TipPayout, error) {
	db, err := s.GetDatabase()
	if err != nil {
		return nil, fmt.Errorf("PayoutTips: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	payout := models.TipPayout{
		EmployeeID:   req.EmployeeID,
		EmployeeName: req.EmployeeName,
		Amount:       req.Amount,
		PayoutDate:   now,
		PayoutMethod: req.PayoutMethod,
		Notes:        req.Notes,
		CreatedAt:    now,
	}

	result, err := db.Collection("tip_payouts").InsertOne(ctx, payout)
	if err != nil {
		return nil, fmt.Errorf("PayoutTips InsertOne: %w", err)
	}

	payout.ID = result.InsertedID.(bson.ObjectID)
	return &payout, nil
}

func (s *TipsService) GetPayouts() ([]models.TipPayout, error) {
	db, err := s.GetDatabase()
	if err != nil {
		return nil, fmt.Errorf("GetPayouts: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.M{"created_at": -1}).SetLimit(100)
	cursor, err := db.Collection("tip_payouts").Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("GetPayouts Find: %w", err)
	}
	defer cursor.Close(ctx)

	var payouts []models.TipPayout
	if err := cursor.All(ctx, &payouts); err != nil {
		return nil, fmt.Errorf("GetPayouts decode: %w", err)
	}

	return payouts, nil
}
