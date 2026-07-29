package services

import (
	"context"
	"math"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	timeclock_models "github.com/nutrixpos/pos/modules/timeclock/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type TimeClockService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *TimeClockService) ClockIn(employeeId, employeeName, notes string) (*timeclock_models.ClockEntry, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("timeclock_entries")

	entry := timeclock_models.ClockEntry{
		Id:           bson.NewObjectID().Hex(),
		EmployeeId:   employeeId,
		EmployeeName: employeeName,
		ClockIn:      time.Now(),
		Status:       "active",
		Notes:        notes,
	}

	_, err = collection.InsertOne(ctx, entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (s *TimeClockService) ClockOut(entryId string) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("timeclock_entries")

	var entry timeclock_models.ClockEntry
	err = collection.FindOne(ctx, bson.M{"id": entryId}).Decode(&entry)
	if err != nil {
		return err
	}

	now := time.Now()
	hoursWorked := now.Sub(entry.ClockIn).Hours()

	_, err = collection.UpdateOne(ctx, bson.M{"id": entryId}, bson.M{
		"$set": bson.M{
			"clock_out":    now,
			"status":       "completed",
			"hours_worked": math.Round(hoursWorked*100) / 100,
		},
	})

	return err
}

func (s *TimeClockService) GetActiveEntries() ([]timeclock_models.ClockEntry, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("timeclock_entries")

	cursor, err := collection.Find(ctx, bson.M{"status": "active"})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	entries := make([]timeclock_models.ClockEntry, 0)
	for cursor.Next(context.Background()) {
		var entry timeclock_models.ClockEntry
		if err := cursor.Decode(&entry); err != nil {
			continue
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (s *TimeClockService) GetEntriesByDate(date time.Time) ([]timeclock_models.ClockEntry, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("timeclock_entries")

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "clock_in", Value: -1}})

	cursor, err := collection.Find(ctx, bson.M{
		"clock_in": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	entries := make([]timeclock_models.ClockEntry, 0)
	for cursor.Next(context.Background()) {
		var entry timeclock_models.ClockEntry
		if err := cursor.Decode(&entry); err != nil {
			continue
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (s *TimeClockService) GetSummary(startDate, endDate time.Time) ([]timeclock_models.TimeClockSummary, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("timeclock_entries")

	cursor, err := collection.Find(ctx, bson.M{
		"clock_in": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
		"status": "completed",
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	employeeMap := make(map[string]*timeclock_models.TimeClockSummary)

	for cursor.Next(context.Background()) {
		var entry timeclock_models.ClockEntry
		if err := cursor.Decode(&entry); err != nil {
			continue
		}

		if _, exists := employeeMap[entry.EmployeeId]; !exists {
			employeeMap[entry.EmployeeId] = &timeclock_models.TimeClockSummary{
				EmployeeId:   entry.EmployeeId,
				EmployeeName: entry.EmployeeName,
			}
		}

		summary := employeeMap[entry.EmployeeId]
		summary.TotalHours += entry.HoursWorked
		summary.ShiftCount++

		if entry.HoursWorked > 8 {
			summary.OvertimeHours += entry.HoursWorked - 8
		}
	}

	summaries := make([]timeclock_models.TimeClockSummary, 0, len(employeeMap))
	for _, s := range employeeMap {
		if s.ShiftCount > 0 {
			s.AvgHoursPerShift = s.TotalHours / float64(s.ShiftCount)
		}
		s.TotalHours = math.Round(s.TotalHours*100) / 100
		s.OvertimeHours = math.Round(s.OvertimeHours*100) / 100
		s.AvgHoursPerShift = math.Round(s.AvgHoursPerShift*100) / 100
		summaries = append(summaries, *s)
	}

	return summaries, nil
}

func (s *TimeClockService) GetDashboard() (*timeclock_models.TimeClockDashboard, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("timeclock_entries")

	cursor, err := collection.Find(ctx, bson.M{"status": "active"})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	active := make([]timeclock_models.ClockEntry, 0)
	for cursor.Next(context.Background()) {
		var entry timeclock_models.ClockEntry
		if err := cursor.Decode(&entry); err != nil {
			continue
		}
		active = append(active, entry)
	}

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekStart := todayStart.AddDate(0, 0, -int(now.Weekday()))

	todaySummary, err := s.GetSummary(todayStart, now)
	if err != nil {
		todaySummary = make([]timeclock_models.TimeClockSummary, 0)
	}

	weekSummary, err := s.GetSummary(weekStart, now)
	if err != nil {
		weekSummary = make([]timeclock_models.TimeClockSummary, 0)
	}

	dashboard := &timeclock_models.TimeClockDashboard{
		CurrentlyClockedIn: active,
		TodaySummary:       todaySummary,
		WeekSummary:        weekSummary,
	}

	return dashboard, nil
}
