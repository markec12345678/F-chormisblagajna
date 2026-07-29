package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type GiftCard struct {
	ID            bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Code          string        `json:"code" bson:"code"`
	InitialAmount float64       `json:"initial_amount" bson:"initial_amount"`
	CurrentAmount float64       `json:"current_amount" bson:"current_amount"`
	Status        string        `json:"status" bson:"status"`
	CustomerID    string        `json:"customer_id" bson:"customer_id"`
	CustomerName  string        `json:"customer_name" bson:"customer_name"`
	IssuedAt      time.Time     `json:"issued_at" bson:"issued_at"`
	ExpiryDate    *time.Time    `json:"expiry_date" bson:"expiry_date"`
	Notes         string        `json:"notes" bson:"notes"`
	CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at" bson:"updated_at"`
}

type GiftCardTransaction struct {
	ID            bson.ObjectID `json:"id" bson:"_id,omitempty"`
	GiftCardID    bson.ObjectID `json:"gift_card_id" bson:"gift_card_id"`
	GiftCardCode  string        `json:"gift_card_code" bson:"gift_card_code"`
	Type          string        `json:"type" bson:"type"`
	Amount        float64       `json:"amount" bson:"amount"`
	BalanceAfter  float64       `json:"balance_after" bson:"balance_after"`
	OrderID       string        `json:"order_id" bson:"order_id"`
	Notes         string        `json:"notes" bson:"notes"`
	CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
}
