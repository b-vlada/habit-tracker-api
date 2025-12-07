package models

import "time"

type Goal struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	TargetDate  time.Time `json:"target_date"`
	CreatedAt   time.Time `json:"created_at"`
	Completed   bool      `json:"completed"`
	CompletedAt time.Time `json:"completed_at"`
	HabitIDs    []int     `json:"habit_ids"`
}
