package models

import "time"

// Task represents a task in the system
type Task struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Status      string    `json:"status"`
}

// TaskDetail represents detailed information about a task
type TaskDetail struct {
	Task
	Priority    string   `json:"priority"`
	Tags        []string `json:"tags,omitempty"`
	EstimatedHours int   `json:"estimated_hours,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}
