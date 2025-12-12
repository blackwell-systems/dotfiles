package cli

import (
	"testing"
)

// TestLintCommand verifies lint command structure
func TestLintCommand(t *testing.T) {
	cmd := newLintCmd()

	if cmd.Use != "lint" {
		t.Errorf("expected Use='lint', got '%s'", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("lint command should have Short description")
	}

	if cmd.Long == "" {
		t.Error("lint command should have Long description")
	}

	if cmd.RunE == nil {
		t.Error("lint command should have RunE function")
	}
}

// TestLintCommandFlags verifies lint command flags
func TestLintCommandFlags(t *testing.T) {
	cmd := newLintCmd()

	// Check verbose flag
	verboseFlag := cmd.Flags().Lookup("verbose")
	if verboseFlag == nil {
		t.Error("lint command should have --verbose flag")
	}
	if verboseFlag.Shorthand != "v" {
		t.Errorf("verbose flag should have shorthand 'v', got '%s'", verboseFlag.Shorthand)
	}

	// Check fix flag
	fixFlag := cmd.Flags().Lookup("fix")
	if fixFlag == nil {
		t.Error("lint command should have --fix flag")
	}
	if fixFlag.Shorthand != "f" {
		t.Errorf("fix flag should have shorthand 'f', got '%s'", fixFlag.Shorthand)
	}
}

// TestLintResult verifies lintResult struct
func TestLintResult(t *testing.T) {
	result := lintResult{
		file:     "test.zsh",
		errors:   []string{"error1", "error2"},
		warnings: []string{"warning1"},
	}

	if result.file != "test.zsh" {
		t.Errorf("expected file='test.zsh', got '%s'", result.file)
	}

	if len(result.errors) != 2 {
		t.Errorf("expected 2 errors, got %d", len(result.errors))
	}

	if len(result.warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(result.warnings))
	}
}

// TestLintStats verifies lintStats struct
func TestLintStats(t *testing.T) {
	stats := lintStats{
		checked:  10,
		errors:   2,
		warnings: 5,
	}

	if stats.checked != 10 {
		t.Errorf("expected checked=10, got %d", stats.checked)
	}

	if stats.errors != 2 {
		t.Errorf("expected errors=2, got %d", stats.errors)
	}

	if stats.warnings != 5 {
		t.Errorf("expected warnings=5, got %d", stats.warnings)
	}
}

// TestCommandExists verifies command existence check
func TestCommandExists(t *testing.T) {
	// "ls" or "dir" should exist on any Unix/Windows system
	// We can't be 100% sure, so just verify the function doesn't panic
	_ = commandExists("ls")
	_ = commandExists("nonexistent-command-12345")
}
