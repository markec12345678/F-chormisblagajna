package models

type MenuItemAnalysis struct {
	ProductId      string  `json:"product_id"`
	ProductName    string  `json:"product_name"`
	TotalSold      int     `json:"total_sold"`
	TotalRevenue   float64 `json:"total_revenue"`
	TotalCost      float64 `json:"total_cost"`
	ProfitMargin   float64 `json:"profit_margin"`
	ProfitPerItem  float64 `json:"profit_per_item"`
	PopularityRank int     `json:"popularity_rank"`
	ProfitRank     int     `json:"profit_rank"`
	Category       string  `json:"category"`
	Quadrant       string  `json:"quadrant"` // star, plowhorse, puzzle, dog
}

type MenuEngineeringSummary struct {
	TotalItems    int                `json:"total_items"`
	TotalRevenue  float64            `json:"total_revenue"`
	AvgProfit     float64            `json:"avg_profit"`
	TopStars      []MenuItemAnalysis `json:"top_stars"`
	TopPlowhorses []MenuItemAnalysis `json:"top_plowhorses"`
	TopPuzzles    []MenuItemAnalysis `json:"top_puzzles"`
	TopDogs       []MenuItemAnalysis `json:"top_dogs"`
}
