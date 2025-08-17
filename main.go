package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var tasks = []Task{{ID: 1, Name: "Learn DevOps on Azure"}}

func main() {
	app := fiber.New()

	app.Get("/tasks", func(c *fiber.Ctx) error {
		return c.JSON(tasks)
	})

	app.Post("/tasks", func(c *fiber.Ctx) error {
		var t Task
		if err := c.BodyParser(&t); err != nil {
			return fiber.ErrBadRequest
		}
		t.ID = len(tasks) + 1
		tasks = append(tasks, t)
		return c.Status(fiber.StatusCreated).JSON(t)
	})

	// PUT /tasks/:id - Update a task
	app.Put("/tasks/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return fiber.ErrBadRequest
		}

		var updatedTask Task
		if err := c.BodyParser(&updatedTask); err != nil {
			return fiber.ErrBadRequest
		}

		// Find and update the task
		for i, task := range tasks {
			if task.ID == id {
				tasks[i].Name = updatedTask.Name
				tasks[i].ID = id // Keep the original ID
				return c.JSON(tasks[i])
			}
		}
		
		return fiber.ErrNotFound
	})

	// DELETE /tasks/:id - Delete a task
	app.Delete("/tasks/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return fiber.ErrBadRequest
		}

		// Find and delete the task
		for i, task := range tasks {
			if task.ID == id {
				tasks = append(tasks[:i], tasks[i+1:]...)
				return c.SendStatus(fiber.StatusNoContent)
			}
		}
		
		return fiber.ErrNotFound
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s\n", port)
	log.Fatal(app.Listen(":" + port))
}