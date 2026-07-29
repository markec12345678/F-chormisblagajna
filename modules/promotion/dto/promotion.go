package dto

type CreatePromotionRequest struct {
	Name           string   `json:"name"`
	Code           string   `json:"code"`
	Type           string   `json:"type"`
	Value          float64  `json:"value"`
	MinOrder       float64  `json:"min_order"`
	MaxDiscount    float64  `json:"max_discount"`
	StartDate      string   `json:"start_date"`
	EndDate        string   `json:"end_date"`
	UsageLimit     int      `json:"usage_limit"`
	ApplicableDays []string `json:"applicable_days"`
	HappyHourStart string   `json:"happy_hour_start"`
	HappyHourEnd   string   `json:"happy_hour_end"`
	IsActive       bool     `json:"is_active"`
}

type UpdatePromotionRequest struct {
	Name           string   `json:"name,omitempty"`
	Code           string   `json:"code,omitempty"`
	Type           string   `json:"type,omitempty"`
	Value          *float64 `json:"value,omitempty"`
	MinOrder       *float64 `json:"min_order,omitempty"`
	MaxDiscount    *float64 `json:"max_discount,omitempty"`
	StartDate      string   `json:"start_date,omitempty"`
	EndDate        string   `json:"end_date,omitempty"`
	UsageLimit     *int     `json:"usage_limit,omitempty"`
	ApplicableDays []string `json:"applicable_days,omitempty"`
	HappyHourStart string   `json:"happy_hour_start,omitempty"`
	HappyHourEnd   string   `json:"happy_hour_end,omitempty"`
	IsActive       *bool    `json:"is_active,omitempty"`
}
