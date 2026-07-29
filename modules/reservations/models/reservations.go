package models

import "time"

type Reservation struct {
	Id              string    `json:"id" bson:"id" mapstructure:"id"`
	CustomerName    string    `json:"customer_name" bson:"customer_name" mapstructure:"customer_name"`
	CustomerPhone   string    `json:"customer_phone" bson:"customer_phone" mapstructure:"customer_phone"`
	CustomerEmail   string    `json:"customer_email" bson:"customer_email" mapstructure:"customer_email"`
	GuestCount      int       `json:"guest_count" bson:"guest_count" mapstructure:"guest_count"`
	ReservationDate string    `json:"reservation_date" bson:"reservation_date" mapstructure:"reservation_date"`
	ReservationTime string    `json:"reservation_time" bson:"reservation_time" mapstructure:"reservation_time"`
	Notes           string    `json:"notes" bson:"notes" mapstructure:"notes"`
	Status          string    `json:"status" bson:"status" mapstructure:"status"`
	TableAssignment string    `json:"table_assignment,omitempty" bson:"table_assignment,omitempty" mapstructure:"table_assignment,omitempty"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at" mapstructure:"created_at"`
}

type ReservationSlot struct {
	Date      string `json:"date"`
	Time      string `json:"time"`
	Available int    `json:"available"`
	Total     int    `json:"total"`
}
