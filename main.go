package main

import (
	"log"
	"os"
	"strconv"

	"github.com/chawatvish/go-task-api/service"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	taskService := service.NewTaskService()

	// GET /tasks - Get all tasks
	app.Get("/tasks", func(c *fiber.Ctx) error {
		tasks := taskService.GetAllTasks()
		return c.JSON(tasks)
	})

	// POST /tasks - Create a new task
	app.Post("/tasks", func(c *fiber.Ctx) error {
		var req struct {
			Name string `json:"name"`
		}
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}
		task := taskService.CreateTask(req.Name)
		return c.Status(fiber.StatusCreated).JSON(task)
	})

	// PUT /tasks/:id - Update a task
	app.Put("/tasks/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return fiber.ErrBadRequest
		}

		var req struct {
			Name string `json:"name"`
		}
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		task, err := taskService.UpdateTask(id, req.Name)
		if err != nil {
			return fiber.ErrNotFound
		}
		
		return c.JSON(task)
	})

	// DELETE /tasks/:id - Delete a task
	app.Delete("/tasks/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return fiber.ErrBadRequest
		}

		err = taskService.DeleteTask(id)
		if err != nil {
			return fiber.ErrNotFound
		}
		
		return c.SendStatus(fiber.StatusNoContent)
	})

	// GET /tasks/:id/detail - Get detailed information about a task
	app.Get("/tasks/:id/detail", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return fiber.ErrBadRequest
		}

		detail, err := taskService.GetTaskDetail(id)
		if err != nil {
			return fiber.ErrNotFound
		}
		
		return c.JSON(detail)
	})

	// PUT /tasks/:id/detail - Update detailed information about a task
	app.Put("/tasks/:id/detail", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return fiber.ErrBadRequest
		}

		var req struct {
			Priority       string   `json:"priority"`
			Tags           []string `json:"tags"`
			EstimatedHours int      `json:"estimated_hours"`
		}
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		detail, err := taskService.UpdateTaskDetail(id, req.Priority, req.Tags, req.EstimatedHours)
		if err != nil {
			return fiber.ErrNotFound
		}
		
		return c.JSON(detail)
	})

	// POST /tasks/:id/complete - Mark a task as completed
	app.Post("/tasks/:id/complete", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return fiber.ErrBadRequest
		}

		err = taskService.MarkTaskCompleted(id)
		if err != nil {
			return fiber.ErrNotFound
		}
		
		return c.SendStatus(fiber.StatusOK)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s\n", port)
	log.Fatal(app.Listen(":" + port))
}