package session

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const sessionPrefix = "gcool-"

// Session represents a tmux session
type Session struct {
	Name    string
	Branch  string
	Active  bool
	Windows int
}

// Manager handles tmux session operations
type Manager struct{}

// NewManager creates a new session manager
func NewManager() *Manager {
	return &Manager{}
}

// SanitizeName sanitizes a branch name for use as a tmux session name
func (m *Manager) SanitizeName(branch string) string {
	// Replace invalid characters with hyphens
	reg := regexp.MustCompile(`[^a-zA-Z0-9\-_]`)
	sanitized := reg.ReplaceAllString(branch, "-")

	// Remove consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	sanitized = reg.ReplaceAllString(sanitized, "-")

	// Trim hyphens from start/end
	sanitized = strings.Trim(sanitized, "-")

	return sessionPrefix + sanitized
}

// SessionExists checks if a tmux session with the given name exists
func (m *Manager) SessionExists(sessionName string) bool {
	cmd := exec.Command("tmux", "has-session", "-t", sessionName)
	err := cmd.Run()
	return err == nil
}

// CreateOrAttach creates a new session or attaches to existing one
func (m *Manager) CreateOrAttach(path, branch string, autoStartClaude bool) error {
	sessionName := m.SanitizeName(branch)

	if m.SessionExists(sessionName) {
		// Attach to existing session
		return m.Attach(sessionName)
	}

	// Create new session
	return m.Create(sessionName, path, autoStartClaude)
}

// Create creates a new tmux session in detached mode
func (m *Manager) Create(sessionName, path string, autoStartClaude bool) error {
	var cmd *exec.Cmd

	if autoStartClaude {
		// Create detached session and start claude
		cmd = exec.Command("tmux", "new-session", "-d", "-s", sessionName, "-c", path, "claude")
	} else {
		// Create detached session with just a shell
		cmd = exec.Command("tmux", "new-session", "-d", "-s", sessionName, "-c", path)
	}

	// Run and wait for session to be created
	return cmd.Run()
}

// Attach attaches to an existing tmux session
func (m *Manager) Attach(sessionName string) error {
	cmd := exec.Command("tmux", "attach-session", "-t", sessionName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run and wait for the command to complete (user detaches from tmux)
	return cmd.Run()
}

// NewWindowAndAttach creates a new window in existing session and attaches to it
func (m *Manager) NewWindowAndAttach(sessionName, path string) error {
	// Create a new window in the existing session with the specified path
	// and attach to the session
	cmd := exec.Command("tmux", "new-window", "-t", sessionName, "-c", path)
	if err := cmd.Run(); err != nil {
		return err
	}

	// Now attach to the session
	return m.Attach(sessionName)
}

// List returns all gcool tmux sessions
func (m *Manager) List() ([]Session, error) {
	// List all sessions with format: name:windows:attached
	cmd := exec.Command("tmux", "list-sessions", "-F", "#{session_name}:#{session_windows}:#{session_attached}")
	output, err := cmd.Output()
	if err != nil {
		// No sessions exist
		return []Session{}, nil
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var sessions []Session

	for _, line := range lines {
		if !strings.HasPrefix(line, sessionPrefix) {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) < 3 {
			continue
		}

		name := parts[0]
		branch := strings.TrimPrefix(name, sessionPrefix)
		active := parts[2] == "1"

		// Parse window count
		windows := 1
		fmt.Sscanf(parts[1], "%d", &windows)

		sessions = append(sessions, Session{
			Name:    name,
			Branch:  branch,
			Active:  active,
			Windows: windows,
		})
	}

	return sessions, nil
}

// Kill terminates a tmux session
func (m *Manager) Kill(sessionName string) error {
	cmd := exec.Command("tmux", "kill-session", "-t", sessionName)
	return cmd.Run()
}

// IsTmuxAvailable checks if tmux is installed
func (m *Manager) IsTmuxAvailable() bool {
	cmd := exec.Command("tmux", "-V")
	err := cmd.Run()
	return err == nil
}
