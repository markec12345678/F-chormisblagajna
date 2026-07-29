package models

import "time"

type OnlineOrderItem struct {
	ProductId   string  `json:"product_id" bson:"product_id" mapstructure:"product_id"`
	ProductName string  `json:"product_name" bson:"product_name" mapstructure:"product_name"`
	Quantity    float64 `json:"quantity" bson:"quantity" mapstructure:"quantity"`
	Price       float64 `json:"price" bson:"price" mapstructure:"price"`
	Comment     string  `json:"comment" bson:"comment" mapstructure:"comment"`
}

type OnlineOrder struct {
	Id            string            `json:"id" bson:"id" mapstructure:"id"`
	DisplayId     string            `json:"display_id" bson:"display_id" mapstructure:"display_id"`
	CustomerName  string            `json:"customer_name" bson:"customer_name" mapstructure:"customer_name"`
	CustomerPhone string            `json:"customer_phone" bson:"customer_phone" mapstructure:"customer_phone"`
	CustomerEmail string            `json:"customer_email" bson:"customer_email" mapstructure:"customer_email"`
	Items         []OnlineOrderItem `json:"items" bson:"items" mapstructure:"items"`
	Subtotal      float64           `json:"subtotal" bson:"subtotal" mapstructure:"subtotal"`
	DeliveryFee   float64           `json:"delivery_fee" bson:"delivery_fee" mapstructure:"delivery_fee"`
	Total         float64           `json:"total" bson:"total" mapstructure:"total"`
	OrderType     string            `json:"order_type" bson:"order_type" mapstructure:"order_type"` // delivery, takeaway, dine_in
	DeliveryAddr  string            `json:"delivery_addr" bson:"delivery_addr" mapstructure:"delivery_addr"`
	Status        string            `json:"status" bson:"status" mapstructure:"status"` // pending, confirmed, preparing, ready, delivered, cancelled
	Notes         string            `json:"notes" bson:"notes" mapstructure:"notes"`
	CreatedAt     time.Time         `json:"created_at" bson:"created_at" mapstructure:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at" bson:"updated_at" mapstructure:"updated_at"`
}

type MenuProduct struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	ImageURL  string  `json:"image_url"`
	Category  string  `json:"category"`
	Available bool    `json:"available"`
}

type MenuCategory struct {
	Id       string         `json:"id"`
	Name     string         `json:"name"`
	Products []MenuProduct  `json:"products"`
}
