package main

import (
	"habit-tracker-api/handlers"
	"habit-tracker-api/storage"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	storage, err := storage.NewJSONStorage("habits.json")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	habitHandler := handlers.NewHabitHandler(storage)
	goalHandler := handlers.NewGoalHandler(storage)
	trackHandler := handlers.NewTrackHandler(storage)

	app := fiber.New(fiber.Config{
		AppName: "Habit Tracker API",
	})

	app.Use(logger.New())
	app.Use(recover.New())

	api := app.Group("/api/v1")

	habits := api.Group("/habits")
	{
		habits.Get("/", habitHandler.GetAllHabits)
		habits.Get("/:id", habitHandler.GetHabitByID)
		habits.Post("/", habitHandler.CreateHabit)
		habits.Put("/:id", habitHandler.UpdateHabit)
		habits.Delete("/:id", habitHandler.DeleteHabit)
		habits.Put("/:id/complete", habitHandler.CompleteHabit)
	}

	goals := api.Group("/goals")
	{
		goals.Get("/", goalHandler.GetAllGoals)
		goals.Get("/:id", goalHandler.GetGoalByID)
		goals.Post("/", goalHandler.CreateGoal)
		goals.Put("/:id", goalHandler.UpdateGoal)
		goals.Delete("/:id", goalHandler.DeleteGoal)
		goals.Put("/:id/complete", goalHandler.CompleteGoal)
	}

	tracks := api.Group("/tracks")
	{
		tracks.Get("/", trackHandler.GetAllTracks)
		tracks.Get("/:id", trackHandler.GetTrackByID)
		tracks.Post("/", trackHandler.CreateTrack)
		tracks.Put("/:id", trackHandler.UpdateTrack)
		tracks.Delete("/:id", trackHandler.DeleteTrack)
	}

	api.Get("/statistics", func(c *fiber.Ctx) error {
		stats, err := storage.GetStatistics()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get statistics",
			})
		}
		return c.JSON(stats)
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Habit Tracker API is running",
			"version": "1.0.0",
			"endpoints": []string{
				"/api/v1/habits",
				"/api/v1/goals",
				"/api/v1/tracks",
				"/api/v1/statistics",
			},
		})
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Endpoint not found",
		})
	})

	log.Println("Server starting on :3000")
	log.Fatal(app.Listen(":3000"))
}
