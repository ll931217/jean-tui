# gcool - Git Worktree TUI Manager

A beautiful terminal user interface for managing Git worktrees, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Features

- **Full CRUD Operations**: Create, list, switch, and delete worktrees
- **Tmux Session Management**: Persistent Claude CLI sessions per worktree
- **Auto-Generated Names**: Random branch and workspace names pre-filled (editable)
- **Organized Workspaces**: All worktrees are created in `.workspaces/` directory
- **Fun Naming**: Names like `happy-panda-42`, `swift-dragon-17`, `brave-falcon-89`
- **Session Persistence**: Detach and return to your work anytime
- **Intuitive Panels UI**: Split-panel interface showing worktrees and detailed information
- **Create from New or Existing Branches**: Choose to create a new branch or use an existing one
- **Shell Integration**: Seamlessly switch directories using shell wrappers
- **Keyboard-First**: Vim-style navigation and shortcuts
- **Beautiful Styling**: Colorful, modern terminal UI using Lipgloss

## Installation

### Using Go Install

```bash
go install github.com/coollabsio/gcool@latest
```

### From Source

```bash
git clone https://github.com/coollabsio/gcool
cd gcool
go build -o gcool
sudo mv gcool /usr/local/bin/
```

## Prerequisites

- **tmux**: Required for persistent session management
  ```bash
  # macOS
  brew install tmux

  # Ubuntu/Debian
  sudo apt install tmux

  # Arch
  sudo pacman -S tmux
  ```

## Shell Integration Setup

The shell wrapper enables:
1. Automatic tmux session creation/attachment
2. Claude CLI auto-start in each worktree
3. Session persistence (detach with `Ctrl+B D`, return anytime)

### Bash / Zsh

Source the provided wrapper in your `~/.bashrc` or `~/.zshrc`:

```bash
source /path/to/gcool/shell/gcool-wrapper.sh
```

Or copy the function manually (see `shell/gcool-wrapper.sh` for full code).

### Fish

Source the provided wrapper in your `~/.config/fish/config.fish`:

```fish
source /path/to/gcool/shell/gcool-wrapper.fish
```

Or copy the function manually (see `shell/gcool-wrapper.fish` for full code).

## Usage

### Basic Usage

Run `gcool` in any Git repository:

```bash
cd /path/to/your/repo
gcool
```

### With Custom Path (for Development)

Test on a different repository without navigating to it:

```bash
gcool -path /path/to/other/repo
```

## Keybindings

### Main View - Navigation
- `↑` / `k` - Move cursor up in worktree list
- `↓` / `j` - Move cursor down in worktree list
- `Enter` / `Space` - Switch to selected worktree (with Claude)
- `t` - Open terminal in worktree (without Claude)

### Main View - Worktree Management
- `n` - Create new worktree with a **new branch** (random name pre-filled)
- `a` - Create worktree from an **existing branch**
- `d` / `x` - Delete selected worktree
- `r` - Refresh worktree list

### Main View - Branch Operations
- `R` (Shift+R) - Rename current branch
- `C` (Shift+C) - Checkout/switch branch in main repository
- `c` - Change base branch for new worktrees

### Main View - Application
- `s` - Open settings menu
- `S` (Shift+S) - View/manage tmux sessions
- `o` - Open worktree in configured editor
- `q` / `Ctrl+C` - Quit application

### Modal Navigation (All Modals)
- `Tab` - Cycle through inputs/buttons
- `Enter` - Confirm action
- `Esc` - Cancel/close modal

### Session List Modal (Press `s`)
- `↑` / `↓` / `j` / `k` - Navigate through sessions
- `Enter` - Attach to selected session
- `k` - Kill selected session
- `Esc` / `q` - Close modal

### Branch Selection Modals (Press `a`, `C`, or `c`)
- `↑` / `↓` - Navigate through branch list
- `Enter` - Select branch
- `Esc` - Cancel

## How It Works

All worktrees are created inside a `.workspaces/` directory in your repository root with randomly generated names like:
- `happy-panda-42`
- `swift-dragon-17`
- `brave-falcon-89`

This keeps your workspace organized and makes it easy to manage multiple feature branches without cluttering your file system.

## Tmux Sessions & Claude CLI

When you switch to a worktree, `gcool` automatically:

1. **Creates or attaches to a tmux session** named `gcool-<branch-name>`
2. **Starts Claude CLI** in the session (by default)
3. **Persists your work** - detach anytime with `Ctrl+B D`

### Session Management

**View all sessions**: Press `s` in the TUI to see active sessions

**Manual session control**:
```bash
# List all gcool sessions
tmux ls | grep gcool-

# Attach to a specific session
tmux attach -t gcool-feature-auth

# Kill a session
tmux kill-session -t gcool-feature-auth
```

**Detach from session**: `Ctrl+B D` (tmux default)

**Disable auto-Claude**: Use the `--no-claude` flag
```bash
gcool --no-claude
```

### How Sessions Work

```
┌─────────────────────────────────────────────────────────┐
│  You: gcool (select "feature-auth")                      │
├─────────────────────────────────────────────────────────┤
│  Shell checks: tmux session "gcool-feature-auth" exists? │
│  ├─ YES → Attach to existing session                   │
│  └─ NO  → Create new session + start Claude CLI        │
└─────────────────────────────────────────────────────────┘
```

**Benefits**:
- Each worktree has its own isolated Claude session
- Work persists across terminal restarts
- Context is maintained per branch
- Easy to switch between multiple features

## UI Overview

```
┌──────────────────────────────┬──────────────────────────────┐
│  📁 Worktrees                │  ℹ️  Details                  │
│                              │                              │
│  ➜ main (current)            │  Branch: main                │
│     └─ my-repo               │  Path: /path/to/my-repo      │
│                              │  Commit: abc1234             │
│  › feature-branch            │  Status: Available           │
│     └─ happy-panda-42        │                              │
│                              │  Press Enter to switch       │
│  bug-fix                     │                              │
│     └─ swift-dragon-17       │                              │
│                              │                              │
└──────────────────────────────┴──────────────────────────────┘
 ↑/↓ navigate • n new • a existing • d delete • enter switch • q quit
```

## Workflow Examples

### Create a New Worktree with a New Branch

1. Press `n` to open the create modal
2. A random branch name is pre-filled (e.g., `happy-panda-42`)
3. You can edit the name or press `Enter` to use it as-is
4. Worktree is created in `.workspaces/` with another random name

### Create a Worktree from an Existing Branch

1. Press `a` to open the branch selection modal
2. Navigate with `↑`/`↓` to select a branch
3. Press `Enter` to confirm
4. Worktree is created instantly with a random name

### Switch to a Worktree

1. Navigate to the desired worktree with `↑`/`↓`
2. Press `Enter` or `Space`
3. Your shell will automatically `cd` to that worktree

### Delete a Worktree

1. Navigate to the worktree you want to delete
2. Press `d`
3. Confirm deletion in the modal
4. The worktree directory will be removed

## Development

### Prerequisites

- **Go 1.21+**: For building and development
- **Git**: Required for all worktree operations
- **tmux**: Required for session management

### Development Commands

```bash
# Run locally
go run main.go

# Run with custom repository path (for testing)
go run main.go -path /path/to/test/repo

# Build binary
go build -o gcool

# Install to system
sudo cp gcool /usr/local/bin/

# Initialize/update dependencies
go mod tidy

# Verify the build
go build -o gcool

# Test with different flags
./gcool --version
./gcool --help
./gcool --no-claude
```

### Project Structure

For detailed codebase documentation, architecture patterns, and development guidelines, see [CLAUDE.md](./CLAUDE.md).

Key areas documented in CLAUDE.md:
- Complete keybinding reference with implementation locations
- Adding new features (keybindings, git operations, modals)
- Message flow and async operation patterns
- File structure with line number references
- Extension points and future enhancements

### Adding New Features

See [CLAUDE.md](./CLAUDE.md) for detailed guides on:
- **Adding a new keybinding**: Step-by-step with code examples
- **Adding a new git operation**: Pattern for extending git functionality
- **Adding a new modal**: Pattern for creating modal dialogs
- **Message flow pattern**: Understanding async operations

## Configuration

### Settings Menu

Press `s` to open the settings menu, where you can configure:

1. **Editor** - Default editor for opening worktrees
   - Press `Enter` on this option to select from available editors
   - Editors: code, cursor, nvim, vim, subl, atom, zed
   - Default: VS Code (`code`)

2. **Base Branch** - Base branch for creating new worktrees
   - Press `Enter` to select from available branches
   - Used when creating new branches with `n` key

All settings are saved per-repository in `~/.config/gcool/config.json`:

```json
{
  "repositories": {
    "/path/to/repo": {
      "base_branch": "main",
      "editor": "code"
    }
  }
}
```

### Tmux Configuration

gcool uses a custom tmux configuration that enhances the default experience without modifying your `~/.tmux.conf`.

**Configuration file**: `~/.config/gcool/tmux.conf`

This config:
- **Sources your `~/.tmux.conf` first** - Your personal settings are preserved
- **Adds gcool-specific enhancements**:
  - Mouse scrolling enabled
  - 10,000 line scrollback buffer
  - 256 color support
  - Better status bar with gcool branding
  - Nice pane border colors

**Customizing**:
You can edit `~/.config/gcool/tmux.conf` to customize gcool's tmux behavior. Changes will apply to new sessions.

**Key features**:
- Mouse wheel scrolling works like a normal terminal
- Click to select panes
- Drag to resize panes
- `Prefix + r` to reload gcool tmux config

### Base Branch

The base branch is used when creating new worktrees with new branches. gcool automatically determines the base branch:

1. Check saved config for repository
2. Fall back to current branch
3. Fall back to default branch (main/master)
4. Fall back to empty string (user must set manually with `c` key)

You can change the base branch at any time by pressing `c` in the main view.

### Editor Integration

gcool includes a built-in editor selection menu for opening worktrees in your IDE.

**Setting your preferred editor:**
1. Press `e` in the main view to open the editor selection modal
2. Use `↑`/`↓` or `j`/`k` to navigate through available editors
3. Press `Enter` to select and save your preference

**Available editors:**
- `code` - VS Code (default)
- `cursor` - Cursor IDE
- `nvim` - Neovim
- `vim` - Vim
- `subl` - Sublime Text
- `atom` - Atom
- `zed` - Zed

**Opening a worktree:**
- Press `o` on any worktree to open it in your configured editor
- The editor preference is saved per repository in `~/.config/gcool/config.json`

## Architecture

### Directory Structure

```
gcool/
├── main.go              # CLI entry point, handles flags and shell integration
├── CLAUDE.md            # Development guide and codebase documentation
├── go.mod               # Module: github.com/coollabsio/gcool
├── config/              # Configuration management
│   └── config.go        # Manages ~/.config/gcool/config.json
├── git/                 # Git operations wrapper
│   └── worktree.go      # Worktree CRUD, branch management, random names
├── session/             # Tmux session management
│   └── tmux.go          # Session creation, attachment, listing, cleanup
├── tui/                 # Bubble Tea TUI (Elm Architecture / MVC)
│   ├── model.go         # State management, data structures, Tea commands
│   ├── update.go        # Event handling, keybindings, state transitions
│   ├── view.go          # UI rendering, modal renderers
│   └── styles.go        # Lipgloss styling definitions
└── shell/               # Shell integration wrappers
    ├── gcool-wrapper.sh   # Bash/Zsh wrapper for directory switching
    └── gcool-wrapper.fish # Fish wrapper for directory switching
```

### Key Architectural Patterns

**Bubble Tea MVC**: The TUI follows the Elm Architecture pattern via Bubble Tea:
- **Model**: Holds all application state (worktrees, branches, sessions, UI state, modals)
- **Update**: Handles messages (keyboard input, async operation results)
- **View**: Renders the UI based on current model state

**Async Operations**: Git and tmux operations are wrapped in Tea commands:
- Operations run asynchronously and return typed messages
- Results are handled in the Update function to update state
- Examples: `worktreesLoadedMsg`, `worktreeCreatedMsg`, `branchRenamedMsg`

**Modal System**: The TUI uses a modal system for different operations:
- Create worktree, delete confirmation, branch selection, session list, rename branch, change base branch
- All modals support Tab navigation, Enter to confirm, Esc to cancel

**Shell Integration Protocol**: Communication with shell wrappers via:
- `GCOOL_SWITCH_FILE` environment variable (preferred): Write switch data to file
- Stdout (legacy): Print switch data in format `path|branch|auto-claude|terminal-only`

**Worktree Organization**: All worktrees are created in `.workspaces/` directory at repository root with randomly generated names (adjective-noun-number pattern)

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT

## Acknowledgments

- Inspired by [git-worktree-tui](https://github.com/FredrikMWold/git-worktree-tui)
- Built with the amazing [Charm](https://charm.sh/) ecosystem
