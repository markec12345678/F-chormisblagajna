package models

type LoyaltyCard struct {
	Id           string `json:"id" bson:"id" mapstructure:"id"`
	CustomerId   string `json:"customer_id" bson:"customer_id" mapstructure:"customer_id"`
	CustomerName string `json:"customer_name" bson:"customer_name" mapstructure:"customer_name"`
	Points       int    `json:"points" bson:"points" mapstructure:"points"`
	Tier         string `json:"tier" bson:"tier" mapstructure:"tier"`
	TotalSpent   float64 `json:"total_spent" bson:"total_spent" mapstructure:"total_spent"`
	VisitCount   int    `json:"visit_count" bson:"visit_count" mapstructure:"visit_count"`
	Active       bool   `json:"active" bson:"active" mapstructure:"active"`
}

type Reward struct {
	Id           string `json:"id" bson:"id" mapstructure:"id"`
	Name         string `json:"name" bson:"name" mapstructure:"name"`
	Description  string `json:"description" bson:"description" mapstructure:"description"`
	PointsCost   int    `json:"points_cost" bson:"points_cost" mapstructure:"points_cost"`
	RewardType   string `json:"reward_type" bson:"reward_type" mapstructure:"reward_type"`
	DiscountPct  float64 `json:"discount_pct,omitempty" bson:"discount_pct,omitempty" mapstructure:"discount_pct,omitempty"`
	FreeItemId   string `json:"free_item_id,omitempty" bson:"free_item_id,omitempty" mapstructure:"free_item_id,omitempty"`
	Active       bool   `json:"active" bson:"active" mapstructure:"active"`
}

type Redemption struct {
	Id          string `json:"id" bson:"id" mapstructure:"id"`
	CardId      string `json:"card_id" bson:"card_id" mapstructure:"card_id"`
	RewardId    string `json:"reward_id" bson:"reward_id" mapstructure:"reward_id"`
	RewardName  string `json:"reward_name" bson:"reward_name" mapstructure:"reward_name"`
	PointsSpent int    `json:"points_spent" bson:"points_spent" mapstructure:"points_spent"`
	RedeemedAt  string `json:"redeemed_at" bson:"redeemed_at" mapstructure:"redeemed_at"`
}

type LoyaltySettings struct {
	PointsPerEuro    float64 `json:"points_per_euro"`
	EuroPerPoint     float64 `json:"euro_per_point"`
	WelcomePoints    int     `json:"welcome_points"`
	TierThresholds   map[string]float64 `json:"tier_thresholds"`
}
