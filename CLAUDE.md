# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

gcool is a Terminal User Interface (TUI) for managing Git worktrees with integrated tmux session management. It's built in Go using the Bubble Tea framework for the TUI and provides persistent Claude CLI sessions per worktree.

## Development Commands

### Build and Run
```bash
# Run locally
go run main.go

# Run with custom repository path (for testing)
go run main.go -path /path/to/test/repo

# Build binary
go build -o gcool

# Install to system
sudo cp gcool /usr/local/bin/
```

### Testing
```bash
# Initialize/update dependencies
go mod tidy

# Verify the build
go build -o gcool

# Test with different flags
./gcool --version
./gcool --help
./gcool --no-claude
```

## Architecture

The application follows a clean separation of concerns:

### Package Structure
- **main.go**: CLI entry point, handles flags and shell integration output
- **tui/**: Bubble Tea TUI implementation (MVC pattern)
  - `model.go`: State management, data structures, and Tea commands
  - `update.go`: Event handling and state transitions
  - `view.go`: UI rendering logic
  - `styles.go`: Lipgloss styling definitions
- **git/**: Git worktree operations wrapper
  - `worktree.go`: All git worktree CRUD operations, branch management, and random name generation
- **session/**: Tmux session management
  - `tmux.go`: Session creation, attachment, listing, and lifecycle management
- **config/**: User configuration persistence
  - `config.go`: Manages base branch settings per repository in `~/.config/gcool/config.json`
- **shell/**: Shell integration wrappers for bash/zsh/fish

### Detailed File Structure

#### tui/model.go
Key components:
- `Model` struct: Main application state (lines 13-33)
  - `worktrees []git.Worktree`: List of all worktrees
  - `cursor int`: Selected worktree index
  - `modal modalType`: Current active modal
  - `sessions []tmux.Session`: Active tmux sessions
  - Various modal-specific state fields
- `modalType` enum: Defines all modal types (lines 36-44)
- Message types: All async operation results (e.g., `worktreesLoadedMsg`, `worktreeCreatedMsg`)
- Tea command functions: Wrap async operations (e.g., `loadWorktrees()`, `createWorktree()`)

#### tui/update.go
Event handling logic:
- `Update(msg tea.Msg)`: Main message router (lines 12-148)
- `handleMainInput(msg tea.KeyMsg)`: Main view keybindings (lines 150-262)
- `handleModalInput(msg tea.KeyMsg)`: Routes to modal-specific handlers (lines 264-695)
- Modal handlers:
  - `handleCreateModalInput()`: New worktree creation
  - `handleDeleteModalInput()`: Deletion confirmation
  - `handleBranchSelectModalInput()`: Branch selection for worktree
  - `handleCheckoutBranchModalInput()`: Checkout branch in main repo
  - `handleSessionListModalInput()`: Tmux session management
  - `handleRenameModalInput()`: Branch renaming
  - `handleChangeBaseBranchModalInput()`: Base branch selection

#### tui/view.go
UI rendering:
- `View()`: Main render function, delegates to modal or main view
- `renderMainView()`: Worktree list and status
- Modal renderers: `renderCreateModal()`, `renderDeleteModal()`, etc.
- `renderHelpBar()`: Bottom help text with keybindings
- Uses lipgloss styles from `styles.go`

#### git/worktree.go
Git operations wrapper:
- `Manager` struct: Handles all git worktree operations
- `ListWorktrees()`: Get all worktrees and their status
- `AddWorktree()`: Create new worktree
- `RemoveWorktree()`: Delete worktree and optionally its branch
- `GetBranches()`: List all local and remote branches
- `RenameBranch()`: Rename current branch
- `CheckoutBranch()`: Switch branch in main repository
- `GenerateRandomName()`: Create random branch names (adjective-noun-number)

#### session/tmux.go
Tmux integration:
- `CreateOrAttachSession()`: Create new tmux session or attach to existing
- `ListSessions()`: Get all active tmux sessions
- `KillSession()`: Terminate a tmux session
- `SessionExists()`: Check if session is running
- `SanitizeSessionName()`: Convert branch names to valid tmux session names

#### config/config.go
Configuration management:
- `Config` struct: Stores repository-specific settings
- `LoadConfig()`: Read from `~/.config/gcool/config.json`
- `SaveConfig()`: Persist configuration changes
- `GetBaseBranch()`: Get base branch for repository
- `SetBaseBranch()`: Update base branch setting

### Key Architectural Patterns

**Bubble Tea MVC**: The TUI follows the Bubble Tea pattern:
- Model holds all state (worktrees, branches, sessions, UI state, modals)
- Update handles messages (keyboard input, async operation results)
- View renders the UI based on current model state

**Async Operations**: Git and tmux operations are wrapped in Tea commands that return messages:
- Operations like `loadWorktrees()`, `createWorktree()`, `deleteWorktree()` run asynchronously
- Results are delivered via typed messages (`worktreesLoadedMsg`, `worktreeCreatedMsg`, etc.)
- The Update function handles these messages and updates state accordingly

**Shell Integration Protocol**: The app communicates with shell wrappers via:
- Environment variable `GCOOL_SWITCH_FILE` (preferred): Write switch data to file
- Stdout (legacy): Print switch data in format `path|branch|auto-claude|terminal-only`
- Shell wrappers read this data to perform `cd` and tmux session management

**Modal Management**: The TUI uses a modal system (`modalType` enum) for different operations:
- `createModal`: Create new worktree with new branch
- `deleteModal`: Confirm worktree deletion
- `branchSelectModal`: Select existing branch for worktree
- `sessionListModal`: View and manage tmux sessions
- `renameModal`: Rename current branch
- `changeBaseBranchModal`: Change base branch for new worktrees

**Session Naming**: Tmux session names are sanitized from branch names:
- Prefix: `gcool-`
- Invalid characters replaced with hyphens
- Format: `gcool-<sanitized-branch-name>`

**Worktree Organization**: All worktrees are created in `.workspaces/` directory at repository root:
- Random names generated from adjectives + nouns + numbers (e.g., `happy-panda-42`)
- Keeps workspace organized and prevents directory conflicts

## Configuration

**User Config Location**: `~/.config/gcool/config.json`
- Stores base branch per repository
- JSON structure: `{"repositories": {"<repo-path>": {"base_branch": "main"}}}`

**Base Branch Logic**:
1. Check saved config for repository
2. Fall back to current branch
3. Fall back to default branch (main/master)
4. Fall back to empty string (user must set manually)

## Dependencies

Key external dependencies:
- `github.com/charmbracelet/bubbletea`: TUI framework
- `github.com/charmbracelet/lipgloss`: Terminal styling
- `github.com/charmbracelet/bubbles`: TUI components (textinput)

## Extension Points

### Potential Feature Additions

**Editor Integration**:
- Add keybinding (e.g., `o`) to open worktree in default IDE
- Detect editor from `$EDITOR` environment variable
- Support multiple editors: VSCode (`code`), Cursor (`cursor`), Neovim (`nvim`), Vim (`vim`)
- Configuration option to set preferred editor per repository

**Enhanced Configuration**:
- Store editor preferences in `~/.config/gcool/config.json`
- Per-repository settings: default branch naming scheme, preferred editor, tmux preferences
- Global defaults with repository overrides

**Worktree Management**:
- Bulk operations: Delete multiple worktrees at once
- Archive old worktrees instead of deleting
- Search/filter worktrees by name or branch
- Sort worktrees by last modified, creation date, or alphabetically

**Branch Management**:
- Create branch from specific commit or tag
- Interactive rebase support
- Merge/rebase branches from TUI
- Show branch history and commits

**Session Management**:
- Multiple session types per worktree (Claude, terminal, editor)
- Session templates (e.g., "start with vim + tmux split + claude")
- Persistent session layouts
- Integration with other tools beyond Claude CLI

**UI Enhancements**:
- Color themes and customization
- Show git status (dirty/clean) per worktree
- Display last commit message for each worktree
- Show active sessions indicator on worktree list

### Adding New External Integrations

The current tmux integration pattern (`session/tmux.go`) can be extended to support other tools:

1. Create new package (e.g., `editor/`) with similar structure to `session/`
2. Define interface for editor operations (open, close, check if running)
3. Implement adapters for different editors (VSCode, Vim, etc.)
4. Add configuration options to select preferred editor
5. Create Tea commands in `tui/model.go` to invoke editor operations
6. Add keybindings in `tui/update.go`

### Testing Strategies

**Unit Tests**:
- Git operations: Mock `exec.Command` for git commands
- Session management: Test session name sanitization, mock tmux commands
- Configuration: Test JSON serialization, file I/O with temp directories

**Integration Tests**:
- Set up test git repository with worktrees
- Verify worktree creation/deletion
- Test session creation and cleanup
- Verify configuration persistence

**TUI Testing**:
- Use Bubble Tea's test utilities for message handling
- Test keyboard input handling with mock messages
- Verify state transitions through modal flows
- Test async operation message handling

## Module Information

**Module Name**: `github.com/coollabsio/gcool`

All internal imports use `github.com/coollabsio/gcool` as the import path. When adding new packages, use this as the base path:
- `github.com/coollabsio/gcool/tui`
- `github.com/coollabsio/gcool/git`
- `github.com/coollabsio/gcool/config`
- `github.com/coollabsio/gcool/session`

## Prerequisites

- **Git**: Required for all worktree operations
- **tmux**: Required for persistent session management
- **Go 1.21+**: For development

## Keybindings

All keybindings are defined in `tui/update.go`. The application uses Bubble Tea's native `tea.KeyMsg` system with string-based matching.

### Main View Keybindings (tui/update.go:150-262)

**Navigation**:
- `↑`, `k` - Move cursor up
- `↓`, `j` - Move cursor down
- `enter`, `space` - Switch to selected worktree (with Claude)
- `t` - Open terminal in worktree (without Claude)

**Worktree Management**:
- `n` - Create new worktree with random branch name
- `a` - Create worktree from existing branch
- `d`, `x` - Delete selected worktree
- `r` - Refresh worktree list

**Branch Operations**:
- `R` (Shift+R) - Rename current branch
- `C` (Shift+C) - Checkout/switch branch in main repository
- `c` - Change base branch for new worktrees

**Application**:
- `q`, `ctrl+c` - Quit application
- `s` - View/manage tmux sessions

### Modal Keybindings

All modals support:
- `esc` - Close modal without action
- `enter` - Confirm action
- `tab` - Navigate between inputs/options (where applicable)

Specific modal handlers are implemented in `tui/update.go`:
- `handleCreateModalInput()` - Create worktree modal
- `handleDeleteModalInput()` - Deletion confirmation
- `handleBranchSelectModalInput()` - Branch selection (uses `↑`/`↓` navigation)
- `handleCheckoutBranchModalInput()` - Checkout modal (uses `↑`/`↓` navigation)
- `handleSessionListModalInput()` - Session list (uses `↑`/`↓`, `k` to kill sessions)
- `handleRenameModalInput()` - Rename modal (text input)
- `handleChangeBaseBranchModalInput()` - Base branch modal (uses `↑`/`↓` navigation)

## Common Patterns

### Adding a New Keybinding

**Example: Adding "o" to open workspace in editor**

1. **Add message type** in `tui/model.go`:
```go
type editorOpenedMsg struct {
    err error
}
```

2. **Add keybinding** in `tui/update.go` in `handleMainInput()`:
```go
case "o":
    if wt := m.selectedWorktree(); wt != nil {
        return m, m.openInEditor(wt.Path)
    }
```

3. **Create command function** in `tui/model.go`:
```go
func (m Model) openInEditor(path string) tea.Cmd {
    return func() tea.Msg {
        editor := os.Getenv("EDITOR")
        if editor == "" {
            editor = "code"  // fallback
        }
        cmd := exec.Command(editor, path)
        err := cmd.Start()
        return editorOpenedMsg{err: err}
    }
}
```

4. **Handle message** in `tui/update.go` in `Update()`:
```go
case editorOpenedMsg:
    if msg.err != nil {
        m.status = "Failed to open editor"
    } else {
        m.status = "Opened in editor"
    }
    return m, nil
```

5. **Update help bar** in `tui/view.go` in `renderHelpBar()`:
```go
row1 := []string{
    "↑/↓ navigate",
    "o open editor",  // Add this
    // ... rest
}
```

### Adding a New Git Operation

1. Add method to `git.Manager` in `git/worktree.go`
2. Create a Tea command function in `tui/model.go` that calls the git method
3. Define a message type for the result
4. Handle the message in `tui/update.go`
5. Update the view in `tui/view.go` if needed

### Adding a New Modal

1. Add new `modalType` constant in `tui/model.go`
2. Add modal state fields to `Model` struct
3. Add keybinding to open modal in `tui/update.go`
4. Implement modal rendering in `tui/view.go`
5. Handle modal interactions (Tab, Enter, Esc) in `tui/update.go`

### Message Flow Pattern

Async operations in gcool follow this pattern:

1. **User Action** → Keybinding triggers command
2. **Command Function** → Returns `tea.Cmd` that executes async operation
3. **Operation Result** → Wrapped in typed message (e.g., `worktreeCreatedMsg`)
4. **Update Handler** → Receives message, updates model state
5. **View Render** → UI reflects new state

Example flow for creating worktree:
```
User presses 'n'
  → handleMainInput() opens createModal
  → User enters branch name and presses Enter
  → handleCreateModalInput() calls createWorktree() command
  → createWorktree() runs git operations asynchronously
  → Returns worktreeCreatedMsg with result
  → Update() handles worktreeCreatedMsg
  → Updates worktree list and closes modal
  → View() renders updated list
```

### Error Handling

Errors from async operations are stored in `model.err` and displayed in the status bar at the bottom of the UI. The status bar shows:
- Success messages (e.g., "Worktree created")
- Error messages (e.g., "Failed to create worktree: ...")
- Current operation status
