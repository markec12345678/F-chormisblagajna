package models

import "time"

type WasteEntry struct {
	Id           string    `json:"id" bson:"id" mapstructure:"id"`
	MaterialId   string    `json:"material_id" bson:"material_id" mapstructure:"material_id"`
	MaterialName string    `json:"material_name" bson:"material_name" mapstructure:"material_name"`
	Quantity     float64   `json:"quantity" bson:"quantity" mapstructure:"quantity"`
	Unit         string    `json:"unit" bson:"unit" mapstructure:"unit"`
	Reason       string    `json:"reason" bson:"reason" mapstructure:"reason"` // expired, damaged, overcooked, other
	Cost         float64   `json:"cost" bson:"cost" mapstructure:"cost"`
	Date         time.Time `json:"date" bson:"date" mapstructure:"date"`
	RecordedBy   string    `json:"recorded_by" bson:"recorded_by" mapstructure:"recorded_by"`
	Notes        string    `json:"notes" bson:"notes" mapstructure:"notes"`
}

type WasteSummary struct {
	TotalWasteCost  float64            `json:"total_waste_cost"`
	TotalEntries    int                `json:"total_entries"`
	ByReason        []ReasonSummary    `json:"by_reason"`
	ByMaterial      []MaterialSummary  `json:"by_material"`
	DailyWaste      []DailyWaste       `json:"daily_waste"`
}

type ReasonSummary struct {
	Reason string  `json:"reason"`
	Total  float64 `json:"total"`
	Count  int     `json:"count"`
}

type MaterialSummary struct {
	MaterialId   string  `json:"material_id"`
	MaterialName string  `json:"material_name"`
	TotalCost    float64 `json:"total_cost"`
	TotalQty     float64 `json:"total_qty"`
	Count        int     `json:"count"`
}

type DailyWaste struct {
	Date  string  `json:"date"`
	Total float64 `json:"total"`
}
