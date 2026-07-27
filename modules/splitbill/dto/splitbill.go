package dto

type CreateSplitBillRequest struct {
	OrderId       string      `json:"order_id"`
	SplitType     string      `json:"split_type"`
	SplitCount    int         `json:"split_count,omitempty"`
	CustomAmounts []float64   `json:"custom_amounts,omitempty"`
	ItemSplits    []ItemSplit `json:"item_splits,omitempty"`
}

type ItemSplit struct {
	ItemId    string `json:"item_id"`
	PartIndex int    `json:"part_index"`
}

type PaySplitPartRequest struct {
	PartId        string  `json:"part_id"`
	PaymentMethod string  `json:"payment_method"`
	Amount        float64 `json:"amount"`
}
