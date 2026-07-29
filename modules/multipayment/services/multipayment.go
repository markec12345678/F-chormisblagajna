package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	mp_models "github.com/nutrixpos/pos/modules/multipayment/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MultiPaymentService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *MultiPaymentService) AddPayment(payment *mp_models.PaymentPart) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("payment_parts")

	payment.Id = bson.NewObjectID().Hex()
	payment.CreatedAt = time.Now()

	_, err = collection.InsertOne(ctx, payment)
	return err
}

func (s *MultiPaymentService) GetPaymentsByOrder(orderId string) ([]mp_models.PaymentPart, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("payment_parts")

	cursor, err := collection.Find(ctx, bson.M{"order_id": orderId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	payments := make([]mp_models.PaymentPart, 0)
	for cursor.Next(ctx) {
		var p mp_models.PaymentPart
		if err := cursor.Decode(&p); err != nil {
			continue
		}
		payments = append(payments, p)
	}

	return payments, nil
}

func (s *MultiPaymentService) GetPaymentSummary(orderId string, totalDue float64) (*mp_models.PaymentSummary, error) {
	payments, err := s.GetPaymentsByOrder(orderId)
	if err != nil {
		return nil, err
	}

	var totalPaid float64
	methodMap := make(map[string]*mp_models.MethodAmount)

	for _, p := range payments {
		totalPaid += p.Amount
		if _, exists := methodMap[p.PaymentMethod]; !exists {
			methodMap[p.PaymentMethod] = &mp_models.MethodAmount{Method: p.PaymentMethod}
		}
		methodMap[p.PaymentMethod].Total += p.Amount
		methodMap[p.PaymentMethod].Count++
	}

	methodBreakdown := make([]mp_models.MethodAmount, 0, len(methodMap))
	for _, m := range methodMap {
		methodBreakdown = append(methodBreakdown, *m)
	}

	summary := &mp_models.PaymentSummary{
		OrderId:       orderId,
		TotalDue:      totalDue,
		TotalPaid:     totalPaid,
		Remaining:     totalDue - totalPaid,
		IsFullyPaid:   totalPaid >= totalDue,
		Payments:      payments,
		MethodBreakdown: methodBreakdown,
	}

	return summary, nil
}

func (s *MultiPaymentService) GetDailyPayments(startDate, endDate time.Time) ([]mp_models.DailyPayments, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("payment_parts")

	filter := bson.M{
		"created_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	cursor, err := collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	dailyMap := make(map[string]*mp_models.DailyPayments)

	for cursor.Next(context.Background()) {
		var p mp_models.PaymentPart
		if err := cursor.Decode(&p); err != nil {
			continue
		}

		dateKey := p.CreatedAt.Format("2006-01-02")
		if _, exists := dailyMap[dateKey]; !exists {
			dailyMap[dateKey] = &mp_models.DailyPayments{Date: dateKey}
		}

		daily := dailyMap[dateKey]
		daily.GrandTotal += p.Amount
		daily.Count++

		switch p.PaymentMethod {
		case "cash":
			daily.TotalCash += p.Amount
		case "card":
			daily.TotalCard += p.Amount
		case "voucher":
			daily.TotalVoucher += p.Amount
		case "mobile":
			daily.TotalMobile += p.Amount
		case "gift_card":
			daily.TotalGift += p.Amount
		}
	}

	result := make([]mp_models.DailyPayments, 0, len(dailyMap))
	for _, d := range dailyMap {
		result = append(result, *d)
	}

	return result, nil
}

// MarkOrderPaid updates the core order's is_paid flag and records the primary payment source.
func (s *MultiPaymentService) MarkOrderPaid(orderId string, paymentSource string) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return fmt.Errorf("MarkOrderPaid: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	coll := client.Database(s.Config.Databases[0].Database).Collection("orders")
	_, err = coll.UpdateOne(ctx, bson.M{"id": orderId}, bson.M{
		"$set": bson.M{"is_paid": true, "payment_source": paymentSource},
	})
	return err
}

// SettleOrder marks an order as paid if the sum of PaymentParts covers the total due.
func (s *MultiPaymentService) SettleOrder(orderId string, totalDue float64) error {
	summary, err := s.GetPaymentSummary(orderId, totalDue)
	if err != nil {
		return fmt.Errorf("SettleOrder: %w", err)
	}
	if summary.IsFullyPaid {
		return s.MarkOrderPaid(orderId, summary.MethodBreakdown[0].Method)
	}
	return fmt.Errorf("order %s is not fully paid (paid %.2f / due %.2f)", orderId, summary.TotalPaid, totalDue)
}

// GetPaymentPartsCollection returns the payment_parts mongo collection.
func (s *MultiPaymentService) GetPaymentPartsCollection() (*mongo.Collection, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}
	return client.Database(s.Config.Databases[0].Database).Collection("payment_parts"), nil
}

func (s *MultiPaymentService) RefundPayment(paymentId string) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("payment_parts")

	_, err = collection.DeleteOne(ctx, bson.M{"id": paymentId})
	return err
}
