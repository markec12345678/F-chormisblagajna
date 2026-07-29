package models

import "time"

type SalesReport struct {
	Period        string  `json:"period" bson:"period"`
	TotalRevenue  float64 `json:"total_revenue" bson:"total_revenue"`
	TotalOrders   int     `json:"total_orders" bson:"total_orders"`
	AverageOrder  float64 `json:"average_order" bson:"average_order"`
	TotalItems    int     `json:"total_items" bson:"total_items"`
	RefundAmount  float64 `json:"refund_amount" bson:"refund_amount"`
	NetRevenue    float64 `json:"net_revenue" bson:"net_revenue"`
	TopProducts   []ProductStat `json:"top_products" bson:"top_products"`
}

type ProductStat struct {
	Name     string  `json:"name" bson:"name"`
	Quantity int     `json:"quantity" bson:"quantity"`
	Revenue  float64 `json:"revenue" bson:"revenue"`
}

type InventoryReport struct {
	TotalMaterials  int     `json:"total_materials" bson:"total_materials"`
	LowStockCount   int     `json:"low_stock_count" bson:"low_stock_count"`
	OutOfStockCount int     `json:"out_of_stock_count" bson:"out_of_stock_count"`
	TotalValue      float64 `json:"total_value" bson:"total_value"`
	LowStockItems   []LowStockItem `json:"low_stock_items" bson:"low_stock_items"`
}

type LowStockItem struct {
	Name     string  `json:"name" bson:"name"`
	Quantity float64 `json:"quantity" bson:"quantity"`
	Unit     string  `json:"unit" bson:"unit"`
	Value    float64 `json:"value" bson:"value"`
}

type TaxReport struct {
	Period        string      `json:"period" bson:"period"`
	TotalTax      float64     `json:"total_tax" bson:"total_tax"`
	TaxByRate     []TaxEntry  `json:"tax_by_rate" bson:"tax_by_rate"`
	TotalRevenue  float64     `json:"total_revenue" bson:"total_revenue"`
}

type TaxEntry struct {
	Rate    float64 `json:"rate" bson:"rate"`
	Label   string  `json:"label" bson:"label"`
	Tax     float64 `json:"tax" bson:"tax"`
	Base    float64 `json:"base" bson:"base"`
}

type DashboardStats struct {
	TodayRevenue   float64 `json:"today_revenue" bson:"today_revenue"`
	TodayOrders    int     `json:"today_orders" bson:"today_orders"`
	TodayCustomers int     `json:"today_customers" bson:"today_customers"`
	WeekRevenue    float64 `json:"week_revenue" bson:"week_revenue"`
	MonthRevenue   float64 `json:"month_revenue" bson:"month_revenue"`
	AverageOrder   float64 `json:"average_order" bson:"average_order"`
	TopProduct     string  `json:"top_product" bson:"top_product"`
	GeneratedAt    time.Time `json:"generated_at" bson:"generated_at"`
}
