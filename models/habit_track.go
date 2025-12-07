package models

import "time"

type HabitTrack struct {
	ID        int       `json:"id"`
	HabitID   int       `json:"habit_id"`
	Date      time.Time `json:"date"`
	Completed bool      `json:"completed"`
	Notes     string    `json:"notes"`
}
