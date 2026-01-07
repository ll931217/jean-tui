package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ConfigEditor handles external editor invocation for config file editing
type ConfigEditor struct {
	editorCommand string
	configPath    string
	tempPath      string
	backupPath    string
}

// ValidationError represents a JSON validation error with location details
type ValidationError struct {
	Message    string
	Line       int
	Column     int
	Offset     int64
	Context    string
	Underlying error
}

func (e *ValidationError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("%s at line %d, column %d", e.Message, e.Line, e.Column)
	}
	return e.Message
}

// Unwrap returns the underlying error for error unwrapping
func (e *ValidationError) Unwrap() error {
	return e.Underlying
}

// NewConfigEditor creates a new ConfigEditor with automatic editor detection
func NewConfigEditor(configPath string) (*ConfigEditor, error) {
	editor, err := DetectEditor()
	if err != nil {
		return nil, fmt.Errorf("failed to detect editor: %w", err)
	}

	return &ConfigEditor{
		editorCommand: editor,
		configPath:    configPath,
	}, nil
}

// DetectEditor finds the appropriate editor using environment variables and fallbacks
func DetectEditor() (string, error) {
	// Check $EDITOR first
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor, nil
	}

	// Check $VISUAL as fallback
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor, nil
	}

	// Fallback to sensible defaults
	editors := []string{"vi", "vim", "nano", "code"}
	for _, editor := range editors {
		if _, err := exec.LookPath(editor); err == nil {
			return editor, nil
		}
	}

	// Last resort - vi should always exist on Unix systems
	return "vi", nil
}

// CreateTempCopy creates a temporary copy of the config file for editing
func (ce *ConfigEditor) CreateTempCopy() error {
	// Read original config
	data, err := os.ReadFile(ce.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	// Create temp file in /tmp with unique name
	timestamp := time.Now().Format("20060102-150405")
	baseName := filepath.Base(ce.configPath)
	tempName := fmt.Sprintf(".%s.%s.%s", baseName, timestamp, "tmp")

	ce.tempPath = filepath.Join(os.TempDir(), tempName)

	// Write temp file with restricted permissions
	if err := os.WriteFile(ce.tempPath, data, 0600); err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	return nil
}

// CleanupTemp removes the temporary file
func (ce *ConfigEditor) CleanupTemp() error {
	if ce.tempPath == "" {
		return nil
	}

	if err := os.Remove(ce.tempPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to cleanup temp file: %w", err)
	}

	ce.tempPath = ""
	return nil
}

// ValidateJSON validates JSON syntax and returns detailed error information
func (ce *ConfigEditor) ValidateJSON() *ValidationError {
	if ce.tempPath == "" {
		return &ValidationError{
			Message: "no temp file to validate",
		}
	}

	data, err := os.ReadFile(ce.tempPath)
	if err != nil {
		return &ValidationError{
			Message:    fmt.Sprintf("failed to read temp file: %v", err),
			Underlying: err,
		}
	}

	var config interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		// Try to extract detailed error information
		var syntaxErr *json.SyntaxError
		if errors.As(err, &syntaxErr) {
			// Calculate line and column from offset
			line, column := offsetToLineColumn(data, syntaxErr.Offset)

			// Extract context snippet
			context := extractErrorContext(data, syntaxErr.Offset)

			return &ValidationError{
				Message:    "Invalid JSON syntax",
				Line:       line,
				Column:     column,
				Offset:     syntaxErr.Offset,
				Context:    context,
				Underlying: err,
			}
		}

		return &ValidationError{
			Message:    fmt.Sprintf("JSON parse error: %v", err),
			Underlying: err,
		}
	}

	return nil
}

// offsetToLineColumn converts a byte offset to line and column numbers
func offsetToLineColumn(data []byte, offset int64) (line, column int) {
	lines := int64(0)
	currentOffset := int64(0)

	for _, b := range data {
		if b == '\n' {
			lines++
			currentOffset++
		} else {
			currentOffset++
		}

		if currentOffset >= offset {
			return int(lines) + 1, int(offset - (currentOffset - int64(column)))
		}
	}

	return int(lines) + 1, int(column)
}

// extractErrorContext extracts a snippet of text around the error location
func extractErrorContext(data []byte, offset int64) string {
	const contextRadius = 20

	start := offset - contextRadius
	if start < 0 {
		start = 0
	}

	end := offset + contextRadius
	if end > int64(len(data)) {
		end = int64(len(data))
	}

	// Extract up to 2 lines of context
	context := string(data[start:end])
	lines := strings.Split(context, "\n")
	if len(lines) > 2 {
		return strings.Join(lines[:2], "\n")
	}

	return context
}

// CreateBackup creates a backup of the original config file
func (ce *ConfigEditor) CreateBackup() error {
	timestamp := time.Now().Format("20060102-150405")
	baseName := filepath.Base(ce.configPath)
	backupName := fmt.Sprintf(".%s.backup.%s", baseName, timestamp)

	ce.backupPath = filepath.Join(filepath.Dir(ce.configPath), backupName)

	// Read original
	data, err := os.ReadFile(ce.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config for backup: %w", err)
	}

	// Write backup
	if err := os.WriteFile(ce.backupPath, data, 0600); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	return nil
}

// RestoreFromBackup restores the config from the backup file
func (ce *ConfigEditor) RestoreFromBackup() error {
	if ce.backupPath == "" {
		return fmt.Errorf("no backup file available")
	}

	data, err := os.ReadFile(ce.backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup: %w", err)
	}

	if err := os.WriteFile(ce.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to restore backup: %w", err)
	}

	// Clean up temp file
	ce.CleanupTemp()

	return nil
}

// ReplaceOriginal replaces the original config with the edited temp file
func (ce *ConfigEditor) ReplaceOriginal() error {
	if ce.tempPath == "" {
		return fmt.Errorf("no temp file to apply")
	}

	// Read temp file content
	data, err := os.ReadFile(ce.tempPath)
	if err != nil {
		return fmt.Errorf("failed to read temp file: %w", err)
	}

	// Write to original location
	if err := os.WriteFile(ce.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	// Clean up temp file
	ce.CleanupTemp()

	// Clean up backup
	ce.cleanupBackup()

	return nil
}

// cleanupBackup removes the backup file
func (ce *ConfigEditor) cleanupBackup() {
	if ce.backupPath != "" {
		os.Remove(ce.backupPath)
		ce.backupPath = ""
	}
}

// GetTempPath returns the path to the temporary file
func (ce *ConfigEditor) GetTempPath() string {
	return ce.tempPath
}

// GetConfigPath returns the path to the config file
func (ce *ConfigEditor) GetConfigPath() string {
	return ce.configPath
}

// ConfigScope represents the scope of configuration to edit
type ConfigScope string

const (
	ScopeGlobal  ConfigScope = "global"
	ScopeRepo    ConfigScope = "repository"
)

// GetConfigPathForScope returns the config file path for the given scope
// Note: Current architecture uses a single config file (~/.config/jean/config.json)
// with repository-specific settings namespaced under the "repositories" key.
func GetConfigPathForScope(scope ConfigScope) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Currently both scopes use the same config file
	// Future enhancement: support per-repo .jean/config.json files
	configPath := filepath.Join(home, ".config", "jean", "config.json")

	return configPath, nil
}

// DetectAvailableScopes returns which config scopes are available in the current context
func DetectAvailableScopes(repoPath string) []ConfigScope {
	scopes := []ConfigScope{ScopeGlobal} // Global is always available

	// Add repo scope if we're in a worktree/repository
	if repoPath != "" {
		scopes = append(scopes, ScopeRepo)
	}

	return scopes
}

// LaunchEdit launches the external editor and blocks until it closes
// Returns the edit result (success, validation error, or launch error)
func (ce *ConfigEditor) LaunchEdit() error {
	if ce.tempPath == "" {
		return fmt.Errorf("no temp file available for editing")
	}

	// Create command to launch editor
	cmd := exec.Command(ce.editorCommand, ce.tempPath)

	// Set up stdin/stdout/stderr to the current process
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run editor and wait for it to complete (blocking)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("editor command failed: %w", err)
	}

	return nil
}
