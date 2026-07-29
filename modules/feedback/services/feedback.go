package services

import (
	"context"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	fb_models "github.com/nutrixpos/pos/modules/feedback/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type FeedbackService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *FeedbackService) getCollection() (*mongo.Collection, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}
	return client.Database(s.Config.Databases[0].Database).Collection("feedbacks"), nil
}

func (s *FeedbackService) SubmitFeedback(fb *fb_models.Feedback) error {
	coll, err := s.getCollection()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fb.Id = bson.NewObjectID().Hex()
	fb.CreatedAt = time.Now()
	fb.Responded = false

	_, err = coll.InsertOne(ctx, fb)
	return err
}

func (s *FeedbackService) GetAllFeedbacks() ([]fb_models.Feedback, error) {
	coll, err := s.getCollection()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := coll.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	feedbacks := make([]fb_models.Feedback, 0)
	for cursor.Next(ctx) {
		var f fb_models.Feedback
		if err := cursor.Decode(&f); err != nil {
			continue
		}
		feedbacks = append(feedbacks, f)
	}

	return feedbacks, nil
}

func (s *FeedbackService) GetFeedbackSummary() (*fb_models.FeedbackSummary, error) {
	coll, err := s.getCollection()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	feedbacks := make([]fb_models.Feedback, 0)
	cursor, err := coll.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(50))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var f fb_models.Feedback
		if err := cursor.Decode(&f); err != nil {
			continue
		}
		feedbacks = append(feedbacks, f)
	}

	summary := &fb_models.FeedbackSummary{
		TotalFeedbacks:  len(feedbacks),
		RatingDist:      make(map[int]int),
		CategoryAvg:     make(map[string]float64),
		RecentFeedbacks: feedbacks,
	}

	var totalRating int
	catSums := make(map[string]float64)
	catCounts := make(map[string]int)

	for _, f := range feedbacks {
		totalRating += f.Rating
		summary.RatingDist[f.Rating]++

		if f.Category != "" {
			catSums[f.Category] += float64(f.Rating)
			catCounts[f.Category]++
		}
	}

	if len(feedbacks) > 0 {
		summary.AverageRating = float64(totalRating) / float64(len(feedbacks))
	}

	for cat, sum := range catSums {
		if catCounts[cat] > 0 {
			summary.CategoryAvg[cat] = sum / float64(catCounts[cat])
		}
	}

	return summary, nil
}

func (s *FeedbackService) RespondToFeedback(id, response string) error {
	coll, err := s.getCollection()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = coll.UpdateOne(ctx, bson.M{"id": id}, bson.M{
		"$set": bson.M{"responded": true, "response": response},
	})
	return err
}

func (s *FeedbackService) DeleteFeedback(id string) error {
	coll, err := s.getCollection()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = coll.DeleteOne(ctx, bson.M{"id": id})
	return err
}
