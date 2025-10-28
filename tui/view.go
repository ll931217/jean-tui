package tui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the TUI
func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	// Show modal if active
	if m.modal != noModal {
		return m.renderModal()
	}

	// Render main view with panels
	leftPanel := m.renderWorktreeList()
	rightPanel := m.renderDetails()

	// Calculate panel dimensions
	// Allocate more space for help bar (now 2 rows + possible status)
	helpBarHeight := 3 // 2 rows of help + 1 for spacing
	if m.status != "" {
		helpBarHeight = 4 // status + 2 rows of help + spacing
	}

	panelWidth := (m.width - 6) / 2
	panelHeight := m.height - helpBarHeight - 2

	// Style panels
	leftPanelStyled := activePanelStyle.
		Width(panelWidth).
		Height(panelHeight).
		Render(leftPanel)

	rightPanelStyled := panelStyle.
		Width(panelWidth).
		Height(panelHeight).
		Render(rightPanel)

	// Combine panels
	panels := lipgloss.JoinHorizontal(lipgloss.Top, leftPanelStyled, rightPanelStyled)

	// Render help bar
	helpBar := m.renderHelpBar()

	// Combine everything
	return lipgloss.JoinVertical(lipgloss.Left, panels, helpBar)
}

func (m Model) renderWorktreeList() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("📁 Worktrees"))
	b.WriteString("\n\n")

	if len(m.worktrees) == 0 {
		b.WriteString(normalItemStyle.Render("No worktrees found"))
		return b.String()
	}

	for i, wt := range m.worktrees {
		var style lipgloss.Style
		icon := "  "

		if wt.IsCurrent {
			style = currentWorktreeStyle
			icon = "➜ "
		} else if i == m.selectedIndex {
			style = selectedItemStyle
			icon = "› "
		} else {
			style = normalItemStyle
			icon = "  "
		}

		// Show branch name and shortened path
		branch := wt.Branch
		if branch == "" {
			branch = "(no branch)"
		}

		dirName := filepath.Base(wt.Path)

		// For current worktree, show it's the main repo
		var line string
		if wt.IsCurrent {
			line = fmt.Sprintf("%smain (branch: %s)", icon, branch)
		} else {
			line = fmt.Sprintf("%s%s", icon, branch)
		}

		b.WriteString(style.Render(line))
		b.WriteString("\n")
		b.WriteString(normalItemStyle.Copy().Foreground(mutedColor).Render(fmt.Sprintf("   └─ %s", dirName)))
		b.WriteString("\n")
	}

	return b.String()
}

func (m Model) renderDetails() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("ℹ️  Details"))
	b.WriteString("\n\n")

	// Show base branch at the top
	if m.baseBranch != "" {
		b.WriteString(detailKeyStyle.Render("Base Branch: "))
		b.WriteString(detailValueStyle.Render(m.baseBranch))
		b.WriteString(normalItemStyle.Copy().Foreground(mutedColor).Render(" (for new worktrees)"))
		b.WriteString("\n\n")
	}

	wt := m.selectedWorktree()
	if wt == nil {
		b.WriteString(normalItemStyle.Render("No worktree selected"))
		return b.String()
	}

	// Render details in a nice format
	details := []struct {
		key   string
		value string
	}{
		{"Branch", wt.Branch},
		{"Path", wt.Path},
		{"Commit", wt.Commit[:min(7, len(wt.Commit))]},
		{"Status", func() string {
			if wt.IsCurrent {
				return "Current"
			}
			return "Available"
		}()},
	}

	for _, d := range details {
		b.WriteString(detailKeyStyle.Render(d.key + ": "))
		b.WriteString(detailValueStyle.Render(d.value))
		b.WriteString("\n")
	}

	// Add extra info
	b.WriteString("\n")
	if wt.IsCurrent {
		b.WriteString(normalItemStyle.Copy().Foreground(mutedColor).Render("You are currently in this worktree"))
	} else {
		b.WriteString(normalItemStyle.Copy().Foreground(accentColor).Render("Press Enter to switch to this worktree"))
	}

	return b.String()
}

func (m Model) renderHelpBar() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf(" Error: %s ", m.err.Error()))
	}

	var statusLine string
	if m.status != "" {
		statusLine = helpStyle.Render(" " + statusStyle.Render(m.status) + " ")
	}

	// Split keybindings into two rows
	row1 := []string{
		"↑/↓ navigate",
		"n new branch",
		"a existing branch",
		"c change base",
		"C checkout",
		"t terminal",
	}

	row2 := []string{
		"R rename",
		"d delete",
		"r refresh",
		"enter switch",
		"q quit",
	}

	help1 := helpStyle.Render(" " + strings.Join(row1, " • ") + " ")
	help2 := helpStyle.Render(" " + strings.Join(row2, " • ") + " ")

	if statusLine != "" {
		return statusLine + "\n" + help1 + "\n" + help2
	}
	return help1 + "\n" + help2
}

func (m Model) renderModal() string {
	switch m.modal {
	case createModal:
		return m.renderCreateModal()
	case deleteModal:
		return m.renderDeleteModal()
	case branchSelectModal:
		return m.renderBranchSelectModal()
	case checkoutBranchModal:
		return m.renderCheckoutBranchModal()
	case sessionListModal:
		return m.renderSessionListModal()
	case renameModal:
		return m.renderRenameModal()
	case changeBaseBranchModal:
		return m.renderChangeBaseBranchModal()
	}
	return ""
}

func (m Model) renderCreateModal() string {
	var b strings.Builder

	title := "Create New Worktree"
	if !m.createNewBranch {
		title = "Create Worktree from Branch"
	}

	b.WriteString(modalTitleStyle.Render(title))
	b.WriteString("\n\n")

	// Branch name input
	if m.createNewBranch {
		b.WriteString(inputLabelStyle.Render("Branch Name:"))
		b.WriteString("\n")
		b.WriteString(m.nameInput.View())
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("(edit the name or press Enter to use generated name)"))
		b.WriteString("\n\n")
	} else {
		b.WriteString(inputLabelStyle.Render("Branch: "))
		b.WriteString(detailValueStyle.Render(m.nameInput.Value()))
		b.WriteString("\n\n")
	}

	// Show info about auto-generated workspace location
	b.WriteString(helpStyle.Render("Workspace location: .workspaces/<random-name>"))
	b.WriteString("\n\n")

	// Buttons (now only 2 buttons: Create and Cancel)
	createBtn := "Create"
	cancelBtn := "Cancel"

	if m.modalFocused == 1 || m.modalFocused == 0 {
		b.WriteString(selectedButtonStyle.Render(createBtn))
	} else {
		b.WriteString(buttonStyle.Render(createBtn))
	}

	if m.modalFocused == 2 {
		b.WriteString(selectedCancelButtonStyle.Render(cancelBtn))
	} else {
		b.WriteString(cancelButtonStyle.Render(cancelBtn))
	}

	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("Enter to confirm • Esc to cancel"))

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		modalStyle.Render(b.String()),
	)
}

func (m Model) renderDeleteModal() string {
	var b strings.Builder

	wt := m.selectedWorktree()
	if wt == nil {
		return ""
	}

	b.WriteString(modalTitleStyle.Render("Delete Worktree"))
	b.WriteString("\n\n")
	b.WriteString(normalItemStyle.Render(fmt.Sprintf("Are you sure you want to delete worktree:")))
	b.WriteString("\n\n")
	b.WriteString(detailValueStyle.Render(fmt.Sprintf("  Branch: %s", wt.Branch)))
	b.WriteString("\n")
	b.WriteString(detailValueStyle.Render(fmt.Sprintf("  Path: %s", wt.Path)))
	b.WriteString("\n\n")
	b.WriteString(errorStyle.Render("This will remove the worktree directory!"))
	b.WriteString("\n\n")

	// Buttons
	yesBtn := "Yes, Delete"
	noBtn := "Cancel"

	if m.modalFocused == 0 {
		b.WriteString(selectedButtonStyle.Copy().Background(errorColor).Render(yesBtn))
	} else {
		b.WriteString(buttonStyle.Copy().Background(errorColor).Render(yesBtn))
	}

	if m.modalFocused == 1 {
		b.WriteString(selectedButtonStyle.Render(noBtn))
	} else {
		b.WriteString(buttonStyle.Render(noBtn))
	}

	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("Tab/←→ to switch • Enter/Y to confirm • Esc/N to cancel"))

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		modalStyle.Render(b.String()),
	)
}

func (m Model) renderBranchSelectModal() string {
	var b strings.Builder

	b.WriteString(modalTitleStyle.Render("Select Branch"))
	b.WriteString("\n\n")

	if len(m.branches) == 0 {
		b.WriteString(normalItemStyle.Render("No branches found"))
		b.WriteString("\n\n")
	} else {
		// Show scrollable branch list
		maxVisible := 10
		start := m.branchIndex - maxVisible/2
		if start < 0 {
			start = 0
		}
		end := start + maxVisible
		if end > len(m.branches) {
			end = len(m.branches)
			start = end - maxVisible
			if start < 0 {
				start = 0
			}
		}

		for i := start; i < end; i++ {
			branch := m.branches[i]
			if i == m.branchIndex {
				b.WriteString(selectedItemStyle.Render(fmt.Sprintf("› %s", branch)))
			} else {
				b.WriteString(normalItemStyle.Render(fmt.Sprintf("  %s", branch)))
			}
			b.WriteString("\n")
		}

		b.WriteString("\n")
		b.WriteString(helpStyle.Render(fmt.Sprintf("Showing %d-%d of %d branches", start+1, end, len(m.branches))))
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Buttons
	okBtn := "OK"
	cancelBtn := "Cancel"

	if m.modalFocused == 1 {
		b.WriteString(selectedButtonStyle.Render(okBtn))
	} else {
		b.WriteString(buttonStyle.Render(okBtn))
	}

	if m.modalFocused == 2 {
		b.WriteString(selectedCancelButtonStyle.Render(cancelBtn))
	} else {
		b.WriteString(cancelButtonStyle.Render(cancelBtn))
	}

	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("↑↓ navigate • Tab to switch • Enter to select • Esc to cancel"))

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		modalStyle.Render(b.String()),
	)
}

func (m Model) renderCheckoutBranchModal() string {
	var b strings.Builder

	b.WriteString(modalTitleStyle.Render("Checkout Branch (Main Repo)"))
	b.WriteString("\n\n")

	// Search input
	b.WriteString(inputLabelStyle.Render("Search:"))
	b.WriteString("\n")
	b.WriteString(m.searchInput.View())
	b.WriteString("\n\n")

	// Use filtered branches if search is active
	branches := m.branches
	if m.searchInput.Value() != "" {
		branches = m.filteredBranches
	}

	if len(branches) == 0 {
		b.WriteString(normalItemStyle.Render("No branches found"))
		b.WriteString("\n\n")
	} else {
		// Show scrollable branch list
		maxVisible := 10
		start := m.branchIndex - maxVisible/2
		if start < 0 {
			start = 0
		}
		end := start + maxVisible
		if end > len(branches) {
			end = len(branches)
			start = end - maxVisible
			if start < 0 {
				start = 0
			}
		}

		for i := start; i < end; i++ {
			branch := branches[i]
			if i == m.branchIndex {
				b.WriteString(selectedItemStyle.Render(fmt.Sprintf("› %s", branch)))
			} else {
				b.WriteString(normalItemStyle.Render(fmt.Sprintf("  %s", branch)))
			}
			b.WriteString("\n")
		}

		b.WriteString("\n")
		b.WriteString(helpStyle.Render(fmt.Sprintf("Showing %d-%d of %d branches", start+1, end, len(branches))))
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Buttons
	checkoutBtn := "Checkout"
	cancelBtn := "Cancel"

	if m.modalFocused == 2 {
		b.WriteString(selectedButtonStyle.Render(checkoutBtn))
	} else {
		b.WriteString(buttonStyle.Render(checkoutBtn))
	}

	if m.modalFocused == 3 {
		b.WriteString(selectedCancelButtonStyle.Render(cancelBtn))
	} else {
		b.WriteString(cancelButtonStyle.Render(cancelBtn))
	}

	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("Type to search • ↑↓ navigate • Tab to switch • Enter to checkout • Esc to cancel"))

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		modalStyle.Render(b.String()),
	)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (m Model) renderSessionListModal() string {
	var b strings.Builder

	b.WriteString(modalTitleStyle.Render("Active Tmux Sessions"))
	b.WriteString("\n\n")

	if len(m.sessions) == 0 {
		b.WriteString(normalItemStyle.Render("No active gwt sessions found"))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Press Esc to close"))
	} else {
		// Show sessions
		maxVisible := 10
		start := m.sessionIndex - maxVisible/2
		if start < 0 {
			start = 0
		}
		end := start + maxVisible
		if end > len(m.sessions) {
			end = len(m.sessions)
			start = end - maxVisible
			if start < 0 {
				start = 0
			}
		}

		for i := start; i < end; i++ {
			sess := m.sessions[i]
			statusIcon := "○"
			statusText := ""
			if sess.Active {
				statusIcon = "●"
				statusText = " (attached)"
			}

			var style lipgloss.Style
			if i == m.sessionIndex {
				style = selectedItemStyle
			} else {
				style = normalItemStyle
			}

			line := fmt.Sprintf("%s %s%s", statusIcon, sess.Branch, statusText)
			b.WriteString(style.Render(line))
			b.WriteString("\n")
		}

		b.WriteString("\n")
		b.WriteString(helpStyle.Render(fmt.Sprintf("Showing %d-%d of %d sessions", start+1, end, len(m.sessions))))
		b.WriteString("\n\n")

		b.WriteString(helpStyle.Render("↑↓ navigate • Enter attach • x/d kill • Esc close"))
	}

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		modalStyle.Render(b.String()),
	)
}

func (m Model) renderRenameModal() string {
	var b strings.Builder

	title := "Rename Branch"
	b.WriteString(modalTitleStyle.Render(title))
	b.WriteString("\n\n")

	// Branch name input
	b.WriteString("New branch name:\n")
	inputStyle := normalItemStyle
	if m.modalFocused == 0 {
		inputStyle = selectedItemStyle
	}
	b.WriteString(inputStyle.Render(m.nameInput.View()))
	b.WriteString("\n\n")

	// Buttons
	renameStyle := normalItemStyle
	cancelStyle := normalItemStyle

	if m.modalFocused == 1 {
		renameStyle = selectedItemStyle
	} else if m.modalFocused == 2 {
		cancelStyle = selectedItemStyle
	}

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Left,
		renameStyle.Render("[ Rename ]"),
		"  ",
		cancelStyle.Render("[ Cancel ]"),
	)
	b.WriteString(buttons)

	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("Tab: cycle • Enter: confirm • Esc: cancel"))

	// Center the modal
	modalContent := b.String()
	modalBox := modalStyle.Render(modalContent)

	// Calculate centering
	modalWidth := lipgloss.Width(modalBox)
	modalHeight := lipgloss.Height(modalBox)

	verticalMargin := (m.height - modalHeight) / 2
	horizontalMargin := (m.width - modalWidth) / 2

	if verticalMargin < 0 {
		verticalMargin = 0
	}
	if horizontalMargin < 0 {
		horizontalMargin = 0
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modalBox,
	)
}

func (m Model) renderChangeBaseBranchModal() string {
	var b strings.Builder

	title := "Change Base Branch"
	b.WriteString(modalTitleStyle.Render(title))
	b.WriteString("\n\n")

	b.WriteString("All new worktrees will branch from this base branch.\n")
	b.WriteString("Current base: ")
	b.WriteString(selectedItemStyle.Render(m.baseBranch))
	b.WriteString("\n\n")

	// Search input
	b.WriteString(inputLabelStyle.Render("Search:"))
	b.WriteString("\n")
	b.WriteString(m.searchInput.View())
	b.WriteString("\n\n")

	// Use filtered branches if search is active
	branches := m.branches
	if m.searchInput.Value() != "" {
		branches = m.filteredBranches
	}

	if len(branches) == 0 {
		b.WriteString(normalItemStyle.Render("No branches found"))
		b.WriteString("\n\n")
	} else {
		// Show scrollable branch list
		maxVisible := 10
		start := m.branchIndex - maxVisible/2
		if start < 0 {
			start = 0
		}
		end := start + maxVisible
		if end > len(branches) {
			end = len(branches)
			start = end - maxVisible
			if start < 0 {
				start = 0
			}
		}

		for i := start; i < end; i++ {
			branch := branches[i]
			if i == m.branchIndex {
				b.WriteString(selectedItemStyle.Render(fmt.Sprintf("› %s", branch)))
			} else {
				b.WriteString(normalItemStyle.Render(fmt.Sprintf("  %s", branch)))
			}
			b.WriteString("\n")
		}

		b.WriteString("\n")
		b.WriteString(helpStyle.Render(fmt.Sprintf("Showing %d-%d of %d branches", start+1, end, len(branches))))
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Buttons
	setBtn := "Set"
	cancelBtn := "Cancel"

	if m.modalFocused == 2 {
		b.WriteString(selectedButtonStyle.Render(setBtn))
	} else {
		b.WriteString(buttonStyle.Render(setBtn))
	}

	if m.modalFocused == 3 {
		b.WriteString(selectedCancelButtonStyle.Render(cancelBtn))
	} else {
		b.WriteString(cancelButtonStyle.Render(cancelBtn))
	}

	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("Type to search • ↑↓ navigate • Tab to switch • Enter to set • Esc to cancel"))

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modalStyle.Render(b.String()),
	)
}
