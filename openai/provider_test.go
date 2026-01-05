package openai

import (
	"encoding/json"
	"testing"
)

// TestDefaultBaseURL tests default base URL for each provider type
func TestDefaultBaseURL(t *testing.T) {
	tests := []struct {
		name            string
		providerType    ProviderType
		expectedBaseURL string
	}{
		{
			name:            "OpenAI provider",
			providerType:    ProviderOpenAI,
			expectedBaseURL: "https://api.openai.com/v1",
		},
		{
			name:            "Azure provider",
			providerType:    ProviderAzure,
			expectedBaseURL: "",
		},
		{
			name:            "Custom provider",
			providerType:    ProviderCustom,
			expectedBaseURL: "",
		},
		{
			name:            "Unknown provider defaults to OpenAI",
			providerType:    ProviderType("unknown"),
			expectedBaseURL: "https://api.openai.com/v1",
		},
		{
			name:            "Empty provider type",
			providerType:    ProviderType(""),
			expectedBaseURL: "https://api.openai.com/v1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DefaultBaseURL(tt.providerType)
			if result != tt.expectedBaseURL {
				t.Errorf("DefaultBaseURL(%q) = %q, want %q", tt.providerType, result, tt.expectedBaseURL)
			}
		})
	}
}

// TestIsValidProvider tests provider type validation
func TestIsValidProvider(t *testing.T) {
	tests := []struct {
		name        string
		provider    ProviderType
		expectedValid bool
	}{
		{
			name:        "OpenAI is valid",
			provider:    ProviderOpenAI,
			expectedValid: true,
		},
		{
			name:        "Azure is valid",
			provider:    ProviderAzure,
			expectedValid: true,
		},
		{
			name:        "Custom is valid",
			provider:    ProviderCustom,
			expectedValid: true,
		},
		{
			name:        "Unknown provider is invalid",
			provider:    ProviderType("unknown"),
			expectedValid: false,
		},
		{
			name:        "Empty provider is invalid",
			provider:    ProviderType(""),
			expectedValid: false,
		},
		{
			name:        "Arbitrary string is invalid",
			provider:    ProviderType("something-else"),
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidProvider(tt.provider)
			if result != tt.expectedValid {
				t.Errorf("IsValidProvider(%q) = %v, want %v", tt.provider, result, tt.expectedValid)
			}
		})
	}
}

// TestNewProviderConfig tests creating provider config with defaults
func TestNewProviderConfig(t *testing.T) {
	tests := []struct {
		name            string
		providerType    ProviderType
		expectedType    ProviderType
		expectedBaseURL string
	}{
		{
			name:            "OpenAI config",
			providerType:    ProviderOpenAI,
			expectedType:    ProviderOpenAI,
			expectedBaseURL: "https://api.openai.com/v1",
		},
		{
			name:            "Azure config",
			providerType:    ProviderAzure,
			expectedType:    ProviderAzure,
			expectedBaseURL: "",
		},
		{
			name:            "Custom config",
			providerType:    ProviderCustom,
			expectedType:    ProviderCustom,
			expectedBaseURL: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewProviderConfig(tt.providerType)

			if config.Type != tt.expectedType {
				t.Errorf("config.Type = %q, want %q", config.Type, tt.expectedType)
			}

			if config.BaseURL != tt.expectedBaseURL {
				t.Errorf("config.BaseURL = %q, want %q", config.BaseURL, tt.expectedBaseURL)
			}

			// Other fields should be empty
			if config.APIKey != "" {
				t.Errorf("config.APIKey = %q, want empty string", config.APIKey)
			}

			if config.Model != "" {
				t.Errorf("config.Model = %q, want empty string", config.Model)
			}
		})
	}
}

// TestGetEffectiveBaseURL tests base URL resolution with and without override
func TestGetEffectiveBaseURL(t *testing.T) {
	tests := []struct {
		name            string
		configType      ProviderType
		configBaseURL   string
		expectedURL     string
	}{
		{
			name:            "OpenAI with default URL",
			configType:      ProviderOpenAI,
			configBaseURL:   "",
			expectedURL:     "https://api.openai.com/v1",
		},
		{
			name:            "OpenAI with custom override",
			configType:      ProviderOpenAI,
			configBaseURL:   "https://custom.openai.com/v1",
			expectedURL:     "https://custom.openai.com/v1",
		},
		{
			name:            "Azure with custom URL",
			configType:      ProviderAzure,
			configBaseURL:   "https://azure-resource.openai.azure.com",
			expectedURL:     "https://azure-resource.openai.azure.com",
		},
		{
			name:            "Azure with empty URL (invalid but tests fallback)",
			configType:      ProviderAzure,
			configBaseURL:   "",
			expectedURL:     "",
		},
		{
			name:            "Custom with URL",
			configType:      ProviderCustom,
			configBaseURL:   "http://localhost:11434/v1",
			expectedURL:     "http://localhost:11434/v1",
		},
		{
			name:            "Custom with empty URL",
			configType:      ProviderCustom,
			configBaseURL:   "",
			expectedURL:     "",
		},
		{
			name:            "Unknown provider defaults to OpenAI",
			configType:      ProviderType("unknown"),
			configBaseURL:   "",
			expectedURL:     "https://api.openai.com/v1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &ProviderConfig{
				Type:    tt.configType,
				BaseURL: tt.configBaseURL,
			}

			result := config.GetEffectiveBaseURL()
			if result != tt.expectedURL {
				t.Errorf("GetEffectiveBaseURL() = %q, want %q", result, tt.expectedURL)
			}
		})
	}
}

// TestDefaultModels tests default model lists for each provider type
func TestDefaultModels(t *testing.T) {
	tests := []struct {
		name                string
		providerType        string
		expectedModelCount  int
		shouldContain       []string
	}{
		{
			name:                "OpenAI models",
			providerType:        "openai",
			expectedModelCount:  5,
			shouldContain:       []string{ModelGPT4, ModelGPT4Turbo, ModelGPT4o, ModelGPT4oMini, ModelGPT35Turbo},
		},
		{
			name:                "Azure models",
			providerType:        "azure",
			expectedModelCount:  2,
			shouldContain:       []string{ModelAzureGPT4, ModelAzureGPT35Turbo},
		},
		{
			name:                "Custom models (empty list)",
			providerType:        "custom",
			expectedModelCount:  0,
			shouldContain:       []string{},
		},
		{
			name:                "Unknown provider defaults to safe list",
			providerType:        "unknown",
			expectedModelCount:  2,
			shouldContain:       []string{ModelGPT4, ModelGPT35Turbo},
		},
		{
			name:                "Empty provider type",
			providerType:        "",
			expectedModelCount:  2,
			shouldContain:       []string{ModelGPT4, ModelGPT35Turbo},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			models := DefaultModels(tt.providerType)

			if len(models) != tt.expectedModelCount {
				t.Errorf("DefaultModels(%q) returned %d models, want %d", tt.providerType, len(models), tt.expectedModelCount)
			}

			for _, expected := range tt.shouldContain {
				found := false
				for _, model := range models {
					if model == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("DefaultModels(%q) should contain %q", tt.providerType, expected)
				}
			}
		})
	}
}

// TestIsValidModel tests model validation
func TestIsValidModel(t *testing.T) {
	tests := []struct {
		name            string
		providerType    string
		model           string
		expectedValid   bool
	}{
		// Valid OpenAI models
		{
			name:          "Valid OpenAI GPT-4",
			providerType:  "openai",
			model:         ModelGPT4,
			expectedValid: true,
		},
		{
			name:          "Valid OpenAI GPT-4 Turbo",
			providerType:  "openai",
			model:         ModelGPT4Turbo,
			expectedValid: true,
		},
		{
			name:          "Valid OpenAI GPT-4o",
			providerType:  "openai",
			model:         ModelGPT4o,
			expectedValid: true,
		},
		{
			name:          "Valid OpenAI GPT-4o Mini",
			providerType:  "openai",
			model:         ModelGPT4oMini,
			expectedValid: true,
		},
		{
			name:          "Valid OpenAI GPT-3.5 Turbo",
			providerType:  "openai",
			model:         ModelGPT35Turbo,
			expectedValid: true,
		},
		// Valid Azure models
		{
			name:          "Valid Azure GPT-4",
			providerType:  "azure",
			model:         ModelAzureGPT4,
			expectedValid: true,
		},
		{
			name:          "Valid Azure GPT-3.5 Turbo",
			providerType:  "azure",
			model:         ModelAzureGPT35Turbo,
			expectedValid: true,
		},
		// Custom models are always valid if non-empty
		{
			name:          "Custom model name for OpenAI",
			providerType:  "openai",
			model:         "gpt-4-turbo-preview",
			expectedValid: true,
		},
		{
			name:          "Custom model name for Azure",
			providerType:  "azure",
			model:         "my-custom-deployment",
			expectedValid: true,
		},
		{
			name:          "Custom model for custom provider",
			providerType:  "custom",
			model:         "llama2-13b",
			expectedValid: true,
		},
		{
			name:          "Ollama model",
			providerType:  "custom",
			model:         "ollama/llama2",
			expectedValid: true,
		},
		// Invalid cases
		{
			name:          "Empty model string",
			providerType:  "openai",
			model:         "",
			expectedValid: false,
		},
		{
			name:          "Model not in predefined list but non-empty",
			providerType:  "openai",
			model:         "some-random-model",
			expectedValid: true, // Custom model names are allowed
		},
		{
			name:          "Unknown provider with non-empty model",
			providerType:  "unknown",
			model:         "any-model",
			expectedValid: true,
		},
		{
			name:          "Unknown provider with empty model",
			providerType:  "unknown",
			model:         "",
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidModel(tt.providerType, tt.model)
			if result != tt.expectedValid {
				t.Errorf("IsValidModel(%q, %q) = %v, want %v", tt.providerType, tt.model, result, tt.expectedValid)
			}
		})
	}
}

// TestGetDefaultCommitPrompt tests retrieving the default commit prompt
func TestGetDefaultCommitPrompt(t *testing.T) {
	result := GetDefaultCommitPrompt()

	if result == "" {
		t.Error("GetDefaultCommitPrompt() returned empty string")
	}

	// Verify it contains expected placeholders
	expectedPlaceholders := []string{"{status}", "{diff}", "{branch}", "{log}"}
	for _, placeholder := range expectedPlaceholders {
		if !contains(result, placeholder) {
			t.Errorf("GetDefaultCommitPrompt() should contain placeholder %q", placeholder)
		}
	}

	// Verify it mentions conventional commits
	if !contains(result, "Conventional Commits") {
		t.Error("GetDefaultCommitPrompt() should mention Conventional Commits")
	}

	// Verify it has examples
	if !contains(result, "feat:") || !contains(result, "fix:") {
		t.Error("GetDefaultCommitPrompt() should have examples")
	}
}

// TestGetDefaultBranchNamePrompt tests retrieving the default branch name prompt
func TestGetDefaultBranchNamePrompt(t *testing.T) {
	result := GetDefaultBranchNamePrompt()

	if result == "" {
		t.Error("GetDefaultBranchNamePrompt() returned empty string")
	}

	// Verify it contains expected placeholders
	if !contains(result, "{diff}") {
		t.Error("GetDefaultBranchNamePrompt() should contain {diff} placeholder")
	}

	// Verify it mentions constraints
	if !contains(result, "kebab-case") || !contains(result, "40 characters") {
		t.Error("GetDefaultBranchNamePrompt() should mention formatting constraints")
	}

	// Verify it has examples
	if !contains(result, "fix-login-bug") {
		t.Error("GetDefaultBranchNamePrompt() should have examples")
	}
}

// TestGetDefaultPRPrompt tests retrieving the default PR prompt
func TestGetDefaultPRPrompt(t *testing.T) {
	result := GetDefaultPRPrompt()

	if result == "" {
		t.Error("GetDefaultPRPrompt() returned empty string")
	}

	// Verify it contains expected placeholders
	if !contains(result, "{diff}") {
		t.Error("GetDefaultPRPrompt() should contain {diff} placeholder")
	}

	// Verify it specifies JSON format
	if !contains(result, `{"title": "...", "description": "..."}`) {
		t.Error("GetDefaultPRPrompt() should specify JSON format")
	}

	// Verify it mentions title length limit
	if !contains(result, "72 characters") {
		t.Error("GetDefaultPRPrompt() should mention title length limit")
	}

	// Verify it mentions "What's Changed" format
	if !contains(result, "What's Changed") {
		t.Error("GetDefaultPRPrompt() should mention 'What's Changed' format")
	}

	// Verify it has example
	if !contains(result, "Add dark mode support") {
		t.Error("GetDefaultPRPrompt() should have example")
	}
}

// TestDefaultPrompts_NotEmpty tests that all default prompts are defined
func TestDefaultPrompts_NotEmpty(t *testing.T) {
	if DefaultCommitPrompt == "" {
		t.Error("DefaultCommitPrompt should not be empty")
	}

	if DefaultBranchNamePrompt == "" {
		t.Error("DefaultBranchNamePrompt should not be empty")
	}

	if DefaultPRPrompt == "" {
		t.Error("DefaultPRPrompt should not be empty")
	}
}

// TestProviderTypeString tests provider type string representations
func TestProviderTypeString(t *testing.T) {
	tests := []struct {
		name     string
		provider ProviderType
		expected string
	}{
		{
			name:     "OpenAI string",
			provider: ProviderOpenAI,
			expected: "openai",
		},
		{
			name:     "Azure string",
			provider: ProviderAzure,
			expected: "azure",
		},
		{
			name:     "Custom string",
			provider: ProviderCustom,
			expected: "custom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := string(tt.provider)
			if result != tt.expected {
				t.Errorf("ProviderType string = %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestModelConstants tests that model constants are defined
func TestModelConstants(t *testing.T) {
	tests := []struct {
		name  string
		model string
	}{
		{"GPT-4", ModelGPT4},
		{"GPT-4 Turbo", ModelGPT4Turbo},
		{"GPT-4o", ModelGPT4o},
		{"GPT-4o Mini", ModelGPT4oMini},
		{"GPT-3.5 Turbo", ModelGPT35Turbo},
		{"Azure GPT-4", ModelAzureGPT4},
		{"Azure GPT-3.5 Turbo", ModelAzureGPT35Turbo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.model == "" {
				t.Errorf("Model constant %q should not be empty", tt.name)
			}
		})
	}
}

// TestProviderConfigJSON tests JSON serialization of ProviderConfig
func TestProviderConfigJSON(t *testing.T) {
	config := &ProviderConfig{
		Type:    ProviderOpenAI,
		BaseURL: "https://api.openai.com/v1",
		APIKey:  "sk-test",
		Model:   "gpt-4",
	}

	// Test JSON tags are correct by marshaling
	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	expectedFields := []string{"type", "base_url", "api_key", "model"}
	for _, field := range expectedFields {
		if !contains(string(data), `"`+field+`"`) {
			t.Errorf("JSON should contain field %q", field)
		}
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestAPIError_Error tests the APIError Error() method
func TestAPIError_Error(t *testing.T) {
	err := &APIError{
		Type:    "test_type",
		Message: "test message",
	}

	expected := "API error (test_type): test message"
	if err.Error() != expected {
		t.Errorf("APIError.Error() = %q, want %q", err.Error(), expected)
	}
}

// TestConfigError_Error tests the ConfigError Error() method
func TestConfigError_Error(t *testing.T) {
	err := &ConfigError{
		Field:   "api_key",
		Message: "is required",
	}

	expected := "configuration error for api_key: is required"
	if err.Error() != expected {
		t.Errorf("ConfigError.Error() = %q, want %q", err.Error(), expected)
	}
}

// TestRequestError_Error tests the RequestError Error() method
func TestRequestError_Error(t *testing.T) {
	tests := []struct {
		name     string
		statusCode int
		message    string
		expected   string
	}{
		{
			name:       "With status code",
			statusCode: 401,
			message:    "Unauthorized",
			expected:   "request failed with status 401: Unauthorized",
		},
		{
			name:       "Without status code",
			statusCode: 0,
			message:    "Connection failed",
			expected:   "request failed: Connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &RequestError{
				StatusCode: tt.statusCode,
				Message:    tt.message,
			}

			if err.Error() != tt.expected {
				t.Errorf("RequestError.Error() = %q, want %q", err.Error(), tt.expected)
			}
		})
	}
}

// TestNewConfigError tests the NewConfigError constructor
func TestNewConfigError(t *testing.T) {
	err := NewConfigError("model", "is invalid")

	if err.Field != "model" {
		t.Errorf("Field = %q, want 'model'", err.Field)
	}

	if err.Message != "is invalid" {
		t.Errorf("Message = %q, want 'is invalid'", err.Message)
	}
}

// TestNewRequestError tests the NewRequestError constructor
func TestNewRequestError(t *testing.T) {
	err := NewRequestError(500, "Server error")

	if err.StatusCode != 500 {
		t.Errorf("StatusCode = %d, want 500", err.StatusCode)
	}

	if err.Message != "Server error" {
		t.Errorf("Message = %q, want 'Server error'", err.Message)
	}
}

