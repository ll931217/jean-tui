package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/coollabsio/jean-tui/config"
)

// setupTestRepo creates a temporary git repository for testing
func setupTestRepo(t *testing.T) (string, func()) {
	t.Helper()

	tempDir := t.TempDir()

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// Configure git
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git: %v", err)
	}

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git: %v", err)
	}

	// Create initial commit
	readmePath := filepath.Join(tempDir, "README.md")
	if err := os.WriteFile(readmePath, []byte("# Test Repo\n"), 0644); err != nil {
		t.Fatalf("Failed to create README: %v", err)
	}

	cmd = exec.Command("git", "add", "README.md")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to add README: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	cleanup := func() {
		// tempDir is automatically cleaned up by t.TempDir()
	}

	return tempDir, cleanup
}

// createTestScript creates a test shell script that records its execution
func createTestScript(t *testing.T, dir, name, marker string) string {
	t.Helper()

	scriptPath := filepath.Join(dir, name)
	scriptContent := `#!/bin/sh
echo "` + marker + `" >> ` + filepath.Join(dir, "hooks-executed.txt")

	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	return scriptPath
}

// TestPreCreateHooksExecution tests that pre-create hooks are executed before worktree creation
func TestPreCreateHooksExecution(t *testing.T) {
	repoPath, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create config manager
	cfgMgr, err := config.NewManager()
	if err != nil {
		t.Fatalf("Failed to create config manager: %v", err)
	}

	// Create test scripts
	preCreateScript := createTestScript(t, repoPath, "pre-create.sh", "pre-create-executed")

	// Configure hooks
	hook := config.Hook{
		Name:    "test-pre-create",
		Command: preCreateScript,
		Enabled: true,
	}

	if err := cfgMgr.AddHook(repoPath, "pre_create", hook); err != nil {
		t.Fatalf("Failed to add hook: %v", err)
	}

	// Create git manager with hooks
	gitMgr := NewManager(repoPath)
	gitMgr.SetConfigManager(cfgMgr)

	// Create worktree
	workspacePath := filepath.Join(repoPath, ".workspaces", "test-worktree")
	if err := gitMgr.Create(workspacePath, "test-branch", true, ""); err != nil {
		t.Fatalf("Failed to create worktree: %v", err)
	}

	// Check that pre-create hook was executed
	hooksLogPath := filepath.Join(repoPath, "hooks-executed.txt")
	content, err := os.ReadFile(hooksLogPath)
	if err != nil {
		t.Fatalf("Failed to read hooks log: %v", err)
	}

	logContent := string(content)
	if !strings.Contains(logContent, "pre-create-executed") {
		t.Error("Pre-create hook was not executed")
	}

	// Verify worktree was created (hooks should not block unless they error)
	if err != nil {
		t.Errorf("AddWorktree failed: %v", err)
	}
}

// TestPostCreateHooksExecution tests that post-create hooks are executed after worktree creation
func TestPostCreateHooksExecution(t *testing.T) {
	repoPath, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create config manager
	cfgMgr, err := config.NewManager()
	if err != nil {
		t.Fatalf("Failed to create config manager: %v", err)
	}

	// Create test scripts
	postCreateScript := createTestScript(t, repoPath, "post-create.sh", "post-create-executed")

	// Configure hooks
	hook := config.Hook{
		Name:    "test-post-create",
		Command: postCreateScript,
		Enabled: true,
	}

	if err := cfgMgr.AddHook(repoPath, "post_create", hook); err != nil {
		t.Fatalf("Failed to add hook: %v", err)
	}

	// Create git manager with hooks
	gitMgr := NewManager(repoPath)
	gitMgr.SetConfigManager(cfgMgr)

	// Create worktree
	workspacePath := filepath.Join(repoPath, ".workspaces", "test-worktree")
	err = gitMgr.Create(workspacePath, "test-branch", true, "")

	// Check that post-create hook was executed
	hooksLogPath := filepath.Join(repoPath, "hooks-executed.txt")
	content, err := os.ReadFile(hooksLogPath)
	if err != nil {
		t.Fatalf("Failed to read hooks log: %v", err)
	}

	logContent := string(content)
	if !strings.Contains(logContent, "post-create-executed") {
		t.Error("Post-create hook was not executed")
	}
}

// TestPreDeleteHooksExecution tests that pre-delete hooks are executed before worktree deletion
func TestPreDeleteHooksExecution(t *testing.T) {
	repoPath, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create config manager
	cfgMgr, err := config.NewManager()
	if err != nil {
		t.Fatalf("Failed to create config manager: %v", err)
	}

	// Create test scripts
	preDeleteScript := createTestScript(t, repoPath, "pre-delete.sh", "pre-delete-executed")

	// Configure hooks
	hook := config.Hook{
		Name:    "test-pre-delete",
		Command: preDeleteScript,
		Enabled: true,
	}

	if err := cfgMgr.AddHook(repoPath, "pre_delete", hook); err != nil {
		t.Fatalf("Failed to add hook: %v", err)
	}

	// Create git manager with hooks
	gitMgr := NewManager(repoPath)
	gitMgr.SetConfigManager(cfgMgr)

	// Create a worktree first
	workspacePath := filepath.Join(repoPath, ".workspaces", "test-worktree")
	if err := gitMgr.Create(workspacePath, "test-branch", true, ""); err != nil {
		t.Fatalf("Failed to create worktree: %v", err)
	}

	// Delete worktree
	err = gitMgr.Remove(workspacePath, false)

	// Check that pre-delete hook was executed
	hooksLogPath := filepath.Join(repoPath, "hooks-executed.txt")
	content, err := os.ReadFile(hooksLogPath)
	if err != nil {
		t.Fatalf("Failed to read hooks log: %v", err)
	}

	logContent := string(content)
	if !strings.Contains(logContent, "pre-delete-executed") {
		t.Error("Pre-delete hook was not executed")
	}

	// Verify worktree was deleted
	if err != nil {
		t.Errorf("RemoveWorktree failed: %v", err)
	}
}

// TestOnSwitchHooksExecution tests that on-switch hooks are executed when switching worktrees
func TestOnSwitchHooksExecution(t *testing.T) {
	repoPath, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create config manager
	cfgMgr, err := config.NewManager()
	if err != nil {
		t.Fatalf("Failed to create config manager: %v", err)
	}

	// Create test scripts
	onSwitchScript := createTestScript(t, repoPath, "on-switch.sh", "on-switch-executed")

	// Configure hooks
	hook := config.Hook{
		Name:    "test-on-switch",
		Command: onSwitchScript,
		Enabled: true,
	}

	if err := cfgMgr.AddHook(repoPath, "on_switch", hook); err != nil {
		t.Fatalf("Failed to add hook: %v", err)
	}

	// Create git manager with hooks
	gitMgr := NewManager(repoPath)
	gitMgr.SetConfigManager(cfgMgr)

	// Create a worktree
	workspacePath := filepath.Join(repoPath, ".workspaces", "test-worktree")
	if err := gitMgr.Create(workspacePath, "test-branch", true, ""); err != nil {
		t.Fatalf("Failed to create worktree: %v", err)
	}

	// Execute on-switch hooks
	err = gitMgr.ExecuteOnSwitchHooks(workspacePath, "test-branch")

	// Check that on-switch hook was executed
	hooksLogPath := filepath.Join(repoPath, "hooks-executed.txt")
	content, err := os.ReadFile(hooksLogPath)
	if err != nil {
		t.Fatalf("Failed to read hooks log: %v", err)
	}

	logContent := string(content)
	if !strings.Contains(logContent, "on-switch-executed") {
		t.Error("On-switch hook was not executed")
	}
}

// TestHookEnvironmentVariables tests that hooks receive correct environment variables
func TestHookEnvironmentVariables(t *testing.T) {
	repoPath, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create config manager
	cfgMgr, err := config.NewManager()
	if err != nil {
		t.Fatalf("Failed to create config manager: %v", err)
	}

	// Create test script that writes environment variables to a file
	envVarsPath := filepath.Join(repoPath, "env-vars.txt")
	scriptContent := "#!/bin/sh\n"
	scriptContent += "echo \"JEAN_WORKSPACE_PATH=$JEAN_WORKSPACE_PATH\" > " + envVarsPath + "\n"
	scriptContent += "echo \"JEAN_ROOT_PATH=$JEAN_ROOT_PATH\" >> " + envVarsPath + "\n"
	scriptContent += "echo \"JEAN_BRANCH=$JEAN_BRANCH\" >> " + envVarsPath + "\n"
	scriptContent += "echo \"JEAN_WORKTREE=$JEAN_WORKTREE\" >> " + envVarsPath + "\n"

	scriptPath := filepath.Join(repoPath, "test-env.sh")
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	// Configure hook
	hook := config.Hook{
		Name:    "test-env",
		Command: scriptPath,
		Enabled: true,
	}

	if err := cfgMgr.AddHook(repoPath, "pre_create", hook); err != nil {
		t.Fatalf("Failed to add hook: %v", err)
	}

	// Create git manager with hooks
	gitMgr := NewManager(repoPath)
	gitMgr.SetConfigManager(cfgMgr)

	// Create worktree
	workspacePath := filepath.Join(repoPath, ".workspaces", "test-worktree")
	branchName := "test-branch"
	worktreeName := filepath.Base(workspacePath)

	if err := gitMgr.Create(workspacePath, branchName, true, ""); err != nil {
		t.Fatalf("Failed to create worktree: %v", err)
	}

	// Read environment variables from file
	content, err := os.ReadFile(envVarsPath)
	if err != nil {
		t.Fatalf("Failed to read env vars: %v", err)
	}

	envLines := strings.Split(string(content), "\n")
	envMap := make(map[string]string)
	for _, line := range envLines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}

	// Verify environment variables
	tests := []struct {
		key      string
		expected string
	}{
		{"JEAN_WORKSPACE_PATH", workspacePath},
		{"JEAN_ROOT_PATH", repoPath},
		{"JEAN_BRANCH", branchName},
		{"JEAN_WORKTREE", worktreeName},
	}

	for _, tt := range tests {
		if got, ok := envMap[tt.key]; !ok {
			t.Errorf("Environment variable %s not set", tt.key)
		} else if got != tt.expected {
			t.Errorf("Environment variable %s = %q, want %q", tt.key, got, tt.expected)
		}
	}
}

// TestAsyncHookExecution tests that async hooks run in background
func TestAsyncHookExecution(t *testing.T) {
	repoPath, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create config manager
	cfgMgr, err := config.NewManager()
	if err != nil {
		t.Fatalf("Failed to create config manager: %v", err)
	}

	// Create test script that sleeps before writing
	hooksLogPath := filepath.Join(repoPath, "hooks-executed.txt")
	scriptContent := "#!/bin/sh\nsleep 0.1\necho \"async-executed\" >> " + hooksLogPath

	scriptPath := filepath.Join(repoPath, "async-hook.sh")
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	// Configure async hook
	hook := config.Hook{
		Name:     "test-async",
		Command:  scriptPath,
		Enabled:  true,
		RunAsync: true,
	}

	if err := cfgMgr.AddHook(repoPath, "post_create", hook); err != nil {
		t.Fatalf("Failed to add hook: %v", err)
	}

	// Create git manager with hooks
	gitMgr := NewManager(repoPath)
	gitMgr.SetConfigManager(cfgMgr)

	// Create worktree (should return immediately, not wait for sleep)
	workspacePath := filepath.Join(repoPath, ".workspaces", "test-worktree")
	start := time.Now()
	if err := gitMgr.Create(workspacePath, "test-branch", true, ""); err != nil {
		t.Fatalf("Failed to create worktree: %v", err)
	}
	elapsed := time.Since(start)

	// Should return quickly (not wait for sleep)
	if elapsed > 500*time.Millisecond {
		t.Errorf("AddWorktree took too long (%v), likely waited for async hook", elapsed)
	}

	// Wait for async hook to complete
	time.Sleep(500 * time.Millisecond)

	// Check that async hook was executed
	content, err := os.ReadFile(hooksLogPath)
	if err != nil {
		t.Fatalf("Failed to read hooks log: %v", err)
	}

	logContent := string(content)
	if !strings.Contains(logContent, "async-executed") {
		t.Error("Async hook was not executed")
	}
}

// TestDisabledHooksNotExecuted tests that disabled hooks are not executed
func TestDisabledHooksNotExecuted(t *testing.T) {
	repoPath, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create config manager
	cfgMgr, err := config.NewManager()
	if err != nil {
		t.Fatalf("Failed to create config manager: %v", err)
	}

	// Create test script
	scriptPath := createTestScript(t, repoPath, "disabled-hook.sh", "disabled-hook-executed")

	// Configure disabled hook
	hook := config.Hook{
		Name:    "test-disabled",
		Command: scriptPath,
		Enabled: false,
	}

	if err := cfgMgr.AddHook(repoPath, "pre_create", hook); err != nil {
		t.Fatalf("Failed to add hook: %v", err)
	}

	// Create git manager with hooks
	gitMgr := NewManager(repoPath)
	gitMgr.SetConfigManager(cfgMgr)

	// Create worktree
	workspacePath := filepath.Join(repoPath, ".workspaces", "test-worktree")
	if err := gitMgr.Create(workspacePath, "test-branch", true, ""); err != nil {
		t.Fatalf("Failed to create worktree: %v", err)
	}

	// Check that disabled hook was NOT executed
	hooksLogPath := filepath.Join(repoPath, "hooks-executed.txt")
	if _, err := os.Stat(hooksLogPath); err == nil {
		content, _ := os.ReadFile(hooksLogPath)
		logContent := string(content)
		if strings.Contains(logContent, "disabled-hook-executed") {
			t.Error("Disabled hook was executed")
		}
	}
}

// TestPreHookErrorBlocksOperation tests that pre-hook errors block the operation
func TestPreHookErrorBlocksOperation(t *testing.T) {
	repoPath, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create config manager
	cfgMgr, err := config.NewManager()
	if err != nil {
		t.Fatalf("Failed to create config manager: %v", err)
	}

	// Create test script that fails
	scriptContent := `#!/bin/sh
echo "Hook failed" >&2
exit 1`

	scriptPath := filepath.Join(repoPath, "failing-hook.sh")
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	// Configure failing pre-create hook
	hook := config.Hook{
		Name:    "test-failing",
		Command: scriptPath,
		Enabled: true,
	}

	if err := cfgMgr.AddHook(repoPath, "pre_create", hook); err != nil {
		t.Fatalf("Failed to add hook: %v", err)
	}

	// Create git manager with hooks
	gitMgr := NewManager(repoPath)
	gitMgr.SetConfigManager(cfgMgr)

	// Try to create worktree (should fail due to pre-create hook)
	workspacePath := filepath.Join(repoPath, ".workspaces", "test-worktree")
	err = gitMgr.Create(workspacePath, "test-branch", true, "")

	// Verify operation was blocked
	if err == nil {
		t.Error("Expected AddWorktree to fail due to pre-create hook error, but it succeeded")
	}

	// Verify worktree was NOT created
	if _, err := os.Stat(workspacePath); err == nil {
		t.Error("Worktree was created despite pre-create hook error")
	}
}

// TestPostHookErrorDoesNotBlock tests that post-hook errors don't block the operation
func TestPostHookErrorDoesNotBlock(t *testing.T) {
	repoPath, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create config manager
	cfgMgr, err := config.NewManager()
	if err != nil {
		t.Fatalf("Failed to create config manager: %v", err)
	}

	// Create test script that fails
	scriptContent := `#!/bin/sh
echo "Hook failed" >&2
exit 1`

	scriptPath := filepath.Join(repoPath, "failing-post-hook.sh")
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	// Configure failing post-create hook
	hook := config.Hook{
		Name:    "test-failing-post",
		Command: scriptPath,
		Enabled: true,
	}

	if err := cfgMgr.AddHook(repoPath, "post_create", hook); err != nil {
		t.Fatalf("Failed to add hook: %v", err)
	}

	// Create git manager with hooks
	gitMgr := NewManager(repoPath)
	gitMgr.SetConfigManager(cfgMgr)

	// Create worktree (should succeed despite post-create hook error)
	workspacePath := filepath.Join(repoPath, ".workspaces", "test-worktree")
	err = gitMgr.Create(workspacePath, "test-branch", true, "")

	// Verify operation succeeded
	if err != nil {
		t.Errorf("AddWorktree failed due to post-create hook error (should be non-blocking): %v", err)
	}

	// Verify worktree WAS created
	if _, err := os.Stat(workspacePath); os.IsNotExist(err) {
		t.Error("Worktree was not created")
	}
}
