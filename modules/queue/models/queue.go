package models

type QueueEntry struct {
	Id           string `json:"id" bson:"id"`
	CustomerName string `json:"customer_name" bson:"customer_name"`
	Phone        string `json:"phone" bson:"phone"`
	PartySize    int    `json:"party_size" bson:"party_size"`
	Position     int    `json:"position" bson:"position"`
	Status       string `json:"status" bson:"status"`
	EstimatedMin int    `json:"estimated_min" bson:"estimated_min"`
	AddedAt      string `json:"added_at" bson:"added_at"`
	Notes        string `json:"notes,omitempty" bson:"notes,omitempty"`
}
