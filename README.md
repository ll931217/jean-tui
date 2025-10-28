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
go install github.com/heyandras/gcool@latest
```

### From Source

```bash
git clone https://github.com/heyandras/gcool
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

### Navigation
- `↑` / `k` - Move up
- `↓` / `j` - Move down
- `Enter` / `Space` - Switch to selected worktree

### Actions
- `n` - Create new worktree with a **new branch**
- `a` - Create worktree from an **existing branch**
- `d` / `x` - Delete selected worktree
- `s` - Show active tmux sessions
- `r` - Refresh worktree list
- `q` / `Ctrl+C` - Quit

### Modal Navigation
- `Tab` - Cycle through inputs and buttons
- `Enter` - Confirm action
- `Esc` - Cancel/close modal

### Session List (Press `s`)
- `↑` / `↓` / `j` / `k` - Navigate sessions
- `Enter` - Attach to selected session
- `x` / `d` - Kill selected session
- `Esc` / `q` - Close modal

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

- Go 1.21 or later
- Git

### Running Locally

```bash
go run main.go
```

### Testing on Another Repo

```bash
go run main.go -path /path/to/test/repo
```

### Building

```bash
go build -o gcool
```

## Architecture

```
gcool/
├── main.go           # CLI entry point
├── git/              # Git operations wrapper
│   └── worktree.go
├── tui/              # Bubble Tea TUI
│   ├── model.go      # State and data structures
│   ├── update.go     # Event handling and state updates
│   ├── view.go       # UI rendering
│   └── styles.go     # Lipgloss styling
└── shell/            # Shell integration wrappers
    ├── gcool-wrapper.sh
    └── gcool-wrapper.fish
```

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
