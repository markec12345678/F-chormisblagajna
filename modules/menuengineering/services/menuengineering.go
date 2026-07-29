package services

import (
	"context"
	"sort"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/core/models"
	menu_models "github.com/nutrixpos/pos/modules/menuengineering/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MenuEngineeringService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *MenuEngineeringService) AnalyzeMenu(startDate, endDate time.Time) (*menu_models.MenuEngineeringSummary, error) {
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

	itemMap := make(map[string]*menu_models.MenuItemAnalysis)

	for cursor.Next(context.Background()) {
		var order models.Order
		if err := cursor.Decode(&order); err != nil {
			continue
		}

		for _, item := range order.Items {
			id := item.Product.Id
			if _, exists := itemMap[id]; !exists {
				itemMap[id] = &menu_models.MenuItemAnalysis{
					ProductId:   id,
					ProductName: item.Product.Name,
					Category:    "",
				}
			}

			analysis := itemMap[id]
			qty := int(item.Quantity)
			analysis.TotalSold += qty
			analysis.TotalRevenue += item.Price * item.Quantity
			analysis.TotalCost += item.Cost * item.Quantity
		}
	}

	items := make([]menu_models.MenuItemAnalysis, 0, len(itemMap))
	for _, item := range itemMap {
		item.ProfitMargin = 0
		if item.TotalRevenue > 0 {
			item.ProfitMargin = ((item.TotalRevenue - item.TotalCost) / item.TotalRevenue) * 100
		}
		item.ProfitPerItem = 0
		if item.TotalSold > 0 {
			item.ProfitPerItem = (item.TotalRevenue - item.TotalCost) / float64(item.TotalSold)
		}
		items = append(items, *item)
	}

	if len(items) == 0 {
		return &menu_models.MenuEngineeringSummary{
			TotalItems:   0,
			TotalRevenue: 0,
		}, nil
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].TotalSold > items[j].TotalSold
	})
	for i := range items {
		items[i].PopularityRank = i + 1
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].ProfitPerItem > items[j].ProfitPerItem
	})
	for i := range items {
		items[i].ProfitRank = i + 1
	}

	var totalProfit float64
	for _, item := range items {
		totalProfit += item.ProfitPerItem
	}
	avgProfit := totalProfit / float64(len(items))

	avgPopularity := float64(len(items)) / 2.0

	for i := range items {
		highProfit := items[i].ProfitPerItem >= avgProfit
		highPopularity := float64(items[i].TotalSold) >= avgPopularity

		switch {
		case highProfit && highPopularity:
			items[i].Quadrant = "star"
		case !highProfit && highPopularity:
			items[i].Quadrant = "plowhorse"
		case highProfit && !highPopularity:
			items[i].Quadrant = "puzzle"
		default:
			items[i].Quadrant = "dog"
		}
	}

	summary := &menu_models.MenuEngineeringSummary{
		TotalItems:   len(items),
		TotalRevenue: 0,
		AvgProfit:    avgProfit,
	}

	for _, item := range items {
		summary.TotalRevenue += item.TotalRevenue
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].ProfitPerItem > items[j].ProfitPerItem
	})

	for _, item := range items {
		switch item.Quadrant {
		case "star":
			if len(summary.TopStars) < 5 {
				summary.TopStars = append(summary.TopStars, item)
			}
		case "plowhorse":
			if len(summary.TopPlowhorses) < 5 {
				summary.TopPlowhorses = append(summary.TopPlowhorses, item)
			}
		case "puzzle":
			if len(summary.TopPuzzles) < 5 {
				summary.TopPuzzles = append(summary.TopPuzzles, item)
			}
		case "dog":
			if len(summary.TopDogs) < 5 {
				summary.TopDogs = append(summary.TopDogs, item)
			}
		}
	}

	return summary, nil
}
