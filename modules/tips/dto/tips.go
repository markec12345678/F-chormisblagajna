package dto

type RecordTipRequest struct {
	OrderID       string  `json:"order_id"`
	EmployeeID    string  `json:"employee_id"`
	EmployeeName  string  `json:"employee_name"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
	BranchID      string  `json:"branch_id"`
}

type PayoutTipsRequest struct {
	EmployeeID   string  `json:"employee_id"`
	EmployeeName string  `json:"employee_name"`
	Amount       float64 `json:"amount"`
	PayoutMethod string  `json:"payout_method"`
	Notes        string  `json:"notes"`
}
