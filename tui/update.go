package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Update handles all state updates
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case tea.KeyMsg:
		// Modal is open - handle modal input
		if m.modal != noModal {
			return m.handleModalInput(msg)
		}

		// Normal mode - handle main UI input
		return m.handleMainInput(msg)

	case worktreesLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.status = "Failed to load worktrees"
			// Auto-clear error after 4 seconds
			return m, tea.Tick(4*time.Second, func(t time.Time) tea.Msg {
				return clearErrorMsg{}
			})
		} else {
			m.worktrees = msg.worktrees
			m.err = nil
			m.status = ""
		}
		return m, nil

	case branchesLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.status = "Failed to load branches"
			// Auto-clear error after 4 seconds
			return m, tea.Tick(4*time.Second, func(t time.Time) tea.Msg {
				return clearErrorMsg{}
			})
		} else {
			m.branches = msg.branches
			m.err = nil
		}
		return m, nil

	case worktreeCreatedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.status = "Failed to create worktree"
			// Auto-clear error after 4 seconds
			return m, tea.Tick(4*time.Second, func(t time.Time) tea.Msg {
				return clearErrorMsg{}
			})
		} else {
			m.status = "Worktree created successfully - switching..."
			m.modal = noModal

			// Automatically switch to the newly created worktree
			m.switchInfo = SwitchInfo{
				Path:       msg.path,
				Branch:     msg.branch,
				AutoClaude: m.autoClaude,
			}
			return m, tea.Quit
		}
		return m, cmd

	case worktreeDeletedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.status = "Failed to delete worktree"
			// Auto-clear error after 4 seconds
			return m, tea.Tick(4*time.Second, func(t time.Time) tea.Msg {
				return clearErrorMsg{}
			})
		} else {
			m.status = "Worktree deleted successfully"
			m.modal = noModal
			if m.selectedIndex >= len(m.worktrees)-1 {
				m.selectedIndex = len(m.worktrees) - 2
				if m.selectedIndex < 0 {
					m.selectedIndex = 0
				}
			}
			cmd = m.loadWorktrees
		}
		return m, cmd

	case branchRenamedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.status = "Failed to rename branch"
			// Auto-clear error after 4 seconds
			return m, tea.Tick(4*time.Second, func(t time.Time) tea.Msg {
				return clearErrorMsg{}
			})
		} else {
			m.status = "Branch renamed successfully"
			cmd = m.loadWorktrees
		}
		return m, cmd

	case branchCheckedOutMsg:
		if msg.err != nil {
			m.err = msg.err
			m.status = "Failed to checkout branch"
			// Auto-clear error after 4 seconds
			return m, tea.Tick(4*time.Second, func(t time.Time) tea.Msg {
				return clearErrorMsg{}
			})
		} else {
			m.status = "Branch checked out successfully"
			cmd = m.loadWorktrees
		}
		return m, cmd

	case baseBranchLoadedMsg:
		m.baseBranch = msg.branch
		return m, nil

	case clearErrorMsg:
		m.err = nil
		m.status = ""
		return m, nil

	case statusMsg:
		m.status = string(msg)
		return m, nil

	case sessionsLoadedMsg:
		m.sessions = msg.sessions
		return m, nil
	}

	return m, cmd
}

func (m Model) handleMainInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "up", "k":
		if m.selectedIndex > 0 {
			m.selectedIndex--
		}

	case "down", "j":
		if m.selectedIndex < len(m.worktrees)-1 {
			m.selectedIndex++
		}

	case "r":
		m.status = "Refreshing..."
		return m, m.loadWorktrees

	case "n":
		// Instantly create worktree with random branch name from base branch
		randomName, err := m.gitManager.GenerateRandomName()
		if err != nil {
			m.status = "Failed to generate random name"
			return m, nil
		}

		// Generate random path
		path, err := m.gitManager.GetDefaultPath(randomName)
		if err != nil {
			m.status = "Failed to generate workspace path"
			return m, nil
		}

		m.status = "Creating worktree with branch: " + randomName
		return m, m.createWorktree(path, randomName, true)

	case "c":
		// Open change base branch modal
		m.modal = changeBaseBranchModal
		m.modalFocused = 0
		m.branchIndex = 0
		m.searchInput.SetValue("")
		m.searchInput.Focus()
		m.filteredBranches = nil
		// Try to find current base branch in the list
		return m, m.loadBranches

	case "a":
		// Open create from existing branch modal
		m.modal = branchSelectModal
		m.modalFocused = 0
		m.createNewBranch = false
		m.branchIndex = 0
		return m, m.loadBranches

	case "d", "x":
		// Open delete modal
		if wt := m.selectedWorktree(); wt != nil && !wt.IsCurrent {
			m.modal = deleteModal
			m.modalFocused = 0
			return m, nil
		} else if wt != nil && wt.IsCurrent {
			m.status = "Cannot delete current worktree"
		}

	case "enter", " ":
		// Switch to selected worktree
		if wt := m.selectedWorktree(); wt != nil && !wt.IsCurrent {
			m.switchInfo = SwitchInfo{
				Path:       wt.Path,
				Branch:     wt.Branch,
				AutoClaude: m.autoClaude,
			}
			return m, tea.Quit
		}

	case "R":
		// Rename current branch (Shift+R)
		if wt := m.selectedWorktree(); wt != nil {
			m.modal = renameModal
			m.modalFocused = 0
			m.nameInput.SetValue(wt.Branch)
			m.nameInput.Focus()
			m.nameInput.CursorEnd()
			return m, nil
		}

	case "C":
		// Checkout/switch branch in main repository (Shift+C)
		m.modal = checkoutBranchModal
		m.modalFocused = 0
		m.branchIndex = 0
		m.searchInput.SetValue("")
		m.searchInput.Focus()
		m.filteredBranches = nil
		return m, m.loadBranches

	case "t":
		// Open terminal in a separate tmux session (not the Claude session)
		if wt := m.selectedWorktree(); wt != nil {
			m.switchInfo = SwitchInfo{
				Path:         wt.Path,
				Branch:       wt.Branch,
				AutoClaude:   false,        // Never auto-start Claude for terminal
				TerminalOnly: true,         // Signal this is a terminal session
			}
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) handleModalInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.modal {
	case createModal:
		return m.handleCreateModalInput(msg)

	case deleteModal:
		return m.handleDeleteModalInput(msg)

	case branchSelectModal:
		return m.handleBranchSelectModalInput(msg)

	case checkoutBranchModal:
		return m.handleCheckoutBranchModalInput(msg)

	case sessionListModal:
		return m.handleSessionListModalInput(msg)

	case renameModal:
		return m.handleRenameModalInput(msg)

	case changeBaseBranchModal:
		return m.handleChangeBaseBranchModalInput(msg)
	}

	return m, cmd
}

func (m Model) handleCreateModalInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.modal = noModal
		m.nameInput.Blur()
		return m, nil

	case "tab", "shift+tab":
		// Cycle through: nameInput -> create button -> cancel button (if new branch)
		// or just: create button -> cancel button (if existing branch)
		if m.createNewBranch {
			m.modalFocused = (m.modalFocused + 1) % 3
			if m.modalFocused == 0 {
				m.nameInput.Focus()
			} else {
				m.nameInput.Blur()
			}
		} else {
			// For existing branch, just toggle between buttons
			if m.modalFocused == 0 {
				m.modalFocused = 1
			} else if m.modalFocused == 1 {
				m.modalFocused = 2
			} else {
				m.modalFocused = 1
			}
		}
		return m, nil

	case "enter":
		if m.modalFocused <= 1 {
			// Create button or enter in input
			name := m.nameInput.Value()

			if name == "" {
				m.status = "Branch name is required"
				return m, nil
			}

			// Always generate random path
			path, err := m.gitManager.GetDefaultPath(name)
			if err != nil {
				m.status = "Failed to generate workspace path"
				return m, nil
			}

			return m, m.createWorktree(path, name, m.createNewBranch)
		} else if m.modalFocused == 2 {
			// Cancel button
			m.modal = noModal
			m.nameInput.Blur()
			return m, nil
		}
	}

	// Handle text input for branch name
	var cmd tea.Cmd
	if m.modalFocused == 0 && m.createNewBranch {
		m.nameInput, cmd = m.nameInput.Update(msg)
	}

	return m, cmd
}

func (m Model) handleDeleteModalInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "n":
		m.modal = noModal
		return m, nil

	case "tab", "left", "right", "h", "l":
		m.modalFocused = (m.modalFocused + 1) % 2

	case "enter", "y":
		if m.modalFocused == 0 || msg.String() == "y" {
			// Confirm delete
			if wt := m.selectedWorktree(); wt != nil {
				return m, m.deleteWorktree(wt.Path, wt.Branch, false)
			}
		}
		m.modal = noModal
		return m, nil
	}

	return m, nil
}

func (m Model) handleBranchSelectModalInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.modal = noModal
		return m, nil

	case "up", "k":
		if m.branchIndex > 0 {
			m.branchIndex--
		}

	case "down", "j":
		if m.branchIndex < len(m.branches)-1 {
			m.branchIndex++
		}

	case "tab":
		// Toggle between branch list and buttons
		if m.modalFocused == 0 {
			m.modalFocused = 1 // Move to OK button
		} else if m.modalFocused == 1 {
			m.modalFocused = 2 // Move to Cancel button
		} else {
			m.modalFocused = 0 // Back to list
		}

	case "enter":
		if m.modalFocused == 0 || m.modalFocused == 1 {
			// Select branch and create worktree directly with random name
			branch := m.selectedBranch()
			if branch == "" {
				return m, nil
			}

			// Generate random path
			path, err := m.gitManager.GetDefaultPath(branch)
			if err != nil {
				m.status = "Failed to generate workspace path"
				return m, nil
			}

			m.modal = noModal
			return m, m.createWorktree(path, branch, false)
		} else if m.modalFocused == 2 {
			// Cancel
			m.modal = noModal
			return m, nil
		}
	}

	return m, nil
}

func (m Model) handleCheckoutBranchModalInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "esc":
		m.modal = noModal
		m.searchInput.Blur()
		return m, nil

	case "up", "k":
		if m.modalFocused == 0 {
			// In search input, move focus down to list
			m.modalFocused = 1
			m.searchInput.Blur()
		} else if m.modalFocused == 1 && m.branchIndex > 0 {
			// In list, move selection up
			m.branchIndex--
		}

	case "down", "j":
		if m.modalFocused == 0 {
			// In search input, move focus down to list
			m.modalFocused = 1
			m.searchInput.Blur()
		} else if m.modalFocused == 1 {
			// In list, move selection down
			branches := m.filteredBranches
			if len(branches) == 0 {
				branches = m.branches
			}
			if m.branchIndex < len(branches)-1 {
				m.branchIndex++
			}
		}

	case "tab":
		// Cycle: search -> list -> OK button -> Cancel button -> search
		m.modalFocused = (m.modalFocused + 1) % 4
		if m.modalFocused == 0 {
			m.searchInput.Focus()
		} else {
			m.searchInput.Blur()
		}

	case "enter":
		if m.modalFocused == 0 || m.modalFocused == 1 || m.modalFocused == 2 {
			// Checkout the selected branch in main repository
			branch := m.selectedBranch()
			if branch == "" {
				return m, nil
			}

			m.modal = noModal
			m.searchInput.Blur()
			m.status = "Checking out branch: " + branch
			return m, m.checkoutBranch(branch)
		} else if m.modalFocused == 3 {
			// Cancel
			m.modal = noModal
			m.searchInput.Blur()
			return m, nil
		}
	}

	// Handle search input typing
	if m.modalFocused == 0 {
		m.searchInput, cmd = m.searchInput.Update(msg)
		// Filter branches based on search
		m.filteredBranches = m.filterBranches(m.searchInput.Value())
		// Reset branch index when search changes
		m.branchIndex = 0
	}

	return m, cmd
}

func (m Model) handleSessionListModalInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.modal = noModal
		return m, nil

	case "up", "k":
		if m.sessionIndex > 0 {
			m.sessionIndex--
		}

	case "down", "j":
		if m.sessionIndex < len(m.sessions)-1 {
			m.sessionIndex++
		}

	case "enter":
		// Attach to selected session
		if m.sessionIndex >= 0 && m.sessionIndex < len(m.sessions) {
			sess := m.sessions[m.sessionIndex]
			// Attach via tmux
			if err := m.sessionManager.Attach(sess.Name); err != nil {
				m.status = "Failed to attach to session"
			}
			return m, tea.Quit
		}

	case "x", "d":
		// Kill selected session
		if m.sessionIndex >= 0 && m.sessionIndex < len(m.sessions) {
			sess := m.sessions[m.sessionIndex]
			if err := m.sessionManager.Kill(sess.Name); err != nil {
				m.status = "Failed to kill session"
			} else {
				m.status = "Session killed"
				// Reload sessions
				return m, m.loadSessions
			}
		}
	}

	return m, nil
}

func (m Model) handleRenameModalInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.modal = noModal
		m.nameInput.Blur()
		return m, nil

	case "tab", "shift+tab":
		// Toggle between input and buttons
		m.modalFocused = (m.modalFocused + 1) % 3
		if m.modalFocused == 0 {
			m.nameInput.Focus()
		} else {
			m.nameInput.Blur()
		}
		return m, nil

	case "enter":
		if m.modalFocused <= 1 {
			// Rename button or enter in input
			newName := m.nameInput.Value()
			if newName == "" {
				m.status = "Branch name cannot be empty"
				return m, nil
			}

			if wt := m.selectedWorktree(); wt != nil {
				if newName == wt.Branch {
					m.status = "Branch name unchanged"
					m.modal = noModal
					m.nameInput.Blur()
					return m, nil
				}

				m.status = "Renaming branch..."
				m.modal = noModal
				m.nameInput.Blur()
				return m, m.renameBranch(wt.Branch, newName)
			}
		} else if m.modalFocused == 2 {
			// Cancel button
			m.modal = noModal
			m.nameInput.Blur()
			return m, nil
		}
	}

	// Handle text input
	var cmd tea.Cmd
	if m.modalFocused == 0 {
		m.nameInput, cmd = m.nameInput.Update(msg)
	}

	return m, cmd
}

func (m Model) handleChangeBaseBranchModalInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "esc":
		m.modal = noModal
		m.searchInput.Blur()
		return m, nil

	case "up", "k":
		if m.modalFocused == 0 {
			// In search input, move focus down to list
			m.modalFocused = 1
			m.searchInput.Blur()
		} else if m.modalFocused == 1 && m.branchIndex > 0 {
			// In list, move selection up
			m.branchIndex--
		}

	case "down", "j":
		if m.modalFocused == 0 {
			// In search input, move focus down to list
			m.modalFocused = 1
			m.searchInput.Blur()
		} else if m.modalFocused == 1 {
			// In list, move selection down
			branches := m.filteredBranches
			if len(branches) == 0 {
				branches = m.branches
			}
			if m.branchIndex < len(branches)-1 {
				m.branchIndex++
			}
		}

	case "tab":
		// Cycle: search -> list -> Set button -> Cancel button -> search
		m.modalFocused = (m.modalFocused + 1) % 4
		if m.modalFocused == 0 {
			m.searchInput.Focus()
		} else {
			m.searchInput.Blur()
		}

	case "enter":
		if m.modalFocused == 0 || m.modalFocused == 1 || m.modalFocused == 2 {
			// Set the selected branch as base branch
			newBaseBranch := m.selectedBranch()
			if newBaseBranch == "" {
				return m, nil
			}

			m.baseBranch = newBaseBranch

			// Save to config
			if m.configManager != nil {
				if err := m.configManager.SetBaseBranch(m.repoPath, newBaseBranch); err != nil {
					m.status = "Base branch set to: " + newBaseBranch + " (warning: failed to save)"
				} else {
					m.status = "Base branch set to: " + newBaseBranch + " (saved)"
				}
			} else {
				m.status = "Base branch set to: " + newBaseBranch
			}

			m.modal = noModal
			m.searchInput.Blur()
			return m, nil
		} else if m.modalFocused == 3 {
			// Cancel button
			m.modal = noModal
			m.searchInput.Blur()
			return m, nil
		}
	}

	// Handle search input typing
	if m.modalFocused == 0 {
		m.searchInput, cmd = m.searchInput.Update(msg)
		// Filter branches based on search
		m.filteredBranches = m.filterBranches(m.searchInput.Value())
		// Reset branch index when search changes
		m.branchIndex = 0
	}

	return m, cmd
}
