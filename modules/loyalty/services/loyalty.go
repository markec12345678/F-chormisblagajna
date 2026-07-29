package services

import (
	"context"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	ly_models "github.com/nutrixpos/pos/modules/loyalty/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type LoyaltyService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *LoyaltyService) cardCollection() (*mongo.Collection, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}
	return client.Database(s.Config.Databases[0].Database).Collection("loyalty_cards"), nil
}

func (s *LoyaltyService) rewardCollection() (*mongo.Collection, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}
	return client.Database(s.Config.Databases[0].Database).Collection("loyalty_rewards"), nil
}

func (s *LoyaltyService) redemptionCollection() (*mongo.Collection, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}
	return client.Database(s.Config.Databases[0].Database).Collection("loyalty_redemptions"), nil
}

func (s *LoyaltyService) GetDefaultSettings() ly_models.LoyaltySettings {
	return ly_models.LoyaltySettings{
		PointsPerEuro:  10,
		EuroPerPoint:   0.01,
		WelcomePoints:  50,
		TierThresholds: map[string]float64{"bronze": 0, "silver": 500, "gold": 2000, "platinum": 5000},
	}
}

func (s *LoyaltyService) CreateCard(card *ly_models.LoyaltyCard) error {
	coll, err := s.cardCollection()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	card.Id = bson.NewObjectID().Hex()
	card.Points = s.GetDefaultSettings().WelcomePoints
	card.Tier = "bronze"
	card.Active = true

	_, err = coll.InsertOne(ctx, card)
	return err
}

func (s *LoyaltyService) GetAllCards() ([]ly_models.LoyaltyCard, error) {
	coll, err := s.cardCollection()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	cards := make([]ly_models.LoyaltyCard, 0)
	for cursor.Next(ctx) {
		var c ly_models.LoyaltyCard
		if err := cursor.Decode(&c); err != nil {
			continue
		}
		cards = append(cards, c)
	}

	return cards, nil
}

func (s *LoyaltyService) AddPoints(cardId string, amount float64) error {
	coll, err := s.cardCollection()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	settings := s.GetDefaultSettings()
	points := int(amount * settings.PointsPerEuro)

	_, err = coll.UpdateOne(ctx, bson.M{"id": cardId}, bson.M{
		"$inc": bson.M{"points": points, "total_spent": amount, "visit_count": 1},
	})
	return err
}

func (s *LoyaltyService) RedeemPoints(cardId, rewardId string) (*ly_models.Redemption, error) {
	cardColl, err := s.cardCollection()
	if err != nil {
		return nil, err
	}

	rewardColl, err := s.rewardCollection()
	if err != nil {
		return nil, err
	}

	redColl, err := s.redemptionCollection()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var reward ly_models.Reward
	err = rewardColl.FindOne(ctx, bson.M{"id": rewardId, "active": true}).Decode(&reward)
	if err != nil {
		return nil, err
	}

	var card ly_models.LoyaltyCard
	err = cardColl.FindOne(ctx, bson.M{"id": cardId, "active": true}).Decode(&card)
	if err != nil {
		return nil, err
	}

	if card.Points < reward.PointsCost {
		return nil, nil
	}

	redemption := &ly_models.Redemption{
		Id:          bson.NewObjectID().Hex(),
		CardId:      cardId,
		RewardId:    rewardId,
		RewardName:  reward.Name,
		PointsSpent: reward.PointsCost,
		RedeemedAt:  time.Now().Format("2006-01-02 15:04"),
	}

	_, err = cardColl.UpdateOne(ctx, bson.M{"id": cardId}, bson.M{
		"$inc": bson.M{"points": -reward.PointsCost},
	})
	if err != nil {
		return nil, err
	}

	_, err = redColl.InsertOne(ctx, redemption)
	return redemption, err
}

func (s *LoyaltyService) CreateReward(reward *ly_models.Reward) error {
	coll, err := s.rewardCollection()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	reward.Id = bson.NewObjectID().Hex()
	_, err = coll.InsertOne(ctx, reward)
	return err
}

func (s *LoyaltyService) GetAllRewards() ([]ly_models.Reward, error) {
	coll, err := s.rewardCollection()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	rewards := make([]ly_models.Reward, 0)
	for cursor.Next(ctx) {
		var r ly_models.Reward
		if err := cursor.Decode(&r); err != nil {
			continue
		}
		rewards = append(rewards, r)
	}

	return rewards, nil
}
