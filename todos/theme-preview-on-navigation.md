# Theme Preview on Navigation

## Feature Request
Enable real-time theme preview when navigating up/down in the theme selection modal. The theme should change instantly as you move through the list, and revert to the original if you press Esc.

## Current Behavior
- User opens Settings → Theme
- User navigates up/down through theme list (only the selection cursor moves)
- User presses Enter to apply and save the theme
- User presses Esc to cancel without applying

## Desired Behavior
1. Open settings → Select theme option
2. As you press ↑/↓, the entire UI immediately updates with that theme (live preview)
3. Press Enter → Theme is saved and you stay in settings
4. Press Esc → Theme reverts to what it was before, return to settings

## Implementation Plan

### 1. Add Original Theme Tracking
**File**: `tui/model.go`

Add new field to Model struct:
```go
originalTheme string  // Stores theme name before entering theme selection modal
```

### 2. Store Original Theme on Modal Open
**File**: `tui/update.go` in `handleSettingsModalInput()`

When opening theme modal (case 1):
```go
case 1:
    // Theme setting - open theme select modal
    m.modal = themeSelectModal
    m.modalFocused = 0
    m.themeIndex = 0

    // Store original theme for preview/revert
    if m.configManager != nil {
        m.originalTheme = m.configManager.GetTheme(m.repoPath)
    }

    // Find current theme in the list
    if m.configManager != nil {
        currentTheme := m.configManager.GetTheme(m.repoPath)
        for i, theme := range m.availableThemes {
            if theme.Name == strings.ToTitle(currentTheme) || theme.Name == strings.Title(currentTheme) {
                m.themeIndex = i
                break
            }
        }
    }
    return m, nil
```

### 3. Create Custom Theme Select Handler
**File**: `tui/update.go`

Replace the current `handleThemeSelectModalInput()` with custom implementation:

```go
func (m Model) handleThemeSelectModalInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "esc", "q":
        // Revert to original theme
        if m.originalTheme != "" {
            ApplyTheme(m.originalTheme)
        }
        m.modal = settingsModal
        m.settingsIndex = 1
        return m, nil

    case "up", "k":
        if m.themeIndex > 0 {
            m.themeIndex--
            // Apply theme preview
            selectedTheme := m.availableThemes[m.themeIndex]
            themeName := strings.ToLower(selectedTheme.Name)
            ApplyTheme(themeName)
        }
        return m, nil

    case "down", "j":
        if m.themeIndex < len(m.availableThemes)-1 {
            m.themeIndex++
            // Apply theme preview
            selectedTheme := m.availableThemes[m.themeIndex]
            themeName := strings.ToLower(selectedTheme.Name)
            ApplyTheme(themeName)
        }
        return m, nil

    case "enter":
        if m.themeIndex >= 0 && m.themeIndex < len(m.availableThemes) {
            selectedTheme := m.availableThemes[m.themeIndex]
            themeName := strings.ToLower(selectedTheme.Name)
            cmd := m.changeTheme(themeName)
            m.modal = settingsModal
            m.settingsIndex = 1
            return m, cmd
        }
        m.modal = settingsModal
        m.settingsIndex = 1
        return m, nil
    }

    return m, nil
}
```

### 4. Edge Cases to Handle
- If `ApplyTheme()` fails during preview, show error but don't crash
- Ensure preview works smoothly without lag
- Clear `m.originalTheme` after modal closes (optional, can reuse)
- Handle case where `m.originalTheme` is empty (fallback to "matrix")

## Technical Details

### Current Theme Application Flow
1. `ApplyTheme(themeName)` in `tui/themes.go` (lines 153-175):
   - Looks up theme colors from `themes` map
   - Updates module-level color variables
   - Calls `rebuildStyles()` to recreate all lipgloss styles
   - Changes are immediate and visible

2. Available themes:
   - `matrix` (green on black - default)
   - `coolify` (purple accents)
   - `dracula` (pink/purple)
   - `nord` (arctic blues)
   - `solarized` (warm blue/yellow)

### Files to Modify
- `tui/model.go` - Add `originalTheme` field
- `tui/update.go` - Modify theme modal handlers for preview functionality

## Benefits
- Instant visual feedback when browsing themes
- No commitment until Enter is pressed
- Safe to experiment with different themes
- Matches modern UI/UX patterns (like color pickers in design tools)

## Testing Checklist
- [ ] Navigate up/down through themes - UI updates instantly
- [ ] Press Enter - theme is saved and persists after restart
- [ ] Press Esc - theme reverts to original
- [ ] Open theme modal with different starting themes
- [ ] Test with all 5 available themes
- [ ] Verify no performance issues with rapid navigation
- [ ] Test error handling if invalid theme somehow selected
