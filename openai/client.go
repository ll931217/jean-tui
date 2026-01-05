package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client represents an OpenAI-compatible API client
type Client struct {
	apiKey  string
	model   string
	baseURL string
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
}

// ChatMessage represents a message in the chat
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *APIError `json:"error"`
}

// PRContent represents structured PR content
type PRContent struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// NewClient creates a new OpenAI-compatible API client
func NewClient(apiKey, baseURL, model string) (*Client, error) {
	if apiKey == "" {
		return nil, NewConfigError("api_key", "API key is required")
	}
	if baseURL == "" {
		return nil, NewConfigError("base_url", "base URL is required")
	}
	if model == "" {
		return nil, NewConfigError("model", "model is required")
	}

	return &Client{
		apiKey:  apiKey,
		model:   model,
		baseURL: baseURL,
	}, nil
}

// GenerateCommitMessage generates a one-line conventional commit message based on git context
func (c *Client) GenerateCommitMessage(status, diff, branch, log, customPrompt string) (string, error) {
	// Limit diff to reasonable size to avoid token limits
	if len(diff) > 5000 {
		diff = diff[:5000]
	}

	// Use custom prompt if provided, otherwise use default
	prompt := customPrompt
	if prompt == "" {
		prompt = DefaultCommitPrompt
	}

	// Replace all placeholders with actual context
	prompt = strings.ReplaceAll(prompt, "{status}", status)
	prompt = strings.ReplaceAll(prompt, "{diff}", diff)
	prompt = strings.ReplaceAll(prompt, "{branch}", branch)
	prompt = strings.ReplaceAll(prompt, "{log}", log)

	response, err := c.callAPI(prompt)
	if err != nil {
		return "", WrapError("generate commit message", err)
	}

	// Parse plain text response (no JSON)
	subject := strings.TrimSpace(response)
	if subject == "" {
		return "", fmt.Errorf("AI generated empty commit subject")
	}

	return subject, nil
}

// GenerateBranchName generates a semantic branch name based on git diff
func (c *Client) GenerateBranchName(diff, customPrompt string) (string, error) {
	// Limit diff to reasonable size
	if len(diff) > 3000 {
		diff = diff[:3000]
	}

	// Use custom prompt if provided, otherwise use default
	prompt := customPrompt
	if prompt == "" {
		prompt = DefaultBranchNamePrompt
	}
	// Replace {diff} placeholder with actual diff
	prompt = strings.ReplaceAll(prompt, "{diff}", diff)

	name, err := c.callAPI(prompt)
	if err != nil {
		return "", WrapError("generate branch name", err)
	}

	// Clean up response
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")

	// Remove non-alphanumeric except hyphens
	var result []rune
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result = append(result, r)
		}
	}
	name = string(result)

	// Remove leading/trailing hyphens
	name = strings.Trim(name, "-")

	// Limit to 40 chars
	if len(name) > 40 {
		name = name[:40]
	}

	if name == "" {
		return "", fmt.Errorf("AI generated invalid branch name")
	}

	return name, nil
}

// GeneratePRContent generates a PR title and description from a git diff
func (c *Client) GeneratePRContent(diff, customPrompt string) (title, description string, err error) {
	// Limit diff to reasonable size
	if len(diff) > 5000 {
		diff = diff[:5000]
	}

	// Use custom prompt if provided, otherwise use default
	prompt := customPrompt
	if prompt == "" {
		prompt = DefaultPRPrompt
	}
	// Replace {diff} placeholder with actual diff
	prompt = strings.ReplaceAll(prompt, "{diff}", diff)

	response, err := c.callAPI(prompt)
	if err != nil {
		return "", "", WrapError("generate PR content", err)
	}

	// Parse JSON response
	var content PRContent
	if err := json.Unmarshal([]byte(response), &content); err != nil {
		return "", "", fmt.Errorf("failed to parse AI response: %w", err)
	}

	// Validate and clean title
	content.Title = strings.TrimSpace(content.Title)
	if content.Title == "" {
		return "", "", fmt.Errorf("AI generated empty PR title")
	}

	// Clean description (optional)
	content.Description = strings.TrimSpace(content.Description)

	return content.Title, content.Description, nil
}

// TestConnection makes a simple API call to verify credentials work
func (c *Client) TestConnection() error {
	// Make a simple test request
	testPrompt := "Respond with exactly: OK"
	_, err := c.callAPI(testPrompt)
	return err
}

// callAPI makes a request to the OpenAI-compatible API
func (c *Client) callAPI(prompt string) (string, error) {
	req := ChatRequest{
		Model: c.model,
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.3, // Low temperature for deterministic output
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/chat/completions", c.baseURL),
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return "", WrapError("API request failed", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API error
	if chatResp.Error != nil {
		return "", chatResp.Error
	}

	// Check for HTTP error status
	if resp.StatusCode != http.StatusOK {
		return "", NewRequestError(resp.StatusCode, string(body))
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	// Clean up response: remove markdown code block formatting if present
	content := chatResp.Choices[0].Message.Content
	content = strings.TrimSpace(content)

	// Remove markdown code block delimiters (```json ... ``` or ``` ... ```)
	if strings.HasPrefix(content, "```") {
		// Remove opening ``` with optional language specifier
		if idx := strings.Index(content, "\n"); idx != -1 {
			content = content[idx+1:]
		} else {
			content = strings.TrimPrefix(content, "```json")
			content = strings.TrimPrefix(content, "```")
		}
		// Remove closing ```
		if strings.HasSuffix(content, "```") {
			content = strings.TrimSuffix(content, "```")
		}
		content = strings.TrimSpace(content)
	}

	return content, nil
}
