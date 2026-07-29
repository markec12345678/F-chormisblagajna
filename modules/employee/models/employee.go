package models

import "time"

type EmployeePerformance struct {
	EmployeeId    string  `json:"employee_id"`
	EmployeeName  string  `json:"employee_name"`
	TotalSales    float64 `json:"total_sales"`
	OrderCount    int     `json:"order_count"`
	AvgOrderValue float64 `json:"avg_order_value"`
	TotalTips     float64 `json:"totalTips"`
	TotalHours    float64 `json:"total_hours"`
	SalesPerHour  float64 `json:"sales_per_hour"`
	TopProducts   []ProductStat `json:"top_products"`
	Rank          int     `json:"rank"`
}

type ProductStat struct {
	ProductId   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Revenue     float64 `json:"revenue"`
}

type PerformanceSummary struct {
	TotalEmployees int                `json:"totalEmployees"`
	TopPerformers  []EmployeePerformance `json:"top_performers"`
	TotalRevenue   float64            `json:"total_revenue"`
	AvgSalesPerHour float64           `json:"avg_sales_per_hour"`
}

type EmployeeShift struct {
	Id         string     `json:"id"`
	EmployeeId string     `json:"employee_id"`
	StartTime  time.Time  `json:"start_time"`
	EndTime    *time.Time `json:"end_time"`
	Status     string     `json:"status"` // active, completed, cancelled
}
