package hooks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestExpandTemplates tests the ExpandTemplates method
func TestExpandTemplates(t *testing.T) {
	executor := NewExecutor("/test/repo")

	ctx := HookContext{
		WorkspacePath: "/test/workspace",
		RootPath:      "/test/repo",
		BranchName:    "feature-branch",
		WorktreeName:  "feature-branch-123",
		Timestamp:     "2024-01-15T10:30:00Z",
		User:          "testuser",
	}

	tests := []struct {
		name     string
		template string
		want     string
	}{
		{
			name:     "simple workspace path",
			template: "{{.WorkspacePath}}",
			want:     "/test/workspace",
		},
		{
			name:     "simple branch name",
			template: "{{.BranchName}}",
			want:     "feature-branch",
		},
		{
			name:     "complex template",
			template: "cd {{.WorkspacePath}} && git checkout {{.BranchName}}",
			want:     "cd /test/workspace && git checkout feature-branch",
		},
		{
			name:     "timestamp",
			template: "echo {{.Timestamp}}",
			want:     "echo 2024-01-15T10:30:00Z",
		},
		{
			name:     "user",
			template: "echo {{.User}}",
			want:     "echo testuser",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := executor.ExpandTemplates(tt.template, ctx)
			if err != nil {
				t.Errorf("ExpandTemplates() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("ExpandTemplates() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestBuildEnv tests the buildEnv method
func TestBuildEnv(t *testing.T) {
	executor := NewExecutor("/test/repo")

	ctx := HookContext{
		WorkspacePath: "/test/workspace",
		RootPath:      "/test/repo",
		BranchName:    "feature-branch",
		WorktreeName:  "feature-branch-123",
		Timestamp:     "2024-01-15T10:30:00Z",
		User:          "testuser",
	}

	env := executor.buildEnv(ctx)

	// Check that JEAN_* variables are set
	expectedVars := map[string]string{
		"JEAN_WORKSPACE_PATH": "/test/workspace",
		"JEAN_ROOT_PATH":      "/test/repo",
		"JEAN_BRANCH":         "feature-branch",
		"JEAN_WORKTREE":       "feature-branch-123",
		"JEAN_TIMESTAMP":      "2024-01-15T10:30:00Z",
	}

	for key, want := range expectedVars {
		found := false
		for _, envVar := range env {
			if strings.HasPrefix(envVar, key+"=") {
				found = true
				got := strings.TrimPrefix(envVar, key+"=")
				if got != want {
					t.Errorf("buildEnv() %s = %q, want %q", key, got, want)
				}
				break
			}
		}
		if !found {
			t.Errorf("buildEnv() missing environment variable: %s", key)
		}
	}
}

// TestBuildEnvWithOldBranch tests the buildEnv method with rename context
func TestBuildEnvWithOldBranch(t *testing.T) {
	executor := NewExecutor("/test/repo")

	ctx := HookContext{
		WorkspacePath: "/test/workspace",
		RootPath:      "/test/repo",
		BranchName:    "new-branch",
		WorktreeName:  "new-branch-123",
		OldBranchName: "old-branch",
		Timestamp:     "2024-01-15T10:30:00Z",
		User:          "testuser",
	}

	env := executor.buildEnv(ctx)

	// Check that JEAN_OLD_BRANCH is set
	found := false
	for _, envVar := range env {
		if strings.HasPrefix(envVar, "JEAN_OLD_BRANCH=") {
			found = true
			got := strings.TrimPrefix(envVar, "JEAN_OLD_BRANCH=")
			if got != "old-branch" {
				t.Errorf("buildEnv() JEAN_OLD_BRANCH = %q, want %q", got, "old-branch")
			}
			break
		}
	}
	if !found {
		t.Error("buildEnv() missing environment variable: JEAN_OLD_BRANCH")
	}
}

// TestExecuteHook tests the ExecuteHook method
func TestExecuteHook(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create a test script that writes to a file
	testScript := filepath.Join(tempDir, "test.sh")
	scriptContent := "#!/bin/sh\necho 'test' > " + filepath.Join(tempDir, "output.txt")
	if err := os.WriteFile(testScript, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	executor := NewExecutor(tempDir)

	hook := Hook{
		Name:    "test-hook",
		Command: filepath.Join(tempDir, "test.sh"),
		Enabled: true,
	}

	ctx := HookContext{
		WorkspacePath: tempDir,
		RootPath:      tempDir,
		BranchName:    "test-branch",
		WorktreeName:  "test-branch",
		Timestamp:     time.Now().Format(time.RFC3339),
		User:          "testuser",
	}

	// Execute the hook
	if err := executor.ExecuteHook(hook, ctx); err != nil {
		t.Fatalf("ExecuteHook() error = %v", err)
	}

	// Check that the output file was created
	outputFile := filepath.Join(tempDir, "output.txt")
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("ExecuteHook() failed to create output file")
	}
}

// TestExecuteHookDisabled tests that disabled hooks are not executed
func TestExecuteHookDisabled(t *testing.T) {
	executor := NewExecutor("/tmp")

	hook := Hook{
		Name:    "disabled-hook",
		Command: "echo 'should not run'",
		Enabled: false,
	}

	ctx := HookContext{
		WorkspacePath: "/tmp",
		RootPath:      "/tmp",
		BranchName:    "test-branch",
		WorktreeName:  "test-branch",
		Timestamp:     time.Now().Format(time.RFC3339),
		User:          "testuser",
	}

	// ExecuteHook should not fail for disabled hooks
	if err := executor.ExecuteHook(hook, ctx); err != nil {
		t.Errorf("ExecuteHook() on disabled hook should not error, got %v", err)
	}
}

// TestExecuteHooksAsync tests the ExecuteHooksAsync method
func TestExecuteHooksAsync(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create a test script that writes to a file
	testScript := filepath.Join(tempDir, "async-test.sh")
	scriptContent := "#!/bin/sh\necho 'async test' > " + filepath.Join(tempDir, "async-output.txt")
	if err := os.WriteFile(testScript, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	executor := NewExecutor(tempDir)

	hooks := []Hook{
		{
			Name:     "async-hook",
			Command:  filepath.Join(tempDir, "async-test.sh"),
			Enabled:  true,
			RunAsync: true,
		},
	}

	ctx := HookContext{
		WorkspacePath: tempDir,
		RootPath:      tempDir,
		BranchName:    "test-branch",
		WorktreeName:  "test-branch",
		Timestamp:     time.Now().Format(time.RFC3339),
		User:          "testuser",
	}

	// Execute hooks asynchronously
	executor.ExecuteHooksAsync(hooks, ctx)

	// Wait for async execution to complete
	time.Sleep(500 * time.Millisecond)

	// Check that the output file was created
	outputFile := filepath.Join(tempDir, "async-output.txt")
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("ExecuteHooksAsync() failed to create output file")
	}
}
