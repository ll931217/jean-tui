package beads

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Manager handles beads operations
type Manager struct {
	repoPath string
}

// NewManager creates a new beads manager
func NewManager(repoPath string) *Manager {
	return &Manager{repoPath: repoPath}
}

// IsInitialized checks if beads is initialized in the repository
func (m *Manager) IsInitialized() bool {
	beadsDir := filepath.Join(m.repoPath, ".beads")
	info, err := os.Stat(beadsDir)
	return err == nil && info.IsDir()
}

// Initialize runs 'bd init' to initialize beads in the repository
func (m *Manager) Initialize() error {
	cmd := exec.Command("bd", "init")
	cmd.Dir = m.repoPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("bd init failed: %s", string(output))
	}
	return nil
}

// GetIssuesForBranch retrieves issues related to a specific branch
// Filters issues by matching branch name in ID, title, or description (case-insensitive)
func (m *Manager) GetIssuesForBranch(branch string) ([]Issue, error) {
	if !m.IsInitialized() {
		return []Issue{}, nil
	}

	issuesPath := filepath.Join(m.repoPath, ".beads", "issues.jsonl")
	file, err := os.Open(issuesPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Issue{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var issues []Issue
	scanner := bufio.NewScanner(file)
	branchLower := strings.ToLower(branch)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var issue Issue
		if err := json.Unmarshal([]byte(line), &issue); err != nil {
			continue
		}

		// Filter: match branch in title, description, or ID
		if strings.Contains(strings.ToLower(issue.Title), branchLower) ||
			strings.Contains(strings.ToLower(issue.Description), branchLower) ||
			strings.Contains(strings.ToLower(issue.ID), branchLower) {
			issues = append(issues, issue)
		}
	}

	return issues, scanner.Err()
}

// GetIssueSummary returns a summary of issues for a specific branch
func (m *Manager) GetIssueSummary(branch string) (IssueSummary, error) {
	issues, err := m.GetIssuesForBranch(branch)
	if err != nil {
		return IssueSummary{}, err
	}

	summary := IssueSummary{}
	for _, issue := range issues {
		summary.TotalCount++
		if issue.Status == "open" || issue.Status == "in_progress" {
			summary.OpenCount++
		} else if issue.Status == "closed" {
			summary.ClosedCount++
		}
	}
	return summary, nil
}
