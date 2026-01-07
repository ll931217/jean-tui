---
prd:
  version: v1
  feature_name: beads-hooks-integration
  status: approved
git:
  branch: feat/beads-hooks-integration
  branch_type: feature
  created_at_commit: 816409cef2c91d0efd3474e1c8cb9abdb69fa3cc
  updated_at_commit: 816409cef2c91d0efd3474e1c8cb9abdb69fa3cc
worktree:
  is_worktree: true
  name: feat-beads-hooks-integration
  path: /home/ll931217/GitHub/jean-tui/.git/worktrees/jean-tui.feat-beads-hooks-integration
  repo_root: /home/ll931217/GitHub/jean-tui.feat-beads-hooks-integration
metadata:
  created_at: 2026-01-06T11:02:20Z
  updated_at: 2026-01-06T12:42:56Z
  created_by: Liang-Shih Lin <liangshihlin@gmail.com>
  filename: prd-beads-hooks-integration-v1.md
beads:
  related_issues: []
  related_epics: []
  note: "Using TodoWrite for task tracking (beads not installed)"
code_references:
  - path: "session/tmux.go"
    lines: "26-32"
    reason: "Manager pattern for external tool integration"
  - path: "github/pr.go"
    lines: "10-11"
    reason: "Manager pattern for external tool integration"
  - path: "config/config.go"
    lines: "34-45"
    reason: "Configuration patterns for per-repository settings"
  - path: "git/worktree.go"
    lines: "266-316"
    reason: "Worktree creation flow with setup script integration"
  - path: "tui/model.go"
    lines: "85-100"
    reason: "TUI integration patterns for external managers"
  - path: "tui/update.go"
    lines: "1514-1526"
    reason: "Worktree switching flow and message handlers"
  - path: "tui/view.go"
    lines: "87-142"
    reason: "Worktree list rendering with status indicators"
  - path: "config/scripts.go"
    lines: "1-50"
    reason: "Existing script integration pattern for hooks"
priorities:
  enabled: true
  default: P2
  inference_method: ai_inference_with_review
  requirements:
    - id: FR-1
      text: "Auto-initialize beads on first jean run if .beads directory doesn't exist"
      priority: P2
      confidence: high
      inferred_from: "user confirmed auto-init preference"
      user_confirmed: true
    - id: FR-2
      text: "Read beads issues from .beads/issues.jsonl and .beads/beads.db"
      priority: P2
      confidence: high
      inferred_from: "core beads functionality for task tracking"
      user_confirmed: true
    - id: FR-3
      text: "Display beads issue counts in worktree list items [open/closed]"
      priority: P2
      confidence: high
      inferred_from: "user selected list item count display"
      user_confirmed: true
    - id: FR-4
      text: "Show detailed beads issues in details panel with titles and status"
      priority: P2
      confidence: high
      inferred_from: "user selected details panel display"
      user_confirmed: true
    - id: FR-5
      text: "Provide keyboard shortcut to open beads for selected worktree"
      priority: P2
      confidence: medium
      inferred_from: "user selected quick access option"
      user_confirmed: true
    - id: FR-6
      text: "Implement pre_create and post_create hooks for worktree creation"
      priority: P2
      confidence: high
      inferred_from: "core hooks system functionality"
      user_confirmed: true
    - id: FR-7
      text: "Implement pre_delete and post_delete hooks for worktree deletion"
      priority: P2
      confidence: high
      inferred_from: "core hooks system functionality"
      user_confirmed: true
    - id: FR-8
      text: "Implement on_switch hooks for worktree switching"
      priority: P2
      confidence: high
      inferred_from: "core hooks system functionality"
      user_confirmed: true
    - id: FR-9
      text: "Support template variable expansion in hook commands ({{.WorkspacePath}}, {{.RootPath}}, {{.BranchName}}, etc.)"
      priority: P2
      confidence: high
      inferred_from: "required for hooks to be useful"
      user_confirmed: true
    - id: FR-10
      text: "Configure hooks via JSON config file (~/.config/jean/config.json)"
      priority: P2
      confidence: high
      inferred_from: "user selected TUI + JSON config approach"
      user_confirmed: true
    - id: FR-11
      text: "Configure hooks via TUI settings menu (user-friendly for simple cases)"
      priority: P2
      confidence: high
      inferred_from: "user selected TUI + JSON config approach"
      user_confirmed: true
    - id: FR-12
      text: "Pre-hook failures should block the operation; post-hook failures are warnings"
      priority: P2
      confidence: high
      inferred_from: "user confirmed pre-hooks block behavior"
      user_confirmed: true
    - id: FR-13
      text: "Support async hook execution via run_async flag"
      priority: P2
      confidence: medium
      inferred_from: "useful for long-running post-create hooks like npm install"
      user_confirmed: true
    - id: FR-14
      text: "Parse ports from package.json (config.port and script regex)"
      priority: P2
      confidence: high
      inferred_from: "user selected all three port parsing sources"
      user_confirmed: true
    - id: FR-15
      text: "Parse ports from .env file (PORT, API_PORT, SERVER_PORT variables)"
      priority: P2
      confidence: high
      inferred_from: "user selected all three port parsing sources"
      user_confirmed: true
    - id: FR-16
      text: "Parse ports from docker-compose.yml (exposed host ports)"
      priority: P2
      confidence: high
      inferred_from: "user selected all three port parsing sources"
      user_confirmed: true
    - id: FR-17
      text: "Display active ports in details panel as clickable localhost URLs"
      priority: P2
      confidence: high
      inferred_from: "enhanced worktree info display"
      user_confirmed: true
    - id: FR-18
      text: "Detect AI session waiting status via generic hook mechanism"
      priority: P2
      confidence: medium
      inferred_from: "user selected generic hooks for AI detection"
      user_confirmed: true
---

# Product Requirements Document: Beads Integration & Hooks System

## 1. Introduction/Overview

This PRD defines the integration of **beads** (task tracking system) and a **lifecycle hooks system** into jean-tui, along with enhanced worktree information display. The primary goal is to provide developers with visibility into task tracking per worktree and enable automation through custom hooks on worktree lifecycle events.

### Problem Statement

Developers working with multiple worktrees lose context on which tasks belong to which branches. Additionally, repetitive setup tasks (installing dependencies, copying config files, starting services) must be manually performed each time a worktree is created or switched to.

### Solution

1. **Beads Integration**: Automatically read and display beads issues per worktree, showing open/closed task counts
2. **Hooks System**: Execute custom commands at key lifecycle points (pre/post create, pre/post delete, on-switch) with template variable support
3. **Enhanced Info**: Display active ports parsed from config files and AI session status

---

## 2. Goals

- Provide visibility into beads task tracking for each worktree
- Enable automation through configurable lifecycle hooks
- Support both JSON configuration and TUI-based hook management
- Parse and display active ports from common config formats
- Auto-initialize beads for seamless first-run experience

---

## 3. User Stories

- **As a developer**, I want to see beads issue counts in the worktree list so I know which branches have pending tasks
- **As a developer**, I want to view detailed beads issues in the details panel so I can see task titles and status
- **As a developer**, I want to define hooks that run automatically when creating worktrees so dependencies are installed and config files are copied
- **As a developer**, I want to configure hooks via JSON for version control and via TUI for quick changes
- **As a developer**, I want hooks to support template variables so I can reference worktree paths, branch names, and timestamps
- **As a developer**, I want pre-create hooks to block on failure so invalid configurations prevent worktree creation
- **As a developer**, I want post-create hooks to be non-blocking so worktree creation succeeds even if setup tasks fail
- **As a developer**, I want to see active ports in the details panel so I can quickly open services in my browser
- **As a developer**, I want hooks to run async when needed so long-running tasks don't block the UI

---

## 4. Functional Requirements

| ID   | Requirement                                                         | Priority | Notes                        |
|------|---------------------------------------------------------------------|----------|------------------------------|
| FR-1 | Auto-initialize beads on first jean run if .beads doesn't exist     | P2       | Non-blocking, shows notification |
| FR-2 | Read beads issues from .beads/issues.jsonl and .beads/beads.db      | P2       | Filter by branch name match  |
| FR-3 | Display beads issue counts in worktree list items [open/closed]     | P2       | Format: [3/5]                |
| FR-4 | Show detailed beads issues in details panel with titles and status  | P2       | Grouped by open/closed       |
| FR-5 | Provide keyboard shortcut to open beads for selected worktree       | P2       | Execute `bd` in worktree dir |
| FR-6 | Implement pre_create hooks for worktree creation                    | P2       | Blocking on failure          |
| FR-7  | Implement post_create hooks for worktree creation                   | P2       | Non-blocking warnings         |
| FR-8  | Implement pre_delete hooks for worktree deletion                    | P2       | Blocking on failure          |
| FR-9  | Implement post_delete hooks for worktree deletion                   | P2       | Non-blocking warnings         |
| FR-10 | Implement on_switch hooks for worktree switching                    | P2       | Async execution supported    |
| FR-11 | Support template variable expansion ({{.WorkspacePath}}, etc.)      | P2       | Uses Go text/template        |
| FR-12 | Configure hooks via JSON config file                                | P2       | Stored in ~/.config/jean/    |
| FR-13 | Configure hooks via TUI settings menu                               | P2       | User-friendly forms          |
| FR-14 | Pre-hook failures block operation; post-hook failures are warnings  | P2       | User confirmed behavior       |
| FR-15 | Support async hook execution via run_async flag                     | P2       | Runs in background goroutine  |
| FR-16 | Parse ports from package.json (config.port and script regex)        | P2       | Matches patterns like --port 3000 |
| FR-17 | Parse ports from .env file (PORT, API_PORT, SERVER_PORT)           | P2       | Standard env variable names   |
| FR-18 | Parse ports from docker-compose.yml (exposed host ports)           | P2       | Matches "3000:80" format     |
| FR-19 | Display active ports in details panel as clickable localhost URLs   | P2       | OSC 8 hyperlink codes         |
| FR-20 | Detect AI session waiting status via generic hook mechanism         | P2       | Extensible for any AI tool    |

---

## 5. Non-Goals (Out of Scope)

- **Beads CRUD operations**: Creating, editing, deleting beads issues via jean-tui (use `bd` CLI directly)
- **Hook script editor**: No built-in editor for hook scripts (use external editors)
- **Hook execution logs**: No persistent log storage for hook outputs (shown as notifications only)
- **Port conflict detection**: Not checking if ports are already in use
- **AI session control**: Not starting/stopping AI sessions, only detecting status
- **Real-time beads sync**: Not watching for changes to issues.jsonl (load on refresh only)

---

## 6. Assumptions

- User has `bd` CLI installed or jean-tui can install it automatically
- User has write permissions to ~/.config/jean/ for configuration
- Template variables use Go's text/template syntax
- Hook commands are executed via `sh -c` (shell syntax)
- Ports are parsed as integers; duplicates are removed
- Beads issues are filtered by case-insensitive substring match on branch name, title, description, or ID

---

## 7. Dependencies

- **External tools**:
  - `bd` CLI (beads) - must be installed or auto-installable
  - `sh` - required for hook command execution
- **Internal packages**:
  - `git/` - for worktree lifecycle integration points
  - `config/` - for hook configuration storage
  - `tui/` - for UI rendering and user interactions
- **Go stdlib**:
  - `text/template` - for template variable expansion
  - `encoding/json` - for config parsing
  - `os/exec` - for hook command execution

---

## 8. Acceptance Criteria

### FR-1: Auto-initialize beads
- [ ] On jean-tui startup, check if `.beads` directory exists
- [ ] If not found, execute `bd init` in repository root
- [ ] Show success notification if init succeeds
- [ ] Show warning notification if init fails (non-blocking)
- [ ] Only attempt init once per repository run

### FR-2: Read beads issues
- [ ] Open `.beads/issues.jsonl` and parse JSONL format
- [ ] Filter issues by case-insensitive match on branch name, title, description, or ID
- [ ] Return empty array if file doesn't exist or beads not initialized
- [ ] Handle JSON parse errors gracefully (skip malformed lines)

### FR-3: Display issue counts in list
- [ ] Append `[open/closed]` badge to worktree list items
- [ ] Only show badge if issues exist for the branch
- [ ] Use warning color for open count, success color for closed count
- [ ] Badge appears after uncommitted indicator (●) and behind count (↓3)

### FR-4: Show detailed issues in details panel
- [ ] Add "Beads Issues" section to details panel
- [ ] Display open issues with warning color
- [ ] Display closed issues with success color
- [ ] Show issue title and status for each issue
- [ ] Group by open/closed status

### FR-5: Keyboard shortcut for beads
- [ ] Add keybinding (e.g., `b` key) to open beads for selected worktree
- [ ] Execute `bd` command in worktree directory
- [ ] Use shell integration protocol to switch to worktree and run beads
- [ ] Show error notification if `bd` not found

### FR-6 to FR-10: Lifecycle hooks
- [ ] pre_create hooks execute before `git worktree add`
- [ ] post_create hooks execute after successful worktree creation
- [ ] pre_delete hooks execute before `git worktree remove`
- [ ] post_delete hooks execute after successful worktree deletion
- [ ] on_switch hooks execute when switching to a worktree

### FR-11: Template variables
- [ ] Support `{{.WorkspacePath}}` - path to worktree directory
- [ ] Support `{{.RootPath}}` - path to repository root
- [ ] Support `{{.BranchName}}` - git branch name
- [ ] Support `{{.WorktreeName}}` - worktree directory name
- [ ] Support `{{.Timestamp}}` - ISO 8601 timestamp
- [ ] Support `{{.User}}` - current username
- [ ] Support `{{.OldBranchName}}` - for rename hooks
- [ ] Support `{{.OldWorktreePath}}` - for move hooks

### FR-12 to FR-13: Configuration
- [ ] Hooks stored in `~/.config/jean/config.json` under `repositories.<repo>.hooks`
- [ ] TUI settings menu has "Hooks" option to open hooks management modal
- [ ] Hooks modal shows list of hook types with add/edit/delete actions
- [ ] Hook form fields: name, command, enabled (bool), run_async (bool)
- [ ] Changes to hooks config require saving

### FR-14: Blocking behavior
- [ ] pre_create hook failure aborts worktree creation and shows error
- [ ] pre_delete hook failure aborts worktree deletion and shows error
- [ ] post_create hook failure shows warning but worktree is usable
- [ ] post_delete hook failure shows warning but worktree is deleted
- [ ] on_switch hook failure shows warning but switch completes

### FR-15: Async execution
- [ ] Hooks with `run_async: true` execute in background goroutine
- [ ] Async hook failures written to stderr with "Hook warning:" prefix
- [ ] Async hooks don't block UI or worktree operations

### FR-16 to FR-18: Port parsing
- [ ] Parse `config.port` from package.json
- [ ] Regex match scripts for patterns like `--port 3000`, `-p 8080`, `port=3000`
- [ ] Parse PORT, API_PORT, SERVER_PORT from .env file
- [ ] Parse docker-compose.yml for exposed host ports (matches "3000:80")
- [ ] Remove duplicate ports from final list
- [ ] Display ports in details panel as `http://localhost:PORT`

### FR-19: Port display
- [ ] Add "Active Ports" section to details panel
- [ ] Each port rendered as clickable OSC 8 hyperlink
- [ ] Hyperlink format: `\033]8;;http://localhost:PORT\033\\http://localhost:PORT\033]8;;\033\\`

### FR-20: AI status detection
- [ ] Generic hook mechanism for detecting AI session status
- [ ] Hook defined in config with custom command
- [ ] Hook output parsed for boolean status (waiting/not waiting)
- [ ] Display "⏳ Claude waiting for input" in details panel if true

---

## 9. Design Considerations

### UI Layout Changes

**Worktree List Item:**
```
› feature/auth ● ↓3 [2/5]
                    ↑   ↑  ↑
                    │   │  └─ Beads issue count badge
                    │   └──── Behind count indicator
                    └──────── Uncommitted indicator
```

**Details Panel - New Sections:**
```
Beads Issues:
  2 open
  5 closed

Active Ports:
  http://localhost:3000
  http://localhost:8080

AI Status:
  ⏳ Claude waiting for input
```

### Hook Configuration Schema

```json
{
  "repositories": {
    "/path/to/repo": {
      "hooks": {
        "pre_create": [
          {
            "name": "validate-env",
            "command": "test -f {{.RootPath}}/.env.example",
            "enabled": true
          }
        ],
        "post_create": [
          {
            "name": "install-deps",
            "command": "cd {{.WorkspacePath}} && npm install",
            "enabled": true,
            "run_async": true
          },
          {
            "name": "copy-env",
            "command": "cp {{.RootPath}}/.env.example {{.WorkspacePath}}/.env",
            "enabled": true
          }
        ],
        "on_switch": [
          {
            "name": "notify-switch",
            "command": "echo 'Switched to {{.BranchName}} at {{.Timestamp}}'",
            "enabled": true,
            "run_async": true
          }
        ]
      }
    }
  }
}
```

### Beads Issue Filtering

Issues are included if the branch name matches (case-insensitive) any of:
- Issue ID (e.g., "feat-123" matches issue ID "feat-123")
- Issue title substring (e.g., "auth" matches "Add user authentication")
- Issue description substring (e.g., "auth" matches description containing "authentication")

---

## 10. Technical Considerations

### Package Structure

**New package: `beads/`**
- `issue.go` - Issue and IssueSummary types
- `beads.go` - Manager with IsInitialized(), Initialize(), GetIssuesForBranch(), GetIssueSummary()

**New package: `hooks/`**
- `executor.go` - Executor with ExpandTemplates(), ExecuteHook(), ExecuteHooksAsync()

**New package: `util/`**
- `ports.go` - PortParser with ParsePorts(), parsePackageJSON(), parseEnvFile(), parseDockerCompose()

**New package: `session/`** (or extend existing)
- `ai_status.go` - AIStatusDetector with DetectAISessionState()

### Integration Points

**git/worktree.go**
- Add fields to Worktree: OpenIssues, ClosedIssues, HasBeads, AIWaiting, Ports
- Modify Create() to run pre_create/post_create hooks
- Modify Remove() to run pre_delete/post_delete hooks

**config/config.go**
- Add HooksConfig and Hook types to RepoConfig
- Add hooks config loading/saving

**tui/model.go**
- Add beadsManager, hooksExecutor fields to Model
- Add beadsInitializedMsg message type
- Modify Init() to auto-initialize beads
- Modify loadWorktrees() to load beads data

**tui/update.go**
- Add beadsInitializedMsg handler
- Add on-switch hook execution to worktreeEnsuredMsg handler
- Add hook management modal handlers

**tui/view.go**
- Modify renderWorktreeList() to show beads counts
- Modify renderDetails() to show beads issues, ports, AI status

### Error Handling

- Beads init failures: Warning notification (non-blocking)
- Hook parse errors: Skip invalid hooks, log to debug
- Hook execution errors: Blocking for pre-hooks, warning for post-hooks
- Port parse errors: Silently skip (port parsing is best-effort)
- Beads read errors: Return empty issue list (graceful degradation)

---

## 11. Architecture Patterns

### SOLID Principles

**Single Responsibility:**
- `beads.Manager` only handles beads operations (init, read issues)
- `hooks.Executor` only handles hook execution
- `util.PortParser` only handles port parsing

**Open/Closed:**
- Hook types are extensible (new hook types can be added without modifying executor)
- Port parsers are pluggable (new file formats can be added)

**Dependency Inversion:**
- TUI depends on Manager interfaces, not concrete implementations
- Hook executor depends on HookContext interface

### Creational Patterns

**Factory Pattern:**
- `NewManager()` for beads.Manager
- `NewExecutor()` for hooks.Executor
- `NewPortParser()` for util.PortParser

**Builder Pattern:**
- HookContext constructed incrementally with optional fields (OldBranchName, OldWorktreePath)

### Structural Patterns

**Adapter Pattern:**
- Hook executor adapts shell commands to Go execution context
- Template expansion adapts Go templates to shell command strings

**Decorator Pattern:**
- Hook execution wrapped with: template expansion → env var building → async wrapper → error handling

### Message Pattern (Existing)

**Bubble Tea Messages:**
- `beadsInitializedMsg` - beads auto-init complete
- `hooksExecutedMsg` - hook execution complete
- Extended `worktreeStatusUpdatedMsg` includes AI status and ports

---

## 12. Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| `bd` CLI not installed | Beads features unavailable | Auto-init fails gracefully with notification; add docs for manual install |
| Hook command injection | Security vulnerability | Document security best practices; warn about untrusted hooks; validate template syntax |
| Hook execution hangs | UI freezes | Enforce timeout; recommend `run_async` for long-running hooks |
| Port parsing false positives | Wrong ports displayed | Document expected formats; allow user to disable port parsing |
| Beads performance issues | UI lag on load | Cache issues per worktree; load asynchronously; show loading indicator |
| Hook failures on every operation | Repeated errors | Add hook disable flag; show error with hook name for debugging |

---

## 13. Success Metrics

- **Adoption**: 80% of users enable hooks within 1 week of feature availability
- **Usage**: 50% of worktrees show beads issue counts (indicating beads adoption)
- **Reliability**: <5% hook failure rate for default hooks
- **Performance**: Beads load time <100ms for repos with <100 issues
- **User Satisfaction**: Net Promoter Score >4/5 for hooks feature

---

## 14. Priority/Timeline

- **Priority**: P2 (Normal) - Standard feature with iterative value
- **Target Release**: Next minor version bump
- **Implementation Phases**:
  1. Foundation (beads, hooks, util packages)
  2. Core integration (beads init, issue display, basic hooks)
  3. Advanced features (TUI hook management, port parsing, AI status)
  4. Testing and documentation

---

## 15. Open Questions

- Should hook execution logs be stored persistently or just shown as notifications? (Current: notifications only)
- Should there be a maximum limit on hook execution time? (Current: no timeout, recommend async for long hooks)
- Should port parsing be configurable per-repository? (Current: hardcoded patterns)
- Should there be a global hooks directory that applies to all repositories? (Current: per-repo only)
- How should beads handle multiple issues with same ID across branches? (Current: substring match, first match wins)

---

## 16. Glossary

- **Beads**: Task tracking system using `bd` CLI with issues stored in `.beads/issues.jsonl`
- **Hook**: Custom command executed at specific lifecycle points (pre/post create, pre/post delete, on-switch)
- **Template Variable**: Placeholder in hook commands (e.g., `{{.BranchName}}`) expanded at runtime
- **Worktree**: Git worktree managed by jean-tui, typically in `.workspaces/` directory
- **Pre-hook**: Hook executed before an operation (can block the operation)
- **Post-hook**: Hook executed after an operation (non-blocking)
- **Async Hook**: Hook that runs in background without blocking UI or operation

---

## 17. Changelog

| Version | Date             | Summary of Changes                    |
| ------- | ---------------- | ------------------------------------- |
| 1       | 2026-01-06 11:02 | Initial PRD created from implementation plan |
| 1       | 2026-01-06 11:05 | PRD approved - ready for task generation |

---

## 18. Relevant Code References

| File Path | Lines | Purpose |
|-----------|-------|---------|
| `session/tmux.go` | 26-32 | Manager pattern for external tool integration |
| `github/pr.go` | 10-11 | Manager pattern for external tool integration |
| `config/config.go` | 34-45 | Configuration patterns for per-repository settings |
| `git/worktree.go` | 266-316 | Worktree creation flow with setup script integration |
| `git/worktree.go` | 357-392 | Worktree deletion flow with session cleanup |
| `tui/model.go` | 85-100 | TUI integration patterns for external managers |
| `tui/update.go` | 1514-1526 | Worktree switching flow and message handlers |
| `tui/view.go` | 87-142 | Worktree list rendering with status indicators |
| `tui/view.go` | 147-272 | Details panel rendering for additional information |
| `config/scripts.go` | 1-50 | Existing script integration pattern for hooks |
| `tui/model.go` | 443-692 | Message-based event system for async operations |
