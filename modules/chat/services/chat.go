package services

import (
	"context"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	chat_models "github.com/nutrixpos/pos/modules/chat/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ChatService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *ChatService) GetChannels() ([]chat_models.ChatChannel, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("chat_channels")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	channels := make([]chat_models.ChatChannel, 0)
	for cursor.Next(ctx) {
		var ch chat_models.ChatChannel
		if err := cursor.Decode(&ch); err != nil {
			continue
		}
		channels = append(channels, ch)
	}

	if len(channels) == 0 {
		defaults := []chat_models.ChatChannel{
			{Id: "general", Name: "General", Description: "General discussion", IsDefault: true},
			{Id: "kitchen", Name: "Kitchen", Description: "Kitchen team chat", IsDefault: false},
			{Id: "service", Name: "Service", Description: "Front of house team", IsDefault: false},
			{Id: "management", Name: "Management", Description: "Management only", IsDefault: false},
		}
		for _, ch := range defaults {
			collection.InsertOne(ctx, ch)
		}
		return defaults, nil
	}

	return channels, nil
}

func (s *ChatService) CreateChannel(channel *chat_models.ChatChannel) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("chat_channels")

	_, err = collection.InsertOne(ctx, channel)
	return err
}

func (s *ChatService) GetMessages(channel string, limit int) ([]chat_models.ChatMessage, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("chat_messages")

	filter := bson.M{}
	if channel != "" {
		filter["channel"] = channel
	}

	if limit <= 0 {
		limit = 100
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	findOptions.SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	messages := make([]chat_models.ChatMessage, 0)
	for cursor.Next(ctx) {
		var msg chat_models.ChatMessage
		if err := cursor.Decode(&msg); err != nil {
			continue
		}
		messages = append(messages, msg)
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (s *ChatService) SendMessage(msg *chat_models.ChatMessage) error {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(s.Config.Databases[0].Database).Collection("chat_messages")

	msg.Id = bson.NewObjectID().Hex()
	msg.CreatedAt = time.Now()

	_, err = collection.InsertOne(ctx, msg)
	return err
}
