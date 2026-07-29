package services

import (
	"context"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	supplier_models "github.com/nutrixpos/pos/modules/supplier/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type SupplierService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *SupplierService) GetAllSuppliers() ([]supplier_models.Supplier, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("suppliers")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	suppliers := make([]supplier_models.Supplier, 0)
	for cursor.Next(context.Background()) {
		var supplier supplier_models.Supplier
		if err := cursor.Decode(&supplier); err != nil {
			continue
		}
		suppliers = append(suppliers, supplier)
	}

	return suppliers, nil
}

func (s *SupplierService) GetSupplier(id string) (*supplier_models.Supplier, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("suppliers")

	var supplier supplier_models.Supplier
	err = collection.FindOne(ctx, bson.M{"id": id}).Decode(&supplier)
	if err != nil {
		return nil, err
	}

	return &supplier, nil
}

func (s *SupplierService) CreateSupplier(supplier *supplier_models.Supplier) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("suppliers")

	supplier.Id = bson.NewObjectID().Hex()
	supplier.CreatedAt = time.Now()
	supplier.IsActive = true

	_, err = collection.InsertOne(ctx, supplier)
	return err
}

func (s *SupplierService) UpdateSupplier(id string, supplier *supplier_models.Supplier) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("suppliers")

	update := bson.M{}
	if supplier.Name != "" {
		update["name"] = supplier.Name
	}
	if supplier.ContactName != "" {
		update["contact_name"] = supplier.ContactName
	}
	if supplier.Email != "" {
		update["email"] = supplier.Email
	}
	if supplier.Phone != "" {
		update["phone"] = supplier.Phone
	}
	if supplier.Address != "" {
		update["address"] = supplier.Address
	}
	if supplier.Website != "" {
		update["website"] = supplier.Website
	}
	if supplier.Notes != "" {
		update["notes"] = supplier.Notes
	}

	_, err = collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": update})
	return err
}

func (s *SupplierService) DeleteSupplier(id string) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("suppliers")

	_, err = collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (s *SupplierService) GetSupplierOrders(supplierId string) ([]supplier_models.SupplierOrder, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("supplier_orders")

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "order_date", Value: -1}})

	filter := bson.M{}
	if supplierId != "" {
		filter["supplier_id"] = supplierId
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	orders := make([]supplier_models.SupplierOrder, 0)
	for cursor.Next(context.Background()) {
		var order supplier_models.SupplierOrder
		if err := cursor.Decode(&order); err != nil {
			continue
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (s *SupplierService) CreateSupplierOrder(order *supplier_models.SupplierOrder) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("supplier_orders")

	order.Id = bson.NewObjectID().Hex()
	order.OrderDate = time.Now()
	order.Status = "pending"

	_, err = collection.InsertOne(ctx, order)
	return err
}

func (s *SupplierService) UpdateSupplierOrderStatus(id string, status string) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("supplier_orders")

	_, err = collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.M{"status": status}})
	return err
}
