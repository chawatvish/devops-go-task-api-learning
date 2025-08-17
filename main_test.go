package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// setupApp creates a new Fiber app for testing
func setupApp() *fiber.App {
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

	return app
}

// resetTasks resets the tasks slice to its initial state for testing
func resetTasks() {
	tasks = []Task{{ID: 1, Name: "Learn DevOps on Azure"}}
}

func TestGetTasks(t *testing.T) {
	resetTasks()
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

	var result []Task
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
	resetTasks()
	app := setupApp()

	newTask := Task{Name: "New Test Task"}
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

	var result Task
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
	resetTasks()
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
	resetTasks()
	app := setupApp()

	updatedTask := Task{Name: "Updated Task Name"}
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

	var result Task
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
	resetTasks()
	app := setupApp()

	updatedTask := Task{Name: "Updated Task Name"}
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
	resetTasks()
	app := setupApp()

	updatedTask := Task{Name: "Updated Task Name"}
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
	resetTasks()
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
	resetTasks()
	app := setupApp()

	req := httptest.NewRequest("DELETE", "/tasks/1", nil)
	
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}

	// Verify task was actually deleted
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks after deletion, got %d", len(tasks))
	}
}

func TestDeleteNonExistentTask(t *testing.T) {
	resetTasks()
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
	resetTasks()
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

func TestMultipleOperations(t *testing.T) {
	resetTasks()
	app := setupApp()

	// Create a new task
	newTask := Task{Name: "Integration Test Task"}
	taskJSON, _ := json.Marshal(newTask)

	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 201 {
		t.Errorf("Expected status 201 for POST, got %d", resp.StatusCode)
	}

	// Update the created task
	updatedTask := Task{Name: "Updated Integration Test Task"}
	updatedJSON, _ := json.Marshal(updatedTask)

	req = httptest.NewRequest("PUT", "/tasks/2", bytes.NewReader(updatedJSON))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200 for PUT, got %d", resp.StatusCode)
	}

	// Verify we have 2 tasks
	req = httptest.NewRequest("GET", "/tasks", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var result []Task
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(result))
	}

	// Delete the second task
	req = httptest.NewRequest("DELETE", "/tasks/2", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("Expected status 204 for DELETE, got %d", resp.StatusCode)
	}

	// Verify we have 1 task remaining
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task after deletion, got %d", len(tasks))
	}
}