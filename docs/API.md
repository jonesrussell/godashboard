# API Documentation

## Task Management API

### Configuration
The Tasks widget connects to a Godo API server. Configure the endpoint using the environment variable:
```bash
GODO_API_URL="http://localhost:8080"  # Default: http://host.docker.internal:8080
```

### Endpoints

#### List Tasks
```
GET /api/v1/tasks
```
Returns list of tasks with their status.

#### Create Task
```
POST /api/v1/tasks
Content-Type: application/json

{
    "title": "Task title",
    "description": "Optional description"
}
```

#### Update Task
```
PUT /api/v1/tasks/{id}
Content-Type: application/json

{
    "title": "Updated title",
    "description": "Updated description"
}
```

#### Delete Task
```
DELETE /api/v1/tasks/{id}
```

### Data Types

#### Task
```go
type Task struct {
    ID          string     `json:"id"`
    Title       string     `json:"title"`
    Description string     `json:"description,omitempty"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
    CompletedAt *time.Time `json:"completed_at,omitempty"`
}
```

### Time Handling
- `CompletedAt` is a pointer to support three states:
  - `nil`: Task has never been completed
  - Zero time value (`0001-01-01T00:00:00Z`): Treated as incomplete
  - Valid time: Task is completed at that time
- When toggling completion:
  - Incomplete to complete: Sets to current time
  - Complete to incomplete: Sets to nil

### Error Handling
- 400: Bad Request - Invalid input
- 404: Not Found - Task doesn't exist
- 500: Internal Server Error

### Rate Limiting
- Default timeout: 10 seconds
- Configurable via `WithTimeout` option 