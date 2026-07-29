package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Notification struct {
	ID        bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Type      string        `json:"type" bson:"type"`
	Title     string        `json:"title" bson:"title"`
	Message   string        `json:"message" bson:"message"`
	Severity  string        `json:"severity" bson:"severity"`
	Reference string        `json:"reference" bson:"reference"`
	Read      bool          `json:"read" bson:"read"`
	UserID    string        `json:"user_id" bson:"user_id"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
}
