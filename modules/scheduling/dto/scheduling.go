package dto

type CreateShiftRequest struct {
	EmployeeId string `json:"employee_id"`
	BranchId   string `json:"branch_id"`
	Date       string `json:"date"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	Role       string `json:"role"`
	Status     string `json:"status"`
	Notes      string `json:"notes"`
}

type UpdateShiftRequest struct {
	EmployeeId string `json:"employee_id,omitempty"`
	BranchId   string `json:"branch_id,omitempty"`
	Date       string `json:"date,omitempty"`
	StartTime  string `json:"start_time,omitempty"`
	EndTime    string `json:"end_time,omitempty"`
	Role       string `json:"role,omitempty"`
	Status     string `json:"status,omitempty"`
	Notes      string `json:"notes,omitempty"`
}
