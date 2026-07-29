package dto

type CreateLoyaltyAccountRequest struct {
	CustomerId string `json:"customer_id"`
	BranchId   string `json:"branch_id"`
}

type RedeemPointsRequest struct {
	CustomerId string `json:"customer_id"`
	Points     int    `json:"points"`
	OrderId    string `json:"order_id"`
}

type AdjustPointsRequest struct {
	CustomerId  string `json:"customer_id"`
	Points      int    `json:"points"`
	Description string `json:"description"`
}

type SetTierRequest struct {
	CustomerId string `json:"customer_id"`
	Tier       string `json:"tier"`
}
