package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	ts_models "github.com/nutrixpos/pos/modules/tableside/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type TablesideService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *TablesideService) getSessionCollection() (*mongo.Collection, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}
	return client.Database(s.Config.Databases[0].Database).Collection("table_sessions"), nil
}

func (s *TablesideService) getOrderCollection() (*mongo.Collection, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}
	return client.Database(s.Config.Databases[0].Database).Collection("table_orders"), nil
}

func (s *TablesideService) GetAllSessions() ([]ts_models.TableSession, error) {
	coll, err := s.getSessionCollection()
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

	sessions := make([]ts_models.TableSession, 0)
	for cursor.Next(ctx) {
		var sess ts_models.TableSession
		if err := cursor.Decode(&sess); err != nil {
			continue
		}
		sessions = append(sessions, sess)
	}

	return sessions, nil
}

func (s *TablesideService) CreateSession(session *ts_models.TableSession) error {
	coll, err := s.getSessionCollection()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	session.Id = bson.NewObjectID().Hex()
	session.QrToken = bson.NewObjectID().Hex()
	session.Active = true
	session.OpenedAt = time.Now()

	_, err = coll.InsertOne(ctx, session)
	return err
}

func (s *TablesideService) CloseSession(id string) error {
	coll, err := s.getSessionCollection()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	_, err = coll.UpdateOne(ctx, bson.M{"id": id, "active": true}, bson.M{
		"$set": bson.M{"active": false, "closed_at": now},
	})
	return err
}

func (s *TablesideService) GetSessionByToken(token string) (*ts_models.TableSession, error) {
	coll, err := s.getSessionCollection()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var sess ts_models.TableSession
	err = coll.FindOne(ctx, bson.M{"qr_token": token, "active": true}).Decode(&sess)
	if err != nil {
		return nil, err
	}

	return &sess, nil
}

func (s *TablesideService) PlaceTableOrder(order *ts_models.TableOrder) error {
	coll, err := s.getOrderCollection()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	order.Id = bson.NewObjectID().Hex()
	order.Status = "pending"
	order.PlacedAt = time.Now()

	var subtotal float64
	items := make([]bson.M, len(order.Items))
	for i, item := range order.Items {
		subtotal += item.UnitPrice * float64(item.Quantity)
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
	order.Subtotal = subtotal

	_, err = coll.InsertOne(ctx, order)
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
		"display_id":   fmt.Sprintf("TS-%s", order.Id[:6]),
		"items":        items,
		"state":        "pending",
		"sale_price":   subtotal,
		"cost":         0,
		"submitted_at": time.Now(),
		"is_paid":      false,
		"is_dine_in":   true,
		"custom_data": bson.M{
			"source":         "tableside",
			"session_id":     order.SessionId,
			"table_order_id": order.Id,
		},
	}
	_, err = client.Database(s.Config.Databases[0].Database).Collection("orders").InsertOne(ctx, coreOrder)
	return err
}

func (s *TablesideService) GetOrdersBySession(sessionId string) ([]ts_models.TableOrder, error) {
	coll, err := s.getOrderCollection()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := coll.Find(ctx, bson.M{"session_id": sessionId}, options.Find().SetSort(bson.D{{Key: "placed_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	orders := make([]ts_models.TableOrder, 0)
	for cursor.Next(ctx) {
		var o ts_models.TableOrder
		if err := cursor.Decode(&o); err != nil {
			continue
		}
		orders = append(orders, o)
	}

	return orders, nil
}

func (s *TablesideService) UpdateOrderStatus(orderId, status string) error {
	coll, err := s.getOrderCollection()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = coll.UpdateOne(ctx, bson.M{"id": orderId}, bson.M{"$set": bson.M{"status": status}})
	return err
}

func (s *TablesideService) GetQrUrl(sessionId, host string) (*ts_models.QrInfo, error) {
	sessions, err := s.GetAllSessions()
	if err != nil {
		return nil, err
	}

	for _, sess := range sessions {
		if sess.Id == sessionId {
			return &ts_models.QrInfo{
				TableLabel: sess.TableLabel,
				Token:      sess.QrToken,
				Url:        fmt.Sprintf("%s/tableside/menu/%s", host, sess.QrToken),
				Host:       host,
			}, nil
		}
	}

	return nil, fmt.Errorf("session not found")
}
