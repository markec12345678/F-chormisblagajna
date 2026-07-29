package models

import "time"

type InventoryAlertRule struct {
	Id            string  `json:"id" bson:"id" mapstructure:"id"`
	MaterialId    string  `json:"material_id" bson:"material_id" mapstructure:"material_id"`
	MaterialName  string  `json:"material_name" bson:"material_name" mapstructure:"material_name"`
	ThresholdLow  float64 `json:"threshold_low" bson:"threshold_low" mapstructure:"threshold_low"`
	ThresholdCrit float64 `json:"threshold_critical" bson:"threshold_critical" mapstructure:"threshold_critical"`
	NotifyEmail   bool    `json:"notify_email" bson:"notify_email" mapstructure:"notify_email"`
	IsActive      bool    `json:"is_active" bson:"is_active" mapstructure:"is_active"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at" mapstructure:"created_at"`
}

type InventoryAlert struct {
	Id           string    `json:"id" bson:"id" mapstructure:"id"`
	RuleId       string    `json:"rule_id" bson:"rule_id" mapstructure:"rule_id"`
	MaterialId   string    `json:"material_id" bson:"material_id" mapstructure:"material_id"`
	MaterialName string    `json:"material_name" bson:"material_name" mapstructure:"material_name"`
	CurrentQty   float64   `json:"current_qty" bson:"current_qty" mapstructure:"current_qty"`
	Threshold    float64   `json:"threshold" bson:"threshold" mapstructure:"threshold"`
	Severity     string    `json:"severity" bson:"severity" mapstructure:"severity"` // low, critical
	IsRead       bool      `json:"is_read" bson:"is_read" mapstructure:"is_read"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at" mapstructure:"created_at"`
}

type AlertSummary struct {
	TotalActive   int `json:"total_active"`
	UnreadCount   int `json:"unread_count"`
	CriticalCount int `json:"critical_count"`
	LowCount      int `json:"low_count"`
}
