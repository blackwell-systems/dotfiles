package cli

import (
	"testing"
)

// TestFeaturesCommand verifies features command structure
func TestFeaturesCommand(t *testing.T) {
	cmd := newFeaturesCmd()

	if cmd.Use != "features" {
		t.Errorf("expected Use='features', got '%s'", cmd.Use)
	}

	// Check aliases
	expectedAliases := map[string]bool{"feature": true, "feat": true}
	for _, alias := range cmd.Aliases {
		if !expectedAliases[alias] {
			t.Errorf("unexpected alias: %s", alias)
		}
		delete(expectedAliases, alias)
	}
	if len(expectedAliases) > 0 {
		for alias := range expectedAliases {
			t.Errorf("missing alias: %s", alias)
		}
	}

	if cmd.Short == "" {
		t.Error("features command should have Short description")
	}

	if cmd.Long == "" {
		t.Error("features command should have Long description")
	}
}

// TestFeaturesSubcommands verifies all expected subcommands exist
func TestFeaturesSubcommands(t *testing.T) {
	cmd := newFeaturesCmd()

	expectedSubcommands := []string{
		"list",
		"enable",
		"disable",
		"preset",
		"check",
		"show",
		"validate",
	}

	subcommands := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		subcommands[sub.Name()] = true
	}

	for _, expected := range expectedSubcommands {
		if !subcommands[expected] {
			t.Errorf("expected subcommand '%s' not found", expected)
		}
	}
}

// TestFeaturesListCmd verifies list subcommand
func TestFeaturesListCmd(t *testing.T) {
	cmd := newFeaturesListCmd()

	if cmd.Use != "list [category]" {
		t.Errorf("expected Use='list [category]', got '%s'", cmd.Use)
	}

	// Check flags
	allFlag := cmd.Flags().Lookup("all")
	if allFlag == nil {
		t.Error("list command should have --all flag")
	}

	jsonFlag := cmd.Flags().Lookup("json")
	if jsonFlag == nil {
		t.Error("list command should have --json flag")
	}
}

// TestFeaturesEnableCmd verifies enable subcommand
func TestFeaturesEnableCmd(t *testing.T) {
	cmd := newFeaturesEnableCmd()

	if cmd.Use != "enable <feature>" {
		t.Errorf("expected Use='enable <feature>', got '%s'", cmd.Use)
	}

	// Check persist flag
	persistFlag := cmd.Flags().Lookup("persist")
	if persistFlag == nil {
		t.Error("enable command should have --persist flag")
	}
	if persistFlag.Shorthand != "p" {
		t.Errorf("persist flag should have shorthand 'p', got '%s'", persistFlag.Shorthand)
	}
}

// TestFeaturesDisableCmd verifies disable subcommand
func TestFeaturesDisableCmd(t *testing.T) {
	cmd := newFeaturesDisableCmd()

	if cmd.Use != "disable <feature>" {
		t.Errorf("expected Use='disable <feature>', got '%s'", cmd.Use)
	}

	// Check persist flag
	persistFlag := cmd.Flags().Lookup("persist")
	if persistFlag == nil {
		t.Error("disable command should have --persist flag")
	}
}

// TestFeaturesPresetCmd verifies preset subcommand
func TestFeaturesPresetCmd(t *testing.T) {
	cmd := newFeaturesPresetCmd()

	if cmd.Use != "preset [name]" {
		t.Errorf("expected Use='preset [name]', got '%s'", cmd.Use)
	}

	// Check flags
	listFlag := cmd.Flags().Lookup("list")
	if listFlag == nil {
		t.Error("preset command should have --list flag")
	}

	persistFlag := cmd.Flags().Lookup("persist")
	if persistFlag == nil {
		t.Error("preset command should have --persist flag")
	}
}

// TestFeaturesCheckCmd verifies check subcommand
func TestFeaturesCheckCmd(t *testing.T) {
	cmd := newFeaturesCheckCmd()

	if cmd.Use != "check <feature>" {
		t.Errorf("expected Use='check <feature>', got '%s'", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("check command should have Short description")
	}
}

// TestFeaturesValidateCmd verifies validate subcommand
func TestFeaturesValidateCmd(t *testing.T) {
	cmd := newFeaturesValidateCmd()

	if cmd.Use != "validate" {
		t.Errorf("expected Use='validate', got '%s'", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("validate command should have Short description")
	}
}

// TestInitRegistry verifies registry initialization
func TestInitRegistry(t *testing.T) {
	// Reset global registry
	registry = nil

	// First call should create new registry
	r1 := initRegistry()
	if r1 == nil {
		t.Fatal("initRegistry should return non-nil registry")
	}

	// Second call should return same instance
	r2 := initRegistry()
	if r1 != r2 {
		t.Error("initRegistry should return same instance on subsequent calls")
	}
}
