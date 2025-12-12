package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNewEngine verifies engine creation
func TestNewEngine(t *testing.T) {
	e := NewEngine("/tmp/templates")
	if e == nil {
		t.Fatal("NewEngine should return non-nil engine")
	}

	if e.vars == nil {
		t.Error("engine.vars should be initialized")
	}
	if e.arrays == nil {
		t.Error("engine.arrays should be initialized")
	}
	if e.filters == nil {
		t.Error("engine.filters should be initialized")
	}
	if e.templateDir != "/tmp/templates" {
		t.Errorf("templateDir should be '/tmp/templates', got '%s'", e.templateDir)
	}
}

// TestSetAndGetVar verifies variable setting and getting
func TestSetAndGetVar(t *testing.T) {
	e := NewEngine("")

	e.SetVar("name", "test value")

	val, ok := e.GetVar("name")
	if !ok {
		t.Error("GetVar should return ok=true for set variable")
	}
	if val != "test value" {
		t.Errorf("GetVar should return 'test value', got '%s'", val)
	}

	// Test undefined variable
	_, ok = e.GetVar("undefined")
	if ok {
		t.Error("GetVar should return ok=false for undefined variable")
	}
}

// TestGetVarWithEnvOverride verifies environment variable override
func TestGetVarWithEnvOverride(t *testing.T) {
	e := NewEngine("")
	e.SetVar("myvar", "original")

	// Set environment override
	os.Setenv("BLACKDOT_TMPL_MYVAR", "overridden")
	defer os.Unsetenv("BLACKDOT_TMPL_MYVAR")

	val, ok := e.GetVar("myvar")
	if !ok {
		t.Error("GetVar should return ok=true")
	}
	if val != "overridden" {
		t.Errorf("GetVar should return 'overridden' from env, got '%s'", val)
	}
}

// TestSetArray verifies array setting
func TestSetArray(t *testing.T) {
	e := NewEngine("")

	values := []string{"one", "two", "three"}
	e.SetArray("items", values)

	if len(e.arrays["items"]) != 3 {
		t.Errorf("expected 3 items, got %d", len(e.arrays["items"]))
	}
}

// TestRenderVariables verifies variable substitution
func TestRenderVariables(t *testing.T) {
	e := NewEngine("")
	e.SetVar("name", "World")

	result, err := e.Render("Hello, {{ name }}!")
	if err != nil {
		t.Fatalf("Render should not return error: %v", err)
	}
	if result != "Hello, World!" {
		t.Errorf("expected 'Hello, World!', got '%s'", result)
	}
}

// TestRenderUndefinedVariable verifies undefined variables render as empty
func TestRenderUndefinedVariable(t *testing.T) {
	e := NewEngine("")

	result, err := e.Render("Hello, {{ undefined }}!")
	if err != nil {
		t.Fatalf("Render should not return error: %v", err)
	}
	if result != "Hello, !" {
		t.Errorf("expected 'Hello, !', got '%s'", result)
	}
}

// TestBuiltInFilters verifies built-in filter functions
func TestBuiltInFilters(t *testing.T) {
	e := NewEngine("")
	e.SetVar("text", " Hello World ")

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{"upper filter", "{{ text | upper }}", " HELLO WORLD "},
		{"lower filter", "{{ text | lower }}", " hello world "},
		{"trim filter", "{{ text | trim }}", "Hello World"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := e.Render(tt.template)
			if err != nil {
				t.Fatalf("Render should not return error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestPathFilters verifies path filter functions
func TestPathFilters(t *testing.T) {
	e := NewEngine("")
	e.SetVar("path", "/home/user/file.txt")

	result, err := e.Render("{{ path | dirname }}")
	if err != nil {
		t.Fatalf("Render should not return error: %v", err)
	}
	if result != "/home/user" {
		t.Errorf("dirname expected '/home/user', got '%s'", result)
	}

	result, err = e.Render("{{ path | basename }}")
	if err != nil {
		t.Fatalf("Render should not return error: %v", err)
	}
	if result != "file.txt" {
		t.Errorf("basename expected 'file.txt', got '%s'", result)
	}
}

// TestQuoteFilter verifies quote filter
func TestQuoteFilter(t *testing.T) {
	e := NewEngine("")
	e.SetVar("text", "hello")

	result, err := e.Render("{{ text | quote }}")
	if err != nil {
		t.Fatalf("Render should not return error: %v", err)
	}
	if result != `"hello"` {
		t.Errorf("quote expected '\"hello\"', got '%s'", result)
	}
}

// TestIfBlock verifies conditional blocks
func TestIfBlock(t *testing.T) {
	tests := []struct {
		name     string
		vars     map[string]string
		template string
		expected string
	}{
		{
			"truthy string",
			map[string]string{"show": "yes"},
			"{{#if show}}visible{{/if}}",
			"visible",
		},
		{
			"falsy empty",
			map[string]string{},
			"{{#if show}}visible{{/if}}",
			"",
		},
		{
			"falsy false",
			map[string]string{"show": "false"},
			"{{#if show}}visible{{/if}}",
			"",
		},
		{
			"falsy zero",
			map[string]string{"show": "0"},
			"{{#if show}}visible{{/if}}",
			"",
		},
		{
			"if-else true",
			map[string]string{"show": "yes"},
			"{{#if show}}yes{{else}}no{{/if}}",
			"yes",
		},
		{
			"if-else false",
			map[string]string{"show": ""},
			"{{#if show}}yes{{else}}no{{/if}}",
			"no",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEngine("")
			for k, v := range tt.vars {
				e.SetVar(k, v)
			}
			result, err := e.Render(tt.template)
			if err != nil {
				t.Fatalf("Render should not return error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestUnlessBlock verifies unless blocks
func TestUnlessBlock(t *testing.T) {
	tests := []struct {
		name     string
		vars     map[string]string
		template string
		expected string
	}{
		{
			"unless with undefined var",
			map[string]string{},
			"{{#unless show}}hidden{{/unless}}",
			"hidden",
		},
		{
			"unless with empty var",
			map[string]string{"show": ""},
			"{{#unless show}}hidden{{/unless}}",
			"hidden",
		},
		{
			"unless with truthy var",
			map[string]string{"show": "yes"},
			"{{#unless show}}hidden{{/unless}}",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEngine("")
			for k, v := range tt.vars {
				e.SetVar(k, v)
			}
			result, err := e.Render(tt.template)
			if err != nil {
				t.Fatalf("Render should not return error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestEachBlock verifies iteration blocks
func TestEachBlock(t *testing.T) {
	e := NewEngine("")
	e.SetArray("items", []string{"a", "b", "c"})

	result, err := e.Render("{{#each items}}[{{this}}]{{/each}}")
	if err != nil {
		t.Fatalf("Render should not return error: %v", err)
	}
	if result != "[a][b][c]" {
		t.Errorf("expected '[a][b][c]', got '%s'", result)
	}
}

// TestEachBlockIndex verifies @index in each blocks
func TestEachBlockIndex(t *testing.T) {
	e := NewEngine("")
	e.SetArray("items", []string{"a", "b", "c"})

	result, err := e.Render("{{#each items}}{{@index}}:{{this}} {{/each}}")
	if err != nil {
		t.Fatalf("Render should not return error: %v", err)
	}
	if result != "0:a 1:b 2:c " {
		t.Errorf("expected '0:a 1:b 2:c ', got '%s'", result)
	}
}

// TestEachBlockFirstLast verifies @first and @last in each blocks
func TestEachBlockFirstLast(t *testing.T) {
	e := NewEngine("")
	e.SetArray("items", []string{"a", "b", "c"})

	result, err := e.Render("{{#each items}}{{@first}}{{this}}{{@last}} {{/each}}")
	if err != nil {
		t.Fatalf("Render should not return error: %v", err)
	}
	// First item: @first=true, @last=""
	// Middle item: @first="", @last=""
	// Last item: @first="", @last=true
	if !strings.Contains(result, "truea") {
		t.Errorf("expected first item to have @first=true, got '%s'", result)
	}
	if !strings.Contains(result, "ctrue") {
		t.Errorf("expected last item to have @last=true, got '%s'", result)
	}
}

// TestEachBlockEmpty verifies empty array handling
func TestEachBlockEmpty(t *testing.T) {
	e := NewEngine("")
	// Don't set the array

	result, err := e.Render("{{#each items}}[{{this}}]{{/each}}")
	if err != nil {
		t.Fatalf("Render should not return error: %v", err)
	}
	if result != "" {
		t.Errorf("expected empty string for undefined array, got '%s'", result)
	}
}

// TestRegisterFilter verifies custom filter registration
func TestRegisterFilter(t *testing.T) {
	e := NewEngine("")
	e.SetVar("text", "hello")

	// Register custom filter
	e.RegisterFilter("exclaim", func(s string) string {
		return s + "!"
	})

	result, err := e.Render("{{ text | exclaim }}")
	if err != nil {
		t.Fatalf("Render should not return error: %v", err)
	}
	if result != "hello!" {
		t.Errorf("expected 'hello!', got '%s'", result)
	}
}

// TestRenderFile verifies file rendering
func TestRenderFile(t *testing.T) {
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "test.tmpl")

	// Create test template file
	content := "Hello, {{ name }}!"
	if err := os.WriteFile(templatePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	e := NewEngine(tmpDir)
	e.SetVar("name", "World")

	result, err := e.RenderFile(templatePath)
	if err != nil {
		t.Fatalf("RenderFile should not return error: %v", err)
	}
	if result != "Hello, World!" {
		t.Errorf("expected 'Hello, World!', got '%s'", result)
	}
}

// TestRenderFileNotFound verifies error handling for missing files
func TestRenderFileNotFound(t *testing.T) {
	e := NewEngine("")

	_, err := e.RenderFile("/nonexistent/file.tmpl")
	if err == nil {
		t.Error("RenderFile should return error for missing file")
	}
}

// TestLoadVariablesFile verifies variables file loading
func TestLoadVariablesFile(t *testing.T) {
	tmpDir := t.TempDir()
	varsPath := filepath.Join(tmpDir, "variables.sh")

	// Create test variables file
	content := `# Comment line
NAME=John
export EMAIL="john@example.com"
PATH_VAR='/usr/bin'
EMPTY=
`
	if err := os.WriteFile(varsPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	e := NewEngine("")
	if err := e.LoadVariablesFile(varsPath); err != nil {
		t.Fatalf("LoadVariablesFile should not return error: %v", err)
	}

	// Check loaded variables
	if val, ok := e.GetVar("NAME"); !ok || val != "John" {
		t.Errorf("expected NAME='John', got '%s'", val)
	}
	if val, ok := e.GetVar("EMAIL"); !ok || val != "john@example.com" {
		t.Errorf("expected EMAIL='john@example.com', got '%s'", val)
	}
	if val, ok := e.GetVar("PATH_VAR"); !ok || val != "/usr/bin" {
		t.Errorf("expected PATH_VAR='/usr/bin', got '%s'", val)
	}
}

// TestLoadVariablesFileNotFound verifies error handling for missing files
func TestLoadVariablesFileNotFound(t *testing.T) {
	e := NewEngine("")

	err := e.LoadVariablesFile("/nonexistent/variables.sh")
	if err == nil {
		t.Error("LoadVariablesFile should return error for missing file")
	}
}

// TestComplexTemplate verifies a template with multiple features
func TestComplexTemplate(t *testing.T) {
	e := NewEngine("")
	e.SetVar("title", "My App")
	e.SetVar("debug", "true")
	e.SetArray("servers", []string{"web1", "web2", "db1"})

	template := `# {{ title | upper }}
{{#if debug}}DEBUG MODE{{/if}}
Servers:
{{#each servers}}
  - {{this}}
{{/each}}`

	result, err := e.Render(template)
	if err != nil {
		t.Fatalf("Render should not return error: %v", err)
	}

	if !strings.Contains(result, "MY APP") {
		t.Error("expected uppercase title")
	}
	if !strings.Contains(result, "DEBUG MODE") {
		t.Error("expected debug mode message")
	}
	if !strings.Contains(result, "- web1") {
		t.Error("expected server web1")
	}
}
