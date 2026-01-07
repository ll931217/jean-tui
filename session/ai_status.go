package session

import (
	"bytes"
	"os/exec"
	"strings"
)

// AIStatusDetector detects the status of AI sessions (Claude)
type AIStatusDetector struct{}

// NewAIStatusDetector creates a new AI status detector
func NewAIStatusDetector() *AIStatusDetector {
	return &AIStatusDetector{}
}

// DetectAISessionState checks if an AI session (Claude) is waiting for input
// It does this by examining the tmux session contents for Claude-specific prompts
func (d *AIStatusDetector) DetectAISessionState(sessionName string) bool {
	if sessionName == "" {
		return false
	}

	// Check if the session exists
	if !d.sessionExists(sessionName) {
		return false
	}

	// Capture the session output to check for Claude prompts
	output, err := d.captureSessionOutput(sessionName)
	if err != nil {
		return false
	}

	// Look for Claude-specific indicators that it's waiting for input
	return d.isClaudeWaiting(output)
}

// sessionExists checks if a tmux session exists
func (d *AIStatusDetector) sessionExists(sessionName string) bool {
	cmd := exec.Command("tmux", "list-sessions", "-F", "#{session_name}")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	sessions := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, session := range sessions {
		if session == sessionName {
			return true
		}
	}

	return false
}

// captureSessionOutput captures the current output of a tmux session
func (d *AIStatusDetector) captureSessionOutput(sessionName string) (string, error) {
	// Use tmux capture-pane to get the current screen content
	cmd := exec.Command("tmux", "capture-pane", "-t", sessionName, "-p")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return stdout.String(), nil
}

// isClaudeWaiting checks if the output contains indicators that Claude is waiting for input
func (d *AIStatusDetector) isClaudeWaiting(output string) bool {
	// Claude CLI indicators that it's waiting for input
	// These patterns indicate Claude is ready and waiting for user input
	indicators := []string{
		"Anthropic",
		"claude",
		"Claude",
		"User:",
		"USER:",
		"Human:",
		"human:",
		">>>", // Some Claude versions use this prompt
		"│",   // Vertical bar often used in Claude CLI interface
		"└",   // Corner characters used in Claude CLI
		"├",   // Tree characters used in Claude CLI
	}

	// Check for multiple indicators to reduce false positives
	// Claude CLI typically has a distinct visual style
	indicatorCount := 0
	outputLower := strings.ToLower(output)

	for _, indicator := range indicators {
		if strings.Contains(output, indicator) || strings.Contains(outputLower, strings.ToLower(indicator)) {
			indicatorCount++
		}
	}

	// Need at least 2 indicators to consider it a Claude session
	// This helps avoid false positives from regular shell sessions
	if indicatorCount >= 2 {
		// Additional check: look for common shell prompts that would indicate NOT in Claude
		shellPrompts := []string{"$", "#", ">", "%", "→"}
		for _, prompt := range shellPrompts {
			// If we see shell prompts at the end, it's probably not Claude waiting
			lines := strings.Split(output, "\n")
			if len(lines) > 0 {
				lastLine := strings.TrimSpace(lines[len(lines)-1])
				if strings.HasSuffix(lastLine, prompt) {
					// Check if there's a clear shell prompt pattern
					if strings.Count(lastLine, prompt) > 0 && !strings.Contains(lastLine, "│") {
						return false
					}
				}
			}
		}

		return true
	}

	// Check for cursor position or input indicators
	if strings.Contains(output, "█") || strings.Contains(output, "▊") || strings.Contains(output, "▋") {
		// Cursor characters often indicate active input mode
		// Combined with Anthropic/Claude mentions
		if strings.Contains(output, "Anthropic") || strings.Contains(output, "Claude") {
			return true
		}
	}

	return false
}

// IsClaudeSessionName checks if a session name appears to be a Claude session
func IsClaudeSessionName(sessionName string) bool {
	return strings.HasPrefix(sessionName, "jean-") && !strings.HasSuffix(sessionName, "-terminal")
}

// GetClaudeSessionName returns the Claude session name for a given branch
func GetClaudeSessionName(branch string) string {
	// This should match the sanitization logic in tmux.go
	// For now, return empty - the actual implementation should be consistent
	// with the session naming in tmux.go
	return ""
}
