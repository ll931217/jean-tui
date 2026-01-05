package openai

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Helper function to create a simple chat response for testing
func newChatResponse(content string) map[string]interface{} {
	return map[string]interface{}{
		"choices": []map[string]interface{}{
			{
				"message": map[string]string{
					"content": content,
				},
			},
		},
	}
}

// Helper function to create a chat response with error
func newChatErrorResponse(errorType, message string) map[string]interface{} {
	return map[string]interface{}{
		"error": map[string]string{
			"type":    errorType,
			"message": message,
		},
	}
}

// Helper function to create an empty choices response
func newEmptyChoicesResponse() map[string]interface{} {
	return map[string]interface{}{
		"choices": []map[string]interface{}{},
	}
}

// TestNewClient_ValidInputs tests creating a client with valid inputs
func TestNewClient_ValidInputs(t *testing.T) {
	tests := []struct {
		name    string
		apiKey  string
		baseURL string
		model   string
	}{
		{
			name:    "OpenAI client",
			apiKey:  "sk-test-key",
			baseURL: "https://api.openai.com/v1",
			model:   "gpt-4",
		},
		{
			name:    "Custom provider",
			apiKey:  "custom-key",
			baseURL: "http://localhost:11434/v1",
			model:   "llama2",
		},
		{
			name:    "Minimal valid inputs",
			apiKey:  "x",
			baseURL: "http://example.com",
			model:   "m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.apiKey, tt.baseURL, tt.model)
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}

			if client == nil {
				t.Fatal("NewClient() returned nil client")
			}

			if client.apiKey != tt.apiKey {
				t.Errorf("client.apiKey = %q, want %q", client.apiKey, tt.apiKey)
			}

			if client.baseURL != tt.baseURL {
				t.Errorf("client.baseURL = %q, want %q", client.baseURL, tt.baseURL)
			}

			if client.model != tt.model {
				t.Errorf("client.model = %q, want %q", client.model, tt.model)
			}
		})
	}
}

// TestNewClient_InvalidInputs tests creating a client with invalid inputs
func TestNewClient_InvalidInputs(t *testing.T) {
	tests := []struct {
		name        string
		apiKey      string
		baseURL     string
		model       string
		expectedErr *ConfigError
	}{
		{
			name:        "Missing API key",
			apiKey:      "",
			baseURL:     "http://example.com",
			model:       "gpt-4",
			expectedErr: &ConfigError{Field: "api_key", Message: "API key is required"},
		},
		{
			name:        "Missing base URL",
			apiKey:      "sk-test",
			baseURL:     "",
			model:       "gpt-4",
			expectedErr: &ConfigError{Field: "base_url", Message: "base URL is required"},
		},
		{
			name:        "Missing model",
			apiKey:      "sk-test",
			baseURL:     "http://example.com",
			model:       "",
			expectedErr: &ConfigError{Field: "model", Message: "model is required"},
		},
		{
			name:        "All fields missing",
			apiKey:      "",
			baseURL:     "",
			model:       "",
			expectedErr: &ConfigError{Field: "api_key", Message: "API key is required"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.apiKey, tt.baseURL, tt.model)

			if err == nil {
				t.Fatal("NewClient() expected error, got nil")
			}

			configErr, ok := err.(*ConfigError)
			if !ok {
				t.Fatalf("NewClient() error type = %T, want *ConfigError", err)
			}

			if configErr.Field != tt.expectedErr.Field {
				t.Errorf("error.Field = %q, want %q", configErr.Field, tt.expectedErr.Field)
			}

			if configErr.Message != tt.expectedErr.Message {
				t.Errorf("error.Message = %q, want %q", configErr.Message, tt.expectedErr.Message)
			}

			if client != nil {
				t.Error("NewClient() returned non-nil client for invalid inputs")
			}
		})
	}
}

// TestGenerateCommitMessage_Success tests successful commit message generation
func TestGenerateCommitMessage_Success(t *testing.T) {
	tests := []struct {
		name           string
		status         string
		diff           string
		branch         string
		log            string
		customPrompt   string
		apiResponse    string
		expectedCommit string
	}{
		{
			name:           "Standard commit",
			status:         "M file.go",
			diff:           "+ new line",
			branch:         "feature/test",
			log:            "Previous commit",
			apiResponse:    "feat: add new feature",
			expectedCommit: "feat: add new feature",
		},
		{
			name:           "Fix commit",
			status:         "M bug.go",
			diff:           "- broken\n+ fixed",
			branch:         "fix/bug",
			log:            "",
			apiResponse:    "fix: resolve null pointer",
			expectedCommit: "fix: resolve null pointer",
		},
		{
			name:           "With whitespace in response",
			status:         "A file.go",
			diff:           "+ new file",
			branch:         "main",
			log:            "",
			apiResponse:    "  docs: update readme  ",
			expectedCommit: "docs: update readme",
		},
		{
			name:           "Using custom prompt",
			status:         "M test.go",
			diff:           "+ test",
			branch:         "test",
			log:            "",
			customPrompt:   "Generate commit for: {diff}",
			apiResponse:    "test: add unit test",
			expectedCommit: "test: add unit test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request method
				if r.Method != "POST" {
					t.Errorf("Method = %q, want POST", r.Method)
				}

				// Verify request path
				if r.URL.Path != "/chat/completions" {
					t.Errorf("Path = %q, want /chat/completions", r.URL.Path)
				}

				// Verify authorization header
				auth := r.Header.Get("Authorization")
				if !strings.HasPrefix(auth, "Bearer ") {
					t.Errorf("Authorization header = %q, want Bearer token", auth)
				}

				// Verify content type
				ct := r.Header.Get("Content-Type")
				if ct != "application/json" {
					t.Errorf("Content-Type = %q, want application/json", ct)
				}

				// Parse request body
				var req ChatRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Fatalf("Failed to decode request: %v", err)
				}

				if req.Model != "gpt-4" {
					t.Errorf("Model = %q, want gpt-4", req.Model)
				}

				// Send response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(newChatResponse(tt.apiResponse))
			}))
			defer server.Close()

			client, err := NewClient("test-key", server.URL, "gpt-4")
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}

			result, err := client.GenerateCommitMessage(tt.status, tt.diff, tt.branch, tt.log, tt.customPrompt)
			if err != nil {
				t.Fatalf("GenerateCommitMessage() error = %v", err)
			}

			if result != tt.expectedCommit {
				t.Errorf("GenerateCommitMessage() = %q, want %q", result, tt.expectedCommit)
			}
		})
	}
}

// TestGenerateCommitMessage_LongDiff tests that diff is truncated
func TestGenerateCommitMessage_LongDiff(t *testing.T) {
	longDiff := strings.Repeat("a", 6000)

	serverCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverCalled = true

		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		// Verify diff was truncated
		prompt := req.Messages[0].Content
		if !strings.Contains(prompt, "{diff}") {
			// Placeholder should be replaced
			// The diff should be truncated to 5000 chars, but the prompt includes
			// other text, so we just verify it's not the original 6000 chars
			if len(prompt) > 5800 {
				t.Errorf("Prompt length = %d, expect diff to be truncated to 5000", len(prompt))
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(newChatResponse("test: commit"))
	}))
	defer server.Close()

	client, _ := NewClient("test-key", server.URL, "gpt-4")
	_, err := client.GenerateCommitMessage("M file", longDiff, "main", "", "")

	if err != nil {
		t.Fatalf("GenerateCommitMessage() error = %v", err)
	}

	if !serverCalled {
		t.Error("Server was not called")
	}
}

// TestGenerateCommitMessage_EmptyResponse tests error on empty AI response
func TestGenerateCommitMessage_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(newChatResponse("   "))
	}))
	defer server.Close()

	client, _ := NewClient("test-key", server.URL, "gpt-4")
	_, err := client.GenerateCommitMessage("M file", "+ change", "main", "", "")

	if err == nil {
		t.Fatal("GenerateCommitMessage() expected error for empty response, got nil")
	}

	expectedErrMsg := "AI generated empty commit subject"
	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("Error = %v, want containing %q", err, expectedErrMsg)
	}
}

// TestGenerateBranchName_Success tests successful branch name generation
func TestGenerateBranchName_Success(t *testing.T) {
	tests := []struct {
		name           string
		diff           string
		customPrompt   string
		apiResponse    string
		expectedBranch string
	}{
		{
			name:           "Simple feature branch",
			diff:           "+ new feature",
			apiResponse:    "feat-user-auth",
			expectedBranch: "feat-user-auth",
		},
		{
			name:           "Branch with spaces",
			diff:           "+ change",
			apiResponse:    "feat add user login",
			expectedBranch: "feat-add-user-login",
		},
		{
			name:           "Branch with underscores",
			diff:           "+ change",
			apiResponse:    "fix_login_bug",
			expectedBranch: "fix-login-bug",
		},
		{
			name:           "Branch with mixed case",
			diff:           "+ change",
			apiResponse:    "Refactor-API-Client",
			expectedBranch: "refactor-api-client",
		},
		{
			name:           "Branch with special characters",
			diff:           "+ change",
			apiResponse:    "feat: add user authentication!!!",
			expectedBranch: "feat-add-user-authentication",
		},
		{
			name:           "Branch with numbers",
			diff:           "+ change",
			apiResponse:    "fix-issue-123",
			expectedBranch: "fix-issue-123",
		},
		{
			name:           "Long branch name gets truncated",
			diff:           "+ change",
			apiResponse:    strings.Repeat("a", 50),
			expectedBranch: strings.Repeat("a", 40),
		},
		{
			name:           "Leading/trailing hyphens removed",
			diff:           "+ change",
			apiResponse:    "-test-branch-",
			expectedBranch: "test-branch",
		},
		{
			name:           "Using custom prompt",
			diff:           "+ new api",
			customPrompt:   "Name this: {diff}",
			apiResponse:    "api-integration",
			expectedBranch: "api-integration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(newChatResponse(tt.apiResponse))
			}))
			defer server.Close()

			client, _ := NewClient("test-key", server.URL, "gpt-4")
			result, err := client.GenerateBranchName(tt.diff, tt.customPrompt)

			if err != nil {
				t.Fatalf("GenerateBranchName() error = %v", err)
			}

			if result != tt.expectedBranch {
				t.Errorf("GenerateBranchName() = %q, want %q", result, tt.expectedBranch)
			}
		})
	}
}

// TestGenerateBranchName_InvalidResponse tests error on invalid branch name
func TestGenerateBranchName_InvalidResponse(t *testing.T) {
	tests := []struct {
		name        string
		apiResponse string
	}{
		{
			name:        "Empty response",
			apiResponse: "",
		},
		{
			name:        "Only special characters",
			apiResponse: "!!!@@@###",
		},
		{
			name:        "Only whitespace",
			apiResponse: "   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(newChatResponse(tt.apiResponse))
			}))
			defer server.Close()

			client, _ := NewClient("test-key", server.URL, "gpt-4")
			_, err := client.GenerateBranchName("+ change", "")

			if err == nil {
				t.Fatal("GenerateBranchName() expected error, got nil")
			}

			expectedErrMsg := "AI generated invalid branch name"
			if !strings.Contains(err.Error(), expectedErrMsg) {
				t.Errorf("Error = %v, want containing %q", err, expectedErrMsg)
			}
		})
	}
}

// TestGeneratePRContent_Success tests successful PR content generation
func TestGeneratePRContent_Success(t *testing.T) {
	tests := []struct {
		name          string
		diff          string
		customPrompt  string
		apiResponse   string
		expectedTitle string
		expectedDesc  string
	}{
		{
			name:         "Standard PR",
			diff:         "+ new feature",
			apiResponse:  `{"title": "Add user authentication", "description": "## What's Changed\n\n### Features\n- Added login system"}`,
			expectedTitle: "Add user authentication",
			expectedDesc:  "## What's Changed\n\n### Features\n- Added login system",
		},
		{
			name:         "With markdown formatting",
			diff:         "+ code",
			apiResponse:  `{"title": "Fix bug", "description": "Fixed critical bug in production"}`,
			expectedTitle: "Fix bug",
			expectedDesc:  "Fixed critical bug in production",
		},
		{
			name:         "Empty description is valid",
			diff:         "+ change",
			apiResponse:  `{"title": "Update docs", "description": ""}`,
			expectedTitle: "Update docs",
			expectedDesc:  "",
		},
		{
			name:         "Using custom prompt",
			diff:         "+ api",
			customPrompt: "PR for: {diff}",
			apiResponse:  `{"title": "API Integration", "description": "Adds external API support"}`,
			expectedTitle: "API Integration",
			expectedDesc:  "Adds external API support",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(newChatResponse(tt.apiResponse))
			}))
			defer server.Close()

			client, _ := NewClient("test-key", server.URL, "gpt-4")
			title, desc, err := client.GeneratePRContent(tt.diff, tt.customPrompt)

			if err != nil {
				t.Fatalf("GeneratePRContent() error = %v", err)
			}

			if title != tt.expectedTitle {
				t.Errorf("GeneratePRContent() title = %q, want %q", title, tt.expectedTitle)
			}

			if desc != tt.expectedDesc {
				t.Errorf("GeneratePRContent() description = %q, want %q", desc, tt.expectedDesc)
			}
		})
	}
}

// TestGeneratePRContent_InvalidJSON tests error on malformed JSON response
func TestGeneratePRContent_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(newChatResponse("invalid json"))
	}))
	defer server.Close()

	client, _ := NewClient("test-key", server.URL, "gpt-4")
	_, _, err := client.GeneratePRContent("+ change", "")

	if err == nil {
		t.Fatal("GeneratePRContent() expected error for invalid JSON, got nil")
	}

	expectedErrMsg := "failed to parse AI response"
	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("Error = %v, want containing %q", err, expectedErrMsg)
	}
}

// TestGeneratePRContent_EmptyTitle tests error on empty title
func TestGeneratePRContent_EmptyTitle(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(newChatResponse(`{"title": "", "description": "desc"}`))
	}))
	defer server.Close()

	client, _ := NewClient("test-key", server.URL, "gpt-4")
	_, _, err := client.GeneratePRContent("+ change", "")

	if err == nil {
		t.Fatal("GeneratePRContent() expected error for empty title, got nil")
	}

	expectedErrMsg := "AI generated empty PR title"
	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("Error = %v, want containing %q", err, expectedErrMsg)
	}
}

// TestGeneratePRContent_MarkdownCleanup tests that markdown code blocks are removed
func TestGeneratePRContent_MarkdownCleanup(t *testing.T) {
	tests := []struct {
		name        string
		apiResponse string
		expected    string
	}{
		{
			name:        "JSON in markdown code block",
			apiResponse: "```json\n{\"title\": \"Test\", \"description\": \"Desc\"}\n```",
			expected:    `{"title": "Test", "description": "Desc"}`,
		},
		{
			name:        "Plain markdown code block",
			apiResponse: "```\n{\"title\": \"Test\", \"description\": \"Desc\"}\n```",
			expected:    `{"title": "Test", "description": "Desc"}`,
		},
		{
			name:        "Markdown without newline after opening",
			apiResponse: "```json{\"title\": \"Test\", \"description\": \"Desc\"}```",
			expected:    `{"title": "Test", "description": "Desc"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(newChatResponse(tt.apiResponse))
			}))
			defer server.Close()

			client, _ := NewClient("test-key", server.URL, "gpt-4")
			_, _, err := client.GeneratePRContent("+ change", "")

			if err != nil {
				t.Fatalf("GeneratePRContent() error = %v", err)
			}
		})
	}
}

// TestTestConnection_Success tests successful connection test
func TestTestConnection_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(newChatResponse("OK"))
	}))
	defer server.Close()

	client, _ := NewClient("test-key", server.URL, "gpt-4")
	err := client.TestConnection()

	if err != nil {
		t.Errorf("TestConnection() error = %v", err)
	}
}

// TestTestConnection_Failure tests connection test failures
func TestTestConnection_Failure(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   string
	}{
		{
			name:       "HTTP 401 Unauthorized",
			statusCode: http.StatusUnauthorized,
			response:   `{"error": {"type": "authentication_error", "message": "Invalid API key"}}`,
		},
		{
			name:       "HTTP 429 Rate Limit",
			statusCode: http.StatusTooManyRequests,
			response:   `{"error": {"type": "rate_limit_error", "message": "Too many requests"}}`,
		},
		{
			name:       "HTTP 500 Server Error",
			statusCode: http.StatusInternalServerError,
			response:   "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client, _ := NewClient("test-key", server.URL, "gpt-4")
			err := client.TestConnection()

			if err == nil {
				t.Fatal("TestConnection() expected error, got nil")
			}
		})
	}
}

// TestCallAPI_APIError tests API error responses
func TestCallAPI_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(newChatErrorResponse("invalid_request_error", "Invalid model specified"))
	}))
	defer server.Close()

	client, _ := NewClient("test-key", server.URL, "invalid-model")
	_, err := client.callAPI("test")

	if err == nil {
		t.Fatal("callAPI() expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("Error type = %T, want *APIError", err)
	}

	if apiErr.Type != "invalid_request_error" {
		t.Errorf("Error.Type = %q, want 'invalid_request_error'", apiErr.Type)
	}

	if apiErr.Message != "Invalid model specified" {
		t.Errorf("Error.Message = %q, want 'Invalid model specified'", apiErr.Message)
	}
}

// TestCallAPI_NoChoices tests error when API returns no choices
func TestCallAPI_NoChoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(newEmptyChoicesResponse())
	}))
	defer server.Close()

	client, _ := NewClient("test-key", server.URL, "gpt-4")
	_, err := client.callAPI("test")

	if err == nil {
		t.Fatal("callAPI() expected error, got nil")
	}

	expectedErrMsg := "no response from API"
	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("Error = %v, want containing %q", err, expectedErrMsg)
	}
}

// TestCallAPI_NetworkError tests handling of network errors
func TestCallAPI_NetworkError(t *testing.T) {
	// Create a client with an invalid URL that will fail to connect
	client, _ := NewClient("test-key", "http://localhost:9999", "gpt-4")
	_, err := client.callAPI("test")

	if err == nil {
		t.Fatal("callAPI() expected error for network failure, got nil")
	}

	// Error should be wrapped
	if !strings.Contains(err.Error(), "generate commit message") && !strings.Contains(err.Error(), "API request failed") {
		t.Logf("Error = %v", err)
	}
}

// TestWrapError tests error wrapping
func TestWrapError(t *testing.T) {
	tests := []struct {
		name          string
		operation     string
		err           error
		wantContains  string
	}{
		{
			name:         "Wrap standard error",
			operation:    "generate commit",
			err:          fmt.Errorf("base error"),
			wantContains: "generate commit",
		},
		{
			name:         "Nil error",
			operation:    "test",
			err:          nil,
			wantContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapError(tt.operation, tt.err)

			if tt.err == nil {
				if result != nil {
					t.Errorf("WrapError() = %v, want nil for nil error", result)
				}
				return
			}

			if result == nil {
				t.Fatal("WrapError() returned nil")
			}

			if !strings.Contains(result.Error(), tt.operation) {
				t.Errorf("Error = %v, want containing %q", result, tt.wantContains)
			}
		})
	}
}

// TestMarkdownCleanup tests markdown code block removal in various formats
func TestMarkdownCleanup(t *testing.T) {
	tests := []struct {
		name     string
		response string
		expected string
	}{
		{
			name:     "JSON with markdown wrapper",
			response: "```json\n{\"key\": \"value\"}\n```",
			expected: "{\"key\": \"value\"}",
		},
		{
			name:     "Plain markdown without language",
			response: "```\nplain text\n```",
			expected: "plain text",
		},
		{
			name:     "Markdown without newline",
			response: "```json\ntext```",
			expected: "text",
		},
		{
			name:     "No markdown - plain text",
			response: "just plain text",
			expected: "just plain text",
		},
		{
			name:     "Markdown with extra whitespace",
			response: "```json  \n  text  \n  ```  ",
			expected: "text",
		},
		{
			name:     "Multi-line content in markdown",
			response: "```\nline1\nline2\nline3\n```",
			expected: "line1\nline2\nline3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(newChatResponse(tt.response))
			}))
			defer server.Close()

			client, _ := NewClient("test-key", server.URL, "gpt-4")
			result, err := client.callAPI("test")

			if err != nil {
				t.Fatalf("callAPI() error = %v", err)
			}

			if result != tt.expected {
				t.Errorf("callAPI() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestRequestHeaders tests that proper headers are sent
func TestRequestHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check Content-Type
		ct := r.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", ct)
		}

		// Check Authorization
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-api-key" {
			t.Errorf("Authorization = %q, want Bearer test-api-key", auth)
		}

		// Check method
		if r.Method != "POST" {
			t.Errorf("Method = %q, want POST", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(newChatResponse("OK"))
	}))
	defer server.Close()

	client, _ := NewClient("test-api-key", server.URL, "gpt-4")
	_, err := client.callAPI("test")

	if err != nil {
		t.Fatalf("callAPI() error = %v", err)
	}
}

// TestTemperatureSetting tests that temperature is set correctly
func TestTemperatureSetting(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		// Verify low temperature for deterministic output
		if req.Temperature != 0.3 {
			t.Errorf("Temperature = %v, want 0.3", req.Temperature)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(newChatResponse("OK"))
	}))
	defer server.Close()

	client, _ := NewClient("test-key", server.URL, "gpt-4")
	_, err := client.callAPI("test")

	if err != nil {
		t.Fatalf("callAPI() error = %v", err)
	}
}
