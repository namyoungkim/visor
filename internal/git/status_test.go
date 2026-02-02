package git

// NOTE: These tests use os.Chdir which modifies global state.
// Do NOT use t.Parallel() in these tests to avoid race conditions.

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGetStatus_NotGitRepo(t *testing.T) {
	// Save current directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Create temp directory (not a git repo)
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	status := GetStatus()

	if status.IsRepo {
		t.Error("Expected IsRepo to be false for non-git directory")
	}
	if status.Branch != "" {
		t.Errorf("Expected empty branch, got %q", status.Branch)
	}
}

func TestGetStatus_GitRepo(t *testing.T) {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	// Save current directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Create temp git repo
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Initialize git repo
	if err := exec.Command("git", "init").Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// Configure git user for commit
	exec.Command("git", "config", "user.email", "test@test.com").Run()
	exec.Command("git", "config", "user.name", "Test").Run()

	status := GetStatus()

	if !status.IsRepo {
		t.Error("Expected IsRepo to be true for git directory")
	}
}

func TestGetStatus_WithBranch(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Initialize git repo with initial commit
	exec.Command("git", "init").Run()
	exec.Command("git", "config", "user.email", "test@test.com").Run()
	exec.Command("git", "config", "user.name", "Test").Run()

	// Create a file and commit
	os.WriteFile("test.txt", []byte("test"), 0644)
	exec.Command("git", "add", "test.txt").Run()
	exec.Command("git", "commit", "-m", "initial").Run()

	status := GetStatus()

	if status.Branch == "" {
		t.Error("Expected branch name to be set")
	}
}

func TestGetStatus_WithModifiedFiles(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Initialize git repo with initial commit
	exec.Command("git", "init").Run()
	exec.Command("git", "config", "user.email", "test@test.com").Run()
	exec.Command("git", "config", "user.name", "Test").Run()

	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)
	exec.Command("git", "add", "test.txt").Run()
	exec.Command("git", "commit", "-m", "initial").Run()

	// Modify the file
	os.WriteFile(testFile, []byte("modified"), 0644)

	status := GetStatus()

	if !status.IsDirty {
		t.Error("Expected IsDirty to be true with modified files")
	}
	if status.Modified == 0 {
		t.Error("Expected Modified count > 0")
	}
}

func TestGetStatus_WithStagedFiles(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Initialize git repo with initial commit
	exec.Command("git", "init").Run()
	exec.Command("git", "config", "user.email", "test@test.com").Run()
	exec.Command("git", "config", "user.name", "Test").Run()

	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)
	exec.Command("git", "add", "test.txt").Run()
	exec.Command("git", "commit", "-m", "initial").Run()

	// Add a new file and stage it
	newFile := filepath.Join(tmpDir, "new.txt")
	os.WriteFile(newFile, []byte("new"), 0644)
	exec.Command("git", "add", "new.txt").Run()

	status := GetStatus()

	if !status.IsDirty {
		t.Error("Expected IsDirty to be true with staged files")
	}
	if status.Staged == 0 {
		t.Error("Expected Staged count > 0")
	}
}
