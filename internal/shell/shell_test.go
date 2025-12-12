package shell

import (
	"os"
	"strings"
	"testing"
)

// TestShellTypeConstants verifies shell type constants
func TestShellTypeConstants(t *testing.T) {
	if ShellZsh != "zsh" {
		t.Errorf("expected ShellZsh='zsh', got '%s'", ShellZsh)
	}
	if ShellBash != "bash" {
		t.Errorf("expected ShellBash='bash', got '%s'", ShellBash)
	}
	if ShellFish != "fish" {
		t.Errorf("expected ShellFish='fish', got '%s'", ShellFish)
	}
	if ShellUnknown != "unknown" {
		t.Errorf("expected ShellUnknown='unknown', got '%s'", ShellUnknown)
	}
}

// TestDetect verifies shell detection
func TestDetect(t *testing.T) {
	// Save original SHELL
	originalShell := os.Getenv("SHELL")
	defer os.Setenv("SHELL", originalShell)

	tests := []struct {
		name     string
		shell    string
		expected ShellType
	}{
		{"zsh path", "/bin/zsh", ShellZsh},
		{"zsh usr local", "/usr/local/bin/zsh", ShellZsh},
		{"bash path", "/bin/bash", ShellBash},
		{"bash usr", "/usr/bin/bash", ShellBash},
		{"fish path", "/usr/bin/fish", ShellFish},
		{"fish opt", "/opt/homebrew/bin/fish", ShellFish},
		{"unknown shell", "/bin/sh", ShellUnknown},
		{"empty SHELL", "", ShellUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("SHELL", tt.shell)
			got := Detect()
			if got != tt.expected {
				t.Errorf("Detect() with SHELL=%q = %q, want %q", tt.shell, got, tt.expected)
			}
		})
	}
}

// TestExportVar verifies export statement generation
func TestExportVar(t *testing.T) {
	// Save original SHELL
	originalShell := os.Getenv("SHELL")
	defer os.Setenv("SHELL", originalShell)

	// Test with bash/zsh style
	os.Setenv("SHELL", "/bin/bash")
	result := ExportVar("MY_VAR", "my value")
	if !strings.HasPrefix(result, "export MY_VAR=") {
		t.Errorf("ExportVar for bash should start with 'export MY_VAR=', got '%s'", result)
	}
	if !strings.Contains(result, "my value") {
		t.Errorf("ExportVar should contain the value 'my value', got '%s'", result)
	}

	// Test with fish style
	os.Setenv("SHELL", "/usr/bin/fish")
	result = ExportVar("MY_VAR", "my value")
	if !strings.HasPrefix(result, "set -gx MY_VAR") {
		t.Errorf("ExportVar for fish should start with 'set -gx MY_VAR', got '%s'", result)
	}
}

// TestEvalOutput verifies eval output generation
func TestEvalOutput(t *testing.T) {
	// Save original SHELL
	originalShell := os.Getenv("SHELL")
	defer os.Setenv("SHELL", originalShell)

	os.Setenv("SHELL", "/bin/bash")

	vars := map[string]string{
		"VAR1": "value1",
		"VAR2": "value2",
	}

	result := EvalOutput(vars)

	if !strings.Contains(result, "VAR1") {
		t.Error("EvalOutput should contain VAR1")
	}
	if !strings.Contains(result, "VAR2") {
		t.Error("EvalOutput should contain VAR2")
	}
	if !strings.Contains(result, "export") {
		t.Error("EvalOutput should contain 'export' for bash")
	}
}

// TestEvalOutputEmpty verifies eval output with empty map
func TestEvalOutputEmpty(t *testing.T) {
	result := EvalOutput(map[string]string{})
	if result != "" {
		t.Errorf("EvalOutput with empty map should return empty string, got '%s'", result)
	}
}

// TestSourceCommand verifies source command generation
func TestSourceCommand(t *testing.T) {
	// Save original SHELL
	originalShell := os.Getenv("SHELL")
	defer os.Setenv("SHELL", originalShell)

	// Test with bash/zsh
	os.Setenv("SHELL", "/bin/zsh")
	result := SourceCommand("/path/to/file.sh")
	if !strings.HasPrefix(result, "source") {
		t.Errorf("SourceCommand should start with 'source', got '%s'", result)
	}
	if !strings.Contains(result, "/path/to/file.sh") {
		t.Errorf("SourceCommand should contain the path, got '%s'", result)
	}

	// Test with fish
	os.Setenv("SHELL", "/usr/bin/fish")
	result = SourceCommand("/path/to/file.fish")
	if !strings.HasPrefix(result, "source") {
		t.Errorf("SourceCommand for fish should start with 'source', got '%s'", result)
	}
}

// TestCommonIntegrations verifies common integrations are defined
func TestCommonIntegrations(t *testing.T) {
	integrations := CommonIntegrations()

	if len(integrations) == 0 {
		t.Fatal("CommonIntegrations should return at least one integration")
	}

	expectedIntegrations := []string{"direnv", "starship", "zoxide"}
	integrationNames := make(map[string]bool)
	for _, i := range integrations {
		integrationNames[i.Name] = true
	}

	for _, expected := range expectedIntegrations {
		if !integrationNames[expected] {
			t.Errorf("expected integration '%s' not found", expected)
		}
	}
}

// TestIntegrationStructure verifies integration structure
func TestIntegrationStructure(t *testing.T) {
	integrations := CommonIntegrations()

	for _, integration := range integrations {
		t.Run(integration.Name, func(t *testing.T) {
			if integration.Name == "" {
				t.Error("integration should have a Name")
			}
			if integration.Description == "" {
				t.Error("integration should have a Description")
			}
			if integration.Condition == nil {
				t.Error("integration should have a Condition function")
			}
			if integration.Setup == nil {
				t.Error("integration should have a Setup function")
			}

			// Test that Condition doesn't panic
			_ = integration.Condition()

			// Test that Setup doesn't panic and returns valid output
			_, err := integration.Setup()
			if err != nil {
				t.Errorf("Setup should not return error: %v", err)
			}
		})
	}
}

// TestGenerateFeatureCheck verifies feature check generation
func TestGenerateFeatureCheck(t *testing.T) {
	result := GenerateFeatureCheck("vault")

	if !strings.Contains(result, "vault") {
		t.Error("GenerateFeatureCheck should contain feature name 'vault'")
	}
	if !strings.Contains(result, "VAULT") {
		t.Error("GenerateFeatureCheck should contain uppercase feature name 'VAULT'")
	}
	if !strings.Contains(result, "blackdot features check") {
		t.Error("GenerateFeatureCheck should contain 'blackdot features check'")
	}
	if !strings.Contains(result, "_BLACKDOT_FEATURE_") {
		t.Error("GenerateFeatureCheck should contain '_BLACKDOT_FEATURE_' prefix")
	}
}

// TestGenerateFeatureCheckMultiple verifies multiple feature checks
func TestGenerateFeatureCheckMultiple(t *testing.T) {
	features := []string{"vault", "aws_helpers", "git_hooks"}

	for _, feature := range features {
		t.Run(feature, func(t *testing.T) {
			result := GenerateFeatureCheck(feature)
			upperFeature := strings.ToUpper(feature)

			if !strings.Contains(result, feature) {
				t.Errorf("GenerateFeatureCheck should contain feature name '%s'", feature)
			}
			if !strings.Contains(result, upperFeature) {
				t.Errorf("GenerateFeatureCheck should contain uppercase feature name '%s'", upperFeature)
			}
		})
	}
}
