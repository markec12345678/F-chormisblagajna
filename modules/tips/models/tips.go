package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Tip struct {
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	OrderID     string        `json:"order_id" bson:"order_id"`
	EmployeeID  string        `json:"employee_id" bson:"employee_id"`
	EmployeeName string       `json:"employee_name" bson:"employee_name"`
	Amount      float64       `json:"amount" bson:"amount"`
	PaymentMethod string      `json:"payment_method" bson:"payment_method"`
	BranchID    string        `json:"branch_id" bson:"branch_id"`
	Date        time.Time     `json:"date" bson:"date"`
	CreatedAt   time.Time     `json:"created_at" bson:"created_at"`
}

type TipSummary struct {
	EmployeeID    string  `json:"employee_id" bson:"employee_id"`
	EmployeeName  string  `json:"employee_name" bson:"employee_name"`
	TotalTips     float64 `json:"total_tips" bson:"total_tips"`
	TipCount      int     `json:"tip_count" bson:"tip_count"`
	AverageTip    float64 `json:"average_tip" bson:"average_tip"`
	Period        string  `json:"period" bson:"period"`
}

type TipPayout struct {
	ID            bson.ObjectID `json:"id" bson:"_id,omitempty"`
	EmployeeID    string        `json:"employee_id" bson:"employee_id"`
	EmployeeName  string        `json:"employee_name" bson:"employee_name"`
	Amount        float64       `json:"amount" bson:"amount"`
	PayoutDate    time.Time     `json:"payout_date" bson:"payout_date"`
	PayoutMethod  string        `json:"payout_method" bson:"payout_method"`
	Notes         string        `json:"notes" bson:"notes"`
	CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
}
