package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/chawatvish/go-task-api/models"
	"github.com/chawatvish/go-task-api/service"
	"github.com/gofiber/fiber/v2"
)

// setupApp creates a new Fiber app for testing
func setupApp() *fiber.App {
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

	return app
}

func TestGetTasks(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("GET", "/tasks", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var result []models.Task
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 task, got %d", len(result))
	}

	if result[0].Name != "Learn DevOps on Azure" {
		t.Errorf("Expected task name 'Learn DevOps on Azure', got '%s'", result[0].Name)
	}
}

func TestCreateTask(t *testing.T) {
	app := setupApp()

	newTask := struct {
		Name string `json:"name"`
	}{Name: "New Test Task"}
	taskJSON, _ := json.Marshal(newTask)

	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 201 {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var result models.Task
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatal(err)
	}

	if result.Name != "New Test Task" {
		t.Errorf("Expected task name 'New Test Task', got '%s'", result.Name)
	}

	if result.ID != 2 {
		t.Errorf("Expected task ID 2, got %d", result.ID)
	}
}

func TestCreateTaskInvalidJSON(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 400 {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestUpdateTask(t *testing.T) {
	app := setupApp()

	updatedTask := struct {
		Name string `json:"name"`
	}{Name: "Updated Task Name"}
	taskJSON, _ := json.Marshal(updatedTask)

	req := httptest.NewRequest("PUT", "/tasks/1", bytes.NewReader(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var result models.Task
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatal(err)
	}

	if result.Name != "Updated Task Name" {
		t.Errorf("Expected updated task name 'Updated Task Name', got '%s'", result.Name)
	}

	if result.ID != 1 {
		t.Errorf("Expected task ID to remain 1, got %d", result.ID)
	}
}

func TestUpdateNonExistentTask(t *testing.T) {
	app := setupApp()

	updatedTask := struct {
		Name string `json:"name"`
	}{Name: "Updated Task Name"}
	taskJSON, _ := json.Marshal(updatedTask)

	req := httptest.NewRequest("PUT", "/tasks/999", bytes.NewReader(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 404 {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}

func TestUpdateTaskInvalidID(t *testing.T) {
	app := setupApp()

	updatedTask := struct {
		Name string `json:"name"`
	}{Name: "Updated Task Name"}
	taskJSON, _ := json.Marshal(updatedTask)

	req := httptest.NewRequest("PUT", "/tasks/invalid", bytes.NewReader(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 400 {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestUpdateTaskInvalidJSON(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("PUT", "/tasks/1", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 400 {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestDeleteTask(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("DELETE", "/tasks/1", nil)
	
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}

	// Verify task was deleted by trying to get its details
	req = httptest.NewRequest("GET", "/tasks/1/detail", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 404 {
		t.Errorf("Expected status 404 for deleted task detail, got %d", resp.StatusCode)
	}
}

func TestDeleteNonExistentTask(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("DELETE", "/tasks/999", nil)
	
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 404 {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}

func TestDeleteTaskInvalidID(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("DELETE", "/tasks/invalid", nil)
	
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 400 {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestGetTaskDetail(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("GET", "/tasks/1/detail", nil)
	
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var result models.TaskDetail
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatal(err)
	}

	if result.Name != "Learn DevOps on Azure" {
		t.Errorf("Expected task name 'Learn DevOps on Azure', got '%s'", result.Name)
	}

	if result.Priority != "medium" {
		t.Errorf("Expected priority 'medium', got '%s'", result.Priority)
	}
}

func TestGetTaskDetailNonExistent(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("GET", "/tasks/999/detail", nil)
	
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 404 {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}

func TestUpdateTaskDetail(t *testing.T) {
	app := setupApp()

	updateDetail := struct {
		Priority       string   `json:"priority"`
		Tags           []string `json:"tags"`
		EstimatedHours int      `json:"estimated_hours"`
	}{
		Priority:       "high",
		Tags:           []string{"urgent", "devops"},
		EstimatedHours: 60,
	}
	detailJSON, _ := json.Marshal(updateDetail)

	req := httptest.NewRequest("PUT", "/tasks/1/detail", bytes.NewReader(detailJSON))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var result models.TaskDetail
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatal(err)
	}

	if result.Priority != "high" {
		t.Errorf("Expected updated priority 'high', got '%s'", result.Priority)
	}

	if len(result.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(result.Tags))
	}

	if result.EstimatedHours != 60 {
		t.Errorf("Expected estimated hours 60, got %d", result.EstimatedHours)
	}
}

func TestCompleteTask(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("POST", "/tasks/1/complete", nil)
	
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Verify task is marked as completed by getting its details
	req = httptest.NewRequest("GET", "/tasks/1/detail", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var result models.TaskDetail
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatal(err)
	}

	if result.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", result.Status)
	}

	if result.CompletedAt == nil {
		t.Error("Expected CompletedAt to be set, but it was nil")
	}
}