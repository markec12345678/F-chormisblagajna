package services

import (
	"context"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	pc_models "github.com/nutrixpos/pos/modules/purchase/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PurchaseService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *PurchaseService) coll() (*mongo.Collection, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}
	return client.Database(s.Config.Databases[0].Database).Collection("purchase_orders"), nil
}

func (s *PurchaseService) Create(po *pc_models.PurchaseOrder) error {
	c, err := s.coll()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var total float64
	for i, item := range po.Items {
		po.Items[i].TotalPrice = item.Quantity * item.UnitPrice
		total += po.Items[i].TotalPrice
	}
	po.TotalCost = total
	po.Id = bson.NewObjectID().Hex()
	po.Status = "pending"
	po.OrderedAt = time.Now()

	_, err = c.InsertOne(ctx, po)
	return err
}

func (s *PurchaseService) GetAll() ([]pc_models.PurchaseOrder, error) {
	c, err := s.coll()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, _ := c.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	var res []pc_models.PurchaseOrder
	for cursor.Next(ctx) {
		var po pc_models.PurchaseOrder
		if cursor.Decode(&po) == nil {
			res = append(res, po)
		}
	}
	return res, nil
}

func (s *PurchaseService) MarkReceived(id string) error {
	c, err := s.coll()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	now := time.Now()
	_, err = c.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.M{"status": "received", "received_at": now}})
	return err
}

func (s *PurchaseService) Cancel(id string) error {
	c, err := s.coll()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = c.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.M{"status": "cancelled"}})
	return err
}
