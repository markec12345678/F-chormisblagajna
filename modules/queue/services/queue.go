package services

import (
	"context"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	q_models "github.com/nutrixpos/pos/modules/queue/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type QueueService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *QueueService) coll() (*mongo.Collection, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}
	return client.Database(s.Config.Databases[0].Database).Collection("queue"), nil
}

func (s *QueueService) Add(entry *q_models.QueueEntry) error {
	c, err := s.coll()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, _ := c.CountDocuments(ctx, bson.M{"status": "waiting"})
	entry.Id = bson.NewObjectID().Hex()
	entry.Position = int(count) + 1
	entry.Status = "waiting"
	entry.EstimatedMin = int(count) * 15
	entry.AddedAt = time.Now().Format("2006-01-02 15:04")

	_, err = c.InsertOne(ctx, entry)
	return err
}

func (s *QueueService) GetAll() ([]q_models.QueueEntry, error) {
	c, _ := s.coll()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, _ := c.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "position", Value: 1}}))
	defer cursor.Close(ctx)
	var res []q_models.QueueEntry
	for cursor.Next(ctx) {
		var e q_models.QueueEntry
		if cursor.Decode(&e) == nil {
			res = append(res, e)
		}
	}
	return res, nil
}

func (s *QueueService) UpdateStatus(id, status string) error {
	c, _ := s.coll()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := c.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.M{"status": status}})
	return err
}

func (s *QueueService) Remove(id string) error {
	c, _ := s.coll()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := c.DeleteOne(ctx, bson.M{"id": id})
	return err
}
