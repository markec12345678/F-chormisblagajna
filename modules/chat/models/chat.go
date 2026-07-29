package models

import "time"

type ChatMessage struct {
	Id        string    `json:"id" bson:"id" mapstructure:"id"`
	Channel   string    `json:"channel" bson:"channel" mapstructure:"channel"`
	Sender    string    `json:"sender" bson:"sender" mapstructure:"sender"`
	SenderId  string    `json:"sender_id" bson:"sender_id" mapstructure:"sender_id"`
	Content   string    `json:"content" bson:"content" mapstructure:"content"`
	CreatedAt time.Time `json:"created_at" bson:"created_at" mapstructure:"created_at"`
}

type ChatChannel struct {
	Id          string `json:"id" bson:"id" mapstructure:"id"`
	Name        string `json:"name" bson:"name" mapstructure:"name"`
	Description string `json:"description" bson:"description" mapstructure:"description"`
	IsDefault   bool   `json:"is_default" bson:"is_default" mapstructure:"is_default"`
}

type ChatSummary struct {
	TotalMessages int              `json:"total_messages"`
	Channels      []ChatChannel    `json:"channels"`
	OnlineUsers   []OnlineUser     `json:"online_users"`
}

type OnlineUser struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}
