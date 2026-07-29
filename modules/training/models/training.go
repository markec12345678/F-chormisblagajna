package models

import "time"

type TrainingSession struct {
	Id         string    `json:"id" bson:"id" mapstructure:"id"`
	UserId     string    `json:"user_id" bson:"user_id" mapstructure:"user_id"`
	Module     string    `json:"module" bson:"module" mapstructure:"module"` // cashier, kitchen, inventory, admin
	StartedAt  time.Time `json:"started_at" bson:"started_at" mapstructure:"started_at"`
	EndedAt    *time.Time `json:"ended_at,omitempty" bson:"ended_at,omitempty" mapstructure:"ended_at,omitempty"`
	Score      int       `json:"score" bson:"score" mapstructure:"score"`
	MaxScore   int       `json:"max_score" bson:"max_score" mapstructure:"max_score"`
	Completed  bool      `json:"completed" bson:"completed" mapstructure:"completed"`
	StepsDone  int       `json:"steps_done" bson:"steps_done" mapstructure:"steps_done"`
	TotalSteps int       `json:"total_steps" bson:"total_steps" mapstructure:"total_steps"`
}

type TrainingStep struct {
	Id          string `json:"id"`
	Module      string `json:"module"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Action      string `json:"action"`
	ExpectedOut string `json:"expected_out"`
}

type TrainingModule struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Steps       int    `json:"steps"`
}

type TrainingProgress struct {
	SessionId     string  `json:"session_id"`
	Module        string  `json:"module"`
	StartedAt     string  `json:"started_at"`
	StepsDone     int     `json:"steps_done"`
	TotalSteps    int     `json:"total_steps"`
	Score         int     `json:"score"`
	MaxScore      int     `json:"max_score"`
	CompletionPct float64 `json:"completion_pct"`
	Completed     bool    `json:"completed"`
}
