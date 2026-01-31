// Package config handles application configuration including theme settings.
package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Theme.Colors.Primary.Light == "" {
		t.Error("expected default Primary.Light to be set")
	}
	if cfg.Theme.Colors.Primary.Dark == "" {
		t.Error("expected default Primary.Dark to be set")
	}
	if cfg.Theme.Colors.Text.Light == "" {
		t.Error("expected default Text.Light to be set")
	}
	if cfg.Theme.Colors.Text.Dark == "" {
		t.Error("expected default Text.Dark to be set")
	}
}

func TestDefaultConfigPath(t *testing.T) {
	path := DefaultConfigPath()

	if path == "" {
		t.Error("expected non-empty config path")
	}

	// Should contain gwt in the path
	if !contains(path, "gwt") {
		t.Errorf("expected path to contain 'gwt', got: %s", path)
	}

	// Should end with theme.yaml
	if filepath.Base(path) != "theme.yaml" {
		t.Errorf("expected path to end with 'theme.yaml', got: %s", filepath.Base(path))
	}
}

func TestLoadConfigNoFile(t *testing.T) {
	// Load from a non-existent path should return defaults
	cfg, err := LoadConfig("/non/existent/path/theme.yaml")
	if err != nil {
		t.Errorf("expected no error for non-existent file, got: %v", err)
	}

	// Should have default values
	defaultCfg := DefaultConfig()
	if cfg.Theme.Colors.Primary.Light != defaultCfg.Theme.Colors.Primary.Light {
		t.Error("expected default Primary.Light when file doesn't exist")
	}
}

func TestLoadConfigValidYAML(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "theme.yaml")

	yamlContent := `theme:
  colors:
    primary:
      light: "#FF0000"
      dark: "#00FF00"
    text:
      light: "#111111"
      dark: "#EEEEEE"
`

	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.Theme.Colors.Primary.Light != "#FF0000" {
		t.Errorf("expected Primary.Light to be '#FF0000', got: %s", cfg.Theme.Colors.Primary.Light)
	}
	if cfg.Theme.Colors.Primary.Dark != "#00FF00" {
		t.Errorf("expected Primary.Dark to be '#00FF00', got: %s", cfg.Theme.Colors.Primary.Dark)
	}
}

func TestLoadConfigPartialYAML(t *testing.T) {
	// Create a partial config file - only override some values
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "theme.yaml")

	yamlContent := `theme:
  colors:
    primary:
      light: "#CUSTOM"
`

	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Custom value should be applied
	if cfg.Theme.Colors.Primary.Light != "#CUSTOM" {
		t.Errorf("expected Primary.Light to be '#CUSTOM', got: %s", cfg.Theme.Colors.Primary.Light)
	}

	// Default values should still be present for unspecified fields
	defaultCfg := DefaultConfig()
	if cfg.Theme.Colors.Primary.Dark != defaultCfg.Theme.Colors.Primary.Dark {
		t.Errorf("expected Primary.Dark to use default '%s', got: %s",
			defaultCfg.Theme.Colors.Primary.Dark, cfg.Theme.Colors.Primary.Dark)
	}
}

func TestLoadConfigInvalidYAML(t *testing.T) {
	// Create an invalid YAML file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "theme.yaml")

	invalidYAML := `invalid yaml: [[[`

	if err := os.WriteFile(configPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	// Should return defaults with error
	cfg, err := LoadConfig(configPath)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}

	// Should still return valid defaults
	defaultCfg := DefaultConfig()
	if cfg.Theme.Colors.Primary.Light != defaultCfg.Theme.Colors.Primary.Light {
		t.Error("expected default values when YAML is invalid")
	}
}

func TestLoadConfigMalformedColors(t *testing.T) {
	// Test that malformed color values are handled gracefully
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "theme.yaml")

	yamlContent := `theme:
  colors:
    primary:
      light: "not-a-color"
      dark: "#VALID"
`

	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// The value should be loaded as-is (lipgloss will handle it)
	if cfg.Theme.Colors.Primary.Light != "not-a-color" {
		t.Errorf("expected color value to be loaded as-is")
	}
}

func TestAdaptiveColorHasLightAndDark(t *testing.T) {
	cfg := DefaultConfig()
	colors := cfg.Theme.Colors

	colorTests := []struct {
		name  string
		color AdaptiveColor
	}{
		{"Primary", colors.Primary},
		{"Text", colors.Text},
		{"TextMuted", colors.TextMuted},
		{"Border", colors.Border},
		{"Success", colors.Success},
		{"Error", colors.Error},
		{"Info", colors.Info},
		{"OnPrimary", colors.OnPrimary},
		{"OnSuccess", colors.OnSuccess},
		{"OnError", colors.OnError},
		{"OnInfo", colors.OnInfo},
	}

	for _, tc := range colorTests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.color.Light == "" {
				t.Errorf("%s.Light is empty", tc.name)
			}
			if tc.color.Dark == "" {
				t.Errorf("%s.Dark is empty", tc.name)
			}
		})
	}
}

func TestGenerateSampleConfig(t *testing.T) {
	sample := GenerateSampleConfig()

	if sample == "" {
		t.Error("expected non-empty sample config")
	}

	// Should contain theme section
	if !contains(sample, "theme:") {
		t.Error("expected sample to contain 'theme:' section")
	}

	// Should contain colors section
	if !contains(sample, "colors:") {
		t.Error("expected sample to contain 'colors:' section")
	}

	// Should contain primary color
	if !contains(sample, "primary:") {
		t.Error("expected sample to contain 'primary:' color")
	}
}

func TestWriteSampleConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "subdir", "theme.yaml")

	err := WriteSampleConfig(configPath)
	if err != nil {
		t.Fatalf("failed to write sample config: %v", err)
	}

	// File should exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("expected config file to be created")
	}

	// File should be valid YAML
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Errorf("expected generated sample to be valid YAML: %v", err)
	}

	if cfg.Theme.Colors.Primary.Light == "" {
		t.Error("expected loaded sample config to have primary color")
	}
}

// contains checks if substr is in s
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
