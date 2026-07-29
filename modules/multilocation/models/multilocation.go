package models

import "time"

type BranchStats struct {
	BranchID       string    `json:"branch_id" bson:"branch_id"`
	BranchName     string    `json:"branch_name" bson:"branch_name"`
	TodayRevenue   float64   `json:"today_revenue" bson:"today_revenue"`
	TodayOrders    int       `json:"today_orders" bson:"today_orders"`
	AvgOrderValue  float64   `json:"avg_order_value" bson:"avg_order_value"`
	WeekRevenue    float64   `json:"week_revenue" bson:"week_revenue"`
	WeekOrders     int       `json:"week_orders" bson:"week_orders"`
	MonthRevenue   float64   `json:"month_revenue" bson:"month_revenue"`
	MonthOrders    int       `json:"month_orders" bson:"month_orders"`
	ActiveStaff    int       `json:"active_staff" bson:"active_staff"`
	TablesOccupied int       `json:"tables_occupied" bson:"tables_occupied"`
	TotalTables    int       `json:"total_tables" bson:"total_tables"`
	Status         string    `json:"status" bson:"status"`
	LastActivity   time.Time `json:"last_activity" bson:"last_activity"`
}

type LocationDashboard struct {
	TotalRevenue   float64      `json:"total_revenue" bson:"total_revenue"`
	TotalOrders    int          `json:"total_orders" bson:"total_orders"`
	TotalBranches  int          `json:"total_branches" bson:"total_branches"`
	AvgOrderValue  float64      `json:"avg_order_value" bson:"avg_order_value"`
	Branches       []BranchStats `json:"branches" bson:"branches"`
	GeneratedAt    time.Time    `json:"generated_at" bson:"generated_at"`
}

type BranchComparison struct {
	Metric    string             `json:"metric" bson:"metric"`
	Branches  []BranchMetricValue `json:"branches" bson:"branches"`
}

type BranchMetricValue struct {
	BranchID   string  `json:"branch_id" bson:"branch_id"`
	BranchName string  `json:"branch_name" bson:"branch_name"`
	Value      float64 `json:"value" bson:"value"`
}
