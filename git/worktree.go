package git

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Worktree represents a Git worktree
type Worktree struct {
	Path      string
	Branch    string
	Commit    string
	IsCurrent bool
}

// Manager handles Git worktree operations
type Manager struct {
	repoPath string
}

// NewManager creates a new worktree manager
func NewManager(repoPath string) *Manager {
	return &Manager{repoPath: repoPath}
}

// List returns all worktrees in the repository
func (m *Manager) List() ([]Worktree, error) {
	cmd := exec.Command("git", "-C", m.repoPath, "worktree", "list", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list worktrees: %w", err)
	}

	return m.parseWorktrees(string(output))
}

// parseWorktrees parses the output of 'git worktree list --porcelain'
func (m *Manager) parseWorktrees(output string) ([]Worktree, error) {
	var worktrees []Worktree
	var current Worktree

	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		if line == "" {
			if current.Path != "" {
				worktrees = append(worktrees, current)
				current = Worktree{}
			}
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		switch key {
		case "worktree":
			current.Path = value
		case "HEAD":
			current.Commit = value
		case "branch":
			// Remove "refs/heads/" prefix
			current.Branch = strings.TrimPrefix(value, "refs/heads/")
		case "detached":
			current.Branch = fmt.Sprintf("(detached at %s)", current.Commit[:7])
		}
	}

	// Add the last worktree if exists
	if current.Path != "" {
		worktrees = append(worktrees, current)
	}

	// Mark current worktree
	currentPath, err := m.getCurrentPath()
	if err == nil {
		for i := range worktrees {
			if worktrees[i].Path == currentPath {
				worktrees[i].IsCurrent = true
			}
		}
	}

	return worktrees, nil
}

// GetCurrentBranch returns the name of the current branch
func (m *Manager) GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "-C", m.repoPath, "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetDefaultBranch tries to determine the default branch (main, master, etc.)
func (m *Manager) GetDefaultBranch() (string, error) {
	// First try to get the default branch from remote
	cmd := exec.Command("git", "-C", m.repoPath, "symbolic-ref", "refs/remotes/origin/HEAD")
	output, err := cmd.Output()
	if err == nil {
		// Extract branch name from refs/remotes/origin/HEAD -> refs/remotes/origin/main
		branch := strings.TrimSpace(string(output))
		branch = strings.TrimPrefix(branch, "refs/remotes/origin/")
		if branch != "" {
			return branch, nil
		}
	}

	// Fallback: check if main or master exists locally
	for _, branch := range []string{"main", "master"} {
		cmd := exec.Command("git", "-C", m.repoPath, "rev-parse", "--verify", branch)
		if err := cmd.Run(); err == nil {
			return branch, nil
		}
	}

	// Last resort: get the first branch
	cmd = exec.Command("git", "-C", m.repoPath, "branch", "--format=%(refname:short)")
	output, err = cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get any branch: %w", err)
	}

	branches := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(branches) > 0 && branches[0] != "" {
		return branches[0], nil
	}

	return "", fmt.Errorf("no branches found in repository")
}

// getCurrentPath returns the current worktree path
func (m *Manager) getCurrentPath() (string, error) {
	cmd := exec.Command("git", "-C", m.repoPath, "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// Create creates a new worktree
func (m *Manager) Create(path, branch string, newBranch bool, baseBranch string) error {
	// Validate base branch exists if specified
	if newBranch && baseBranch != "" {
		cmd := exec.Command("git", "-C", m.repoPath, "rev-parse", "--verify", baseBranch)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("base branch '%s' does not exist. Use 'c' to change the base branch", baseBranch)
		}
	}

	args := []string{"-C", m.repoPath, "worktree", "add"}

	if newBranch {
		args = append(args, "-b", branch)
	}

	args = append(args, path)

	if !newBranch {
		args = append(args, branch)
	} else if baseBranch != "" {
		// When creating new branch, specify the base branch to start from
		args = append(args, baseBranch)
	}

	cmd := exec.Command("git", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create worktree: %s", string(output))
	}

	return nil
}

// Remove removes a worktree
func (m *Manager) Remove(path string, force bool) error {
	args := []string{"-C", m.repoPath, "worktree", "remove"}

	if force {
		args = append(args, "--force")
	}

	args = append(args, path)

	cmd := exec.Command("git", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to remove worktree: %s", string(output))
	}

	return nil
}

// RenameBranch renames the current branch
func (m *Manager) RenameBranch(oldName, newName string) error {
	cmd := exec.Command("git", "-C", m.repoPath, "branch", "-m", oldName, newName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to rename branch: %s", string(output))
	}
	return nil
}

// CheckoutBranch checks out a branch in the main repository
func (m *Manager) CheckoutBranch(branch string) error {
	cmd := exec.Command("git", "-C", m.repoPath, "checkout", branch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to checkout branch: %s", string(output))
	}
	return nil
}

// ListBranches returns all branches in the repository
func (m *Manager) ListBranches() ([]string, error) {
	cmd := exec.Command("git", "-C", m.repoPath, "branch", "-a", "--format=%(refname:short)")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	branches := strings.Split(strings.TrimSpace(string(output)), "\n")

	// Filter out remote branches and current branch marker
	var filtered []string
	for _, b := range branches {
		b = strings.TrimSpace(b)
		if b != "" && !strings.HasPrefix(b, "origin/HEAD") {
			// Remove "origin/" prefix for remote branches
			b = strings.TrimPrefix(b, "origin/")
			filtered = append(filtered, b)
		}
	}

	return filtered, nil
}

// GetRepoRoot returns the root path of the repository
func (m *Manager) GetRepoRoot() (string, error) {
	cmd := exec.Command("git", "-C", m.repoPath, "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not a git repository or git is not installed")
	}
	return strings.TrimSpace(string(output)), nil
}

// GetDefaultPath returns a default path for a new worktree in .workspaces directory
func (m *Manager) GetDefaultPath(branch string) (string, error) {
	root, err := m.GetRepoRoot()
	if err != nil {
		return "", err
	}

	// Use .workspaces directory inside repo root
	workspacesDir := filepath.Join(root, ".workspaces")

	// Use branch name as the directory name
	return filepath.Join(workspacesDir, branch), nil
}

// GetWorkspacesDir returns the .workspaces directory path
func (m *Manager) GetWorkspacesDir() (string, error) {
	root, err := m.GetRepoRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, ".workspaces"), nil
}

// EnsureWorkspacesDir creates the .workspaces directory if it doesn't exist
func (m *Manager) EnsureWorkspacesDir() error {
	dir, err := m.GetWorkspacesDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create .workspaces directory: %w", err)
	}

	return nil
}

// Random name generator - generates fun, memorable names
var (
	adjectives = []string{
		"happy", "swift", "bright", "calm", "brave", "clever", "gentle", "quick",
		"proud", "wise", "noble", "bold", "cool", "keen", "warm", "kind",
		"fine", "grand", "pure", "true", "vital", "witty", "zesty", "apt",
	}

	nouns = []string{
		"panda", "fox", "wolf", "bear", "lion", "tiger", "eagle", "falcon",
		"hawk", "raven", "owl", "dragon", "phoenix", "wizard", "knight", "warrior",
		"ninja", "samurai", "pirate", "viking", "spartan", "ranger", "hunter", "scout",
	}
)

// GenerateRandomName creates a random name like "happy-panda-42"
func (m *Manager) GenerateRandomName() (string, error) {
	return generateRandomName(), nil
}

// generateRandomName creates a random name like "happy-panda-42"
func generateRandomName() string {
	adjective := adjectives[randomInt(len(adjectives))]
	noun := nouns[randomInt(len(nouns))]
	number := randomInt(100)

	return fmt.Sprintf("%s-%s-%d", adjective, noun, number)
}

// randomInt returns a random integer between 0 and max (exclusive)
func randomInt(max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		// Fallback to a simple number if crypto/rand fails
		return 0
	}
	return int(n.Int64())
}
