package dto

type CreateReservationRequest struct {
	CustomerName  string `json:"customer_name"`
	CustomerPhone string `json:"customer_phone"`
	CustomerEmail string `json:"customer_email"`
	TableId       string `json:"table_id"`
	BranchId      string `json:"branch_id"`
	Date          string `json:"date"`
	Time          string `json:"time"`
	GuestCount    int    `json:"guest_count"`
	Status        string `json:"status"`
	Notes         string `json:"notes"`
}

type UpdateReservationRequest struct {
	CustomerName  string `json:"customer_name,omitempty"`
	CustomerPhone string `json:"customer_phone,omitempty"`
	CustomerEmail string `json:"customer_email,omitempty"`
	TableId       string `json:"table_id,omitempty"`
	BranchId      string `json:"branch_id,omitempty"`
	Date          string `json:"date,omitempty"`
	Time          string `json:"time,omitempty"`
	GuestCount    *int   `json:"guest_count,omitempty"`
	Status        string `json:"status,omitempty"`
	Notes         string `json:"notes,omitempty"`
}
