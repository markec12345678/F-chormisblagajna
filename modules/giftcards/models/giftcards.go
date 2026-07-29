package models

type GiftCard struct {
	Id         string  `json:"id" bson:"id"`
	Code       string  `json:"code" bson:"code"`
	Balance    float64 `json:"balance" bson:"balance"`
	InitialAmt float64 `json:"initial_amt" bson:"initial_amt"`
	IssuedTo   string  `json:"issued_to" bson:"issued_to"`
	IssuedAt   string  `json:"issued_at" bson:"issued_at"`
	ExpiresAt  string  `json:"expires_at,omitempty" bson:"expires_at,omitempty"`
	Active     bool    `json:"active" bson:"active"`
}
