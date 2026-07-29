package models

type Campaign struct {
	Id          string `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Type        string `json:"type" bson:"type"`
	StartDate   string `json:"start_date" bson:"start_date"`
	EndDate     string `json:"end_date" bson:"end_date"`
	DiscountPct float64 `json:"discount_pct,omitempty" bson:"discount_pct,omitempty"`
	TargetAud   string `json:"target_audience" bson:"target_audience"`
	Active      bool   `json:"active" bson:"active"`
	CreatedAt   string `json:"created_at" bson:"created_at"`
}

type CampaignStats struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	TotalSent   int    `json:"total_sent"`
	TotalOpened int    `json:"total_opened"`
	TotalRedeemed int  `json:"total_redeemed"`
	Revenue     float64 `json:"revenue"`
}
