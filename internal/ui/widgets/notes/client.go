package notes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/jonesrussell/dashboard/internal/logger"
	"github.com/jonesrussell/dashboard/internal/logger/types"
)

// Default configuration
const (
	defaultBaseURL    = "http://host.docker.internal:8080"
	defaultAPITimeout = 10 * time.Second
	envGodoAPIBaseURL = "GODO_API_URL"
)

// Client handles communication with the godo API
type Client struct {
	baseURL    string
	httpClient *http.Client
	logger     logger.Logger
}

// ClientOption allows configuring the client
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL for the API
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithTimeout sets a custom timeout for API requests
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithLogger sets the logger for the client
func WithLogger(l logger.Logger) ClientOption {
	return func(c *Client) {
		c.logger = l
	}
}

// NewClient creates a new godo API client
func NewClient(opts ...ClientOption) *Client {
	// Get base URL from environment or use default
	baseURL := os.Getenv(envGodoAPIBaseURL)
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	// Create default logger if none provided
	defaultLogger, err := logger.New(logger.DefaultConfig())
	if err != nil {
		defaultLogger, _ = logger.NewZapLogger(types.Config{
			Level:      "debug",
			OutputPath: "logs/app.log",
		})
	}

	client := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: defaultAPITimeout,
		},
		logger: defaultLogger,
	}

	// Apply options
	for _, opt := range opts {
		opt(client)
	}

	return client
}

// Note represents a task from the godo API
type Note struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// NoteInput represents the input for creating/updating a task
type NoteInput struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

// ListNotes retrieves all tasks
func (c *Client) ListNotes() ([]Note, error) {
	c.logger.Debug("Making notes request",
		logger.NewField("url", c.baseURL),
		logger.NewField("method", "GET"),
		logger.NewField("path", "/api/v1/tasks"),
	)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/tasks", c.baseURL), nil)
	if err != nil {
		c.logger.Error("Failed to create request",
			logger.NewField("error", err),
		)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Failed to make request",
			logger.NewField("error", err),
		)
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body for logging
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Failed to read response body",
			logger.NewField("error", err),
		)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	c.logger.Debug("Received response",
		logger.NewField("body", string(bodyBytes)),
	)

	// Create a new reader with the body bytes for decoding
	resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Unexpected status code",
			logger.NewField("status_code", resp.StatusCode),
		)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Try to decode as a single note first
	var note Note
	if err := json.NewDecoder(resp.Body).Decode(&note); err != nil {
		// Reset the reader
		resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		// Try to decode as array
		var notes []Note
		if err := json.NewDecoder(resp.Body).Decode(&notes); err != nil {
			c.logger.Error("Failed to decode response",
				logger.NewField("error", err),
				logger.NewField("body", string(bodyBytes)),
			)
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return notes, nil
	}

	// If we successfully decoded a single note, return it as a slice
	c.logger.Debug("Successfully decoded single note",
		logger.NewField("id", note.ID),
	)
	return []Note{note}, nil
}

// CreateNote creates a new task
func (c *Client) CreateNote(input NoteInput) (*Note, error) {
	c.logger.Debug("Creating note",
		logger.NewField("title", input.Title),
	)

	body, err := json.Marshal(input)
	if err != nil {
		c.logger.Error("Failed to marshal input",
			logger.NewField("error", err),
		)
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/api/v1/tasks", c.baseURL),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		c.logger.Error("Failed to create note",
			logger.NewField("error", err),
		)
		return nil, fmt.Errorf("failed to create task: %w", err)
	}
	defer resp.Body.Close()

	// Read and log response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Failed to read response body",
			logger.NewField("error", err),
		)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	c.logger.Debug("Received response",
		logger.NewField("status_code", resp.StatusCode),
		logger.NewField("body", string(bodyBytes)),
	)
	resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	if resp.StatusCode != http.StatusCreated {
		c.logger.Error("Unexpected status code",
			logger.NewField("status_code", resp.StatusCode),
			logger.NewField("status", resp.Status),
		)
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var task Note
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		c.logger.Error("Failed to decode response",
			logger.NewField("error", err),
		)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	c.logger.Debug("Successfully created note",
		logger.NewField("id", task.ID),
	)
	return &task, nil
}

// UpdateNote updates an existing note
func (c *Client) UpdateNote(id string, input NoteInput) (*Note, error) {
	c.logger.Debug("Updating note",
		logger.NewField("id", id),
		logger.NewField("title", input.Title),
	)

	body, err := json.Marshal(input)
	if err != nil {
		c.logger.Error("Failed to marshal input",
			logger.NewField("error", err),
		)
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/api/v1/tasks/%s", c.baseURL, id),
		bytes.NewReader(body),
	)
	if err != nil {
		c.logger.Error("Failed to create request",
			logger.NewField("error", err),
		)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Failed to update note",
			logger.NewField("error", err),
		)
		return nil, fmt.Errorf("failed to update task: %w", err)
	}
	defer resp.Body.Close()

	// Read and log response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Failed to read response body",
			logger.NewField("error", err),
		)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	c.logger.Debug("Received response",
		logger.NewField("status_code", resp.StatusCode),
		logger.NewField("body", string(bodyBytes)),
	)
	resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Unexpected status code",
			logger.NewField("status_code", resp.StatusCode),
			logger.NewField("status", resp.Status),
		)
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var task Note
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		c.logger.Error("Failed to decode response",
			logger.NewField("error", err),
		)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	c.logger.Debug("Successfully updated note",
		logger.NewField("id", task.ID),
	)
	return &task, nil
}

// DeleteNote deletes a note
func (c *Client) DeleteNote(id string) error {
	c.logger.Debug("Deleting note",
		logger.NewField("id", id),
	)

	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/api/v1/tasks/%s", c.baseURL, id),
		nil,
	)
	if err != nil {
		c.logger.Error("Failed to create request",
			logger.NewField("error", err),
		)
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Failed to delete note",
			logger.NewField("error", err),
		)
		return fmt.Errorf("failed to delete task: %w", err)
	}
	defer resp.Body.Close()

	// Read and log response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Failed to read response body",
			logger.NewField("error", err),
		)
		return fmt.Errorf("failed to read response body: %w", err)
	}
	c.logger.Debug("Received response",
		logger.NewField("status_code", resp.StatusCode),
		logger.NewField("body", string(bodyBytes)),
	)
	resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	if resp.StatusCode != http.StatusNoContent {
		c.logger.Error("Unexpected status code",
			logger.NewField("status_code", resp.StatusCode),
			logger.NewField("status", resp.Status),
		)
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	c.logger.Debug("Successfully deleted note",
		logger.NewField("id", id),
	)
	return nil
}
