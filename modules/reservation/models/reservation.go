package models

import "time"

type Reservation struct {
	Id           string    `json:"id" bson:"id,omitempty" mapstructure:"id,omitempty"`
	CustomerName string    `json:"customer_name" bson:"customer_name" mapstructure:"customer_name"`
	CustomerPhone string   `json:"customer_phone" bson:"customer_phone" mapstructure:"customer_phone"`
	CustomerEmail string   `json:"customer_email" bson:"customer_email" mapstructure:"customer_email"`
	TableId      string    `json:"table_id" bson:"table_id" mapstructure:"table_id"`
	BranchId     string    `json:"branch_id" bson:"branch_id" mapstructure:"branch_id"`
	Date         string    `json:"date" bson:"date" mapstructure:"date"`
	Time         string    `json:"time" bson:"time" mapstructure:"time"`
	GuestCount   int       `json:"guest_count" bson:"guest_count" mapstructure:"guest_count"`
	Status       string    `json:"status" bson:"status" mapstructure:"status"`
	Notes        string    `json:"notes" bson:"notes" mapstructure:"notes"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at" mapstructure:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at" mapstructure:"updated_at"`
}
