package cli

import (
	"os"
	"path/filepath"
	"testing"
)

// TestIsSymlink verifies symlink detection
func TestIsSymlink(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()

	// Test regular file (not a symlink)
	regularFile := filepath.Join(tmpDir, "regular.txt")
	if err := os.WriteFile(regularFile, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	if isSymlink(regularFile) {
		t.Error("isSymlink should return false for regular file")
	}

	// Test symlink
	symlinkPath := filepath.Join(tmpDir, "symlink")
	if err := os.Symlink(regularFile, symlinkPath); err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}
	if !isSymlink(symlinkPath) {
		t.Error("isSymlink should return true for symlink")
	}

	// Test non-existent path
	if isSymlink(filepath.Join(tmpDir, "nonexistent")) {
		t.Error("isSymlink should return false for non-existent path")
	}

	// Test directory (not a symlink)
	dirPath := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(dirPath, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	if isSymlink(dirPath) {
		t.Error("isSymlink should return false for directory")
	}
}

// TestFileExists verifies file existence check
func TestFileExists(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()

	// Test existing file
	existingFile := filepath.Join(tmpDir, "exists.txt")
	if err := os.WriteFile(existingFile, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	if !fileExists(existingFile) {
		t.Error("fileExists should return true for existing file")
	}

	// Test non-existent file
	if fileExists(filepath.Join(tmpDir, "nonexistent.txt")) {
		t.Error("fileExists should return false for non-existent file")
	}

	// Test existing directory
	dirPath := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(dirPath, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	if !fileExists(dirPath) {
		t.Error("fileExists should return true for existing directory")
	}
}

// TestIsMacOS verifies macOS detection
func TestIsMacOS(t *testing.T) {
	// Save original env
	originalOSType := os.Getenv("OSTYPE")
	defer os.Setenv("OSTYPE", originalOSType)

	// Test with OSTYPE=darwin
	os.Setenv("OSTYPE", "darwin")
	if !isMacOS() {
		t.Error("isMacOS should return true when OSTYPE=darwin")
	}

	// Test with different OSTYPE
	os.Setenv("OSTYPE", "linux-gnu")
	// Note: isMacOS also checks for system file, so result depends on OS
	// On non-macOS, should return false
	// On macOS, might return true because of file check
	// This test primarily verifies the function doesn't panic
	_ = isMacOS()
}

// TestStatusItem verifies statusItem struct
func TestStatusItem(t *testing.T) {
	item := statusItem{
		name: "test",
		ok:   true,
		info: "test info",
		fix:  "test fix",
		skip: false,
	}

	if item.name != "test" {
		t.Errorf("expected name='test', got '%s'", item.name)
	}
	if !item.ok {
		t.Error("expected ok=true")
	}
	if item.info != "test info" {
		t.Errorf("expected info='test info', got '%s'", item.info)
	}
	if item.fix != "test fix" {
		t.Errorf("expected fix='test fix', got '%s'", item.fix)
	}
	if item.skip {
		t.Error("expected skip=false")
	}
}

// TestStatusCommand verifies status command structure
func TestStatusCommand(t *testing.T) {
	cmd := newStatusCmd()

	if cmd.Use != "status" {
		t.Errorf("expected Use='status', got '%s'", cmd.Use)
	}

	// Check alias
	hasAlias := false
	for _, alias := range cmd.Aliases {
		if alias == "s" {
			hasAlias = true
			break
		}
	}
	if !hasAlias {
		t.Error("status command should have 's' alias")
	}

	if cmd.Short == "" {
		t.Error("status command should have Short description")
	}

	if cmd.RunE == nil {
		t.Error("status command should have RunE function")
	}
}
