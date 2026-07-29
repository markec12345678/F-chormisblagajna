package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/report/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ReportService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *ReportService) GetSalesReport(startDate, endDate string) (*models.SalesReport, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, fmt.Errorf("GetSalesReport: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := client.Database(s.Config.Databases[0].Database)

	filter := bson.M{"state": "finished"}
	if startDate != "" {
		filter["finished_at"] = bson.M{"$gte": startDate}
	}
	if endDate != "" {
		if _, ok := filter["finished_at"]; ok {
			filter["finished_at"] = bson.M{"$gte": startDate, "$lte": endDate}
		} else {
			filter["finished_at"] = bson.M{"$lte": endDate}
		}
	}

	ordersCol := db.Collection("orders")

	totalOrders, _ := ordersCol.CountDocuments(ctx, filter)

	report := &models.SalesReport{
		Period:       fmt.Sprintf("%s to %s", startDate, endDate),
		TotalOrders:  int(totalOrders),
		TopProducts:  make([]models.ProductStat, 0),
	}

	cursor, err := ordersCol.Find(ctx, filter)
	if err != nil {
		return report, nil
	}
	defer cursor.Close(ctx)

	totalRevenue := 0.0
	totalItems := 0
	productMap := make(map[string]*models.ProductStat)

	for cursor.Next(ctx) {
		var order map[string]interface{}
		if err := cursor.Decode(&order); err != nil {
			continue
		}

		if salePrice, ok := order["sale_price"].(float64); ok {
			totalRevenue += salePrice
		}

		if items, ok := order["items"].([]interface{}); ok {
			for _, item := range items {
				totalItems++
				if itemMap, ok := item.(map[string]interface{}); ok {
					name := ""
					if n, ok := itemMap["name"].(string); ok {
						name = n
					} else if p, ok := itemMap["product"].(map[string]interface{}); ok {
						if n, ok := p["name"].(string); ok {
							name = n
						}
					}

					qty := 1
					if q, ok := itemMap["quantity"].(float64); ok {
						qty = int(q)
					}

					sp := 0.0
					if s, ok := itemMap["sale_price"].(float64); ok {
						sp = s
					}

					if existing, ok := productMap[name]; ok {
						existing.Quantity += qty
						existing.Revenue += sp
					} else {
						productMap[name] = &models.ProductStat{
							Name:     name,
							Quantity: qty,
							Revenue:  sp,
						}
					}
				}
			}
		}
	}

	report.TotalRevenue = totalRevenue
	report.TotalItems = totalItems
	if totalOrders > 0 {
		report.AverageOrder = totalRevenue / float64(totalOrders)
	}
	report.NetRevenue = totalRevenue

	for _, ps := range productMap {
		report.TopProducts = append(report.TopProducts, *ps)
	}

	return report, nil
}

func (s *ReportService) GetInventoryReport() (*models.InventoryReport, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, fmt.Errorf("GetInventoryReport: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := client.Database(s.Config.Databases[0].Database)
	col := db.Collection("materials")

	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("GetInventoryReport: %w", err)
	}
	defer cursor.Close(ctx)

	report := &models.InventoryReport{
		LowStockItems: make([]models.LowStockItem, 0),
	}

	for cursor.Next(ctx) {
		var material map[string]interface{}
		if err := cursor.Decode(&material); err != nil {
			continue
		}

		report.TotalMaterials++

		qty := 0.0
		if q, ok := material["quantity"].(float64); ok {
			qty = q
		}
		if q, ok := material["ready"].(float64); ok {
			qty = q
		}

		name := ""
		if n, ok := material["name"].(string); ok {
			name = n
		}

		unit := ""
		if u, ok := material["unit"].(string); ok {
			unit = u
		}

		cost := 0.0
		if c, ok := material["cost"].(float64); ok {
			cost = c
		}

		report.TotalValue += qty * cost

		threshold := 10.0
		if t, ok := material["low_stock_threshold"].(float64); ok {
			threshold = t
		}

		if qty <= 0 {
			report.OutOfStockCount++
			report.LowStockItems = append(report.LowStockItems, models.LowStockItem{
				Name: name, Quantity: qty, Unit: unit, Value: qty * cost,
			})
		} else if qty <= threshold {
			report.LowStockCount++
			report.LowStockItems = append(report.LowStockItems, models.LowStockItem{
				Name: name, Quantity: qty, Unit: unit, Value: qty * cost,
			})
		}
	}

	return report, nil
}

func (s *ReportService) GetDashboardStats() (*models.DashboardStats, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, fmt.Errorf("GetDashboardStats: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := client.Database(s.Config.Databases[0].Database)
	ordersCol := db.Collection("orders")

	today := time.Now().Format("2006-01-02")

	todayFilter := bson.M{
		"state":       "finished",
		"finished_at": bson.M{"$gte": today},
	}

	todayOrders, _ := ordersCol.CountDocuments(ctx, todayFilter)

	stats := &models.DashboardStats{
		TodayOrders: int(todayOrders),
		GeneratedAt: time.Now(),
	}

	cursor, err := ordersCol.Find(ctx, todayFilter)
	if err == nil {
		defer cursor.Close(ctx)

		totalRevenue := 0.0
		customerSet := make(map[string]bool)

		for cursor.Next(ctx) {
			var order map[string]interface{}
			if err := cursor.Decode(&order); err != nil {
				continue
			}

			if salePrice, ok := order["sale_price"].(float64); ok {
				totalRevenue += salePrice
			}

			if customer, ok := order["customer"].(map[string]interface{}); ok {
				if id, ok := customer["id"].(string); ok && id != "" {
					customerSet[id] = true
				}
			}
		}

		stats.TodayRevenue = totalRevenue
		stats.TodayCustomers = len(customerSet)
		if todayOrders > 0 {
			stats.AverageOrder = totalRevenue / float64(todayOrders)
		}
	}

	return stats, nil
}
