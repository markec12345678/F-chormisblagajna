package models

import "time"

type Expense struct {
	Id          string    `json:"id" bson:"id" mapstructure:"id"`
	Description string    `json:"description" bson:"description" mapstructure:"description"`
	Amount      float64   `json:"amount" bson:"amount" mapstructure:"amount"`
	Category    string    `json:"category" bson:"category" mapstructure:"category"`
	Date        time.Time `json:"date" bson:"date" mapstructure:"date"`
	Receipt     string    `json:"receipt" bson:"receipt" mapstructure:"receipt"`
	Notes       string    `json:"notes" bson:"notes" mapstructure:"notes"`
	CreatedBy   string    `json:"created_by" bson:"created_by" mapstructure:"created_by"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at" mapstructure:"created_at"`
}

type ExpenseCategory struct {
	Id          string  `json:"id" bson:"id" mapstructure:"id"`
	Name        string  `json:"name" bson:"name" mapstructure:"name"`
	Budget      float64 `json:"budget" bson:"budget" mapstructure:"budget"`
	TotalSpent  float64 `json:"total_spent" bson:"total_spent" mapstructure:"total_spent"`
}

type ExpenseSummary struct {
	TotalExpenses float64           `json:"total_expenses"`
	ByCategory    []CategorySummary `json:"by_category"`
	MonthlyTotal  float64           `json:"monthly_total"`
}

type CategorySummary struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
	Count    int     `json:"count"`
}
