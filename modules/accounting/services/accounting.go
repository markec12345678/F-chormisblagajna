package services

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	core_models "github.com/nutrixpos/pos/modules/core/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type AccountingService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *AccountingService) generateQuickBooksCSV(startDate, endDate time.Time) ([]byte, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ordersCollection := client.Database(s.Config.Databases[0].Database).Collection("orders")

	filter := bson.M{
		"submitted_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
		"is_paid": true,
	}

	cursor, err := ordersCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	writer.Write([]string{"Date", "Type", "Num", "Name", "Memo", "Amount", "Account", "Income Account", "Tax Rate"})

	for cursor.Next(context.Background()) {
		var order core_models.Order
		if err := cursor.Decode(&order); err != nil {
			continue
		}

		amount := order.SalePrice - order.Discount
		taxRate := 0.0
		if order.SalePrice > 0 {
			taxRate = ((order.SalePrice - amount) / amount) * 100
		}

		writer.Write([]string{
			order.SubmittedAt.Format("01/02/2006"),
			"Invoice",
			order.DisplayId,
			order.Customer.Name,
			fmt.Sprintf("POS Order %s", order.DisplayId),
			fmt.Sprintf("%.2f", amount),
			"Sales",
			"Sales Income",
			fmt.Sprintf("%.2f", taxRate),
		})
	}

	writer.Flush()
	return buf.Bytes(), writer.Error()
}

func (s *AccountingService) generateXeroCSV(startDate, endDate time.Time) ([]byte, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ordersCollection := client.Database(s.Config.Databases[0].Database).Collection("orders")

	filter := bson.M{
		"submitted_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
		"is_paid": true,
	}

	cursor, err := ordersCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	writer.Write([]string{"Date", "Reference", "Payee", "Description", "Line Amount", "Tax Amount", "Account Code", "Tax Type", "Category"})

	for cursor.Next(context.Background()) {
		var order core_models.Order
		if err := cursor.Decode(&order); err != nil {
			continue
		}

		lineAmount := order.SalePrice - order.Discount
		taxAmount := order.Discount

		writer.Write([]string{
			order.SubmittedAt.Format("02/01/2006"),
			order.DisplayId,
			order.Customer.Name,
			fmt.Sprintf("POS Sale %s", order.DisplayId),
			fmt.Sprintf("%.2f", lineAmount),
			fmt.Sprintf("%.2f", taxAmount),
			"200",
			"OUTPUT2",
			"Sales",
		})
	}

	writer.Flush()
	return buf.Bytes(), writer.Error()
}

func (s *AccountingService) ExportQuickBooks(startDate, endDate time.Time) ([]byte, error) {
	return s.generateQuickBooksCSV(startDate, endDate)
}

func (s *AccountingService) ExportXero(startDate, endDate time.Time) ([]byte, error) {
	return s.generateXeroCSV(startDate, endDate)
}
