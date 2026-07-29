package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	dv_models "github.com/nutrixpos/pos/modules/delivery/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DeliveryService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *DeliveryService) db(name string) (*mongo.Collection, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}
	return client.Database(s.Config.Databases[0].Database).Collection(name), nil
}

func (s *DeliveryService) GetZones() ([]dv_models.DeliveryZone, error) {
	c, err := s.db("delivery_zones")
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, _ := c.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	var res []dv_models.DeliveryZone
	for cursor.Next(ctx) {
		var z dv_models.DeliveryZone
		if cursor.Decode(&z) == nil {
			res = append(res, z)
		}
	}
	return res, nil
}

func (s *DeliveryService) SaveZone(z *dv_models.DeliveryZone) error {
	c, err := s.db("delivery_zones")
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if z.Id == "" {
		z.Id = bson.NewObjectID().Hex()
	}
	_, err = c.ReplaceOne(ctx, bson.M{"id": z.Id}, z)
	return err
}

func (s *DeliveryService) DeleteZone(id string) error {
	c, err := s.db("delivery_zones")
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = c.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (s *DeliveryService) GetOrders() ([]dv_models.DeliveryOrder, error) {
	c, err := s.db("delivery_orders")
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, _ := c.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	var res []dv_models.DeliveryOrder
	for cursor.Next(ctx) {
		var o dv_models.DeliveryOrder
		if cursor.Decode(&o) == nil {
			res = append(res, o)
		}
	}
	return res, nil
}

func (s *DeliveryService) CreateOrder(o *dv_models.DeliveryOrder) error {
	c, err := s.db("delivery_orders")
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	o.Id = bson.NewObjectID().Hex()
	o.Status = "pending"
	o.PlacedAt = time.Now()
	_, err = c.InsertOne(ctx, o)
	if err != nil {
		return err
	}

	// also create a corresponding core order document
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}
	coreOrder := bson.M{
		"id":           bson.NewObjectID().Hex(),
		"display_id":   fmt.Sprintf("DV-%s", o.Id[:6]),
		"state":        "pending",
		"sale_price":   o.DeliveryFee,
		"cost":         0,
		"submitted_at": time.Now(),
		"is_paid":      false,
		"is_delivery":  true,
		"delivery_info": bson.M{
			"receiver_name": o.CustomerName,
			"address":       o.Address,
			"phone":         o.CustomerPhone,
		},
		"custom_data": bson.M{
			"source":            "delivery",
			"delivery_order_id": o.Id,
			"zone_id":           o.ZoneId,
		},
	}
	_, err = client.Database(s.Config.Databases[0].Database).Collection("orders").InsertOne(ctx, coreOrder)
	return err
}

func (s *DeliveryService) UpdateOrderStatus(id, status string) error {
	c, err := s.db("delivery_orders")
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	set := bson.M{"status": status}
	if status == "delivered" {
		now := time.Now()
		set["delivered_at"] = now
	}
	_, err = c.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": set})
	return err
}
