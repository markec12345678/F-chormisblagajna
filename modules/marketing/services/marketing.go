package services

import (
	"context"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	mk_models "github.com/nutrixpos/pos/modules/marketing/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MarketingService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *MarketingService) coll() (*mongo.Collection, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}
	return client.Database(s.Config.Databases[0].Database).Collection("campaigns"), nil
}

func (s *MarketingService) Create(c *mk_models.Campaign) error {
	coll, err := s.coll()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c.Id = bson.NewObjectID().Hex()
	c.CreatedAt = time.Now().Format("2006-01-02 15:04")
	_, err = coll.InsertOne(ctx, c)
	return err
}

func (s *MarketingService) GetAll() ([]mk_models.Campaign, error) {
	coll, _ := s.coll()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, _ := coll.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	var res []mk_models.Campaign
	for cursor.Next(ctx) {
		var c mk_models.Campaign
		if cursor.Decode(&c) == nil {
			res = append(res, c)
		}
	}
	return res, nil
}

func (s *MarketingService) ToggleActive(id string) error {
	coll, _ := s.coll()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var c mk_models.Campaign
	coll.FindOne(ctx, bson.M{"id": id}).Decode(&c)
	_, err := coll.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.M{"active": !c.Active}})
	return err
}

func (s *MarketingService) Delete(id string) error {
	coll, _ := s.coll()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := coll.DeleteOne(ctx, bson.M{"id": id})
	return err
}
