package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/coollabsio/gcool/tui"
)

const version = "0.1.0"

func main() {
	// Parse flags
	pathFlag := flag.String("path", ".", "Path to git repository (default: current directory)")
	noClaudeFlag := flag.Bool("no-claude", false, "Don't auto-start Claude CLI in tmux session")
	versionFlag := flag.Bool("version", false, "Print version and exit")
	helpFlag := flag.Bool("help", false, "Show help")

	flag.Parse()

	// Handle flags
	if *versionFlag {
		fmt.Printf("gcool version %s\n", version)
		os.Exit(0)
	}

	if *helpFlag {
		printHelp()
		os.Exit(0)
	}

	// Get repo path and auto-claude setting
	repoPath := *pathFlag
	autoClaude := !*noClaudeFlag

	// Create and run TUI
	model := tui.NewModel(repoPath, autoClaude)
	p := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}

	// Check if we need to switch directories
	if m, ok := finalModel.(tui.Model); ok {
		switchInfo := m.GetSwitchInfo()
		if switchInfo.Path != "" {
			// Format: path|branch|auto-claude|terminal-only
			autoCl := "false"
			if switchInfo.AutoClaude {
				autoCl = "true"
			}
			termOnly := "false"
			if switchInfo.TerminalOnly {
				termOnly = "true"
			}
			switchData := fmt.Sprintf("%s|%s|%s|%s", switchInfo.Path, switchInfo.Branch, autoCl, termOnly)

			// Check if we should write to a file (for shell wrapper integration)
			if switchFile := os.Getenv("GCOOL_SWITCH_FILE"); switchFile != "" {
				// Write to file for shell wrapper
				if err := os.WriteFile(switchFile, []byte(switchData), 0600); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: could not write switch file: %v\n", err)
				}
			} else {
				// Print to stdout (legacy behavior)
				fmt.Println(switchData)
			}
		}
	}
}

func printHelp() {
	fmt.Printf(`gcool - Git Worktree TUI Manager v%s

A terminal user interface for managing Git worktrees.

USAGE:
    gcool [OPTIONS]

OPTIONS:
    -path <path>    Path to git repository (default: current directory)
    -no-claude      Don't auto-start Claude CLI in tmux session
    -version        Print version and exit
    -help           Show this help message

KEYBINDINGS:
    Navigation:
        ↑/k         Move up
        ↓/j         Move down

    Actions:
        Enter       Switch to selected worktree
        n           Create new worktree with new branch
        a           Create worktree from existing branch
        d           Delete selected worktree
        r           Refresh worktree list
        q/Ctrl+C    Quit

    Modal Navigation:
        Tab         Cycle through inputs/buttons
        Enter       Confirm action
        Esc         Cancel/close modal

SHELL INTEGRATION:
    For directory switching to work, you need to wrap gwt in a shell function.
    Add this to your shell rc file (~/.bashrc, ~/.zshrc, etc.):

    Bash/Zsh:
        gcool() {
            local output
            output=$(command gcool "$@")
            if [ -n "$output" ]; then
                cd "$output" || return
            fi
        }

    Fish:
        function gcool
            set output (command gcool $argv)
            if test -n "$output"
                cd $output
            end
        end

EXAMPLES:
    # Run in current directory
    gcool

    # Run for a specific repository
    gcool -path /path/to/repo

For more information, visit: https://github.com/coollabsio/gcool
`, version)
}
