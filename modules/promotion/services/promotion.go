package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/promotion/dto"
	"github.com/nutrixpos/pos/modules/promotion/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type PromotionService struct {
	Logger logger.ILogger
	Config config.Config
}

type GetPromotionsParams struct {
	PageNumber int
	PageSize   int
	ActiveOnly bool
}

func (s *PromotionService) GetPromotions(params GetPromotionsParams) ([]models.Promotion, int, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, 0, fmt.Errorf("GetPromotions: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("promotions")

	filter := bson.M{}
	if params.ActiveOnly {
		filter["is_active"] = true
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64((params.PageNumber - 1) * params.PageSize))
	findOptions.SetLimit(int64(params.PageSize))
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("GetPromotions: %w", err)
	}
	defer cursor.Close(ctx)

	promotions := make([]models.Promotion, 0)
	for cursor.Next(ctx) {
		var promotion models.Promotion
		if err := cursor.Decode(&promotion); err != nil {
			return nil, 0, fmt.Errorf("GetPromotions: %w", err)
		}
		promotions = append(promotions, promotion)
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("GetPromotions: %w", err)
	}

	return promotions, int(count), nil
}

func (s *PromotionService) ValidatePromotionCode(code string, orderTotal float64) (*models.Promotion, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, fmt.Errorf("ValidatePromotionCode: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("promotions")

	var promotion models.Promotion
	err = collection.FindOne(ctx, bson.M{"code": code, "is_active": true}).Decode(&promotion)
	if err != nil {
		return nil, fmt.Errorf("promotion not found or inactive")
	}

	now := time.Now()
	today := now.Format("2006-01-02")

	if promotion.StartDate != "" && today < promotion.StartDate {
		return nil, fmt.Errorf("promotion not yet active")
	}
	if promotion.EndDate != "" && today > promotion.EndDate {
		return nil, fmt.Errorf("promotion has expired")
	}
	if promotion.UsageLimit > 0 && promotion.UsageCount >= promotion.UsageLimit {
		return nil, fmt.Errorf("promotion usage limit reached")
	}
	if promotion.MinOrder > 0 && orderTotal < promotion.MinOrder {
		return nil, fmt.Errorf("order total below minimum %.2f", promotion.MinOrder)
	}

	return &promotion, nil
}

func (s *PromotionService) CreatePromotion(req dto.CreatePromotionRequest) (models.Promotion, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return models.Promotion{}, fmt.Errorf("CreatePromotion: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("promotions")

	now := time.Now()
	promotion := models.Promotion{
		Id:             bson.NewObjectID().Hex(),
		Name:           req.Name,
		Code:           req.Code,
		Type:           req.Type,
		Value:          req.Value,
		MinOrder:       req.MinOrder,
		MaxDiscount:    req.MaxDiscount,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		UsageLimit:     req.UsageLimit,
		ApplicableDays: req.ApplicableDays,
		HappyHourStart: req.HappyHourStart,
		HappyHourEnd:   req.HappyHourEnd,
		IsActive:       req.IsActive,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	_, err = collection.InsertOne(ctx, promotion)
	if err != nil {
		return models.Promotion{}, fmt.Errorf("CreatePromotion: %w", err)
	}

	return promotion, nil
}

func (s *PromotionService) UpdatePromotion(id string, req dto.UpdatePromotionRequest) (models.Promotion, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return models.Promotion{}, fmt.Errorf("UpdatePromotion: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("promotions")

	update := bson.M{"updated_at": time.Now()}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Code != "" {
		update["code"] = req.Code
	}
	if req.Type != "" {
		update["type"] = req.Type
	}
	if req.Value != nil {
		update["value"] = *req.Value
	}
	if req.MinOrder != nil {
		update["min_order"] = *req.MinOrder
	}
	if req.MaxDiscount != nil {
		update["max_discount"] = *req.MaxDiscount
	}
	if req.StartDate != "" {
		update["start_date"] = req.StartDate
	}
	if req.EndDate != "" {
		update["end_date"] = req.EndDate
	}
	if req.UsageLimit != nil {
		update["usage_limit"] = *req.UsageLimit
	}
	if req.ApplicableDays != nil {
		update["applicable_days"] = req.ApplicableDays
	}
	if req.HappyHourStart != "" {
		update["happy_hour_start"] = req.HappyHourStart
	}
	if req.HappyHourEnd != "" {
		update["happy_hour_end"] = req.HappyHourEnd
	}
	if req.IsActive != nil {
		update["is_active"] = *req.IsActive
	}

	_, err = collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": update})
	if err != nil {
		return models.Promotion{}, fmt.Errorf("UpdatePromotion: %w", err)
	}

	var updatedPromotion models.Promotion
	err = collection.FindOne(ctx, bson.M{"id": id}).Decode(&updatedPromotion)
	if err != nil {
		return models.Promotion{}, fmt.Errorf("UpdatePromotion: %w", err)
	}

	return updatedPromotion, nil
}

func (s *PromotionService) DeletePromotion(id string) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return fmt.Errorf("DeletePromotion: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("promotions")

	_, err = collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return fmt.Errorf("DeletePromotion: %w", err)
	}

	return nil
}
