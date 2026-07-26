package models

import "time"

type Table struct {
	Id        string    `json:"id" bson:"id,omitempty" mapstructure:"id,omitempty"`
	Number    int       `json:"number" bson:"number" mapstructure:"number"`
	Name      string    `json:"name" bson:"name" mapstructure:"name"`
	Capacity  int       `json:"capacity" bson:"capacity" mapstructure:"capacity"`
	Zone      string    `json:"zone" bson:"zone" mapstructure:"zone"`
	Status    string    `json:"status" bson:"status" mapstructure:"status"`
	QRCode    string    `json:"qr_code" bson:"qr_code" mapstructure:"qr_code"`
	OrderId   string    `json:"order_id" bson:"order_id" mapstructure:"order_id"`
	BranchId  string    `json:"branch_id" bson:"branch_id" mapstructure:"branch_id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at" mapstructure:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" mapstructure:"updated_at"`
}
