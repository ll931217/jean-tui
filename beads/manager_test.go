package beads

import (
	"os"
	"path/filepath"
	"testing"
)

// TestIsInitialized tests the IsInitialized method
func TestIsInitialized(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	manager := NewManager(tempDir)

	// Initially should not be initialized
	if manager.IsInitialized() {
		t.Error("Expected beads to not be initialized initially")
	}

	// Create .beads directory
	if err := os.MkdirAll(filepath.Join(tempDir, ".beads"), 0755); err != nil {
		t.Fatalf("Failed to create .beads directory: %v", err)
	}

	// Now should be initialized
	if !manager.IsInitialized() {
		t.Error("Expected beads to be initialized after creating .beads directory")
	}
}

// TestInitialize tests the Initialize method
func TestInitialize(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	manager := NewManager(tempDir)

	// Initialize beads
	if err := manager.Initialize(); err != nil {
		t.Fatalf("Failed to initialize beads: %v", err)
	}

	// Check that .beads directory was created
	beadsDir := filepath.Join(tempDir, ".beads")
	if _, err := os.Stat(beadsDir); os.IsNotExist(err) {
		t.Error("Expected .beads directory to be created")
	}

	// Should now be initialized
	if !manager.IsInitialized() {
		t.Error("Expected beads to be initialized after Initialize()")
	}

	// Note: bd init is NOT idempotent - it will error if run again
	// This is expected behavior for beads
}

// TestGetIssueSummary tests the GetIssueSummary method
func TestGetIssueSummary(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	manager := NewManager(tempDir)

	// Initialize beads
	if err := manager.Initialize(); err != nil {
		t.Fatalf("Failed to initialize beads: %v", err)
	}

	// Get summary for a branch with no issues
	summary, err := manager.GetIssueSummary("test-branch")
	if err != nil {
		t.Fatalf("GetIssueSummary() error = %v", err)
	}
	if summary.OpenCount != 0 {
		t.Errorf("Expected 0 open issues, got %d", summary.OpenCount)
	}
	if summary.ClosedCount != 0 {
		t.Errorf("Expected 0 closed issues, got %d", summary.ClosedCount)
	}
	if summary.TotalCount != 0 {
		t.Errorf("Expected 0 total issues, got %d", summary.TotalCount)
	}
}

// TestGetIssuesForBranch tests the GetIssuesForBranch method
func TestGetIssuesForBranch(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	manager := NewManager(tempDir)

	// Initialize beads
	if err := manager.Initialize(); err != nil {
		t.Fatalf("Failed to initialize beads: %v", err)
	}

	// Get issues for a branch with no issues
	issues, err := manager.GetIssuesForBranch("test-branch")
	if err != nil {
		t.Fatalf("Failed to get issues for branch: %v", err)
	}
	if len(issues) != 0 {
		t.Errorf("Expected 0 issues, got %d", len(issues))
	}
}
