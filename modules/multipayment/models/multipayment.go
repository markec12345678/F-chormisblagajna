package models

import "time"

type PaymentPart struct {
	Id            string    `json:"id" bson:"id" mapstructure:"id"`
	OrderId       string    `json:"order_id" bson:"order_id" mapstructure:"order_id"`
	Amount        float64   `json:"amount" bson:"amount" mapstructure:"amount"`
	PaymentMethod string    `json:"payment_method" bson:"payment_method" mapstructure:"payment_method"` // cash, card, voucher, mobile, gift_card
	Reference     string    `json:"reference" bson:"reference" mapstructure:"reference"`
	Notes         string    `json:"notes" bson:"notes" mapstructure:"notes"`
	ReceivedBy    string    `json:"received_by" bson:"received_by" mapstructure:"received_by"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at" mapstructure:"created_at"`
}

type PaymentSummary struct {
	OrderId       string  `json:"order_id"`
	TotalDue      float64 `json:"total_due"`
	TotalPaid     float64 `json:"total_paid"`
	Remaining     float64 `json:"remaining"`
	IsFullyPaid   bool    `json:"is_fully_paid"`
	Payments      []PaymentPart `json:"payments"`
	MethodBreakdown []MethodAmount `json:"method_breakdown"`
}

type MethodAmount struct {
	Method string  `json:"method"`
	Total  float64 `json:"total"`
	Count  int     `json:"count"`
}

type DailyPayments struct {
	Date         string  `json:"date"`
	TotalCash    float64 `json:"total_cash"`
	TotalCard    float64 `json:"total_card"`
	TotalVoucher float64 `json:"total_voucher"`
	TotalMobile  float64 `json:"total_mobile"`
	TotalGift    float64 `json:"total_gift_card"`
	GrandTotal   float64 `json:"grand_total"`
	Count        int     `json:"count"`
}
