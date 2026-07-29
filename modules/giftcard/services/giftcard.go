package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/giftcard/dto"
	"github.com/nutrixpos/pos/modules/giftcard/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type GiftCardService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *GiftCardService) GetDatabase() (*mongo.Database, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}
	return client.Database(s.Config.Databases[0].Database), nil
}

func (s *GiftCardService) CreateGiftCard(req dto.CreateGiftCardRequest) (*models.GiftCard, error) {
	db, err := s.GetDatabase()
	if err != nil {
		return nil, fmt.Errorf("CreateGiftCard: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	card := models.GiftCard{
		Code:          req.Code,
		InitialAmount: req.InitialAmount,
		CurrentAmount: req.InitialAmount,
		Status:        "active",
		CustomerID:    req.CustomerID,
		CustomerName:  req.CustomerName,
		IssuedAt:      now,
		ExpiryDate:    req.ExpiryDate,
		Notes:         req.Notes,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	result, err := db.Collection("gift_cards").InsertOne(ctx, card)
	if err != nil {
		return nil, fmt.Errorf("CreateGiftCard InsertOne: %w", err)
	}

	card.ID = result.InsertedID.(bson.ObjectID)

	tx := models.GiftCardTransaction{
		GiftCardID:   card.ID,
		GiftCardCode: card.Code,
		Type:         "issue",
		Amount:       card.InitialAmount,
		BalanceAfter: card.CurrentAmount,
		Notes:        "Gift card issued",
		CreatedAt:    now,
	}
	db.Collection("gift_card_transactions").InsertOne(ctx, tx)

	return &card, nil
}

func (s *GiftCardService) GetAllGiftCards() ([]models.GiftCard, error) {
	db, err := s.GetDatabase()
	if err != nil {
		return nil, fmt.Errorf("GetAllGiftCards: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.M{"created_at": -1})
	cursor, err := db.Collection("gift_cards").Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("GetAllGiftCards Find: %w", err)
	}
	defer cursor.Close(ctx)

	var cards []models.GiftCard
	if err := cursor.All(ctx, &cards); err != nil {
		return nil, fmt.Errorf("GetAllGiftCards decode: %w", err)
	}

	return cards, nil
}

func (s *GiftCardService) GetGiftCardByCode(code string) (*models.GiftCard, error) {
	db, err := s.GetDatabase()
	if err != nil {
		return nil, fmt.Errorf("GetGiftCardByCode: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var card models.GiftCard
	err = db.Collection("gift_cards").FindOne(ctx, bson.M{"code": code}).Decode(&card)
	if err != nil {
		return nil, fmt.Errorf("GetGiftCardByCode: %w", err)
	}

	return &card, nil
}

func (s *GiftCardService) RedeemGiftCard(req dto.RedeemGiftCardRequest) (*models.GiftCard, error) {
	db, err := s.GetDatabase()
	if err != nil {
		return nil, fmt.Errorf("RedeemGiftCard: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var card models.GiftCard
	err = db.Collection("gift_cards").FindOne(ctx, bson.M{"code": req.Code, "status": "active"}).Decode(&card)
	if err != nil {
		return nil, fmt.Errorf("RedeemGiftCard FindOne: %w", err)
	}

	if card.CurrentAmount < req.Amount {
		return nil, fmt.Errorf("insufficient balance: have %.2f, need %.2f", card.CurrentAmount, req.Amount)
	}

	if card.ExpiryDate != nil && time.Now().After(*card.ExpiryDate) {
		return nil, fmt.Errorf("gift card expired")
	}

	newBalance := card.CurrentAmount - req.Amount
	now := time.Now()

	update := bson.M{
		"$set": bson.M{
			"current_amount": newBalance,
			"updated_at":     now,
		},
	}

	if newBalance == 0 {
		update["$set"].(bson.M)["status"] = "redeemed"
	}

	_, err = db.Collection("gift_cards").UpdateOne(ctx, bson.M{"_id": card.ID}, update)
	if err != nil {
		return nil, fmt.Errorf("RedeemGiftCard UpdateOne: %w", err)
	}

	tx := models.GiftCardTransaction{
		GiftCardID:   card.ID,
		GiftCardCode: card.Code,
		Type:         "redeem",
		Amount:       req.Amount,
		BalanceAfter: newBalance,
		OrderID:      req.OrderID,
		Notes:        req.Notes,
		CreatedAt:    now,
	}
	db.Collection("gift_card_transactions").InsertOne(ctx, tx)

	card.CurrentAmount = newBalance
	if newBalance == 0 {
		card.Status = "redeemed"
	}
	card.UpdatedAt = now

	return &card, nil
}

func (s *GiftCardService) GetTransactions(giftCardID string) ([]models.GiftCardTransaction, error) {
	db, err := s.GetDatabase()
	if err != nil {
		return nil, fmt.Errorf("GetTransactions: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(giftCardID)
	if err != nil {
		return nil, fmt.Errorf("GetTransactions invalid ID: %w", err)
	}

	opts := options.Find().SetSort(bson.M{"created_at": -1})
	cursor, err := db.Collection("gift_card_transactions").Find(ctx, bson.M{"gift_card_id": objID}, opts)
	if err != nil {
		return nil, fmt.Errorf("GetTransactions Find: %w", err)
	}
	defer cursor.Close(ctx)

	var txs []models.GiftCardTransaction
	if err := cursor.All(ctx, &txs); err != nil {
		return nil, fmt.Errorf("GetTransactions decode: %w", err)
	}

	return txs, nil
}
