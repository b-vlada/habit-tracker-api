package handlers

import (
	"habit-tracker-api/models"
	"habit-tracker-api/storage"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type TrackHandler struct {
	storage *storage.JSONStorage
}

func NewTrackHandler(storage *storage.JSONStorage) *TrackHandler {
	return &TrackHandler{storage: storage}
}

type CreateTrackRequest struct {
	HabitID   int       `json:"habit_id" validate:"required"`
	Date      time.Time `json:"date" validate:"required"`
	Completed bool      `json:"completed"`
	Notes     string    `json:"notes"`
}

type UpdateTrackRequest struct {
	HabitID   int       `json:"habit_id" validate:"required"`
	Date      time.Time `json:"date" validate:"required"`
	Completed bool      `json:"completed"`
	Notes     string    `json:"notes"`
}

func (h *TrackHandler) GetAllTracks(c *fiber.Ctx) error {
	tracks, err := h.storage.GetAllTracks()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get tracks",
		})
	}

	return c.JSON(fiber.Map{
		"tracks": tracks,
		"count":  len(tracks),
	})
}

func (h *TrackHandler) GetTrackByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid track ID",
		})
	}

	track, err := h.storage.GetTrackByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get track",
		})
	}

	if track == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Track not found",
		})
	}

	return c.JSON(track)
}

func (h *TrackHandler) CreateTrack(c *fiber.Ctx) error {
	var req CreateTrackRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.HabitID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Habit ID is required",
		})
	}

	habit, err := h.storage.GetHabitByID(req.HabitID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to validate habit",
		})
	}
	if habit == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Habit not found",
		})
	}

	track := &models.HabitTrack{
		HabitID:   req.HabitID,
		Date:      req.Date,
		Completed: req.Completed,
		Notes:     req.Notes,
	}

	if err := h.storage.CreateTrack(track); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create track",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(track)
}

func (h *TrackHandler) UpdateTrack(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid track ID",
		})
	}

	var req UpdateTrackRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.HabitID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Habit ID is required",
		})
	}

	habit, err := h.storage.GetHabitByID(req.HabitID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to validate habit",
		})
	}
	if habit == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Habit not found",
		})
	}

	existingTrack, err := h.storage.GetTrackByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get track",
		})
	}

	if existingTrack == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Track not found",
		})
	}

	updatedTrack := &models.HabitTrack{
		ID:        id,
		HabitID:   req.HabitID,
		Date:      req.Date,
		Completed: req.Completed,
		Notes:     req.Notes,
	}

	if err := h.storage.UpdateTrack(id, updatedTrack); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update track",
		})
	}

	return c.JSON(updatedTrack)
}

func (h *TrackHandler) DeleteTrack(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid track ID",
		})
	}

	track, err := h.storage.GetTrackByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get track",
		})
	}

	if track == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Track not found",
		})
	}

	if err := h.storage.DeleteTrack(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete track",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
