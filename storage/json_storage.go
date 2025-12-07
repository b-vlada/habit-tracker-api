package storage

import (
	"encoding/json"
	"habit-tracker-api/models"
	"os"
	"sync"
	"time"
)

type JSONStorage struct {
	filename    string
	mu          sync.RWMutex
	Habits      map[int]models.Habit      `json:"habits"`
	Goals       map[int]models.Goal       `json:"goals"`
	HabitTracks map[int]models.HabitTrack `json:"habit_tracks"`
	NextHabitID int                       `json:"next_habit_id"`
	NextGoalID  int                       `json:"next_goal_id"`
	NextTrackID int                       `json:"next_track_id"`
}

func NewJSONStorage(filename string) (*JSONStorage, error) {
	storage := &JSONStorage{
		filename:    filename,
		Habits:      make(map[int]models.Habit),
		Goals:       make(map[int]models.Goal),
		HabitTracks: make(map[int]models.HabitTrack),
		NextHabitID: 1,
		NextGoalID:  1,
		NextTrackID: 1,
	}

	if err := storage.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return storage, nil
}

func (s *JSONStorage) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, s)
}

func (s *JSONStorage) save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filename, data, 0644)
}

func (s *JSONStorage) GetAllHabits() ([]models.Habit, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var habits []models.Habit
	for _, habit := range s.Habits {
		habits = append(habits, habit)
	}

	return habits, nil
}

func (s *JSONStorage) GetHabitByID(id int) (*models.Habit, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	habit, exists := s.Habits[id]
	if !exists {
		return nil, nil
	}

	return &habit, nil
}

func (s *JSONStorage) CreateHabit(habit *models.Habit) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	habit.ID = s.NextHabitID
	s.NextHabitID++
	s.Habits[habit.ID] = *habit

	return s.save()
}

func (s *JSONStorage) UpdateHabit(id int, habit *models.Habit) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Habits[id]; !exists {
		return nil
	}

	habit.ID = id
	s.Habits[id] = *habit

	return s.save()
}

func (s *JSONStorage) DeleteHabit(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Habits[id]; !exists {
		return nil
	}

	delete(s.Habits, id)
	return s.save()
}

func (s *JSONStorage) CompleteHabit(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	habit, exists := s.Habits[id]
	if !exists {
		return nil
	}

	if habit.Completed {
		return nil
	}

	habit.Completed = true
	s.Habits[id] = habit

	track := models.HabitTrack{
		ID:        s.NextTrackID,
		HabitID:   id,
		Date:      time.Now(),
		Completed: true,
		Notes:     "Marked as completed via API",
	}
	s.NextTrackID++
	s.HabitTracks[track.ID] = track

	return s.save()
}

func (s *JSONStorage) GetAllGoals() ([]models.Goal, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var goals []models.Goal
	for _, goal := range s.Goals {
		goals = append(goals, goal)
	}

	return goals, nil
}

func (s *JSONStorage) GetGoalByID(id int) (*models.Goal, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	goal, exists := s.Goals[id]
	if !exists {
		return nil, nil
	}

	return &goal, nil
}

func (s *JSONStorage) CreateGoal(goal *models.Goal) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	goal.ID = s.NextGoalID
	s.NextGoalID++
	s.Goals[goal.ID] = *goal

	return s.save()
}

func (s *JSONStorage) UpdateGoal(id int, goal *models.Goal) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Goals[id]; !exists {
		return nil
	}

	goal.ID = id
	s.Goals[id] = *goal

	return s.save()
}

func (s *JSONStorage) DeleteGoal(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Goals[id]; !exists {
		return nil
	}

	delete(s.Goals, id)
	return s.save()
}

func (s *JSONStorage) CompleteGoal(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	goal, exists := s.Goals[id]
	if !exists {
		return nil
	}

	if goal.Completed {
		return nil
	}

	goal.Completed = true
	goal.CompletedAt = time.Now()
	s.Goals[id] = goal

	return s.save()
}

func (s *JSONStorage) GetAllTracks() ([]models.HabitTrack, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var tracks []models.HabitTrack
	for _, track := range s.HabitTracks {
		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (s *JSONStorage) GetTrackByID(id int) (*models.HabitTrack, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	track, exists := s.HabitTracks[id]
	if !exists {
		return nil, nil
	}

	return &track, nil
}

func (s *JSONStorage) CreateTrack(track *models.HabitTrack) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	track.ID = s.NextTrackID
	s.NextTrackID++
	s.HabitTracks[track.ID] = *track

	return s.save()
}

func (s *JSONStorage) UpdateTrack(id int, track *models.HabitTrack) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.HabitTracks[id]; !exists {
		return nil
	}

	track.ID = id
	s.HabitTracks[id] = *track

	return s.save()
}

func (s *JSONStorage) DeleteTrack(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.HabitTracks[id]; !exists {
		return nil
	}

	delete(s.HabitTracks, id)
	return s.save()
}

type Statistics struct {
	TotalHabits         int                      `json:"total_habits"`
	CompletedHabits     int                      `json:"completed_habits"`
	HabitCompletionRate float64                  `json:"habit_completion_rate"`
	TotalGoals          int                      `json:"total_goals"`
	CompletedGoals      int                      `json:"completed_goals"`
	GoalCompletionRate  float64                  `json:"goal_completion_rate"`
	OverdueGoals        int                      `json:"overdue_goals"`
	TodayCompleted      int                      `json:"today_completed"`
	TotalItems          int                      `json:"total_items"`
	CompletedItems      int                      `json:"completed_items"`
	OverallProgress     float64                  `json:"overall_progress"`
	Categories          map[string]CategoryStats `json:"categories"`
}

type CategoryStats struct {
	Total      int     `json:"total"`
	Completed  int     `json:"completed"`
	Percentage float64 `json:"percentage"`
}

func (s *JSONStorage) GetStatistics() (*Statistics, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := &Statistics{
		Categories: make(map[string]CategoryStats),
	}

	// Считаем привычки
	totalHabits := len(s.Habits)
	completedHabits := 0
	for _, habit := range s.Habits {
		if habit.Completed {
			completedHabits++
		}

		// Важно: работаем с копией структуры из map
		catStats := stats.Categories[habit.Category]
		catStats.Total++
		if habit.Completed {
			catStats.Completed++
		}
		stats.Categories[habit.Category] = catStats
	}
	if totalHabits > 0 {
		stats.HabitCompletionRate = float64(completedHabits) / float64(totalHabits) * 100
	}

	// Считаем цели
	totalGoals := len(s.Goals)
	completedGoals := 0
	overdueGoals := 0
	now := time.Now()
	for _, goal := range s.Goals {
		if goal.Completed {
			completedGoals++
		} else if goal.TargetDate.Before(now) {
			overdueGoals++
		}
	}
	if totalGoals > 0 {
		stats.GoalCompletionRate = float64(completedGoals) / float64(totalGoals) * 100
	}

	// Сегодняшние выполнения
	today := now.Format("2006-01-02")
	todayCompleted := 0
	for _, track := range s.HabitTracks {
		if track.Date.Format("2006-01-02") == today && track.Completed {
			todayCompleted++
		}
	}

	// Общий прогресс
	totalItems := totalHabits + totalGoals
	completedItems := completedHabits + completedGoals
	if totalItems > 0 {
		stats.OverallProgress = float64(completedItems) / float64(totalItems) * 100
	}

	// Считаем проценты по категориям
	for category, catStats := range stats.Categories {
		if catStats.Total > 0 {
			catStats.Percentage = float64(catStats.Completed) / float64(catStats.Total) * 100
		}
		stats.Categories[category] = catStats
	}

	// Заполняем остальные поля
	stats.TotalHabits = totalHabits
	stats.CompletedHabits = completedHabits
	stats.TotalGoals = totalGoals
	stats.CompletedGoals = completedGoals
	stats.OverdueGoals = overdueGoals
	stats.TodayCompleted = todayCompleted
	stats.TotalItems = totalItems
	stats.CompletedItems = completedItems

	return stats, nil
}
