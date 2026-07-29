package dto

import "time"

type CreateGiftCardRequest struct {
	Code          string     `json:"code"`
	InitialAmount float64    `json:"initial_amount"`
	CustomerID    string     `json:"customer_id"`
	CustomerName  string     `json:"customer_name"`
	ExpiryDate    *time.Time `json:"expiry_date"`
	Notes         string     `json:"notes"`
}

type RedeemGiftCardRequest struct {
	Code   string  `json:"code"`
	Amount float64 `json:"amount"`
	OrderID string `json:"order_id"`
	Notes  string  `json:"notes"`
}
