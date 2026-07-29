package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	kk_models "github.com/nutrixpos/pos/modules/kiosk/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type KioskService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *KioskService) db(name string) (*mongo.Collection, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}
	return client.Database(s.Config.Databases[0].Database).Collection(name), nil
}

func (s *KioskService) GetConfigs() ([]kk_models.KioskConfig, error) {
	c, _ := s.db("kiosk_configs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, _ := c.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	var res []kk_models.KioskConfig
	for cursor.Next(ctx) {
		var k kk_models.KioskConfig
		if cursor.Decode(&k) == nil {
			res = append(res, k)
		}
	}
	return res, nil
}

func (s *KioskService) SaveConfig(k *kk_models.KioskConfig) error {
	c, _ := s.db("kiosk_configs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if k.Id == "" {
		k.Id = bson.NewObjectID().Hex()
	}
	_, err := c.ReplaceOne(ctx, bson.M{"id": k.Id}, k)
	return err
}

func (s *KioskService) PlaceOrder(o *kk_models.KioskOrder) error {
	c, _ := s.db("kiosk_orders")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	o.Id = bson.NewObjectID().Hex()
	o.Status = "pending"
	o.CreatedAt = time.Now().Format("2006-01-02 15:04")
	var total float64
	items := make([]bson.M, len(o.Items))
	for i, item := range o.Items {
		total += item.UnitPrice * float64(item.Quantity)
		items[i] = bson.M{
			"id": bson.NewObjectID().Hex(),
			"product": bson.M{
				"id":    item.ProductId,
				"name":  item.ProductName,
				"price": item.UnitPrice,
			},
			"price":      item.UnitPrice,
			"quantity":   item.Quantity,
			"sale_price": item.UnitPrice * float64(item.Quantity),
			"comment":    item.Notes,
		}
	}
	o.Total = total
	_, err := c.InsertOne(ctx, o)
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
		"display_id":   fmt.Sprintf("KS-%s", o.Id[:6]),
		"items":        items,
		"state":        "pending",
		"sale_price":   total,
		"cost":         0,
		"submitted_at": time.Now(),
		"is_paid":      false,
		"custom_data": bson.M{
			"source":   "kiosk",
			"kiosk_id": o.KioskId,
		},
	}
	_, err = client.Database(s.Config.Databases[0].Database).Collection("orders").InsertOne(ctx, coreOrder)
	return err
}

func (s *KioskService) GetOrders() ([]kk_models.KioskOrder, error) {
	c, _ := s.db("kiosk_orders")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, _ := c.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	var res []kk_models.KioskOrder
	for cursor.Next(ctx) {
		var o kk_models.KioskOrder
		if cursor.Decode(&o) == nil {
			res = append(res, o)
		}
	}
	return res, nil
}
