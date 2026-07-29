package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	gc_models "github.com/nutrixpos/pos/modules/giftcards/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type GiftCardService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *GiftCardService) coll() (*mongo.Collection, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}
	return client.Database(s.Config.Databases[0].Database).Collection("gift_cards"), nil
}

func randomCode() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("GC-%08X", b)
}

func (s *GiftCardService) Issue(card *gc_models.GiftCard) error {
	c, err := s.coll()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	card.Id = bson.NewObjectID().Hex()
	card.Code = randomCode()
	card.Balance = card.InitialAmt
	card.IssuedAt = time.Now().Format("2006-01-02")
	card.Active = true
	_, err = c.InsertOne(ctx, card)
	return err
}

func (s *GiftCardService) GetAll() ([]gc_models.GiftCard, error) {
	c, _ := s.coll()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, _ := c.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	var res []gc_models.GiftCard
	for cursor.Next(ctx) {
		var g gc_models.GiftCard
		if cursor.Decode(&g) == nil {
			res = append(res, g)
		}
	}
	return res, nil
}

func (s *GiftCardService) Redeem(code string, amount float64, orderId string) error {
	c, err := s.coll()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r := c.FindOneAndUpdate(ctx, bson.M{"code": code, "active": true, "balance": bson.M{"$gte": amount}},
		bson.M{"$inc": bson.M{"balance": -amount}})
	if r.Err() != nil {
		return r.Err()
	}

	// also record a payment part for multi-payment tracking
	if orderId != "" {
		client, err := common.GetDatabaseClient(s.Logger, &s.Config)
		if err != nil {
			return err
		}
		paymentPart := bson.M{
			"id":             bson.NewObjectID().Hex(),
			"order_id":       orderId,
			"amount":         amount,
			"payment_method": "gift_card",
			"reference":      code,
			"created_at":     time.Now(),
		}
		_, err = client.Database(s.Config.Databases[0].Database).Collection("payment_parts").InsertOne(ctx, paymentPart)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *GiftCardService) Deactivate(id string) error {
	c, _ := s.coll()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := c.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.M{"active": false}})
	return err
}
