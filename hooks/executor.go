package hooks

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

// Hook represents a single hook configuration
type Hook struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	Enabled  bool   `json:"enabled"`
	RunAsync bool   `json:"run_async"`
}

// HookContext provides context variables for template expansion
type HookContext struct {
	WorkspacePath   string
	RootPath        string
	BranchName      string
	WorktreeName    string
	Timestamp       string
	User            string
	OldBranchName   string // for rename hooks
	OldWorktreePath string // for move hooks
}

// Executor manages hook execution
type Executor struct {
	repoPath string
}

// NewExecutor creates a new hook executor
func NewExecutor(repoPath string) *Executor {
	return &Executor{repoPath: repoPath}
}

// ExpandTemplates expands template variables in the command string
func (e *Executor) ExpandTemplates(command string, ctx HookContext) (string, error) {
	tmpl, err := template.New("hook").Parse(command)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, ctx); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.String(), nil
}

// ExecuteHook executes a single hook with template expansion
func (e *Executor) ExecuteHook(hook Hook, ctx HookContext) error {
	if !hook.Enabled {
		return nil
	}

	expandedCmd, err := e.ExpandTemplates(hook.Command, ctx)
	if err != nil {
		return fmt.Errorf("template expansion failed: %w", err)
	}

	cmd := exec.Command("sh", "-c", expandedCmd)
	// Use workspace path if it exists, otherwise fall back to root path
	// This handles pre_create hooks where workspace doesn't exist yet
	workDir := ctx.WorkspacePath
	if _, err := os.Stat(workDir); os.IsNotExist(err) {
		workDir = ctx.RootPath
	}
	cmd.Dir = workDir
	cmd.Env = e.buildEnv(ctx)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("hook '%s' failed: %s\nOutput: %s", hook.Name, err.Error(), string(output))
	}
	return nil
}

// ExecuteHooksAsync executes hooks asynchronously in background goroutines
func (e *Executor) ExecuteHooksAsync(hooks []Hook, ctx HookContext) {
	go func() {
		for _, hook := range hooks {
			if err := e.ExecuteHook(hook, ctx); err != nil {
				fmt.Fprintf(os.Stderr, "Hook warning: %v\n", err)
			}
		}
	}()
}

// buildEnv builds JEAN_* environment variables for hook execution
func (e *Executor) buildEnv(ctx HookContext) []string {
	env := os.Environ()
	env = append(env, fmt.Sprintf("JEAN_WORKSPACE_PATH=%s", ctx.WorkspacePath))
	env = append(env, fmt.Sprintf("JEAN_ROOT_PATH=%s", ctx.RootPath))
	env = append(env, fmt.Sprintf("JEAN_BRANCH=%s", ctx.BranchName))
	env = append(env, fmt.Sprintf("JEAN_WORKTREE=%s", ctx.WorktreeName))
	env = append(env, fmt.Sprintf("JEAN_TIMESTAMP=%s", ctx.Timestamp))
	if ctx.OldBranchName != "" {
		env = append(env, fmt.Sprintf("JEAN_OLD_BRANCH=%s", ctx.OldBranchName))
	}
	if ctx.OldWorktreePath != "" {
		env = append(env, fmt.Sprintf("JEAN_OLD_PATH=%s", ctx.OldWorktreePath))
	}
	return env
}
