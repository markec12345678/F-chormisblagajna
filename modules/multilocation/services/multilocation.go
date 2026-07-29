package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/multilocation/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MultiLocationService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *MultiLocationService) GetDatabase() (*mongo.Database, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}
	return client.Database(s.Config.Databases[0].Database), nil
}

func (s *MultiLocationService) GetDashboard() (*models.LocationDashboard, error) {
	db, err := s.GetDatabase()
	if err != nil {
		return nil, fmt.Errorf("GetDashboard: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	branchesCol := db.Collection("branches")
	ordersCol := db.Collection("orders")

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekStart := todayStart.AddDate(0, 0, -int(todayStart.Weekday()))
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	cursor, err := branchesCol.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("GetDashboard Find branches: %w", err)
	}
	defer cursor.Close(ctx)

	var branches []struct {
		ID   string `bson:"_id"`
		Name string `bson:"name"`
	}
	if err := cursor.All(ctx, &branches); err != nil {
		return nil, fmt.Errorf("GetDashboard decode branches: %w", err)
	}

	dashboard := &models.LocationDashboard{
		TotalBranches: len(branches),
		GeneratedAt:   now,
	}

	for _, branch := range branches {
		bs := models.BranchStats{
			BranchID:   branch.ID,
			BranchName: branch.Name,
			Status:     "active",
		}

		bs.TodayOrders, _ = s.countOrders(ctx, ordersCol, branch.ID, todayStart, now)
		bs.TodayRevenue, _ = s.sumRevenue(ctx, ordersCol, branch.ID, todayStart, now)
		bs.WeekOrders, _ = s.countOrders(ctx, ordersCol, branch.ID, weekStart, now)
		bs.WeekRevenue, _ = s.sumRevenue(ctx, ordersCol, branch.ID, weekStart, now)
		bs.MonthOrders, _ = s.countOrders(ctx, ordersCol, branch.ID, monthStart, now)
		bs.MonthRevenue, _ = s.sumRevenue(ctx, ordersCol, branch.ID, monthStart, now)

		if bs.TodayOrders > 0 {
			bs.AvgOrderValue = bs.TodayRevenue / float64(bs.TodayOrders)
		}

		dashboard.TotalRevenue += bs.TodayRevenue
		dashboard.TotalOrders += bs.TodayOrders

		dashboard.Branches = append(dashboard.Branches, bs)
	}

	if dashboard.TotalOrders > 0 {
		dashboard.AvgOrderValue = dashboard.TotalRevenue / float64(dashboard.TotalOrders)
	}

	return dashboard, nil
}

func (s *MultiLocationService) countOrders(ctx context.Context, col *mongo.Collection, branchID string, start, end time.Time) (int, error) {
	filter := bson.M{
		"branch_id":  branchID,
		"created_at": bson.M{"$gte": start, "$lte": end},
		"state":      "finished",
	}
	count, err := col.CountDocuments(ctx, filter)
	return int(count), err
}

func (s *MultiLocationService) sumRevenue(ctx context.Context, col *mongo.Collection, branchID string, start, end time.Time) (float64, error) {
	filter := bson.M{
		"branch_id":  branchID,
		"created_at": bson.M{"$gte": start, "$lte": end},
		"state":      "finished",
	}
	pipeline := []bson.M{
		{"$match": filter},
		{"$group": bson.M{"_id": nil, "total": bson.M{"$sum": "$total"}}},
	}
	cursor, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return 0, err
	}
	if len(results) > 0 {
		if v, ok := results[0]["total"].(float64); ok {
			return v, nil
		}
	}
	return 0, nil
}

func (s *MultiLocationService) GetComparison() ([]models.BranchComparison, error) {
	db, err := s.GetDatabase()
	if err != nil {
		return nil, fmt.Errorf("GetComparison: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	branchesCol := db.Collection("branches")
	ordersCol := db.Collection("orders")

	cursor, err := branchesCol.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("GetComparison Find branches: %w", err)
	}
	defer cursor.Close(ctx)

	var branches []struct {
		ID   string `bson:"_id"`
		Name string `bson:"name"`
	}
	if err := cursor.All(ctx, &branches); err != nil {
		return nil, fmt.Errorf("GetComparison decode branches: %w", err)
	}

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekStart := todayStart.AddDate(0, 0, -int(todayStart.Weekday()))
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var comparisons []models.BranchComparison

	comparisons = append(comparisons, s.buildComparison(ctx, ordersCol, branches, "today_revenue", todayStart, now))
	comparisons = append(comparisons, s.buildComparison(ctx, ordersCol, branches, "week_revenue", weekStart, now))
	comparisons = append(comparisons, s.buildComparison(ctx, ordersCol, branches, "month_revenue", monthStart, now))

	return comparisons, nil
}

func (s *MultiLocationService) buildComparison(ctx context.Context, col *mongo.Collection, branches []struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
}, metric string, start, end time.Time) models.BranchComparison {
	comp := models.BranchComparison{Metric: metric}

	for _, branch := range branches {
		value, _ := s.sumRevenue(ctx, col, branch.ID, start, end)
		comp.Branches = append(comp.Branches, models.BranchMetricValue{
			BranchID:   branch.ID,
			BranchName: branch.Name,
			Value:      value,
		})
	}

	return comp
}
