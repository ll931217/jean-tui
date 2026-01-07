---
prd:
  version: v1
  feature_name: config-editor
  status: implementing
git:
  branch: feat-config-editor
  branch_type: feature
  created_at_commit: 6ccf3b241d30bbc25ed5a70354a2e038de9dd7b8
  updated_at_commit: 6ccf3b241d30bbc25ed5a70354a2e038de9dd7b8
worktree:
  is_worktree: true
  name: feat-config-editor
  path: /home/ll931217/GitHub/jean-tui.feat-config-editor/.git/worktrees/feat-config-editor
  repo_root: /home/ll931217/GitHub/jean-tui.feat-config-editor
metadata:
  created_at: 2026-01-07T22:13:00+08:00
  updated_at: 2026-01-07T20:30:26Z
  created_by: Liang-Shih Lin <liangshihlin@gmail.com>
  filename: prd-config-editor-v1.md
beads:
  related_issues: [jean-tui.feat-llm-integration-bzr, jean-tui.feat-llm-integration-bzr.6, jean-tui.feat-llm-integration-bzr.7, jean-tui.feat-llm-integration-bzr.8, jean-tui.feat-llm-integration-bzr.9, jean-tui.feat-llm-integration-bzr.10, jean-tui.feat-llm-integration-0iv, jean-tui.feat-llm-integration-0iv.1, jean-tui.feat-llm-integration-0iv.2, jean-tui.feat-llm-integration-0iv.3, jean-tui.feat-llm-integration-tjt, jean-tui.feat-llm-integration-tjt.1, jean-tui.feat-llm-integration-tjt.2, jean-tui.feat-llm-integration-tjt.3, jean-tui.feat-llm-integration-tjt.4, jean-tui.feat-llm-integration-tjt.5, jean-tui.feat-llm-integration-9wh, jean-tui.feat-llm-integration-9wh.1, jean-tui.feat-llm-integration-9wh.2, jean-tui.feat-llm-integration-9wh.3, jean-tui.feat-llm-integration-9wh.4, jean-tui.feat-llm-integration-9wh.5, jean-tui.feat-llm-integration-8d0, jean-tui.feat-llm-integration-8d0.1, jean-tui.feat-llm-integration-8d0.2, jean-tui.feat-llm-integration-8d0.3, jean-tui.feat-llm-integration-8d0.4, jean-tui.feat-llm-integration-4vx, jean-tui.feat-llm-integration-4vx.1, jean-tui.feat-llm-integration-4vx.2, jean-tui.feat-llm-integration-4vx.3, jean-tui.feat-llm-integration-4vx.4, jean-tui.feat-llm-integration-4vx.5]
  related_epics: [jean-tui.feat-llm-integration-bzr, jean-tui.feat-llm-integration-0iv, jean-tui.feat-llm-integration-tjt, jean-tui.feat-llm-integration-9wh, jean-tui.feat-llm-integration-8d0, jean-tui.feat-llm-integration-4vx]
code_references:
  - path: "config/config.go"
    lines: "89-145"
    reason: "Config Manager struct, load() and save() methods for config file I/O"
  - path: "config/config.go"
    lines: "199-221"
    reason: "GetEditor() and SetEditor() methods for editor preference handling"
  - path: "config/config.go"
    lines: "34-87"
    reason: "Config and RepoConfig structs defining JSON structure"
  - path: "tui/model.go"
    lines: "1081-1095"
    reason: "openInEditor() function showing external editor invocation pattern"
  - path: "tui/model.go"
    lines: "375-383"
    reason: "Hardcoded editors list - may need to be configurable"
  - path: "tui/update.go"
    lines: "1614-1618"
    reason: "Editor selection keybinding (e) - can add config editor here"
priorities:
  enabled: true
  default: P2
  inference_method: ai_inference_with_review
  requirements:
    - id: FR-1
      text: "Launch external editor to modify config files"
      priority: P1
      confidence: high
      inferred_from: "Core feature requested by user"
      user_confirmed: true
    - id: FR-2
      text: "Validate JSON syntax before applying configuration changes"
      priority: P1
      confidence: high
      inferred_from: "User explicitly requested strict validation"
      user_confirmed: true
    - id: FR-3
      text: "Support both global (~/.config/jean/config.json) and repository-specific configuration"
      priority: P1
      confidence: high
      inferred_from: "User explicitly requested both global and repo support"
      user_confirmed: true
    - id: FR-4
      text: "Provide full access to all config.json fields including AI provider profiles"
      priority: P2
      confidence: medium
      inferred_from: "User wants full config access but can be phased"
      user_confirmed: true
    - id: FR-5
      text: "Automatically reload configuration after external editor closes"
      priority: P2
      confidence: high
      inferred_from: "Expected UX behavior for config editing"
      user_confirmed: true
    - id: FR-6
      text: "Display clear error messages and revert on validation failure"
      priority: P2
      confidence: medium
      inferred_from: "Supports strict validation requirement"
      user_confirmed: true
    - id: FR-7
      text: "Handle API keys and other sensitive fields without special masking"
      priority: P3
      confidence: low
      inferred_from: "User chose full access over security"
      user_confirmed: true
---

# Product Requirements Document: External Config Editor

## 1. Introduction

This PRD describes a feature to enable users to edit jean's configuration files directly through their preferred external text editor. Currently, configuration changes are made through TUI modals, which can be limiting for advanced users who want to modify multiple settings at once or access nested configuration options like AI provider profiles.

The feature will launch the user's configured `$EDITOR` (or a sensible default) to edit either the global configuration file (`~/.config/jean/config.json`) or the repository-specific configuration, then validate and apply the changes upon save.

## 2. Goals

- Enable direct editing of configuration files through external editors (vim, nano, VS Code, etc.)
- Support both global and repository-specific configuration editing
- Provide strict JSON validation to prevent configuration corruption
- Automatically reload configuration after successful edits
- Maintain application stability by rejecting invalid configurations

## 3. User Stories

- **As a developer**, I want to edit my config file in vim so that I can quickly modify multiple settings at once
- **As a power user**, I want to access and edit AI provider profiles directly so that I can configure custom endpoints without navigating multiple modals
- **As a user**, I want the application to validate my JSON changes so that I don't corrupt my configuration with a syntax error
- **As a user**, I want to choose between editing global settings or repository-specific settings so that I can control the scope of my changes
- **As a developer**, I want the application to reload my configuration automatically after editing so that changes take effect immediately

## 4. Functional Requirements

| ID   | Requirement | Priority | Notes |
|------|-------------|----------|-------|
| FR-1 | Launch external editor to modify config files | P1 | Core feature |
| FR-2 | Validate JSON syntax before applying configuration changes | P1 | Prevent corruption |
| FR-3 | Support both global (~/.config/jean/config.json) and repository-specific configuration | P1 | User explicitly requested both |
| FR-4 | Provide full access to all config.json fields including AI provider profiles | P2 | Complete config access |
| FR-5 | Automatically reload configuration after external editor closes | P2 | Expected UX behavior |
| FR-6 | Display clear error messages and revert on validation failure | P2 | Supports validation |
| FR-7 | Handle API keys and other sensitive fields without special masking | P3 | User chose full access |

## 5. Non-Goals (Out of Scope)

- **In-editor JSON schema validation** - Schema validation happens after save, not during editing
- **Config file diff/merge** - No conflict resolution for concurrent edits
- **Multiple file editing** - Only one config file at a time
- **Config file history/rollback** - No version control for configs (users can use git)
- **Custom editor detection** - Relies on `$EDITOR` environment variable or defaults
- **Config file templates** - No scaffolding or example configs
- **Remote config editing** - Only local config files

## 6. Assumptions

- User has a text editor installed (or can use defaults like `vi`)
- `$EDITOR` environment variable may be set (will use fallback if not)
- User understands JSON syntax and structure
- Config file exists at expected location (`~/.config/jean/config.json`)
- User is in a git repository when editing repo-specific config

## 7. Dependencies

- **Existing config package** (`config/config.go`) - Config Manager, load/save methods
- **Existing TUI framework** (Bubble Tea) - Modal system, keybinding handling
- **External editor** - User's `$EDITOR` or system default (vi, nano, etc.)
- **JSON marshaling** - Go's `encoding/json` package for validation

## 8. Acceptance Criteria

### FR-1: Launch External Editor
- **Given** the user presses the config editor keybinding
- **When** the external editor closes
- **Then** the application should read the modified config file

### FR-2: JSON Validation
- **Given** the user has modified the config file
- **When** the editor closes and the file is read
- **Then** the application must validate JSON syntax
- **And** reject the file if JSON is malformed
- **And** display an error message indicating the validation failure

### FR-3: Global vs Repository Config
- **Given** the user is in any jean session
- **When** the user triggers config editing
- **Then** they should be able to choose between:
  - Global config (`~/.config/jean/config.json`)
  - Repository-specific config (if in a worktree)
- **And** the appropriate file should be opened

### FR-4: Full Config Access
- **Given** the user opens the config file
- **When** they modify any valid JSON field
- **Then** the change should be applied after validation
- **Including** nested fields like `ai_provider.profiles`

### FR-5: Auto-Reload
- **Given** a valid config file modification
- **When** the editor closes
- **Then** the application should reload the configuration
- **And** changes should take effect immediately

### FR-6: Error Handling
- **Given** an invalid JSON file after editing
- **When** validation fails
- **Then** the application should:
  - Display a specific error message
  - Keep the original configuration active
  - Not write the invalid config to disk

### FR-7: Sensitive Fields
- **Given** the config file contains API keys
- **When** the file is opened in the editor
- **Then** API keys should be displayed in plaintext
- **No** masking or redaction should occur

## 9. Design Considerations

### UI/UX Flow
1. User presses keybinding (e.g., `E` for "Edit config")
2. Quick selection modal appears (if needed):
   - "Global configuration"
   - "Repository configuration"
3. External editor launches with config file
4. User edits and saves, then closes editor
5. Application validates JSON and reloads config
6. Success notification or error message displayed

### Keybinding
- **Recommended**: `E` (Shift+E) for "Edit config"
- **Alternative**: Add to Settings modal as "Edit Config (External)" option

### Notification Messages
- **Success**: "Configuration reloaded successfully"
- **JSON error**: "Invalid JSON: [parse error details]"
- **File error**: "Failed to read config: [error details]"

## 10. Technical Considerations

### Editor Detection
```go
// Priority order for editor selection
1. $EDITOR environment variable
2. $VISUAL environment variable
3. Sensible defaults: "vi", "nano", "code"
```

### File Path Resolution
```go
// Global config path
globalConfigPath := "~/.config/jean/config.json"

// Repo config path (same file, repo settings are namespaced)
// Note: Currently, jean uses a single config file with repository-specific
// settings nested under the "repositories" key. This may need architectural
// consideration if users want truly separate config files per repository.
```

### Validation Strategy
```go
// Two-stage validation
1. JSON syntax validation (json.Unmarshal)
2. Schema validation (optional, can use JSON schema library)
```

### Temp File Approach
- Consider copying config to temp file before editing
- Only replace original if validation passes
- Prevents leaving disk in corrupted state

### Process Waiting
- Use `exec.Command().Run()` to block until editor closes
- This ensures application knows when editing is complete
- Alternative: Use `exec.Command().Start()` and poll file modification time

## 11. Architecture Patterns

### SOLID Principles

**Single Responsibility Principle**:
- Config editor service handles editor invocation and validation
- Config manager handles loading/saving (existing)
- TUI update handler handles keybindings and state

**Open/Closed Principle**:
- Editor detection strategy should be extensible
- New editor preferences can be added without modifying core logic

**Dependency Inversion**:
- TUI should depend on ConfigEditor interface, not concrete implementation
- Allows testing with mock editor

### Creational Patterns

**Factory Pattern** (Recommended):
```go
type ConfigEditorFactory interface {
    CreateEditor(path string) ConfigEditor
}

type ExternalConfigEditor struct {
    editorCommand string
    configPath    string
}
```

### Structural Patterns

**Adapter Pattern** (If needed):
- Adapt different editor invocation styles (GUI vs terminal)
- Handle editor-specific exit codes and behaviors

## 12. Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| User corrupts config with invalid JSON | High | Strict validation before apply; keep backup |
| Editor doesn't exit cleanly (crash, background) | Medium | Implement timeout; add manual reload option |
| Concurrent edits (multiple jean instances) | Low | File locking or last-write-wins with warning |
| Editor not found on system | Low | Graceful fallback to vi/nano; clear error message |
| Large config files difficult to edit | Low | Consider future enhancement for partial editing |

## 13. Success Metrics

- **Functional**: All P1 requirements implemented and passing tests
- **Reliability**: Zero config corruptions in production (caught by validation)
- **Usability**: Average time to edit config < 30 seconds (including validation)
- **Adoption**: Feature is discoverable (documented in help modal)

## 14. Implementation Notes

### Key Files to Modify
- `config/config.go` - Add validation method, backup/restore logic
- `tui/model.go` - Add config editor command, messages
- `tui/update.go` - Add keybinding handler
- `tui/view.go` - Update help bar with new keybinding

### Testing Strategy
- Unit tests for JSON validation logic
- Integration tests for editor invocation (with mock editor)
- Manual testing with real editors (vim, nano, code)

### Future Enhancements
- JSON schema validation with detailed error messages
- Config file diff view before applying changes
- Support for editing multiple config files in one session
- In-editor syntax highlighting hints (comments in JSON)

## 15. Open Questions

1. **Config File Architecture**: Currently, jean uses a single config file with repository-specific settings under a "repositories" key. Should we:
   - Keep the current architecture (single file, repository settings namespaced)?
   - Move to per-repository config files (`.jean/config.json` in each repo)?
   - This is an architectural decision that affects FR-3 implementation.

2. **Validation Strictness**: Should we validate against a JSON schema, or just check for valid JSON syntax? Schema validation would catch type mismatches but requires maintaining a schema definition.

3. **Temp File Strategy**: Should we edit a temp file and replace on success, or edit in-place? Temp files are safer but may confuse users about where the actual config lives.

## 16. Glossary

- **Global config**: Configuration at `~/.config/jean/config.json` that applies across all repositories
- **Repository-specific config**: Configuration settings for a particular git repository, currently namespaced under the "repositories" key
- **External editor**: User's preferred text editor (vim, nano, VS Code, etc.)
- **JSON validation**: Syntax checking to ensure the config file is valid JSON

## 17. Changelog

| Version | Date | Summary of Changes |
|---------|------|-------------------|
| 1 | 2026-01-07 | Initial PRD |
