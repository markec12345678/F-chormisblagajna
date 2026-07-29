package services

import (
	"context"
	"sort"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	expense_models "github.com/nutrixpos/pos/modules/expense/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ExpenseService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *ExpenseService) GetAllExpenses() ([]expense_models.Expense, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("expenses")

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "date", Value: -1}})

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	expenses := make([]expense_models.Expense, 0)
	for cursor.Next(context.Background()) {
		var expense expense_models.Expense
		if err := cursor.Decode(&expense); err != nil {
			continue
		}
		expenses = append(expenses, expense)
	}

	return expenses, nil
}

func (s *ExpenseService) CreateExpense(expense *expense_models.Expense) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("expenses")

	expense.Id = bson.NewObjectID().Hex()
	expense.CreatedAt = time.Now()

	_, err = collection.InsertOne(ctx, expense)
	return err
}

func (s *ExpenseService) UpdateExpense(id string, expense *expense_models.Expense) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("expenses")

	update := bson.M{}
	if expense.Description != "" {
		update["description"] = expense.Description
	}
	if expense.Amount != 0 {
		update["amount"] = expense.Amount
	}
	if expense.Category != "" {
		update["category"] = expense.Category
	}
	if !expense.Date.IsZero() {
		update["date"] = expense.Date
	}
	if expense.Notes != "" {
		update["notes"] = expense.Notes
	}

	_, err = collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": update})
	return err
}

func (s *ExpenseService) DeleteExpense(id string) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("expenses")

	_, err = collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (s *ExpenseService) GetExpenseSummary(startDate, endDate time.Time) (*expense_models.ExpenseSummary, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("expenses")

	filter := bson.M{
		"date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	categoryMap := make(map[string]*expense_models.CategorySummary)
	var totalExpenses float64
	var monthlyTotal float64

	now := time.Now()
	currentMonth := now.Month()
	currentYear := now.Year()

	for cursor.Next(context.Background()) {
		var expense expense_models.Expense
		if err := cursor.Decode(&expense); err != nil {
			continue
		}

		totalExpenses += expense.Amount

		if expense.Date.Month() == currentMonth && expense.Date.Year() == currentYear {
			monthlyTotal += expense.Amount
		}

		if _, exists := categoryMap[expense.Category]; !exists {
			categoryMap[expense.Category] = &expense_models.CategorySummary{
				Category: expense.Category,
			}
		}
		cat := categoryMap[expense.Category]
		cat.Total += expense.Amount
		cat.Count++
	}

	categories := make([]expense_models.CategorySummary, 0, len(categoryMap))
	for _, cat := range categoryMap {
		categories = append(categories, *cat)
	}

	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Total > categories[j].Total
	})

	summary := &expense_models.ExpenseSummary{
		TotalExpenses: totalExpenses,
		ByCategory:    categories,
		MonthlyTotal:  monthlyTotal,
	}

	return summary, nil
}
