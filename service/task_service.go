package service

import (
	"errors"
	"time"

	"github.com/chawatvish/go-task-api/models"
)

// TaskService handles task operations
type TaskService struct {
	tasks       []models.Task
	taskDetails map[int]*models.TaskDetail
	nextID      int
}

// NewTaskService creates a new task service instance
func NewTaskService() *TaskService {
	now := time.Now()
	initialTask := models.Task{
		ID:        1,
		Name:      "Learn DevOps on Azure",
		Status:    "pending",
		CreatedAt: now,
		UpdatedAt: now,
	}

	service := &TaskService{
		tasks:       []models.Task{initialTask},
		taskDetails: make(map[int]*models.TaskDetail),
		nextID:      2,
	}

	// Initialize detailed information for the initial task
	service.taskDetails[1] = &models.TaskDetail{
		Task:           initialTask,
		Priority:       "medium",
		Tags:           []string{"devops", "azure", "learning"},
		EstimatedHours: 40,
	}

	return service
}

// GetAllTasks returns all tasks
func (ts *TaskService) GetAllTasks() []models.Task {
	return ts.tasks
}

// CreateTask creates a new task
func (ts *TaskService) CreateTask(name string) models.Task {
	now := time.Now()
	task := models.Task{
		ID:        ts.nextID,
		Name:      name,
		Status:    "pending",
		CreatedAt: now,
		UpdatedAt: now,
	}

	ts.tasks = append(ts.tasks, task)
	
	// Initialize basic task details
	ts.taskDetails[task.ID] = &models.TaskDetail{
		Task:     task,
		Priority: "medium",
	}
	
	ts.nextID++
	return task
}

// UpdateTask updates an existing task
func (ts *TaskService) UpdateTask(id int, name string) (models.Task, error) {
	for i, task := range ts.tasks {
		if task.ID == id {
			ts.tasks[i].Name = name
			ts.tasks[i].UpdatedAt = time.Now()
			
			// Update the detailed information as well
			if detail, exists := ts.taskDetails[id]; exists {
				detail.Task = ts.tasks[i]
			}
			
			return ts.tasks[i], nil
		}
	}
	return models.Task{}, errors.New("task not found")
}

// DeleteTask deletes a task by ID
func (ts *TaskService) DeleteTask(id int) error {
	for i, task := range ts.tasks {
		if task.ID == id {
			ts.tasks = append(ts.tasks[:i], ts.tasks[i+1:]...)
			delete(ts.taskDetails, id)
			return nil
		}
	}
	return errors.New("task not found")
}

// GetTaskDetail returns detailed information about a task
func (ts *TaskService) GetTaskDetail(id int) (*models.TaskDetail, error) {
	detail, exists := ts.taskDetails[id]
	if !exists {
		return nil, errors.New("task not found")
	}
	return detail, nil
}

// UpdateTaskDetail updates detailed information about a task
func (ts *TaskService) UpdateTaskDetail(id int, priority string, tags []string, estimatedHours int) (*models.TaskDetail, error) {
	detail, exists := ts.taskDetails[id]
	if !exists {
		return nil, errors.New("task not found")
	}

	detail.Priority = priority
	detail.Tags = tags
	detail.EstimatedHours = estimatedHours
	detail.UpdatedAt = time.Now()

	// Also update the task in the main tasks slice
	for i, task := range ts.tasks {
		if task.ID == id {
			ts.tasks[i].UpdatedAt = time.Now()
			detail.Task = ts.tasks[i]
			break
		}
	}

	return detail, nil
}

// MarkTaskCompleted marks a task as completed
func (ts *TaskService) MarkTaskCompleted(id int) error {
	for i, task := range ts.tasks {
		if task.ID == id {
			now := time.Now()
			ts.tasks[i].Status = "completed"
			ts.tasks[i].UpdatedAt = now
			
			if detail, exists := ts.taskDetails[id]; exists {
				detail.Task = ts.tasks[i]
				detail.CompletedAt = &now
			}
			
			return nil
		}
	}
	return errors.New("task not found")
}
