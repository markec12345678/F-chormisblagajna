package services

import (
	"context"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	rs_models "github.com/nutrixpos/pos/modules/reservations/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ReservationService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *ReservationService) coll() (*mongo.Collection, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}
	return client.Database(s.Config.Databases[0].Database).Collection("reservations"), nil
}

func (s *ReservationService) Create(r *rs_models.Reservation) error {
	c, err := s.coll()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r.Id = bson.NewObjectID().Hex()
	r.Status = "pending"
	r.CreatedAt = time.Now()
	_, err = c.InsertOne(ctx, r)
	return err
}

func (s *ReservationService) GetAll() ([]rs_models.Reservation, error) {
	c, err := s.coll()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := c.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "reservation_date", Value: 1}, {Key: "reservation_time", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	res := make([]rs_models.Reservation, 0)
	for cursor.Next(ctx) {
		var r rs_models.Reservation
		if err := cursor.Decode(&r); err != nil {
			continue
		}
		res = append(res, r)
	}
	return res, nil
}

func (s *ReservationService) UpdateStatus(id, status string) error {
	c, err := s.coll()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = c.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.M{"status": status}})
	return err
}

func (s *ReservationService) AssignTable(id, table string) error {
	c, err := s.coll()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = c.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.M{"table_assignment": table, "status": "confirmed"}})
	return err
}

func (s *ReservationService) Delete(id string) error {
	c, err := s.coll()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = c.DeleteOne(ctx, bson.M{"id": id})
	return err
}
