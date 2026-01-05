package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/coollabsio/jean-tui/openai"
)

// Helper function to create a test manager with a temporary config file
func createTestManager(t *testing.T) (*Manager, string) {
	t.Helper()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Create a new manager with the test config path
	m := &Manager{
		configPath: configPath,
		config: &Config{
			Repositories: make(map[string]*RepoConfig),
		},
	}

	return m, tempDir
}

// Helper function to create a test profile
func createTestProfile(name string) *AIProviderProfile {
	return &AIProviderProfile{
		Name:   name,
		Type:   openai.ProviderOpenAI,
		BaseURL: "https://api.openai.com/v1",
		APIKey: "sk-test-key-123",
		Model:  "gpt-4",
	}
}

// TestGetProviderProfiles tests the GetProviderProfiles method
func TestGetProviderProfiles(t *testing.T) {
	t.Run("returns empty map when no profiles configured", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		profiles := m.GetProviderProfiles(repoPath)

		if profiles == nil {
			t.Fatal("Expected non-nil map, got nil")
		}

		if len(profiles) != 0 {
			t.Errorf("Expected empty map, got %d profiles", len(profiles))
		}
	})

	t.Run("returns all profiles when configured", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Add a profile
		profile := createTestProfile("test-profile")
		err := m.AddProviderProfile(repoPath, profile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		profiles := m.GetProviderProfiles(repoPath)

		if len(profiles) != 1 {
			t.Errorf("Expected 1 profile, got %d", len(profiles))
		}

		if _, exists := profiles["test-profile"]; !exists {
			t.Error("Expected 'test-profile' to exist in profiles")
		}
	})

	t.Run("returns empty map for non-existent repository", func(t *testing.T) {
		m, _ := createTestManager(t)

		profiles := m.GetProviderProfiles("/non/existent/repo")

		if profiles == nil {
			t.Fatal("Expected non-nil map, got nil")
		}

		if len(profiles) != 0 {
			t.Errorf("Expected empty map, got %d profiles", len(profiles))
		}
	})
}

// TestAddProviderProfile tests the AddProviderProfile method
func TestAddProviderProfile(t *testing.T) {
	t.Run("adds profile to existing repo config", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Create repo config first
		if m.config.Repositories == nil {
			m.config.Repositories = make(map[string]*RepoConfig)
		}
		m.config.Repositories[repoPath] = &RepoConfig{
			BaseBranch: "main",
		}

		// Add profile
		profile := createTestProfile("openai-primary")
		err := m.AddProviderProfile(repoPath, profile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		// Verify profile was added
		profiles := m.GetProviderProfiles(repoPath)
		if len(profiles) != 1 {
			t.Errorf("Expected 1 profile, got %d", len(profiles))
		}

		retrieved := profiles["openai-primary"]
		if retrieved.Name != "openai-primary" {
			t.Errorf("Expected name 'openai-primary', got '%s'", retrieved.Name)
		}
		if retrieved.APIKey != "sk-test-key-123" {
			t.Errorf("Expected API key 'sk-test-key-123', got '%s'", retrieved.APIKey)
		}
	})

	t.Run("creates AIProvider if nil", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Create repo config without AIProvider
		m.config.Repositories[repoPath] = &RepoConfig{
			BaseBranch: "main",
		}

		// Verify AIProvider is nil
		if m.config.Repositories[repoPath].AIProvider != nil {
			t.Fatal("Expected AIProvider to be nil initially")
		}

		// Add profile
		profile := createTestProfile("test")
		err := m.AddProviderProfile(repoPath, profile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		// Verify AIProvider was created
		if m.config.Repositories[repoPath].AIProvider == nil {
			t.Fatal("Expected AIProvider to be created")
		}
	})

	t.Run("creates profiles map if nil", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Create repo config with AIProvider but nil profiles
		m.config.Repositories[repoPath] = &RepoConfig{
			BaseBranch: "main",
			AIProvider: &AIProviderConfig{
				ActiveProfile:   "",
				FallbackProfile: "",
			},
		}

		// Verify profiles map is nil
		if m.config.Repositories[repoPath].AIProvider.Profiles != nil {
			t.Fatal("Expected Profiles to be nil initially")
		}

		// Add profile
		profile := createTestProfile("test")
		err := m.AddProviderProfile(repoPath, profile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		// Verify profiles map was created
		if m.config.Repositories[repoPath].AIProvider.Profiles == nil {
			t.Fatal("Expected Profiles map to be created")
		}
	})

	t.Run("saves config after adding", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Add profile
		profile := createTestProfile("test")
		err := m.AddProviderProfile(repoPath, profile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		// Verify file was written
		data, err := os.ReadFile(m.configPath)
		if err != nil {
			t.Fatalf("Failed to read config file: %v", err)
		}

		// Verify JSON contains the profile
		jsonStr := string(data)
		if !contains(jsonStr, "test") {
			t.Error("Expected config file to contain profile name")
		}
		if !contains(jsonStr, "sk-test-key-123") {
			t.Error("Expected config file to contain API key")
		}
	})
}

// TestUpdateProviderProfile tests the UpdateProviderProfile method
func TestUpdateProviderProfile(t *testing.T) {
	t.Run("updates existing profile", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Add initial profile
		profile := createTestProfile("test")
		profile.Model = "gpt-3.5-turbo"
		err := m.AddProviderProfile(repoPath, profile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		// Update profile
		updatedProfile := createTestProfile("test")
		updatedProfile.Model = "gpt-4"
		updatedProfile.APIKey = "sk-new-key-456"
		err = m.UpdateProviderProfile(repoPath, updatedProfile)
		if err != nil {
			t.Fatalf("Failed to update profile: %v", err)
		}

		// Verify update
		profiles := m.GetProviderProfiles(repoPath)
		retrieved := profiles["test"]
		if retrieved.Model != "gpt-4" {
			t.Errorf("Expected model 'gpt-4', got '%s'", retrieved.Model)
		}
		if retrieved.APIKey != "sk-new-key-456" {
			t.Errorf("Expected API key 'sk-new-key-456', got '%s'", retrieved.APIKey)
		}
	})

	t.Run("returns error for non-existent profile", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Try to update non-existent profile
		profile := createTestProfile("non-existent")
		err := m.UpdateProviderProfile(repoPath, profile)

		if err == nil {
			t.Fatal("Expected error when updating non-existent profile, got nil")
		}

		expectedErr := "profile 'non-existent' not found"
		if err.Error() != expectedErr {
			t.Errorf("Expected error '%s', got '%s'", expectedErr, err.Error())
		}
	})

	t.Run("returns error for non-existent repository", func(t *testing.T) {
		m, _ := createTestManager(t)

		// Try to update profile in non-existent repo
		profile := createTestProfile("test")
		err := m.UpdateProviderProfile("/non/existent/repo", profile)

		if err == nil {
			t.Fatal("Expected error when updating profile in non-existent repo, got nil")
		}
	})

	t.Run("saves config after updating", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Add profile
		profile := createTestProfile("test")
		err := m.AddProviderProfile(repoPath, profile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		// Update profile
		updatedProfile := createTestProfile("test")
		updatedProfile.Model = "gpt-4-turbo"
		err = m.UpdateProviderProfile(repoPath, updatedProfile)
		if err != nil {
			t.Fatalf("Failed to update profile: %v", err)
		}

		// Verify file was written
		data, err := os.ReadFile(m.configPath)
		if err != nil {
			t.Fatalf("Failed to read config file: %v", err)
		}

		// Verify JSON contains updated value
		jsonStr := string(data)
		if !contains(jsonStr, "gpt-4-turbo") {
			t.Error("Expected config file to contain updated model")
		}
	})
}

// TestDeleteProviderProfile tests the DeleteProviderProfile method
func TestDeleteProviderProfile(t *testing.T) {
	t.Run("deletes existing profile", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Add profile
		profile := createTestProfile("test")
		err := m.AddProviderProfile(repoPath, profile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		// Verify profile exists
		profiles := m.GetProviderProfiles(repoPath)
		if len(profiles) != 1 {
			t.Fatalf("Expected 1 profile before deletion, got %d", len(profiles))
		}

		// Delete profile
		err = m.DeleteProviderProfile(repoPath, "test")
		if err != nil {
			t.Fatalf("Failed to delete profile: %v", err)
		}

		// Verify profile was deleted
		profiles = m.GetProviderProfiles(repoPath)
		if len(profiles) != 0 {
			t.Errorf("Expected 0 profiles after deletion, got %d", len(profiles))
		}

		if _, exists := profiles["test"]; exists {
			t.Error("Expected 'test' profile to be deleted")
		}
	})

	t.Run("clears active profile if deleting active", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Add and set active profile
		profile := createTestProfile("active")
		err := m.AddProviderProfile(repoPath, profile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		err = m.SetActiveProfile(repoPath, "active")
		if err != nil {
			t.Fatalf("Failed to set active profile: %v", err)
		}

		// Verify active profile is set
		if m.GetActiveProfile(repoPath) != "active" {
			t.Fatal("Expected active profile to be 'active'")
		}

		// Delete active profile
		err = m.DeleteProviderProfile(repoPath, "active")
		if err != nil {
			t.Fatalf("Failed to delete profile: %v", err)
		}

		// Verify active profile was cleared
		if m.GetActiveProfile(repoPath) != "" {
			t.Errorf("Expected active profile to be cleared, got '%s'", m.GetActiveProfile(repoPath))
		}
	})

	t.Run("clears fallback profile if deleting fallback", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Add and set fallback profile
		profile := createTestProfile("fallback")
		err := m.AddProviderProfile(repoPath, profile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		err = m.SetFallbackProfile(repoPath, "fallback")
		if err != nil {
			t.Fatalf("Failed to set fallback profile: %v", err)
		}

		// Verify fallback profile is set
		if m.GetFallbackProfile(repoPath) != "fallback" {
			t.Fatal("Expected fallback profile to be 'fallback'")
		}

		// Delete fallback profile
		err = m.DeleteProviderProfile(repoPath, "fallback")
		if err != nil {
			t.Fatalf("Failed to delete profile: %v", err)
		}

		// Verify fallback profile was cleared
		if m.GetFallbackProfile(repoPath) != "" {
			t.Errorf("Expected fallback profile to be cleared, got '%s'", m.GetFallbackProfile(repoPath))
		}
	})

	t.Run("returns error for non-existent profile", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Try to delete non-existent profile
		err := m.DeleteProviderProfile(repoPath, "non-existent")

		if err == nil {
			t.Fatal("Expected error when deleting non-existent profile, got nil")
		}

		expectedErr := "profile 'non-existent' not found"
		if err.Error() != expectedErr {
			t.Errorf("Expected error '%s', got '%s'", expectedErr, err.Error())
		}
	})

	t.Run("returns error for non-existent repository", func(t *testing.T) {
		m, _ := createTestManager(t)

		// Try to delete profile from non-existent repo
		err := m.DeleteProviderProfile("/non/existent/repo", "test")

		if err == nil {
			t.Fatal("Expected error when deleting profile from non-existent repo, got nil")
		}
	})

	t.Run("saves config after deleting", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Add profile
		profile := createTestProfile("test")
		err := m.AddProviderProfile(repoPath, profile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		// Delete profile
		err = m.DeleteProviderProfile(repoPath, "test")
		if err != nil {
			t.Fatalf("Failed to delete profile: %v", err)
		}

		// Verify file was written
		data, err := os.ReadFile(m.configPath)
		if err != nil {
			t.Fatalf("Failed to read config file: %v", err)
		}

		// Verify JSON doesn't contain the deleted profile name
		jsonStr := string(data)
		// The profile name should not appear as a key in profiles
		if contains(jsonStr, `"test":`) && contains(jsonStr, `"profiles"`) {
			// Check if it's actually in the profiles section
			if contains(jsonStr, `"profiles":{`) && contains(jsonStr, `"test":{`) {
				t.Error("Expected config file to not contain deleted profile")
			}
		}
	})
}

// TestGetActiveProfile tests the GetActiveProfile method
func TestGetActiveProfile(t *testing.T) {
	t.Run("returns empty string when not set", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		active := m.GetActiveProfile(repoPath)
		if active != "" {
			t.Errorf("Expected empty string, got '%s'", active)
		}
	})

	t.Run("returns empty string for non-existent repository", func(t *testing.T) {
		m, _ := createTestManager(t)

		active := m.GetActiveProfile("/non/existent/repo")
		if active != "" {
			t.Errorf("Expected empty string, got '%s'", active)
		}
	})

	t.Run("returns set active profile", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Add profile and set as active
		profile := createTestProfile("test")
		err := m.AddProviderProfile(repoPath, profile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		err = m.SetActiveProfile(repoPath, "test")
		if err != nil {
			t.Fatalf("Failed to set active profile: %v", err)
		}

		active := m.GetActiveProfile(repoPath)
		if active != "test" {
			t.Errorf("Expected 'test', got '%s'", active)
		}
	})
}

// TestSetActiveProfile tests the SetActiveProfile method
func TestSetActiveProfile(t *testing.T) {
	t.Run("sets and retrieves active profile", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Add profile
		profile := createTestProfile("test")
		err := m.AddProviderProfile(repoPath, profile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		// Set as active
		err = m.SetActiveProfile(repoPath, "test")
		if err != nil {
			t.Fatalf("Failed to set active profile: %v", err)
		}

		// Verify
		active := m.GetActiveProfile(repoPath)
		if active != "test" {
			t.Errorf("Expected 'test', got '%s'", active)
		}
	})

	t.Run("validates profile exists before setting", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Try to set non-existent profile as active
		err := m.SetActiveProfile(repoPath, "non-existent")

		if err == nil {
			t.Fatal("Expected error when setting non-existent profile as active, got nil")
		}

		expectedErr := "profile 'non-existent' not found"
		if err.Error() != expectedErr {
			t.Errorf("Expected error '%s', got '%s'", expectedErr, err.Error())
		}
	})

	t.Run("allows setting empty string (clearing active profile)", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Add and set active profile
		profile := createTestProfile("test")
		err := m.AddProviderProfile(repoPath, profile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		err = m.SetActiveProfile(repoPath, "test")
		if err != nil {
			t.Fatalf("Failed to set active profile: %v", err)
		}

		// Clear active profile
		err = m.SetActiveProfile(repoPath, "")
		if err != nil {
			t.Fatalf("Failed to clear active profile: %v", err)
		}

		// Verify
		active := m.GetActiveProfile(repoPath)
		if active != "" {
			t.Errorf("Expected empty string, got '%s'", active)
		}
	})

	t.Run("creates AIProvider if nil", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Create repo config without AIProvider
		m.config.Repositories[repoPath] = &RepoConfig{
			BaseBranch: "main",
		}

		// Set empty active profile (should create AIProvider)
		err := m.SetActiveProfile(repoPath, "")
		if err != nil {
			t.Fatalf("Failed to set active profile: %v", err)
		}

		// Verify AIProvider was created
		if m.config.Repositories[repoPath].AIProvider == nil {
			t.Fatal("Expected AIProvider to be created")
		}
	})

	t.Run("saves config after setting", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Add profile
		profile := createTestProfile("test")
		err := m.AddProviderProfile(repoPath, profile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		// Set as active
		err = m.SetActiveProfile(repoPath, "test")
		if err != nil {
			t.Fatalf("Failed to set active profile: %v", err)
		}

		// Verify file was written
		data, err := os.ReadFile(m.configPath)
		if err != nil {
			t.Fatalf("Failed to read config file: %v", err)
		}

		// Verify JSON contains active profile
		jsonStr := string(data)
		if !contains(jsonStr, `"active_profile":"test"`) && !contains(jsonStr, `"active_profile": "test"`) {
			t.Error("Expected config file to contain active_profile")
		}
	})
}

// TestGetFallbackProfile tests the GetFallbackProfile method
func TestGetFallbackProfile(t *testing.T) {
	t.Run("returns empty string when not set", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		fallback := m.GetFallbackProfile(repoPath)
		if fallback != "" {
			t.Errorf("Expected empty string, got '%s'", fallback)
		}
	})

	t.Run("returns empty string for non-existent repository", func(t *testing.T) {
		m, _ := createTestManager(t)

		fallback := m.GetFallbackProfile("/non/existent/repo")
		if fallback != "" {
			t.Errorf("Expected empty string, got '%s'", fallback)
		}
	})

	t.Run("returns set fallback profile", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Add profile and set as fallback
		profile := createTestProfile("test")
		err := m.AddProviderProfile(repoPath, profile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		err = m.SetFallbackProfile(repoPath, "test")
		if err != nil {
			t.Fatalf("Failed to set fallback profile: %v", err)
		}

		fallback := m.GetFallbackProfile(repoPath)
		if fallback != "test" {
			t.Errorf("Expected 'test', got '%s'", fallback)
		}
	})
}

// TestSetFallbackProfile tests the SetFallbackProfile method
func TestSetFallbackProfile(t *testing.T) {
	t.Run("sets and retrieves fallback profile", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Add profile
		profile := createTestProfile("test")
		err := m.AddProviderProfile(repoPath, profile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		// Set as fallback
		err = m.SetFallbackProfile(repoPath, "test")
		if err != nil {
			t.Fatalf("Failed to set fallback profile: %v", err)
		}

		// Verify
		fallback := m.GetFallbackProfile(repoPath)
		if fallback != "test" {
			t.Errorf("Expected 'test', got '%s'", fallback)
		}
	})

	t.Run("validates profile exists before setting", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Try to set non-existent profile as fallback
		err := m.SetFallbackProfile(repoPath, "non-existent")

		if err == nil {
			t.Fatal("Expected error when setting non-existent profile as fallback, got nil")
		}

		expectedErr := "profile 'non-existent' not found"
		if err.Error() != expectedErr {
			t.Errorf("Expected error '%s', got '%s'", expectedErr, err.Error())
		}
	})

	t.Run("allows setting empty string (clearing fallback profile)", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Add and set fallback profile
		profile := createTestProfile("test")
		err := m.AddProviderProfile(repoPath, profile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		err = m.SetFallbackProfile(repoPath, "test")
		if err != nil {
			t.Fatalf("Failed to set fallback profile: %v", err)
		}

		// Clear fallback profile
		err = m.SetFallbackProfile(repoPath, "")
		if err != nil {
			t.Fatalf("Failed to clear fallback profile: %v", err)
		}

		// Verify
		fallback := m.GetFallbackProfile(repoPath)
		if fallback != "" {
			t.Errorf("Expected empty string, got '%s'", fallback)
		}
	})

	t.Run("creates AIProvider if nil", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Create repo config without AIProvider
		m.config.Repositories[repoPath] = &RepoConfig{
			BaseBranch: "main",
		}

		// Set empty fallback profile (should create AIProvider)
		err := m.SetFallbackProfile(repoPath, "")
		if err != nil {
			t.Fatalf("Failed to set fallback profile: %v", err)
		}

		// Verify AIProvider was created
		if m.config.Repositories[repoPath].AIProvider == nil {
			t.Fatal("Expected AIProvider to be created")
		}
	})

	t.Run("saves config after setting", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Add profile
		profile := createTestProfile("test")
		err := m.AddProviderProfile(repoPath, profile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		// Set as fallback
		err = m.SetFallbackProfile(repoPath, "test")
		if err != nil {
			t.Fatalf("Failed to set fallback profile: %v", err)
		}

		// Verify file was written
		data, err := os.ReadFile(m.configPath)
		if err != nil {
			t.Fatalf("Failed to read config file: %v", err)
		}

		// Verify JSON contains fallback profile
		jsonStr := string(data)
		if !contains(jsonStr, `"fallback_profile":"test"`) && !contains(jsonStr, `"fallback_profile": "test"`) {
			t.Error("Expected config file to contain fallback_profile")
		}
	})
}

// TestGetActiveProviderProfile tests the GetActiveProviderProfile method
func TestGetActiveProviderProfile(t *testing.T) {
	t.Run("returns nil when no active profile", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		profile := m.GetActiveProviderProfile(repoPath)
		if profile != nil {
			t.Error("Expected nil profile, got non-nil")
		}
	})

	t.Run("returns nil when active profile is empty string", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Set active profile to empty string
		err := m.SetActiveProfile(repoPath, "")
		if err != nil {
			t.Fatalf("Failed to set active profile: %v", err)
		}

		profile := m.GetActiveProviderProfile(repoPath)
		if profile != nil {
			t.Error("Expected nil profile, got non-nil")
		}
	})

	t.Run("returns active profile with all fields", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Add profile
		testProfile := &AIProviderProfile{
			Name:   "openai-gpt4",
			Type:   openai.ProviderOpenAI,
			BaseURL: "https://api.openai.com/v1",
			APIKey: "sk-test-key-12345",
			Model:  "gpt-4",
		}
		err := m.AddProviderProfile(repoPath, testProfile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		// Set as active
		err = m.SetActiveProfile(repoPath, "openai-gpt4")
		if err != nil {
			t.Fatalf("Failed to set active profile: %v", err)
		}

		// Get active profile
		profile := m.GetActiveProviderProfile(repoPath)
		if profile == nil {
			t.Fatal("Expected non-nil profile, got nil")
		}

		// Verify all fields
		if profile.Name != "openai-gpt4" {
			t.Errorf("Expected name 'openai-gpt4', got '%s'", profile.Name)
		}
		if profile.Type != openai.ProviderOpenAI {
			t.Errorf("Expected type '%s', got '%s'", openai.ProviderOpenAI, profile.Type)
		}
		if profile.BaseURL != "https://api.openai.com/v1" {
			t.Errorf("Expected base URL 'https://api.openai.com/v1', got '%s'", profile.BaseURL)
		}
		if profile.APIKey != "sk-test-key-12345" {
			t.Errorf("Expected API key 'sk-test-key-12345', got '%s'", profile.APIKey)
		}
		if profile.Model != "gpt-4" {
			t.Errorf("Expected model 'gpt-4', got '%s'", profile.Model)
		}
	})

	t.Run("returns nil when active profile is set but profile doesn't exist", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Manually set active profile without adding it
		if m.config.Repositories == nil {
			m.config.Repositories = make(map[string]*RepoConfig)
		}
		m.config.Repositories[repoPath] = &RepoConfig{
			AIProvider: &AIProviderConfig{
				ActiveProfile: "non-existent",
			},
		}

		profile := m.GetActiveProviderProfile(repoPath)
		if profile != nil {
			t.Error("Expected nil profile, got non-nil")
		}
	})
}

// TestGetFallbackProviderProfile tests the GetFallbackProviderProfile method
func TestGetFallbackProviderProfile(t *testing.T) {
	t.Run("returns nil when no fallback profile", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		profile := m.GetFallbackProviderProfile(repoPath)
		if profile != nil {
			t.Error("Expected nil profile, got non-nil")
		}
	})

	t.Run("returns nil when fallback profile is empty string", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Set fallback profile to empty string
		err := m.SetFallbackProfile(repoPath, "")
		if err != nil {
			t.Fatalf("Failed to set fallback profile: %v", err)
		}

		profile := m.GetFallbackProviderProfile(repoPath)
		if profile != nil {
			t.Error("Expected nil profile, got non-nil")
		}
	})

	t.Run("returns fallback profile with all fields", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Add profile
		testProfile := &AIProviderProfile{
			Name:   "azure-openai",
			Type:   openai.ProviderAzure,
			BaseURL: "https://test.azure.openai.com",
			APIKey: "azure-key-67890",
			Model:  "gpt-35-turbo",
		}
		err := m.AddProviderProfile(repoPath, testProfile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		// Set as fallback
		err = m.SetFallbackProfile(repoPath, "azure-openai")
		if err != nil {
			t.Fatalf("Failed to set fallback profile: %v", err)
		}

		// Get fallback profile
		profile := m.GetFallbackProviderProfile(repoPath)
		if profile == nil {
			t.Fatal("Expected non-nil profile, got nil")
		}

		// Verify all fields
		if profile.Name != "azure-openai" {
			t.Errorf("Expected name 'azure-openai', got '%s'", profile.Name)
		}
		if profile.Type != openai.ProviderAzure {
			t.Errorf("Expected type '%s', got '%s'", openai.ProviderAzure, profile.Type)
		}
		if profile.BaseURL != "https://test.azure.openai.com" {
			t.Errorf("Expected base URL 'https://test.azure.openai.com', got '%s'", profile.BaseURL)
		}
		if profile.APIKey != "azure-key-67890" {
			t.Errorf("Expected API key 'azure-key-67890', got '%s'", profile.APIKey)
		}
		if profile.Model != "gpt-35-turbo" {
			t.Errorf("Expected model 'gpt-35-turbo', got '%s'", profile.Model)
		}
	})

	t.Run("returns nil when fallback profile is set but profile doesn't exist", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Manually set fallback profile without adding it
		if m.config.Repositories == nil {
			m.config.Repositories = make(map[string]*RepoConfig)
		}
		m.config.Repositories[repoPath] = &RepoConfig{
			AIProvider: &AIProviderConfig{
				FallbackProfile: "non-existent",
			},
		}

		profile := m.GetFallbackProviderProfile(repoPath)
		if profile != nil {
			t.Error("Expected nil profile, got non-nil")
		}
	})
}

// TestHasActiveAIProvider tests the HasActiveAIProvider method
func TestHasActiveAIProvider(t *testing.T) {
	t.Run("returns false when no active profile", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		hasProvider := m.HasActiveAIProvider(repoPath)
		if hasProvider {
			t.Error("Expected false, got true")
		}
	})

	t.Run("returns false when profile has empty API key", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Add profile with empty API key
		profile := &AIProviderProfile{
			Name:   "test",
			Type:   openai.ProviderOpenAI,
			BaseURL: "https://api.openai.com/v1",
			APIKey: "", // Empty
			Model:  "gpt-4",
		}
		err := m.AddProviderProfile(repoPath, profile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		err = m.SetActiveProfile(repoPath, "test")
		if err != nil {
			t.Fatalf("Failed to set active profile: %v", err)
		}

		hasProvider := m.HasActiveAIProvider(repoPath)
		if hasProvider {
			t.Error("Expected false when API key is empty, got true")
		}
	})

	t.Run("returns false when profile has empty BaseURL", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Add profile with empty BaseURL
		profile := &AIProviderProfile{
			Name:   "test",
			Type:   openai.ProviderCustom,
			BaseURL: "", // Empty
			APIKey: "sk-test-key",
			Model:  "gpt-4",
		}
		err := m.AddProviderProfile(repoPath, profile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		err = m.SetActiveProfile(repoPath, "test")
		if err != nil {
			t.Fatalf("Failed to set active profile: %v", err)
		}

		hasProvider := m.HasActiveAIProvider(repoPath)
		if hasProvider {
			t.Error("Expected false when BaseURL is empty, got true")
		}
	})

	t.Run("returns false when profile has empty Model", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Add profile with empty Model
		profile := &AIProviderProfile{
			Name:   "test",
			Type:   openai.ProviderOpenAI,
			BaseURL: "https://api.openai.com/v1",
			APIKey: "sk-test-key",
			Model:  "", // Empty
		}
		err := m.AddProviderProfile(repoPath, profile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		err = m.SetActiveProfile(repoPath, "test")
		if err != nil {
			t.Fatalf("Failed to set active profile: %v", err)
		}

		hasProvider := m.HasActiveAIProvider(repoPath)
		if hasProvider {
			t.Error("Expected false when Model is empty, got true")
		}
	})

	t.Run("returns true when profile is complete", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Add complete profile
		profile := createTestProfile("complete")
		err := m.AddProviderProfile(repoPath, profile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		err = m.SetActiveProfile(repoPath, "complete")
		if err != nil {
			t.Fatalf("Failed to set active profile: %v", err)
		}

		hasProvider := m.HasActiveAIProvider(repoPath)
		if !hasProvider {
			t.Error("Expected true when profile is complete, got false")
		}
	})

	t.Run("returns true for OpenAI provider with default URL", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Add OpenAI profile (uses default URL)
		profile := &AIProviderProfile{
			Name:   "openai-default",
			Type:   openai.ProviderOpenAI,
			BaseURL: "https://api.openai.com/v1", // Explicit URL
			APIKey: "sk-test-key",
			Model:  "gpt-4",
		}
		err := m.AddProviderProfile(repoPath, profile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		err = m.SetActiveProfile(repoPath, "openai-default")
		if err != nil {
			t.Fatalf("Failed to set active profile: %v", err)
		}

		hasProvider := m.HasActiveAIProvider(repoPath)
		if !hasProvider {
			t.Error("Expected true for OpenAI provider, got false")
		}
	})
}

// TestMultipleProfiles tests managing multiple profiles
func TestMultipleProfiles(t *testing.T) {
	t.Run("can add and manage multiple profiles", func(t *testing.T) {
		m, _ := createTestManager(t)
		repoPath := "/test/repo"

		// Add multiple profiles
		profiles := []*AIProviderProfile{
			{
				Name:   "openai-gpt4",
				Type:   openai.ProviderOpenAI,
				BaseURL: "https://api.openai.com/v1",
				APIKey: "sk-openai-key",
				Model:  "gpt-4",
			},
			{
				Name:   "azure-gpt35",
				Type:   openai.ProviderAzure,
				BaseURL: "https://test.azure.openai.com",
				APIKey: "azure-key",
				Model:  "gpt-35-turbo",
			},
			{
				Name:   "custom-local",
				Type:   openai.ProviderCustom,
				BaseURL: "http://localhost:8080",
				APIKey: "local-key",
				Model:  "llama-2",
			},
		}

		for _, profile := range profiles {
			err := m.AddProviderProfile(repoPath, profile)
			if err != nil {
				t.Fatalf("Failed to add profile %s: %v", profile.Name, err)
			}
		}

		// Verify all profiles were added
		retrievedProfiles := m.GetProviderProfiles(repoPath)
		if len(retrievedProfiles) != 3 {
			t.Errorf("Expected 3 profiles, got %d", len(retrievedProfiles))
		}

		// Set active and fallback
		err := m.SetActiveProfile(repoPath, "openai-gpt4")
		if err != nil {
			t.Fatalf("Failed to set active profile: %v", err)
		}

		err = m.SetFallbackProfile(repoPath, "azure-gpt35")
		if err != nil {
			t.Fatalf("Failed to set fallback profile: %v", err)
		}

		// Verify active and fallback
		active := m.GetActiveProviderProfile(repoPath)
		if active.Name != "openai-gpt4" {
			t.Errorf("Expected active profile 'openai-gpt4', got '%s'", active.Name)
		}

		fallback := m.GetFallbackProviderProfile(repoPath)
		if fallback.Name != "azure-gpt35" {
			t.Errorf("Expected fallback profile 'azure-gpt35', got '%s'", fallback.Name)
		}

		// Delete one profile
		err = m.DeleteProviderProfile(repoPath, "custom-local")
		if err != nil {
			t.Fatalf("Failed to delete profile: %v", err)
		}

		// Verify deletion
		retrievedProfiles = m.GetProviderProfiles(repoPath)
		if len(retrievedProfiles) != 2 {
			t.Errorf("Expected 2 profiles after deletion, got %d", len(retrievedProfiles))
		}

		if _, exists := retrievedProfiles["custom-local"]; exists {
			t.Error("Expected 'custom-local' profile to be deleted")
		}
	})
}

// TestProfilePersistence tests that profiles persist across manager reloads
func TestProfilePersistence(t *testing.T) {
	t.Run("profiles persist when reloading from disk", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.json")
		repoPath := "/test/repo"

		// Create first manager and add profiles
		m1 := &Manager{
			configPath: configPath,
			config: &Config{
				Repositories: make(map[string]*RepoConfig),
			},
		}

		profile := createTestProfile("persistent")
		err := m1.AddProviderProfile(repoPath, profile)
		if err != nil {
			t.Fatalf("Failed to add profile: %v", err)
		}

		err = m1.SetActiveProfile(repoPath, "persistent")
		if err != nil {
			t.Fatalf("Failed to set active profile: %v", err)
		}

		// Create second manager that loads from the same file
		m2 := &Manager{
			configPath: configPath,
		}
		err = m2.load()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		// Verify profile was persisted
		profiles := m2.GetProviderProfiles(repoPath)
		if len(profiles) != 1 {
			t.Errorf("Expected 1 profile after reload, got %d", len(profiles))
		}

		if _, exists := profiles["persistent"]; !exists {
			t.Error("Expected 'persistent' profile to exist after reload")
		}

		// Verify active profile was persisted
		active := m2.GetActiveProfile(repoPath)
		if active != "persistent" {
			t.Errorf("Expected active profile 'persistent' after reload, got '%s'", active)
		}

		// Verify profile details
		retrieved := profiles["persistent"]
		if retrieved.APIKey != "sk-test-key-123" {
			t.Errorf("Expected API key 'sk-test-key-123' after reload, got '%s'", retrieved.APIKey)
		}
		if retrieved.Model != "gpt-4" {
			t.Errorf("Expected model 'gpt-4' after reload, got '%s'", retrieved.Model)
		}
	})
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
