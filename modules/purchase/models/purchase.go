package models

import "time"

type PurchaseOrder struct {
	Id            string    `json:"id" bson:"id"`
	SupplierId    string    `json:"supplier_id" bson:"supplier_id"`
	SupplierName  string    `json:"supplier_name" bson:"supplier_name"`
	Items         []POItem  `json:"items" bson:"items"`
	TotalCost     float64   `json:"total_cost" bson:"total_cost"`
	Status        string    `json:"status" bson:"status"`
	OrderedAt     time.Time `json:"ordered_at" bson:"ordered_at"`
	ExpectedAt    string    `json:"expected_at,omitempty" bson:"expected_at,omitempty"`
	ReceivedAt    *time.Time `json:"received_at,omitempty" bson:"received_at,omitempty"`
	Notes         string    `json:"notes,omitempty" bson:"notes,omitempty"`
}

type POItem struct {
	MaterialId   string  `json:"material_id" bson:"material_id"`
	MaterialName string  `json:"material_name" bson:"material_name"`
	Quantity     float64 `json:"quantity" bson:"quantity"`
	UnitPrice    float64 `json:"unit_price" bson:"unit_price"`
	TotalPrice   float64 `json:"total_price" bson:"total_price"`
}
