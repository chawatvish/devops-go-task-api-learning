package main

import (
	"log"
	"os"

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s\n", port)
	log.Fatal(app.Listen(":" + port))
}