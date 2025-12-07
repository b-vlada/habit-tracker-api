package handlers

import (
	"habit-tracker-api/models"
	"habit-tracker-api/storage"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type GoalHandler struct {
	storage *storage.JSONStorage
}

func NewGoalHandler(storage *storage.JSONStorage) *GoalHandler {
	return &GoalHandler{storage: storage}
}

type CreateGoalRequest struct {
	Title       string    `json:"title" validate:"required,min=1"`
	Description string    `json:"description"`
	TargetDate  time.Time `json:"target_date" validate:"required"`
}

type UpdateGoalRequest struct {
	Title       string    `json:"title" validate:"required,min=1"`
	Description string    `json:"description"`
	TargetDate  time.Time `json:"target_date" validate:"required"`
}

func (h *GoalHandler) GetAllGoals(c *fiber.Ctx) error {
	goals, err := h.storage.GetAllGoals()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get goals",
		})
	}

	return c.JSON(fiber.Map{
		"goals": goals,
		"count": len(goals),
	})
}

func (h *GoalHandler) GetGoalByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid goal ID",
		})
	}

	goal, err := h.storage.GetGoalByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get goal",
		})
	}

	if goal == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Goal not found",
		})
	}

	return c.JSON(goal)
}

func (h *GoalHandler) CreateGoal(c *fiber.Ctx) error {
	var req CreateGoalRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Goal title is required",
		})
	}

	now := time.Now()
	goal := &models.Goal{
		Title:       req.Title,
		Description: req.Description,
		TargetDate:  req.TargetDate,
		CreatedAt:   now,
		Completed:   false,
		CompletedAt: time.Time{},
		HabitIDs:    []int{},
	}

	if err := h.storage.CreateGoal(goal); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create goal",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(goal)
}

func (h *GoalHandler) UpdateGoal(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid goal ID",
		})
	}

	var req UpdateGoalRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Goal title is required",
		})
	}

	existingGoal, err := h.storage.GetGoalByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get goal",
		})
	}

	if existingGoal == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Goal not found",
		})
	}

	updatedGoal := &models.Goal{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		TargetDate:  req.TargetDate,
		CreatedAt:   existingGoal.CreatedAt,
		Completed:   existingGoal.Completed,
		CompletedAt: existingGoal.CompletedAt,
		HabitIDs:    existingGoal.HabitIDs,
	}

	if err := h.storage.UpdateGoal(id, updatedGoal); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update goal",
		})
	}

	return c.JSON(updatedGoal)
}

func (h *GoalHandler) DeleteGoal(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid goal ID",
		})
	}

	goal, err := h.storage.GetGoalByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get goal",
		})
	}

	if goal == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Goal not found",
		})
	}

	if err := h.storage.DeleteGoal(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete goal",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func (h *GoalHandler) CompleteGoal(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid goal ID",
		})
	}

	if err := h.storage.CompleteGoal(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to complete goal",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Goal marked as completed",
		"id":      id,
	})
}
