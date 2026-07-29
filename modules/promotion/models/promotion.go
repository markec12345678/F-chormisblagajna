package models

import "time"

type Promotion struct {
	Id            string    `json:"id" bson:"id,omitempty" mapstructure:"id,omitempty"`
	Name          string    `json:"name" bson:"name" mapstructure:"name"`
	Code          string    `json:"code" bson:"code" mapstructure:"code"`
	Type          string    `json:"type" bson:"type" mapstructure:"type"`
	Value         float64   `json:"value" bson:"value" mapstructure:"value"`
	MinOrder      float64   `json:"min_order" bson:"min_order" mapstructure:"min_order"`
	MaxDiscount   float64   `json:"max_discount" bson:"max_discount" mapstructure:"max_discount"`
	StartDate     string    `json:"start_date" bson:"start_date" mapstructure:"start_date"`
	EndDate       string    `json:"end_date" bson:"end_date" mapstructure:"end_date"`
	UsageLimit    int       `json:"usage_limit" bson:"usage_limit" mapstructure:"usage_limit"`
	UsageCount    int       `json:"usage_count" bson:"usage_count" mapstructure:"usage_count"`
	ApplicableDays []string `json:"applicable_days" bson:"applicable_days" mapstructure:"applicable_days"`
	HappyHourStart string   `json:"happy_hour_start" bson:"happy_hour_start" mapstructure:"happy_hour_start"`
	HappyHourEnd   string   `json:"happy_hour_end" bson:"happy_hour_end" mapstructure:"happy_hour_end"`
	IsActive      bool      `json:"is_active" bson:"is_active" mapstructure:"is_active"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at" mapstructure:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at" mapstructure:"updated_at"`
}
