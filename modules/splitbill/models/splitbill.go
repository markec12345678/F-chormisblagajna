package models

import "time"

type SplitBill struct {
	Id        string      `json:"id" bson:"id,omitempty" mapstructure:"id,omitempty"`
	OrderId   string      `json:"order_id" bson:"order_id" mapstructure:"order_id"`
	SplitType string      `json:"split_type" bson:"split_type" mapstructure:"split_type"`
	Parts     []SplitPart `json:"parts" bson:"parts" mapstructure:"parts"`
	Status    string      `json:"status" bson:"status" mapstructure:"status"`
	CreatedAt time.Time   `json:"created_at" bson:"created_at" mapstructure:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" bson:"updated_at" mapstructure:"updated_at"`
}

type SplitPart struct {
	Id            string     `json:"id" bson:"id" mapstructure:"id"`
	Amount        float64    `json:"amount" bson:"amount" mapstructure:"amount"`
	PaymentMethod string     `json:"payment_method" bson:"payment_method" mapstructure:"payment_method"`
	IsPaid        bool       `json:"is_paid" bson:"is_paid" mapstructure:"is_paid"`
	PaidAt        *time.Time `json:"paid_at,omitempty" bson:"paid_at,omitempty" mapstructure:"paid_at,omitempty"`
}
