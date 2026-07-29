package models

import "time"

type Customer struct {
	Id            string            `json:"id" bson:"id" mapstructure:"id"`
	Name          string            `json:"name" bson:"name" mapstructure:"name"`
	Phone         string            `json:"phone" bson:"phone" mapstructure:"phone"`
	Address       string            `json:"address" bson:"address" mapstructure:"address"`
	Email         string            `json:"email" bson:"email" mapstructure:"email"`
	Notes         string            `json:"notes" bson:"notes" mapstructure:"notes"`
	Tags          []string          `json:"tags" bson:"tags" mapstructure:"tags"`
	LoyaltyPoints int               `json:"loyalty_points" bson:"loyalty_points" mapstructure:"loyalty_points"`
	TotalSpent    float64           `json:"total_spent" bson:"total_spent" mapstructure:"total_spent"`
	OrderCount    int               `json:"order_count" bson:"order_count" mapstructure:"order_count"`
	LastOrderDate *time.Time        `json:"last_order_date" bson:"last_order_date" mapstructure:"last_order_date"`
	Preferences   map[string]string `json:"preferences" bson:"preferences" mapstructure:"preferences"`
	CreatedAt     time.Time         `json:"created_at" bson:"created_at" mapstructure:"created_at"`
}
