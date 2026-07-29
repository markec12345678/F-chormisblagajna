package models

import "time"

type TableSession struct {
	Id         string    `json:"id" bson:"id" mapstructure:"id"`
	TableLabel string    `json:"table_label" bson:"table_label" mapstructure:"table_label"`
	Zone       string    `json:"zone" bson:"zone" mapstructure:"zone"`
	QrToken    string    `json:"qr_token" bson:"qr_token" mapstructure:"qr_token"`
	Active     bool      `json:"active" bson:"active" mapstructure:"active"`
	WaiterId   string    `json:"waiter_id" bson:"waiter_id" mapstructure:"waiter_id"`
	GuestCount int       `json:"guest_count" bson:"guest_count" mapstructure:"guest_count"`
	OpenedAt   time.Time `json:"opened_at" bson:"opened_at" mapstructure:"opened_at"`
	ClosedAt   *time.Time `json:"closed_at,omitempty" bson:"closed_at,omitempty" mapstructure:"closed_at,omitempty"`
}

type TableOrderItem struct {
	ProductId   string  `json:"product_id" bson:"product_id"`
	ProductName string  `json:"product_name" bson:"product_name"`
	Quantity    int     `json:"quantity" bson:"quantity"`
	UnitPrice   float64 `json:"unit_price" bson:"unit_price"`
	Notes       string  `json:"notes" bson:"notes"`
}

type TableOrder struct {
	Id        string          `json:"id" bson:"id" mapstructure:"id"`
	SessionId string          `json:"session_id" bson:"session_id" mapstructure:"session_id"`
	Items     []TableOrderItem `json:"items" bson:"items" mapstructure:"items"`
	Status    string          `json:"status" bson:"status" mapstructure:"status"`
	Subtotal  float64         `json:"subtotal" bson:"subtotal" mapstructure:"subtotal"`
	PlacedAt  time.Time       `json:"placed_at" bson:"placed_at" mapstructure:"placed_at"`
}

type QrInfo struct {
	TableLabel string `json:"table_label"`
	Token      string `json:"token"`
	Url        string `json:"url"`
	Host       string `json:"host"`
}
