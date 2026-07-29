package services

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/core/models"
	onlineorder_models "github.com/nutrixpos/pos/modules/onlineorder/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type OnlineOrderService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *OnlineOrderService) GetMenu() ([]onlineorder_models.MenuCategory, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	categoriesCol := client.Database(s.Config.Databases[0].Database).Collection("categories")
	productsCol := client.Database(s.Config.Databases[0].Database).Collection("recipes")

	catCursor, err := categoriesCol.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer catCursor.Close(ctx)

	var categories []models.Category
	if err := catCursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	menuCategories := make([]onlineorder_models.MenuCategory, 0, len(categories))

	for _, cat := range categories {
		products := make([]onlineorder_models.MenuProduct, 0)

		for _, catProd := range cat.Products {
			var product models.Product
			filter := bson.D{{Key: "id", Value: catProd.Id}}
			err := productsCol.FindOne(ctx, filter).Decode(&product)
			if err != nil {
				continue
			}
			if product.Name == "" {
				continue
			}

			products = append(products, onlineorder_models.MenuProduct{
				Id:        product.Id,
				Name:      product.Name,
				Price:     product.Price,
				ImageURL:  product.ImageURL,
				Category:  cat.Name,
				Available: product.Quantity > 0 || product.Ready > 0,
			})
		}

		if len(products) > 0 {
			menuCategories = append(menuCategories, onlineorder_models.MenuCategory{
				Id:       cat.Id,
				Name:     cat.Name,
				Products: products,
			})
		}
	}

	return menuCategories, nil
}

func (s *OnlineOrderService) CreateOrder(order *onlineorder_models.OnlineOrder) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("online_orders")

	order.Id = bson.NewObjectID().Hex()
	order.DisplayId = fmt.Sprintf("WO-%d", time.Now().UnixMilli()%10000)
	order.Status = "pending"
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	// Calculate subtotal
	var subtotal float64
	for _, item := range order.Items {
		subtotal += item.Price * float64(item.Quantity)
	}
	order.Subtotal = subtotal
	order.Total = subtotal + order.DeliveryFee

	_, err = collection.InsertOne(ctx, order)
	return err
}

func (s *OnlineOrderService) GetOrder(orderId string) (*onlineorder_models.OnlineOrder, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("online_orders")

	var order onlineorder_models.OnlineOrder
	err = collection.FindOne(ctx, bson.M{"id": orderId}).Decode(&order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *OnlineOrderService) TrackOrder(displayId string) (*onlineorder_models.OnlineOrder, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("online_orders")

	var order onlineorder_models.OnlineOrder
	err = collection.FindOne(ctx, bson.M{
		"display_id": bson.M{"$regex": fmt.Sprintf("(?i)%s", regexp.QuoteMeta(displayId))},
	}).Decode(&order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *OnlineOrderService) GetAllOrders(status string) ([]onlineorder_models.OnlineOrder, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("online_orders")

	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	orders := make([]onlineorder_models.OnlineOrder, 0)
	for cursor.Next(context.Background()) {
		var order onlineorder_models.OnlineOrder
		if err := cursor.Decode(&order); err != nil {
			continue
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (s *OnlineOrderService) UpdateOrderStatus(orderId, status string) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("online_orders")

	_, err = collection.UpdateOne(ctx,
		bson.M{"id": orderId},
		bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}},
	)
	return err
}
