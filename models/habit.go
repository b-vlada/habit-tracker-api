package models

import "time"

type Habit struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Frequency   string    `json:"frequency"`
	CreatedAt   time.Time `json:"created_at"`
	Completed   bool      `json:"completed"`
}
