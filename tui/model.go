package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/heyandras/gwt/config"
	"github.com/heyandras/gwt/git"
	"github.com/heyandras/gwt/session"
)

// SwitchInfo contains information about the worktree to switch to
type SwitchInfo struct {
	Path           string
	Branch         string
	AutoClaude     bool
	TerminalOnly   bool // If true, open terminal session instead of Claude session
}

type modalType int

const (
	noModal modalType = iota
	createModal
	deleteModal
	branchSelectModal
	checkoutBranchModal
	sessionListModal
	renameModal
	changeBaseBranchModal
)

// Model represents the TUI state
type Model struct {
	gitManager     *git.Manager
	sessionManager *session.Manager
	configManager  *config.Manager
	worktrees      []git.Worktree
	branches       []string
	sessions       []session.Session
	repoPath       string // Path to the repository

	// UI state
	selectedIndex   int
	sessionIndex    int
	width           int
	height          int
	ready           bool
	err             error
	status          string
	switchInfo      SwitchInfo // Info about worktree to switch to
	autoClaude      bool       // Whether to auto-start Claude
	baseBranch      string     // Base branch for new worktrees

	// Modal state
	modal           modalType
	modalFocused    int // Which input/button is focused in modal
	nameInput       textinput.Model
	pathInput       textinput.Model
	searchInput     textinput.Model
	branchIndex     int
	filteredBranches []string // Filtered list of branches for search
	createNewBranch bool
}

// NewModel creates a new TUI model
func NewModel(repoPath string, autoClaude bool) Model {
	nameInput := textinput.New()
	nameInput.Placeholder = "branch-name"
	nameInput.Focus()
	nameInput.CharLimit = 156
	nameInput.Width = 50

	pathInput := textinput.New()
	pathInput.Placeholder = "/path/to/worktree"
	pathInput.CharLimit = 256
	pathInput.Width = 50

	searchInput := textinput.New()
	searchInput.Placeholder = "Search branches..."
	searchInput.CharLimit = 100
	searchInput.Width = 50

	// Initialize config manager (ignore errors, will use defaults)
	configManager, _ := config.NewManager()

	// Create git manager and get absolute repo root path
	gitManager := git.NewManager(repoPath)
	absoluteRepoPath := repoPath
	if root, err := gitManager.GetRepoRoot(); err == nil {
		absoluteRepoPath = root
	}

	return Model{
		gitManager:     gitManager,
		sessionManager: session.NewManager(),
		configManager:  configManager,
		nameInput:      nameInput,
		pathInput:      pathInput,
		searchInput:    searchInput,
		autoClaude:     autoClaude,
		repoPath:       absoluteRepoPath,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.loadWorktrees,
		m.loadBaseBranch,
		tea.EnterAltScreen,
	)
}

// Messages
type (
	worktreesLoadedMsg struct {
		worktrees []git.Worktree
		err       error
	}

	branchesLoadedMsg struct {
		branches []string
		err      error
	}

	worktreeCreatedMsg struct {
		err    error
		path   string
		branch string
	}

	worktreeDeletedMsg struct {
		err error
	}

	branchRenamedMsg struct {
		err error
	}

	branchCheckedOutMsg struct {
		err error
	}

	baseBranchLoadedMsg struct {
		branch string
	}

	clearErrorMsg struct{}

	statusMsg string
)

// Commands
func (m Model) loadWorktrees() tea.Msg {
	worktrees, err := m.gitManager.List()
	return worktreesLoadedMsg{worktrees: worktrees, err: err}
}

func (m Model) loadBranches() tea.Msg {
	branches, err := m.gitManager.ListBranches()
	return branchesLoadedMsg{branches: branches, err: err}
}

func (m Model) createWorktree(path, branch string, newBranch bool) tea.Cmd {
	return func() tea.Msg {
		// Ensure .workspaces directory exists
		if err := m.gitManager.EnsureWorkspacesDir(); err != nil {
			return worktreeCreatedMsg{err: err, path: path, branch: branch}
		}

		// Use base branch when creating new branch
		baseBranch := ""
		if newBranch {
			baseBranch = m.baseBranch
		}

		err := m.gitManager.Create(path, branch, newBranch, baseBranch)
		return worktreeCreatedMsg{err: err, path: path, branch: branch}
	}
}

func (m Model) deleteWorktree(path, branch string, force bool) tea.Cmd {
	return func() tea.Msg {
		// First remove the worktree
		err := m.gitManager.Remove(path, force)
		if err != nil {
			return worktreeDeletedMsg{err: err}
		}

		// Then kill the associated tmux session if it exists
		sessionName := m.sessionManager.SanitizeName(branch)
		_ = m.sessionManager.Kill(sessionName) // Ignore error if session doesn't exist

		return worktreeDeletedMsg{err: nil}
	}
}

func (m Model) renameBranch(oldName, newName string) tea.Cmd {
	return func() tea.Msg {
		err := m.gitManager.RenameBranch(oldName, newName)
		return branchRenamedMsg{err: err}
	}
}

func (m Model) checkoutBranch(branch string) tea.Cmd {
	return func() tea.Msg {
		err := m.gitManager.CheckoutBranch(branch)
		return branchCheckedOutMsg{err: err}
	}
}

func (m Model) loadBaseBranch() tea.Msg {
	// First, try to load from config
	if m.configManager != nil {
		if savedBranch := m.configManager.GetBaseBranch(m.repoPath); savedBranch != "" {
			return baseBranchLoadedMsg{branch: savedBranch}
		}
	}

	// If not in config, try current branch
	branch, err := m.gitManager.GetCurrentBranch()
	if err != nil || branch == "" {
		// Try to get default branch (main or master)
		defaultBranch, err := m.gitManager.GetDefaultBranch()
		if err != nil {
			// Last resort: empty (user must set manually)
			return baseBranchLoadedMsg{branch: ""}
		}
		return baseBranchLoadedMsg{branch: defaultBranch}
	}
	return baseBranchLoadedMsg{branch: branch}
}

// Helper methods
func (m Model) selectedWorktree() *git.Worktree {
	if m.selectedIndex < 0 || m.selectedIndex >= len(m.worktrees) {
		return nil
	}
	return &m.worktrees[m.selectedIndex]
}

func (m Model) selectedBranch() string {
	// Use filtered branches if search is active
	branches := m.branches
	if len(m.filteredBranches) > 0 || m.searchInput.Value() != "" {
		branches = m.filteredBranches
	}

	if m.branchIndex < 0 || m.branchIndex >= len(branches) {
		return ""
	}
	return branches[m.branchIndex]
}

func (m Model) filterBranches(query string) []string {
	if query == "" {
		return m.branches
	}

	var filtered []string
	queryLower := strings.ToLower(query)
	for _, branch := range m.branches {
		if strings.Contains(strings.ToLower(branch), queryLower) {
			filtered = append(filtered, branch)
		}
	}
	return filtered
}

// GetSwitchInfo returns the switch information (for shell integration)
func (m Model) GetSwitchInfo() SwitchInfo {
	return m.switchInfo
}

// loadSessions loads tmux sessions
func (m Model) loadSessions() tea.Msg {
	sessions, err := m.sessionManager.List()
	if err != nil {
		return statusMsg("Failed to load sessions")
	}
	return sessionsLoadedMsg{sessions: sessions}
}

type sessionsLoadedMsg struct {
	sessions []session.Session
}
