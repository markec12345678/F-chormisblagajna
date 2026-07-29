package models

import "time"

type Shift struct {
	Id         string    `json:"id" bson:"id,omitempty" mapstructure:"id,omitempty"`
	EmployeeId string    `json:"employee_id" bson:"employee_id" mapstructure:"employee_id"`
	BranchId   string    `json:"branch_id" bson:"branch_id" mapstructure:"branch_id"`
	Date       string    `json:"date" bson:"date" mapstructure:"date"`
	StartTime  string    `json:"start_time" bson:"start_time" mapstructure:"start_time"`
	EndTime    string    `json:"end_time" bson:"end_time" mapstructure:"end_time"`
	Role       string    `json:"role" bson:"role" mapstructure:"role"`
	Status     string    `json:"status" bson:"status" mapstructure:"status"`
	Notes      string    `json:"notes" bson:"notes" mapstructure:"notes"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at" mapstructure:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at" mapstructure:"updated_at"`
}
