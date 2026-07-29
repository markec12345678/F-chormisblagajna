package services

import (
	"context"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	auditlog_models "github.com/nutrixpos/pos/modules/auditlog/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type AuditLogService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *AuditLogService) GetAll(action, resource, userId string, limit int) ([]auditlog_models.AuditLogEntry, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("audit_logs")

	filter := bson.M{}
	if action != "" {
		filter["action"] = action
	}
	if resource != "" {
		filter["resource"] = resource
	}
	if userId != "" {
		filter["user_id"] = userId
	}

	if limit <= 0 {
		limit = 200
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	findOptions.SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	entries := make([]auditlog_models.AuditLogEntry, 0)
	for cursor.Next(context.Background()) {
		var entry auditlog_models.AuditLogEntry
		if err := cursor.Decode(&entry); err != nil {
			continue
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (s *AuditLogService) Create(entry *auditlog_models.AuditLogEntry) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("audit_logs")

	entry.Id = bson.NewObjectID().Hex()
	entry.CreatedAt = time.Now()

	_, err = collection.InsertOne(ctx, entry)
	return err
}

func (s *AuditLogService) GetSummary() (*auditlog_models.AuditLogSummary, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("audit_logs")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	actionMap := make(map[string]*auditlog_models.ActionSummary)
	resourceMap := make(map[string]*auditlog_models.ResourceSummary)
	userMap := make(map[string]*auditlog_models.UserSummary)
	totalCount := 0

	for cursor.Next(context.Background()) {
		var entry auditlog_models.AuditLogEntry
		if err := cursor.Decode(&entry); err != nil {
			continue
		}

		totalCount++

		if _, exists := actionMap[entry.Action]; !exists {
			actionMap[entry.Action] = &auditlog_models.ActionSummary{Action: entry.Action}
		}
		actionMap[entry.Action].Count++

		if _, exists := resourceMap[entry.Resource]; !exists {
			resourceMap[entry.Resource] = &auditlog_models.ResourceSummary{Resource: entry.Resource}
		}
		resourceMap[entry.Resource].Count++

		if entry.UserId != "" {
			if _, exists := userMap[entry.UserId]; !exists {
				userMap[entry.UserId] = &auditlog_models.UserSummary{
					UserId:   entry.UserId,
					Username: entry.Username,
				}
			}
			userMap[entry.UserId].Count++
		}
	}

	actions := make([]auditlog_models.ActionSummary, 0, len(actionMap))
	for _, a := range actionMap {
		actions = append(actions, *a)
	}

	resources := make([]auditlog_models.ResourceSummary, 0, len(resourceMap))
	for _, r := range resourceMap {
		resources = append(resources, *r)
	}

	users := make([]auditlog_models.UserSummary, 0, len(userMap))
	for _, u := range userMap {
		users = append(users, *u)
	}

	return &auditlog_models.AuditLogSummary{
		TotalEntries: totalCount,
		ByAction:     actions,
		ByResource:   resources,
		ByUser:       users,
	}, nil
}
