package notes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
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

// NewClient creates a new godo API client
func NewClient(opts ...ClientOption) *Client {
	// Get base URL from environment or use default
	baseURL := os.Getenv(envGodoAPIBaseURL)
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	client := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: defaultAPITimeout,
		},
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
	fmt.Printf("DEBUG: ListNotes - Making request to %s\n", c.baseURL)
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/tasks", c.baseURL), nil)
	if err != nil {
		fmt.Printf("DEBUG: ListNotes - Error creating request: %v\n", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Printf("DEBUG: ListNotes - Error making request: %v\n", err)
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("DEBUG: ListNotes - Unexpected status code: %d\n", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var notes []Note
	if err := json.NewDecoder(resp.Body).Decode(&notes); err != nil {
		fmt.Printf("DEBUG: ListNotes - Error decoding response: %v\n", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Printf("DEBUG: ListNotes - Successfully decoded %d notes\n", len(notes))
	return notes, nil
}

// CreateNote creates a new task
func (c *Client) CreateNote(input NoteInput) (*Note, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/api/v1/tasks", c.baseURL),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var task Note
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &task, nil
}

// UpdateNote updates an existing note
func (c *Client) UpdateNote(id string, input NoteInput) (*Note, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/api/v1/tasks/%s", c.baseURL, id),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var task Note
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &task, nil
}

// DeleteNote deletes a note
func (c *Client) DeleteNote(id string) error {
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/api/v1/tasks/%s", c.baseURL, id),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	return nil
}
