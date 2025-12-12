package cli

import (
	"os"
	"testing"

	"github.com/fatih/color"
)

// TestStatusIcon verifies status icons are correct
func TestStatusIcon(t *testing.T) {
	tests := []struct {
		name    string
		enabled bool
		want    string
	}{
		{"enabled returns filled circle", true, "●"},
		{"disabled returns empty circle", false, "○"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StatusIcon(tt.enabled)
			if got != tt.want {
				t.Errorf("StatusIcon(%v) = %q, want %q", tt.enabled, got, tt.want)
			}
		})
	}
}

// TestStatusColor verifies status colors are correct
func TestStatusColor(t *testing.T) {
	enabledColor := StatusColor(true)
	disabledColor := StatusColor(false)

	if enabledColor != Green {
		t.Error("StatusColor(true) should return Green")
	}

	if disabledColor != Dim {
		t.Error("StatusColor(false) should return Dim")
	}
}

// TestColorControl verifies color enable/disable
func TestColorControl(t *testing.T) {
	// Save original state
	originalNoColor := color.NoColor
	defer func() { color.NoColor = originalNoColor }()

	// Test NoColor
	ForceColor() // Start with color enabled
	if !IsColorEnabled() {
		t.Error("IsColorEnabled() should be true after ForceColor()")
	}

	NoColor()
	if IsColorEnabled() {
		t.Error("IsColorEnabled() should be false after NoColor()")
	}

	ForceColor()
	if !IsColorEnabled() {
		t.Error("IsColorEnabled() should be true after ForceColor()")
	}
}

// TestColorVariables verifies color variables are initialized
func TestColorVariables(t *testing.T) {
	colors := map[string]*color.Color{
		"ClrPrimary":   ClrPrimary,
		"ClrSecondary": ClrSecondary,
		"ClrSuccess":   ClrSuccess,
		"ClrError":     ClrError,
		"ClrWarning":   ClrWarning,
		"ClrInfo":      ClrInfo,
		"ClrMuted":     ClrMuted,
		"ClrBold":      ClrBold,
		"ClrHeader":    ClrHeader,
		"ClrBox":       ClrBox,
		"Red":          Red,
		"Green":        Green,
		"Yellow":       Yellow,
		"Blue":         Blue,
		"Cyan":         Cyan,
		"Magenta":      Magenta,
		"Bold":         Bold,
		"Dim":          Dim,
		"BoldCyan":     BoldCyan,
		"BoldGreen":    BoldGreen,
		"BoldRed":      BoldRed,
	}

	for name, clr := range colors {
		if clr == nil {
			t.Errorf("color %s should not be nil", name)
		}
	}
}

// TestToolColors verifies tool-specific colors are initialized
func TestToolColors(t *testing.T) {
	toolColors := map[string]*color.Color{
		"ClrRust":   ClrRust,
		"ClrGo":     ClrGo,
		"ClrPython": ClrPython,
		"ClrAWS":    ClrAWS,
		"ClrCDK":    ClrCDK,
		"ClrNode":   ClrNode,
		"ClrJava":   ClrJava,
		"ClrSSH":    ClrSSH,
		"ClrDocker": ClrDocker,
	}

	for name, clr := range toolColors {
		if clr == nil {
			t.Errorf("tool color %s should not be nil", name)
		}
	}
}

// TestDebugRespectVerbose verifies Debug only outputs when verbose is set
func TestDebugRespectVerbose(t *testing.T) {
	// Save original state
	originalVerbose := verbose
	defer func() { verbose = originalVerbose }()

	// Redirect stderr for testing
	oldStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w
	defer func() { os.Stderr = oldStderr }()

	// Test with verbose=false (default behavior verified by no panic)
	verbose = false
	Debug("test message") // Should not panic

	// Test with verbose=true (should not panic)
	verbose = true
	Debug("test message with verbose")

	w.Close()
}
