package services

import (
	"context"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/core/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CustomersService struct {
	Logger   logger.ILogger
	Config   config.Config
	Settings models.Settings
}

type GetCustomersParams struct {
	// page_number is used in pagination to set the index of the first record to be returned.
	PageNumber int
	// Rows is to set the desired row count limit.
	PageSize int
}

func (cs CustomersService) GetCustomers(params GetCustomersParams) (customers []models.Customer, customers_count int, err error) {

	client, err := common.GetDatabaseClient(cs.Logger, &cs.Config)
	if err != nil {
		return customers, customers_count, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	customers = make([]models.Customer, 0)

	collection := client.Database(cs.Config.Databases[0].Database).Collection("customers")

	findOptions := options.Find()
	findOptions.SetSkip(int64((params.PageNumber - 1) * params.PageSize))
	findOptions.SetLimit(int64(params.PageSize))

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return customers, customers_count, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(context.Background()) {
		var customer models.Customer
		if err := cursor.Decode(&customer); err != nil {
			return customers, customers_count, err
		}

		customers = append(customers, customer)
	}

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return customers, customers_count, err
	}
	customers_count = int(count)

	return customers, customers_count, err
}

func (cs CustomersService) GetCustomer(customer_id string) (customer models.Customer, err error) {

	client, err := common.GetDatabaseClient(cs.Logger, &cs.Config)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(cs.Config.Databases[0].Database).Collection("customers")

	err = collection.FindOne(ctx, bson.M{"id": customer_id}).Decode(&customer)
	if err != nil {
		return customer, err
	}

	return
}

func (cs CustomersService) InsertNew(customer models.Customer) (afterInsert models.Customer, err error) {

	client, err := common.GetDatabaseClient(cs.Logger, &cs.Config)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(cs.Config.Databases[0].Database).Collection("customers")

	customer.Id = bson.NewObjectID().Hex()
	customer.CreatedAt = time.Now()

	result, err := collection.InsertOne(ctx, customer)
	if err != nil {
		return afterInsert, err
	}

	err = collection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&afterInsert)
	if err != nil {
		return afterInsert, err
	}

	return
}

func (cs CustomersService) UpdateCustomer(customer models.Customer, customer_id string) (afterUpdate models.Customer, err error) {

	client, err := common.GetDatabaseClient(cs.Logger, &cs.Config)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(cs.Config.Databases[0].Database).Collection("customers")

	update := bson.M{}
	if customer.Name != "" {
		update["name"] = customer.Name
	}
	if customer.Phone != "" {
		update["phone"] = customer.Phone
	}
	if customer.Address != "" {
		update["address"] = customer.Address
	}
	if customer.Email != "" {
		update["email"] = customer.Email
	}
	if customer.Notes != "" {
		update["notes"] = customer.Notes
	}
	if customer.Tags != nil {
		update["tags"] = customer.Tags
	}
	if customer.Preferences != nil {
		update["preferences"] = customer.Preferences
	}

	_, err = collection.UpdateOne(ctx, bson.M{"id": customer_id}, bson.M{"$set": update})
	if err != nil {
		return afterUpdate, err
	}

	err = collection.FindOne(ctx, bson.M{"id": customer_id}).Decode(&afterUpdate)
	if err != nil {
		return afterUpdate, err
	}

	return
}

func (cs CustomersService) DeleteCustomer(customer_id string) (err error) {

	client, err := common.GetDatabaseClient(cs.Logger, &cs.Config)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(cs.Config.Databases[0].Database).Collection("customers")

	_, err = collection.DeleteOne(ctx, bson.M{"id": customer_id})
	if err != nil {
		return err
	}

	return
}

func (cs CustomersService) GetCustomerOrders(customer_id string, limit int) (orders []models.Order, err error) {

	client, err := common.GetDatabaseClient(cs.Logger, &cs.Config)
	if err != nil {
		return orders, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orders = make([]models.Order, 0)

	collection := client.Database(cs.Config.Databases[0].Database).Collection("orders")

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.D{{Key: "submitted_at", Value: -1}})

	cursor, err := collection.Find(ctx, bson.M{"customer.id": customer_id}, findOptions)
	if err != nil {
		return orders, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(context.Background()) {
		var order models.Order
		if err := cursor.Decode(&order); err != nil {
			return orders, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (cs CustomersService) UpdateCustomerStats(customer_id string) (err error) {

	client, err := common.GetDatabaseClient(cs.Logger, &cs.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ordersCollection := client.Database(cs.Config.Databases[0].Database).Collection("orders")
	customersCollection := client.Database(cs.Config.Databases[0].Database).Collection("customers")

	var totalSpent float64
	var orderCount int
	var lastOrderDate *time.Time

	cursor, err := ordersCollection.Find(ctx, bson.M{
		"customer.id": customer_id,
		"is_paid":     true,
	})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(context.Background()) {
		var order models.Order
		if err := cursor.Decode(&order); err != nil {
			return err
		}
		totalSpent += order.SalePrice - order.Discount
		orderCount++
		if lastOrderDate == nil || order.SubmittedAt.After(*lastOrderDate) {
			lastOrderDate = &order.SubmittedAt
		}
	}

	update := bson.M{
		"total_spent":     totalSpent,
		"order_count":     orderCount,
		"loyalty_points":  int(totalSpent / 10),
	}
	if lastOrderDate != nil {
		update["last_order_date"] = lastOrderDate
	}

	_, err = customersCollection.UpdateOne(ctx, bson.M{"id": customer_id}, bson.M{"$set": update})
	return err
}
