package services

import (
	"context"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	fp_models "github.com/nutrixpos/pos/modules/floorplan/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type FloorplanService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *FloorplanService) db(name string) (*mongo.Collection, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}
	return client.Database(s.Config.Databases[0].Database).Collection(name), nil
}

func (s *FloorplanService) GetTables() ([]fp_models.FloorTable, error) {
	c, _ := s.db("floor_tables")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, _ := c.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	var res []fp_models.FloorTable
	for cursor.Next(ctx) {
		var t fp_models.FloorTable
		if cursor.Decode(&t) == nil {
			res = append(res, t)
		}
	}
	return res, nil
}

func (s *FloorplanService) SaveTable(t *fp_models.FloorTable) error {
	c, _ := s.db("floor_tables")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if t.Id == "" {
		t.Id = bson.NewObjectID().Hex()
	}
	_, err := c.ReplaceOne(ctx, bson.M{"id": t.Id}, t)
	return err
}

func (s *FloorplanService) DeleteTable(id string) error {
	c, _ := s.db("floor_tables")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := c.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (s *FloorplanService) GetZones() ([]fp_models.FloorZone, error) {
	c, _ := s.db("floor_zones")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, _ := c.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	var res []fp_models.FloorZone
	for cursor.Next(ctx) {
		var z fp_models.FloorZone
		if cursor.Decode(&z) == nil {
			res = append(res, z)
		}
	}
	return res, nil
}

func (s *FloorplanService) SaveZone(z *fp_models.FloorZone) error {
	c, _ := s.db("floor_zones")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if z.Id == "" {
		z.Id = bson.NewObjectID().Hex()
	}
	_, err := c.ReplaceOne(ctx, bson.M{"id": z.Id}, z)
	return err
}
