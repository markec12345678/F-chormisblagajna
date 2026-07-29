package dto

type CreateTransferRequest struct {
	MaterialID   string  `json:"material_id"`
	MaterialName string  `json:"material_name"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
	FromBranchID string  `json:"from_branch_id"`
	ToBranchID   string  `json:"to_branch_id"`
	Notes        string  `json:"notes"`
}

type UpdateTransferStatusRequest struct {
	Status string `json:"status"`
}
