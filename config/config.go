package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"github.com/coollabsio/jean-tui/openai"
)

// AIPrompts represents customizable AI prompts for various generation tasks
type AIPrompts struct {
	CommitMessage string `json:"commit_message,omitempty"` // Custom prompt for commit message generation
	BranchName    string `json:"branch_name,omitempty"`    // Custom prompt for branch name generation
	PRContent     string `json:"pr_content,omitempty"`     // Custom prompt for PR title and description generation
}

// AIProviderProfile represents a single AI provider profile configuration
type AIProviderProfile struct {
	Name   string              `json:"name"`   // Display name for the profile
	Type   openai.ProviderType `json:"type"`   // Provider type (openai, azure, custom)
	BaseURL string              `json:"base_url"` // API base URL (overrides default if set)
	APIKey string              `json:"api_key"` // API key for authentication
	Model  string              `json:"model"`  // Model to use
}

// AIProviderConfig holds all AI provider profiles and settings
type AIProviderConfig struct {
	Profiles        map[string]*AIProviderProfile `json:"profiles"`         // name -> profile
	ActiveProfile   string                        `json:"active_profile"`   // Name of active profile
	FallbackProfile string                        `json:"fallback_profile"` // Name of fallback profile
}

// Config represents the global jean configuration
type Config struct {
	Repositories        map[string]*RepoConfig `json:"repositories"`
	LastUpdateCheckTime string                 `json:"lastUpdateCheckTime"` // RFC3339 format
	DefaultTheme        string                 `json:"default_theme,omitempty"` // Global default theme, "" = matrix
	AICommitEnabled     bool                   `json:"ai_commit_enabled,omitempty"` // Enable AI commit message generation
	AIBranchNameEnabled bool                   `json:"ai_branch_name_enabled,omitempty"` // Enable AI branch name generation
	DebugLoggingEnabled bool                   `json:"debug_logging_enabled"` // Enable debug logging to temp files
	AIPrompts           *AIPrompts             `json:"ai_prompts,omitempty"` // Customizable AI prompts
	WrapperChecksums    map[string]string      `json:"wrapper_checksums,omitempty"` // Shell -> SHA256 checksum of installed wrapper
	Onboarded           bool                   `json:"onboarded"` // Whether the user has completed the onboarding flow
}

// PRInfo represents information about a pull request
type PRInfo struct {
	URL       string `json:"url"`
	Status    string `json:"status"` // "open", "merged", or "closed"
	CreatedAt string `json:"created_at,omitempty"` // RFC3339 format
	Branch    string `json:"branch"`
	PRNumber  int    `json:"pr_number,omitempty"` // GitHub PR number (e.g., 42 from github.com/owner/repo/pull/42)
	Title     string `json:"title,omitempty"`     // PR title for display
	Author    string `json:"author,omitempty"`    // Author login for display
}

// RepoConfig represents configuration for a specific repository
type RepoConfig struct {
	BaseBranch         string                  `json:"base_branch"`
	LastSelectedBranch string                  `json:"last_selected_branch,omitempty"`
	Editor             string                  `json:"editor,omitempty"`
	AutoFetchInterval  int                     `json:"auto_fetch_interval,omitempty"` // in seconds, 0 = use default (10s)
	Theme              string                  `json:"theme,omitempty"`               // Per-repo theme override, "" = use global default
	PRDefaultState     string                  `json:"pr_default_state,omitempty"`    // "draft" or "ready", "" = use default (ready)
	PRs                map[string][]PRInfo      `json:"prs,omitempty"`                // branch -> list of PRs
	InitializedClaudes map[string]bool        `json:"initialized_claudes,omitempty"` // branch -> whether Claude has been started
	AIProvider         *AIProviderConfig       `json:"ai_provider,omitempty"`       // AI provider profiles and settings
	Hooks              *HooksConfig            `json:"hooks,omitempty"`              // Hooks configuration
}

// Hook represents a single hook configuration (duplicated from hooks package for JSON serialization)
type Hook struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	Enabled  bool   `json:"enabled"`
	RunAsync bool   `json:"run_async"`
}

// HooksConfig holds all hook configurations organized by hook type
type HooksConfig struct {
	PreCreate  []Hook `json:"pre_create,omitempty"`  // Hooks before worktree creation
	PostCreate []Hook `json:"post_create,omitempty"` // Hooks after successful worktree creation
	PreDelete  []Hook `json:"pre_delete,omitempty"`  // Hooks before worktree deletion
	PostDelete []Hook `json:"post_delete,omitempty"` // Hooks after successful worktree deletion
	OnSwitch   []Hook `json:"on_switch,omitempty"`   // Hooks when switching worktrees
}

// Manager handles configuration loading and saving
type Manager struct {
	configPath string
	config     *Config
}

// NewManager creates a new configuration manager
func NewManager() (*Manager, error) {
	// Get user home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	// Create config directory: ~/.config/jean
	configDir := filepath.Join(home, ".config", "jean")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "config.json")

	m := &Manager{
		configPath: configPath,
	}

	// Load existing config or create new one
	if err := m.load(); err != nil {
		// If file doesn't exist, create empty config
		m.config = &Config{
			Repositories: make(map[string]*RepoConfig),
		}
	}

	return m, nil
}

// load reads the configuration from disk
func (m *Manager) load() error {
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return err
	}

	m.config = &Config{}
	return json.Unmarshal(data, m.config)
}

// save writes the configuration to disk
func (m *Manager) save() error {
	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(m.configPath, data, 0644)
}

// Reload reloads the configuration from disk
func (m *Manager) Reload() error {
	return m.load()
}

// GetBaseBranch returns the base branch for a repository
func (m *Manager) GetBaseBranch(repoPath string) string {
	if repo, ok := m.config.Repositories[repoPath]; ok {
		return repo.BaseBranch
	}
	return ""
}

// SetBaseBranch sets the base branch for a repository
func (m *Manager) SetBaseBranch(repoPath, branch string) error {
	if m.config.Repositories == nil {
		m.config.Repositories = make(map[string]*RepoConfig)
	}

	if _, ok := m.config.Repositories[repoPath]; !ok {
		m.config.Repositories[repoPath] = &RepoConfig{}
	}

	m.config.Repositories[repoPath].BaseBranch = branch
	return m.save()
}

// GetRepoConfig returns the configuration for a specific repository
func (m *Manager) GetRepoConfig(repoPath string) *RepoConfig {
	if repo, ok := m.config.Repositories[repoPath]; ok {
		return repo
	}
	return &RepoConfig{}
}

// GetLastSelectedBranch returns the last selected branch for a repository
func (m *Manager) GetLastSelectedBranch(repoPath string) string {
	if repo, ok := m.config.Repositories[repoPath]; ok {
		return repo.LastSelectedBranch
	}
	return ""
}

// SetLastSelectedBranch sets the last selected branch for a repository
func (m *Manager) SetLastSelectedBranch(repoPath, branch string) error {
	if m.config.Repositories == nil {
		m.config.Repositories = make(map[string]*RepoConfig)
	}

	if _, ok := m.config.Repositories[repoPath]; !ok {
		m.config.Repositories[repoPath] = &RepoConfig{}
	}

	m.config.Repositories[repoPath].LastSelectedBranch = branch
	return m.save()
}

// GetEditor returns the preferred editor for a repository
func (m *Manager) GetEditor(repoPath string) string {
	if repo, ok := m.config.Repositories[repoPath]; ok {
		if repo.Editor != "" {
			return repo.Editor
		}
	}
	return "code" // Default to VS Code
}

// SetEditor sets the preferred editor for a repository
func (m *Manager) SetEditor(repoPath, editor string) error {
	if m.config.Repositories == nil {
		m.config.Repositories = make(map[string]*RepoConfig)
	}

	if _, ok := m.config.Repositories[repoPath]; !ok {
		m.config.Repositories[repoPath] = &RepoConfig{}
	}

	m.config.Repositories[repoPath].Editor = editor
	return m.save()
}

// GetAutoFetchInterval returns the auto-fetch interval for a repository
// Returns the configured interval in seconds, or 10 if not set
func (m *Manager) GetAutoFetchInterval(repoPath string) int {
	if repo, ok := m.config.Repositories[repoPath]; ok {
		if repo.AutoFetchInterval > 0 {
			return repo.AutoFetchInterval
		}
	}
	return 10 // Default to 10 seconds
}

// SetAutoFetchInterval sets the auto-fetch interval for a repository
func (m *Manager) SetAutoFetchInterval(repoPath string, interval int) error {
	if m.config.Repositories == nil {
		m.config.Repositories = make(map[string]*RepoConfig)
	}

	if _, ok := m.config.Repositories[repoPath]; !ok {
		m.config.Repositories[repoPath] = &RepoConfig{}
	}

	m.config.Repositories[repoPath].AutoFetchInterval = interval
	return m.save()
}

// GetLastUpdateCheckTime returns the last update check time
func (m *Manager) GetLastUpdateCheckTime() string {
	return m.config.LastUpdateCheckTime
}

// SetLastUpdateCheckTime sets the last update check time
func (m *Manager) SetLastUpdateCheckTime(timestamp string) error {
	m.config.LastUpdateCheckTime = timestamp
	return m.save()
}

// GetTheme returns the theme for a repository
// Returns per-repo theme if set, otherwise returns global default theme
// Returns "coolify" if no theme is configured
func (m *Manager) GetTheme(repoPath string) string {
	// Check if repo has a per-repo theme override
	if repo, ok := m.config.Repositories[repoPath]; ok {
		if repo.Theme != "" {
			return repo.Theme
		}
	}

	// Fall back to global default theme
	if m.config.DefaultTheme != "" {
		return m.config.DefaultTheme
	}

	// Default to coolify theme
	return "coolify"
}

// SetTheme sets the theme for a specific repository
// If theme is empty string, it will use the global default
func (m *Manager) SetTheme(repoPath, theme string) error {
	if m.config.Repositories == nil {
		m.config.Repositories = make(map[string]*RepoConfig)
	}

	if _, ok := m.config.Repositories[repoPath]; !ok {
		m.config.Repositories[repoPath] = &RepoConfig{}
	}

	m.config.Repositories[repoPath].Theme = theme
	return m.save()
}

// SetGlobalTheme sets the global default theme for all repositories
func (m *Manager) SetGlobalTheme(theme string) error {
	m.config.DefaultTheme = theme
	return m.save()
}

// GetGlobalTheme returns the global default theme
// Returns "coolify" if not set
func (m *Manager) GetGlobalTheme() string {
	if m.config.DefaultTheme != "" {
		return m.config.DefaultTheme
	}
	return "coolify"
}

// GetProviderProfiles returns all provider profiles for a repository
func (m *Manager) GetProviderProfiles(repoPath string) map[string]*AIProviderProfile {
	if repo, ok := m.config.Repositories[repoPath]; ok {
		if repo.AIProvider != nil && repo.AIProvider.Profiles != nil {
			return repo.AIProvider.Profiles
		}
	}
	return make(map[string]*AIProviderProfile)
}

// AddProviderProfile adds a new provider profile for a repository
func (m *Manager) AddProviderProfile(repoPath string, profile *AIProviderProfile) error {
	if m.config.Repositories == nil {
		m.config.Repositories = make(map[string]*RepoConfig)
	}

	if _, ok := m.config.Repositories[repoPath]; !ok {
		m.config.Repositories[repoPath] = &RepoConfig{}
	}

	repo := m.config.Repositories[repoPath]
	if repo.AIProvider == nil {
		repo.AIProvider = &AIProviderConfig{
			Profiles: make(map[string]*AIProviderProfile),
		}
	}

	if repo.AIProvider.Profiles == nil {
		repo.AIProvider.Profiles = make(map[string]*AIProviderProfile)
	}

	repo.AIProvider.Profiles[profile.Name] = profile
	return m.save()
}

// UpdateProviderProfile updates an existing provider profile
func (m *Manager) UpdateProviderProfile(repoPath string, profile *AIProviderProfile) error {
	if repo, ok := m.config.Repositories[repoPath]; ok {
		if repo.AIProvider != nil && repo.AIProvider.Profiles != nil {
			if _, exists := repo.AIProvider.Profiles[profile.Name]; exists {
				repo.AIProvider.Profiles[profile.Name] = profile
				return m.save()
			}
		}
	}
	return fmt.Errorf("profile '%s' not found", profile.Name)
}

// DeleteProviderProfile deletes a provider profile
func (m *Manager) DeleteProviderProfile(repoPath, name string) error {
	if repo, ok := m.config.Repositories[repoPath]; ok {
		if repo.AIProvider != nil && repo.AIProvider.Profiles != nil {
			if repo.AIProvider.ActiveProfile == name {
				repo.AIProvider.ActiveProfile = ""
			}
			if repo.AIProvider.FallbackProfile == name {
				repo.AIProvider.FallbackProfile = ""
			}
			delete(repo.AIProvider.Profiles, name)
			return m.save()
		}
	}
	return fmt.Errorf("profile '%s' not found", name)
}

// GetActiveProfile returns the active profile name for a repository
func (m *Manager) GetActiveProfile(repoPath string) string {
	if repo, ok := m.config.Repositories[repoPath]; ok {
		if repo.AIProvider != nil {
			return repo.AIProvider.ActiveProfile
		}
	}
	return ""
}

// SetActiveProfile sets the active profile for a repository
func (m *Manager) SetActiveProfile(repoPath, profileName string) error {
	if m.config.Repositories == nil {
		m.config.Repositories = make(map[string]*RepoConfig)
	}

	if _, ok := m.config.Repositories[repoPath]; !ok {
		m.config.Repositories[repoPath] = &RepoConfig{}
	}

	repo := m.config.Repositories[repoPath]
	if repo.AIProvider == nil {
		repo.AIProvider = &AIProviderConfig{
			Profiles: make(map[string]*AIProviderProfile),
		}
	}

	// Verify profile exists
	if repo.AIProvider.Profiles != nil {
		if _, exists := repo.AIProvider.Profiles[profileName]; !exists && profileName != "" {
			return fmt.Errorf("profile '%s' not found", profileName)
		}
	}

	repo.AIProvider.ActiveProfile = profileName
	return m.save()
}

// GetFallbackProfile returns the fallback profile name for a repository
func (m *Manager) GetFallbackProfile(repoPath string) string {
	if repo, ok := m.config.Repositories[repoPath]; ok {
		if repo.AIProvider != nil {
			return repo.AIProvider.FallbackProfile
		}
	}
	return ""
}

// SetFallbackProfile sets the fallback profile for a repository
func (m *Manager) SetFallbackProfile(repoPath, profileName string) error {
	if m.config.Repositories == nil {
		m.config.Repositories = make(map[string]*RepoConfig)
	}

	if _, ok := m.config.Repositories[repoPath]; !ok {
		m.config.Repositories[repoPath] = &RepoConfig{}
	}

	repo := m.config.Repositories[repoPath]
	if repo.AIProvider == nil {
		repo.AIProvider = &AIProviderConfig{
			Profiles: make(map[string]*AIProviderProfile),
		}
	}

	// Verify profile exists
	if repo.AIProvider.Profiles != nil {
		if _, exists := repo.AIProvider.Profiles[profileName]; !exists && profileName != "" {
			return fmt.Errorf("profile '%s' not found", profileName)
		}
	}

	repo.AIProvider.FallbackProfile = profileName
	return m.save()
}

// GetActiveProviderProfile returns the active profile configuration for a repository
func (m *Manager) GetActiveProviderProfile(repoPath string) *AIProviderProfile {
	activeProfileName := m.GetActiveProfile(repoPath)
	if activeProfileName == "" {
		return nil
	}

	profiles := m.GetProviderProfiles(repoPath)
	return profiles[activeProfileName]
}

// GetFallbackProviderProfile returns the fallback profile configuration for a repository
func (m *Manager) GetFallbackProviderProfile(repoPath string) *AIProviderProfile {
	fallbackProfileName := m.GetFallbackProfile(repoPath)
	if fallbackProfileName == "" {
		return nil
	}

	profiles := m.GetProviderProfiles(repoPath)
	return profiles[fallbackProfileName]
}

// HasActiveAIProvider returns whether an AI provider is configured for a repository
func (m *Manager) HasActiveAIProvider(repoPath string) bool {
	profile := m.GetActiveProviderProfile(repoPath)
	return profile != nil && profile.APIKey != "" && profile.BaseURL != "" && profile.Model != ""
}

// GetAICommitEnabled returns whether AI commit message generation is enabled
func (m *Manager) GetAICommitEnabled() bool {
	return m.config.AICommitEnabled
}

// SetAICommitEnabled sets whether AI commit message generation is enabled
func (m *Manager) SetAICommitEnabled(enabled bool) error {
	m.config.AICommitEnabled = enabled
	return m.save()
}

// GetAIBranchNameEnabled returns whether AI branch name generation is enabled
func (m *Manager) GetAIBranchNameEnabled() bool {
	return m.config.AIBranchNameEnabled
}

// SetAIBranchNameEnabled sets whether AI branch name generation is enabled
func (m *Manager) SetAIBranchNameEnabled(enabled bool) error {
	m.config.AIBranchNameEnabled = enabled
	return m.save()
}

// GetDebugLoggingEnabled returns whether debug logging is enabled
func (m *Manager) GetDebugLoggingEnabled() bool {
	return m.config.DebugLoggingEnabled
}

// SetDebugLoggingEnabled sets whether debug logging is enabled
func (m *Manager) SetDebugLoggingEnabled(enabled bool) error {
	m.config.DebugLoggingEnabled = enabled
	return m.save()
}

// GetPRs returns all pull requests for a given branch
func (m *Manager) GetPRs(repoPath, branch string) []PRInfo {
	if repo, ok := m.config.Repositories[repoPath]; ok {
		if repo.PRs != nil {
			if prs, ok := repo.PRs[branch]; ok {
				return prs
			}
		}
	}
	return []PRInfo{}
}

// GetLatestPR returns the most recent pull request for a given branch
func (m *Manager) GetLatestPR(repoPath, branch string) *PRInfo {
	prs := m.GetPRs(repoPath, branch)
	if len(prs) == 0 {
		return nil
	}
	return &prs[len(prs)-1]
}

// AddPR adds a pull request for a given branch
func (m *Manager) AddPR(repoPath, branch, url string, prNumber int, title string, author string) error {
	if m.config.Repositories == nil {
		m.config.Repositories = make(map[string]*RepoConfig)
	}

	if _, ok := m.config.Repositories[repoPath]; !ok {
		m.config.Repositories[repoPath] = &RepoConfig{}
	}

	repo := m.config.Repositories[repoPath]
	if repo.PRs == nil {
		repo.PRs = make(map[string][]PRInfo)
	}

	// Check if PR already exists
	for _, pr := range repo.PRs[branch] {
		if pr.URL == url {
			return nil // Already exists
		}
	}

	// Add new PR with "open" status by default
	prInfo := PRInfo{
		URL:      url,
		Status:   "open",
		Branch:   branch,
		PRNumber: prNumber,
		Title:    title,
		Author:   author,
	}
	repo.PRs[branch] = append(repo.PRs[branch], prInfo)
	return m.save()
}

// UpdatePRStatus updates the status of a pull request
func (m *Manager) UpdatePRStatus(repoPath, branch, url, status string) error {
	if repo, ok := m.config.Repositories[repoPath]; ok {
		if repo.PRs != nil {
			if prs, ok := repo.PRs[branch]; ok {
				for i, pr := range prs {
					if pr.URL == url {
						prs[i].Status = status
						return m.save()
					}
				}
			}
		}
	}
	return nil
}

// RemovePR removes a pull request
func (m *Manager) RemovePR(repoPath, branch, url string) error {
	if repo, ok := m.config.Repositories[repoPath]; ok {
		if repo.PRs != nil {
			if prs, ok := repo.PRs[branch]; ok {
				newPRs := []PRInfo{}
				for _, pr := range prs {
					if pr.URL != url {
						newPRs = append(newPRs, pr)
					}
				}
				repo.PRs[branch] = newPRs
				return m.save()
			}
		}
	}
	return nil
}

// HasPRs checks if there are any pull requests for a given branch
func (m *Manager) HasPRs(repoPath, branch string) bool {
	prs := m.GetPRs(repoPath, branch)
	return len(prs) > 0
}

// IsClaudeInitialized checks if a Claude session has been initialized for a branch
func (m *Manager) IsClaudeInitialized(repoPath, branch string) bool {
	if repo, ok := m.config.Repositories[repoPath]; ok {
		if repo.InitializedClaudes != nil {
			return repo.InitializedClaudes[branch]
		}
	}
	return false
}

// SetClaudeInitialized marks a branch as having an initialized Claude session
func (m *Manager) SetClaudeInitialized(repoPath, branch string) error {
	if m.config.Repositories == nil {
		m.config.Repositories = make(map[string]*RepoConfig)
	}

	if _, ok := m.config.Repositories[repoPath]; !ok {
		m.config.Repositories[repoPath] = &RepoConfig{}
	}

	repo := m.config.Repositories[repoPath]
	if repo.InitializedClaudes == nil {
		repo.InitializedClaudes = make(map[string]bool)
	}

	repo.InitializedClaudes[branch] = true
	return m.save()
}

// CleanupBranch removes all branch-specific data from config when a worktree is deleted
// This includes:
// - All pull requests for the branch
// - Claude initialization flag
// - Last selected branch reference (if it matches the deleted branch)
func (m *Manager) CleanupBranch(repoPath, branch string) error {
	repo, ok := m.config.Repositories[repoPath]
	if !ok {
		return nil // Nothing to clean up
	}

	// Remove all PRs for this branch
	if repo.PRs != nil {
		delete(repo.PRs, branch)
	}

	// Remove Claude initialization flag for this branch
	if repo.InitializedClaudes != nil {
		delete(repo.InitializedClaudes, branch)
	}

	// Clear last selected branch if it matches the deleted branch
	if repo.LastSelectedBranch == branch {
		repo.LastSelectedBranch = ""
	}

	return m.save()
}

// GetCommitPrompt returns the custom commit message prompt
// Returns the custom prompt if set, otherwise returns the default prompt
func (m *Manager) GetCommitPrompt() string {
	if m.config.AIPrompts != nil && m.config.AIPrompts.CommitMessage != "" {
		return m.config.AIPrompts.CommitMessage
	}
	return openai.GetDefaultCommitPrompt()
}

// SetCommitPrompt sets the custom commit message prompt
func (m *Manager) SetCommitPrompt(prompt string) error {
	if m.config.AIPrompts == nil {
		m.config.AIPrompts = &AIPrompts{}
	}
	m.config.AIPrompts.CommitMessage = prompt
	return m.save()
}

// GetBranchNamePrompt returns the custom branch name prompt
// Returns the custom prompt if set, otherwise returns the default prompt
func (m *Manager) GetBranchNamePrompt() string {
	if m.config.AIPrompts != nil && m.config.AIPrompts.BranchName != "" {
		return m.config.AIPrompts.BranchName
	}
	return openai.GetDefaultBranchNamePrompt()
}

// SetBranchNamePrompt sets the custom branch name prompt
func (m *Manager) SetBranchNamePrompt(prompt string) error {
	if m.config.AIPrompts == nil {
		m.config.AIPrompts = &AIPrompts{}
	}
	m.config.AIPrompts.BranchName = prompt
	return m.save()
}

// GetPRPrompt returns the custom PR content prompt
// Returns the custom prompt if set, otherwise returns the default prompt
func (m *Manager) GetPRPrompt() string {
	if m.config.AIPrompts != nil && m.config.AIPrompts.PRContent != "" {
		return m.config.AIPrompts.PRContent
	}
	return openai.GetDefaultPRPrompt()
}

// SetPRPrompt sets the custom PR content prompt
func (m *Manager) SetPRPrompt(prompt string) error {
	if m.config.AIPrompts == nil {
		m.config.AIPrompts = &AIPrompts{}
	}
	m.config.AIPrompts.PRContent = prompt
	return m.save()
}

// ResetAIPromptsToDefaults resets all AI prompts to their default values
func (m *Manager) ResetAIPromptsToDefaults() error {
	m.config.AIPrompts = &AIPrompts{} // Empty AIPrompts means use defaults
	return m.save()
}

// GetWrapperChecksum returns the stored checksum for a shell wrapper
// Returns empty string if no checksum is stored
func (m *Manager) GetWrapperChecksum(shell string) string {
	if m.config.WrapperChecksums == nil {
		return ""
	}
	return m.config.WrapperChecksums[shell]
}

// SetWrapperChecksum stores the checksum for a shell wrapper
func (m *Manager) SetWrapperChecksum(shell, checksum string) error {
	if m.config.WrapperChecksums == nil {
		m.config.WrapperChecksums = make(map[string]string)
	}
	m.config.WrapperChecksums[shell] = checksum
	return m.save()
}

// IsOnboarded returns whether the user has completed the onboarding flow
func (m *Manager) IsOnboarded() bool {
	return m.config.Onboarded
}

// SetOnboarded marks the onboarding flow as completed
func (m *Manager) SetOnboarded() error {
	m.config.Onboarded = true
	return m.save()
}

// GetPRDefaultState returns the default PR state for a repository
// Returns "draft" or "ready", defaults to "ready" if not set
func (m *Manager) GetPRDefaultState(repoPath string) string {
	if repo, ok := m.config.Repositories[repoPath]; ok {
		if repo.PRDefaultState == "draft" || repo.PRDefaultState == "ready" {
			return repo.PRDefaultState
		}
	}
	return "ready" // Default to "ready for review"
}

// SetPRDefaultState sets the default PR state for a repository
func (m *Manager) SetPRDefaultState(repoPath, state string) error {
	if m.config.Repositories == nil {
		m.config.Repositories = make(map[string]*RepoConfig)
	}

	if _, ok := m.config.Repositories[repoPath]; !ok {
		m.config.Repositories[repoPath] = &RepoConfig{}
	}

	m.config.Repositories[repoPath].PRDefaultState = state
	return m.save()
}

// GetHooks returns the hooks configuration for a repository
func (m *Manager) GetHooks(repoPath string) *HooksConfig {
	if repo, ok := m.config.Repositories[repoPath]; ok {
		return repo.Hooks
	}
	return nil
}

// SetHooks sets the entire hooks configuration for a repository
func (m *Manager) SetHooks(repoPath string, hooks *HooksConfig) error {
	if m.config.Repositories == nil {
		m.config.Repositories = make(map[string]*RepoConfig)
	}

	if _, ok := m.config.Repositories[repoPath]; !ok {
		m.config.Repositories[repoPath] = &RepoConfig{}
	}

	m.config.Repositories[repoPath].Hooks = hooks
	return m.save()
}

// AddHook adds a hook to a specific hook type for a repository
func (m *Manager) AddHook(repoPath, hookType string, hook Hook) error {
	if m.config.Repositories == nil {
		m.config.Repositories = make(map[string]*RepoConfig)
	}

	if _, ok := m.config.Repositories[repoPath]; !ok {
		m.config.Repositories[repoPath] = &RepoConfig{}
	}

	repo := m.config.Repositories[repoPath]
	if repo.Hooks == nil {
		repo.Hooks = &HooksConfig{}
	}

	switch hookType {
	case "pre_create":
		repo.Hooks.PreCreate = append(repo.Hooks.PreCreate, hook)
	case "post_create":
		repo.Hooks.PostCreate = append(repo.Hooks.PostCreate, hook)
	case "pre_delete":
		repo.Hooks.PreDelete = append(repo.Hooks.PreDelete, hook)
	case "post_delete":
		repo.Hooks.PostDelete = append(repo.Hooks.PostDelete, hook)
	case "on_switch":
		repo.Hooks.OnSwitch = append(repo.Hooks.OnSwitch, hook)
	default:
		return fmt.Errorf("invalid hook type: %s", hookType)
	}

	return m.save()
}

// RemoveHook removes a hook from a specific hook type by index
func (m *Manager) RemoveHook(repoPath, hookType string, index int) error {
	repo, ok := m.config.Repositories[repoPath]
	if !ok || repo.Hooks == nil {
		return fmt.Errorf("no hooks configured for repository")
	}

	switch hookType {
	case "pre_create":
		if index < 0 || index >= len(repo.Hooks.PreCreate) {
			return fmt.Errorf("invalid hook index")
		}
		repo.Hooks.PreCreate = append(repo.Hooks.PreCreate[:index], repo.Hooks.PreCreate[index+1:]...)
	case "post_create":
		if index < 0 || index >= len(repo.Hooks.PostCreate) {
			return fmt.Errorf("invalid hook index")
		}
		repo.Hooks.PostCreate = append(repo.Hooks.PostCreate[:index], repo.Hooks.PostCreate[index+1:]...)
	case "pre_delete":
		if index < 0 || index >= len(repo.Hooks.PreDelete) {
			return fmt.Errorf("invalid hook index")
		}
		repo.Hooks.PreDelete = append(repo.Hooks.PreDelete[:index], repo.Hooks.PreDelete[index+1:]...)
	case "post_delete":
		if index < 0 || index >= len(repo.Hooks.PostDelete) {
			return fmt.Errorf("invalid hook index")
		}
		repo.Hooks.PostDelete = append(repo.Hooks.PostDelete[:index], repo.Hooks.PostDelete[index+1:]...)
	case "on_switch":
		if index < 0 || index >= len(repo.Hooks.OnSwitch) {
			return fmt.Errorf("invalid hook index")
		}
		repo.Hooks.OnSwitch = append(repo.Hooks.OnSwitch[:index], repo.Hooks.OnSwitch[index+1:]...)
	default:
		return fmt.Errorf("invalid hook type: %s", hookType)
	}

	return m.save()
}

// UpdateHook updates a hook at a specific index for a hook type
func (m *Manager) UpdateHook(repoPath, hookType string, index int, hook Hook) error {
	repo, ok := m.config.Repositories[repoPath]
	if !ok || repo.Hooks == nil {
		return fmt.Errorf("no hooks configured for repository")
	}

	switch hookType {
	case "pre_create":
		if index < 0 || index >= len(repo.Hooks.PreCreate) {
			return fmt.Errorf("invalid hook index")
		}
		repo.Hooks.PreCreate[index] = hook
	case "post_create":
		if index < 0 || index >= len(repo.Hooks.PostCreate) {
			return fmt.Errorf("invalid hook index")
		}
		repo.Hooks.PostCreate[index] = hook
	case "pre_delete":
		if index < 0 || index >= len(repo.Hooks.PreDelete) {
			return fmt.Errorf("invalid hook index")
		}
		repo.Hooks.PreDelete[index] = hook
	case "post_delete":
		if index < 0 || index >= len(repo.Hooks.PostDelete) {
			return fmt.Errorf("invalid hook index")
		}
		repo.Hooks.PostDelete[index] = hook
	case "on_switch":
		if index < 0 || index >= len(repo.Hooks.OnSwitch) {
			return fmt.Errorf("invalid hook index")
		}
		repo.Hooks.OnSwitch[index] = hook
	default:
		return fmt.Errorf("invalid hook type: %s", hookType)
	}

	return m.save()
}
