package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type KitchenStation struct {
	ID       bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string        `json:"name" bson:"name"`
	BranchID string        `json:"branch_id" bson:"branch_id"`
	Active   bool          `json:"active" bson:"active"`
	CreatedAt time.Time    `json:"created_at" bson:"created_at"`
}

type KitchenOrderItem struct {
	OrderID    string `json:"order_id" bson:"order_id"`
	OrderItemIndex int `json:"order_item_index" bson:"order_item_index"`
	ProductName string `json:"product_name" bson:"product_name"`
	Quantity   int    `json:"quantity" bson:"quantity"`
	Notes      string `json:"notes" bson:"notes"`
	Status     string `json:"status" bson:"status"`
	Station    string `json:"station" bson:"station"`
	Priority   string `json:"priority" bson:"priority"`
	StartedAt  *time.Time `json:"started_at" bson:"started_at"`
	ReadyAt    *time.Time `json:"ready_at" bson:"ready_at"`
}
