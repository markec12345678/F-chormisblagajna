package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type NotificationService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *NotificationService) GetDatabase() (*mongo.Database, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}
	return client.Database(s.Config.Databases[0].Database), nil
}

func (s *NotificationService) CreateNotification(notificationType, title, message, severity, reference, userID string) error {
	db, err := s.GetDatabase()
	if err != nil {
		return fmt.Errorf("CreateNotification: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	notification := bson.M{
		"type":       notificationType,
		"title":      title,
		"message":    message,
		"severity":   severity,
		"reference":  reference,
		"read":       false,
		"user_id":    userID,
		"created_at": time.Now(),
	}

	_, err = db.Collection("notifications").InsertOne(ctx, notification)
	return err
}

func (s *NotificationService) GetNotifications(userID string, unreadOnly bool) ([]bson.M, error) {
	db, err := s.GetDatabase()
	if err != nil {
		return nil, fmt.Errorf("GetNotifications: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}
	if userID != "" {
		filter["user_id"] = userID
	}
	if unreadOnly {
		filter["read"] = false
	}

	opts := options.Find().SetSort(bson.M{"created_at": -1}).SetLimit(50)
	cursor, err := db.Collection("notifications").Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("GetNotifications Find: %w", err)
	}
	defer cursor.Close(ctx)

	var notifications []bson.M
	if err := cursor.All(ctx, &notifications); err != nil {
		return nil, fmt.Errorf("GetNotifications decode: %w", err)
	}

	return notifications, nil
}

func (s *NotificationService) MarkAsRead(id string) error {
	db, err := s.GetDatabase()
	if err != nil {
		return fmt.Errorf("MarkAsRead: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("MarkAsRead invalid ID: %w", err)
	}

	_, err = db.Collection("notifications").UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
		"$set": bson.M{"read": true},
	})
	return err
}

func (s *NotificationService) MarkAllAsRead(userID string) error {
	db, err := s.GetDatabase()
	if err != nil {
		return fmt.Errorf("MarkAllAsRead: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"read": false}
	if userID != "" {
		filter["user_id"] = userID
	}

	_, err = db.Collection("notifications").UpdateMany(ctx, filter, bson.M{
		"$set": bson.M{"read": true},
	})
	return err
}

func (s *NotificationService) GetUnreadCount(userID string) (int64, error) {
	db, err := s.GetDatabase()
	if err != nil {
		return 0, fmt.Errorf("GetUnreadCount: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"read": false}
	if userID != "" {
		filter["user_id"] = userID
	}

	count, err := db.Collection("notifications").CountDocuments(ctx, filter)
	return count, err
}
