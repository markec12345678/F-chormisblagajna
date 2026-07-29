package models

import "time"

type AuditLogEntry struct {
	Id        string            `json:"id" bson:"id" mapstructure:"id"`
	Action    string            `json:"action" bson:"action" mapstructure:"action"` // create, update, delete, login, logout, fiscalize
	Resource  string            `json:"resource" bson:"resource" mapstructure:"resource"` // order, product, material, customer, settings, user, etc.
	ResourceId string           `json:"resource_id" bson:"resource_id" mapstructure:"resource_id"`
	UserId    string            `json:"user_id" bson:"user_id" mapstructure:"user_id"`
	Username  string            `json:"username" bson:"username" mapstructure:"username"`
	Details   map[string]string `json:"details" bson:"details" mapstructure:"details"`
	IpAddress string            `json:"ip_address" bson:"ip_address" mapstructure:"ip_address"`
	CreatedAt time.Time         `json:"created_at" bson:"created_at" mapstructure:"created_at"`
}

type AuditLogSummary struct {
	TotalEntries    int                `json:"total_entries"`
	ByAction        []ActionSummary    `json:"by_action"`
	ByResource      []ResourceSummary  `json:"by_resource"`
	ByUser          []UserSummary      `json:"by_user"`
}

type ActionSummary struct {
	Action string `json:"action"`
	Count  int    `json:"count"`
}

type ResourceSummary struct {
	Resource string `json:"resource"`
	Count    int    `json:"count"`
}

type UserSummary struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	Count    int    `json:"count"`
}
