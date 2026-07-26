package dto

type CreateTableRequest struct {
	Number   int    `json:"number"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
	Zone     string `json:"zone"`
	BranchId string `json:"branch_id"`
}

type UpdateTableRequest struct {
	Name     string `json:"name,omitempty"`
	Capacity int    `json:"capacity,omitempty"`
	Zone     string `json:"zone,omitempty"`
	Status   string `json:"status,omitempty"`
}
