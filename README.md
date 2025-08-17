# Task API

A REST API for managing tasks with detailed information support.

## Project Structure

```
├── main.go              # Main application with HTTP handlers
├── main_test.go         # Tests for HTTP handlers
├── models/
│   └── task.go         # Task and TaskDetail model definitions
└── service/
    └── task_service.go # Task service layer with business logic
```

## API Endpoints

### Basic Task Operations

- `GET /tasks` - Get all tasks
- `POST /tasks` - Create a new task
  ```json
  {
    "name": "Task name"
  }
  ```
- `PUT /tasks/{id}` - Update a task
  ```json
  {
    "name": "Updated task name"
  }
  ```
- `DELETE /tasks/{id}` - Delete a task

### Task Detail Operations

- `GET /tasks/{id}/detail` - Get detailed information about a task
- `PUT /tasks/{id}/detail` - Update detailed information about a task
  ```json
  {
    "priority": "high",
    "tags": ["docker", "deployment", "devops"],
    "estimated_hours": 20
  }
  ```
- `POST /tasks/{id}/complete` - Mark a task as completed

## Data Models

### Task

```json
{
  "id": 1,
  "name": "Learn DevOps on Azure",
  "created_at": "2025-08-17T23:59:47.760443+07:00",
  "updated_at": "2025-08-17T23:59:47.760443+07:00",
  "status": "pending"
}
```

### TaskDetail

```json
{
  "id": 2,
  "name": "Test Docker Deployment",
  "created_at": "2025-08-18T00:00:20.977599+07:00",
  "updated_at": "2025-08-18T00:00:39.067526+07:00",
  "status": "completed",
  "priority": "high",
  "tags": ["docker", "deployment", "devops"],
  "estimated_hours": 20,
  "completed_at": "2025-08-18T00:00:39.067526+07:00"
}
```

## Running the Application

1. Build the application:

   ```bash
   go build -o task-api
   ```

2. Run the application:

   ```bash
   ./task-api
   ```

   The server will start on port 8080 by default. You can change the port using the `PORT` environment variable:

   ```bash
   PORT=8081 ./task-api
   ```

## Running Tests

```bash
go test -v
```

## Features

- **Separation of Concerns**: Service layer separated from HTTP handlers
- **Enhanced Task Model**: Tasks now include timestamps, status, description
- **Task Details**: Extended information including priority, tags, estimated hours, completion tracking
- **Comprehensive Testing**: Full test coverage for all endpoints
- **Error Handling**: Proper error responses for invalid requests
- **RESTful Design**: Following REST conventions for API endpoints

## Task Status

Tasks can have the following status values:

- `pending` - Default status for new tasks
- `completed` - Status after calling the complete endpoint

## Priority Levels

Tasks can have the following priority levels:

- `low`
- `medium` (default)
- `high`
