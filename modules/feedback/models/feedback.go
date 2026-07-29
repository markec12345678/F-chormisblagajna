package models

import "time"

type Feedback struct {
	Id         string    `json:"id" bson:"id" mapstructure:"id"`
	OrderId    string    `json:"order_id" bson:"order_id" mapstructure:"order_id"`
	CustomerId string    `json:"customer_id" bson:"customer_id" mapstructure:"customer_id"`
	Rating     int       `json:"rating" bson:"rating" mapstructure:"rating"`
	Comment    string    `json:"comment" bson:"comment" mapstructure:"comment"`
	Category   string    `json:"category" bson:"category" mapstructure:"category"` // food, service, ambiance, overall
	Anonymous  bool      `json:"anonymous" bson:"anonymous" mapstructure:"anonymous"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at" mapstructure:"created_at"`
	Responded  bool      `json:"responded" bson:"responded" mapstructure:"responded"`
	Response   string    `json:"response,omitempty" bson:"response,omitempty" mapstructure:"response,omitempty"`
}

type FeedbackSummary struct {
	TotalFeedbacks  int     `json:"total_feedbacks"`
	AverageRating   float64 `json:"average_rating"`
	RatingDist      map[int]int `json:"rating_distribution"`
	CategoryAvg     map[string]float64 `json:"category_averages"`
	RecentFeedbacks []Feedback `json:"recent_feedbacks"`
}
