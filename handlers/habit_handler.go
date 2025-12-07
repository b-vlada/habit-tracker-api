package handlers

import (
	"habit-tracker-api/models"
	"habit-tracker-api/storage"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type HabitHandler struct {
	storage *storage.JSONStorage
}

func NewHabitHandler(storage *storage.JSONStorage) *HabitHandler {
	return &HabitHandler{storage: storage}
}

type CreateHabitRequest struct {
	Name        string `json:"name" validate:"required,min=1"`
	Description string `json:"description"`
	Category    string `json:"category" validate:"required"`
	Frequency   string `json:"frequency" validate:"required"`
}

type UpdateHabitRequest struct {
	Name        string `json:"name" validate:"required,min=1"`
	Description string `json:"description"`
	Category    string `json:"category" validate:"required"`
	Frequency   string `json:"frequency" validate:"required"`
}

func (h *HabitHandler) GetAllHabits(c *fiber.Ctx) error {
	habits, err := h.storage.GetAllHabits()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get habits",
		})
	}

	return c.JSON(fiber.Map{
		"habits": habits,
		"count":  len(habits),
	})
}

func (h *HabitHandler) GetHabitByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid habit ID",
		})
	}

	habit, err := h.storage.GetHabitByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get habit",
		})
	}

	if habit == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Habit not found",
		})
	}

	return c.JSON(habit)
}

func (h *HabitHandler) CreateHabit(c *fiber.Ctx) error {
	var req CreateHabitRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Habit name is required",
		})
	}

	if req.Category == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category is required",
		})
	}

	if req.Frequency == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Frequency is required",
		})
	}

	now := time.Now()
	habit := &models.Habit{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Frequency:   req.Frequency,
		CreatedAt:   now,
		Completed:   false,
	}

	if err := h.storage.CreateHabit(habit); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create habit",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(habit)
}

func (h *HabitHandler) UpdateHabit(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid habit ID",
		})
	}

	var req UpdateHabitRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Habit name is required",
		})
	}

	if req.Category == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category is required",
		})
	}

	if req.Frequency == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Frequency is required",
		})
	}

	existingHabit, err := h.storage.GetHabitByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get habit",
		})
	}

	if existingHabit == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Habit not found",
		})
	}

	updatedHabit := &models.Habit{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Frequency:   req.Frequency,
		CreatedAt:   existingHabit.CreatedAt,
		Completed:   existingHabit.Completed,
	}

	if err := h.storage.UpdateHabit(id, updatedHabit); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update habit",
		})
	}

	return c.JSON(updatedHabit)
}

func (h *HabitHandler) DeleteHabit(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid habit ID",
		})
	}

	habit, err := h.storage.GetHabitByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get habit",
		})
	}

	if habit == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Habit not found",
		})
	}

	if err := h.storage.DeleteHabit(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete habit",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func (h *HabitHandler) CompleteHabit(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid habit ID",
		})
	}

	if err := h.storage.CompleteHabit(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to complete habit",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Habit marked as completed",
		"id":      id,
	})
}
