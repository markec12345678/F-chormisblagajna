package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type InventoryTransfer struct {
	ID             bson.ObjectID `json:"id" bson:"_id,omitempty"`
	MaterialID     string        `json:"material_id" bson:"material_id"`
	MaterialName   string        `json:"material_name" bson:"material_name"`
	Quantity       float64       `json:"quantity" bson:"quantity"`
	Unit           string        `json:"unit" bson:"unit"`
	FromBranchID   string        `json:"from_branch_id" bson:"from_branch_id"`
	FromBranchName string        `json:"from_branch_name" bson:"from_branch_name"`
	ToBranchID     string        `json:"to_branch_id" bson:"to_branch_id"`
	ToBranchName   string        `json:"to_branch_name" bson:"to_branch_name"`
	Status         string        `json:"status" bson:"status"`
	Notes          string        `json:"notes" bson:"notes"`
	CreatedBy      string        `json:"created_by" bson:"created_by"`
	CreatedAt      time.Time     `json:"created_at" bson:"created_at"`
	CompletedAt    *time.Time    `json:"completed_at" bson:"completed_at"`
}
