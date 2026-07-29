package services

import (
	"context"
	"sort"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	waste_models "github.com/nutrixpos/pos/modules/waste/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type WasteService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *WasteService) GetAllWaste() ([]waste_models.WasteEntry, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("waste_entries")

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "date", Value: -1}})

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	entries := make([]waste_models.WasteEntry, 0)
	for cursor.Next(context.Background()) {
		var entry waste_models.WasteEntry
		if err := cursor.Decode(&entry); err != nil {
			continue
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (s *WasteService) CreateWaste(entry *waste_models.WasteEntry) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("waste_entries")

	entry.Id = bson.NewObjectID().Hex()
	entry.Date = time.Now()

	_, err = collection.InsertOne(ctx, entry)
	return err
}

func (s *WasteService) DeleteWaste(id string) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("waste_entries")

	_, err = collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (s *WasteService) GetSummary(startDate, endDate time.Time) (*waste_models.WasteSummary, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("waste_entries")

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

	reasonMap := make(map[string]*waste_models.ReasonSummary)
	materialMap := make(map[string]*waste_models.MaterialSummary)
	dailyMap := make(map[string]*waste_models.DailyWaste)
	var totalCost float64
	var totalCount int

	for cursor.Next(context.Background()) {
		var entry waste_models.WasteEntry
		if err := cursor.Decode(&entry); err != nil {
			continue
		}

		totalCost += entry.Cost
		totalCount++

		if _, exists := reasonMap[entry.Reason]; !exists {
			reasonMap[entry.Reason] = &waste_models.ReasonSummary{Reason: entry.Reason}
		}
		reasonMap[entry.Reason].Total += entry.Cost
		reasonMap[entry.Reason].Count++

		if _, exists := materialMap[entry.MaterialId]; !exists {
			materialMap[entry.MaterialId] = &waste_models.MaterialSummary{
				MaterialId:   entry.MaterialId,
				MaterialName: entry.MaterialName,
			}
		}
		materialMap[entry.MaterialId].TotalCost += entry.Cost
		materialMap[entry.MaterialId].TotalQty += entry.Quantity
		materialMap[entry.MaterialId].Count++

		dateKey := entry.Date.Format("2006-01-02")
		if _, exists := dailyMap[dateKey]; !exists {
			dailyMap[dateKey] = &waste_models.DailyWaste{Date: dateKey}
		}
		dailyMap[dateKey].Total += entry.Cost
	}

	reasons := make([]waste_models.ReasonSummary, 0, len(reasonMap))
	for _, r := range reasonMap {
		reasons = append(reasons, *r)
	}
	sort.Slice(reasons, func(i, j int) bool { return reasons[i].Total > reasons[j].Total })

	materials := make([]waste_models.MaterialSummary, 0, len(materialMap))
	for _, m := range materialMap {
		materials = append(materials, *m)
	}
	sort.Slice(materials, func(i, j int) bool { return materials[i].TotalCost > materials[j].TotalCost })

	daily := make([]waste_models.DailyWaste, 0, len(dailyMap))
	for _, d := range dailyMap {
		daily = append(daily, *d)
	}
	sort.Slice(daily, func(i, j int) bool { return daily[i].Date > daily[j].Date })

	summary := &waste_models.WasteSummary{
		TotalWasteCost: totalCost,
		TotalEntries:   totalCount,
		ByReason:       reasons,
		ByMaterial:     materials,
		DailyWaste:     daily,
	}

	return summary, nil
}
