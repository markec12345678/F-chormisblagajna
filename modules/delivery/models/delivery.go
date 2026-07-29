package models

import "time"

type DeliveryZone struct {
	Id       string  `json:"id" bson:"id"`
	Name     string  `json:"name" bson:"name"`
	Fee      float64 `json:"fee" bson:"fee"`
	MinOrder float64 `json:"min_order" bson:"min_order"`
	Active   bool    `json:"active" bson:"active"`
}

type DeliveryOrder struct {
	Id            string    `json:"id" bson:"id"`
	OrderId       string    `json:"order_id" bson:"order_id"`
	CustomerName  string    `json:"customer_name" bson:"customer_name"`
	CustomerPhone string    `json:"customer_phone" bson:"customer_phone"`
	Address       string    `json:"address" bson:"address"`
	ZoneId        string    `json:"zone_id" bson:"zone_id"`
	DeliveryFee   float64   `json:"delivery_fee" bson:"delivery_fee"`
	Status        string    `json:"status" bson:"status"`
	CourierId     string    `json:"courier_id,omitempty" bson:"courier_id,omitempty"`
	EstimatedMin  int       `json:"estimated_min" bson:"estimated_min"`
	PlacedAt      time.Time `json:"placed_at" bson:"placed_at"`
	DeliveredAt   *time.Time `json:"delivered_at,omitempty" bson:"delivered_at,omitempty"`
}
