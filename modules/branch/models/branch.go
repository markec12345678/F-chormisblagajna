package models

import "time"

type Branch struct {
	Id        string    `json:"id" bson:"id,omitempty" mapstructure:"id,omitempty"`
	Name      string    `json:"name" bson:"name" mapstructure:"name"`
	Address   string    `json:"address" bson:"address" mapstructure:"address"`
	Phone     string    `json:"phone" bson:"phone" mapstructure:"phone"`
	Email     string    `json:"email" bson:"email" mapstructure:"email"`
	IsActive  bool      `json:"is_active" bson:"is_active" mapstructure:"is_active"`
	CreatedAt time.Time `json:"created_at" bson:"created_at" mapstructure:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" mapstructure:"updated_at"`
}
