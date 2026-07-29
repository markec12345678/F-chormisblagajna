package models

import "time"

type Supplier struct {
	Id          string    `json:"id" bson:"id" mapstructure:"id"`
	Name        string    `json:"name" bson:"name" mapstructure:"name"`
	ContactName string    `json:"contact_name" bson:"contact_name" mapstructure:"contact_name"`
	Email       string    `json:"email" bson:"email" mapstructure:"email"`
	Phone       string    `json:"phone" bson:"phone" mapstructure:"phone"`
	Address     string    `json:"address" bson:"address" mapstructure:"address"`
	Website     string    `json:"website" bson:"website" mapstructure:"website"`
	Notes       string    `json:"notes" bson:"notes" mapstructure:"notes"`
	IsActive    bool      `json:"is_active" bson:"is_active" mapstructure:"is_active"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at" mapstructure:"created_at"`
}

type SupplierOrder struct {
	Id           string    `json:"id" bson:"id" mapstructure:"id"`
	SupplierId   string    `json:"supplier_id" bson:"supplier_id" mapstructure:"supplier_id"`
	SupplierName string    `json:"supplier_name" bson:"supplier_name" mapstructure:"supplier_name"`
	OrderDate    time.Time `json:"order_date" bson:"order_date" mapstructure:"order_date"`
	TotalAmount  float64   `json:"total_amount" bson:"total_amount" mapstructure:"total_amount"`
	Status       string    `json:"status" bson:"status" mapstructure:"status"` // pending, delivered, cancelled
	Items        []SupplierOrderItem `json:"items" bson:"items" mapstructure:"items"`
}

type SupplierOrderItem struct {
	MaterialId   string  `json:"material_id" bson:"material_id" mapstructure:"material_id"`
	MaterialName string  `json:"material_name" bson:"material_name" mapstructure:"material_name"`
	Quantity     float64 `json:"quantity" bson:"quantity" mapstructure:"quantity"`
	UnitPrice    float64 `json:"unit_price" bson:"unit_price" mapstructure:"unit_price"`
}
