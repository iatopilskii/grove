// Package config handles application configuration including theme settings.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// AdaptiveColor represents a color that adapts to light/dark terminal themes.
type AdaptiveColor struct {
	Light string `yaml:"light"`
	Dark  string `yaml:"dark"`
}

// ThemeColors defines the color palette for the application theme.
type ThemeColors struct {
	// Primary colors (accent color for active/selected states)
	Primary   AdaptiveColor `yaml:"primary"`
	OnPrimary AdaptiveColor `yaml:"on_primary"`

	// Text colors
	Text      AdaptiveColor `yaml:"text"`
	TextMuted AdaptiveColor `yaml:"text_muted"`

	// Border colors
	Border AdaptiveColor `yaml:"border"`

	// Semantic colors
	Success   AdaptiveColor `yaml:"success"`
	Error     AdaptiveColor `yaml:"error"`
	Info      AdaptiveColor `yaml:"info"`
	OnSuccess AdaptiveColor `yaml:"on_success"`
	OnError   AdaptiveColor `yaml:"on_error"`
	OnInfo    AdaptiveColor `yaml:"on_info"`
}

// Theme defines the visual theme configuration.
type Theme struct {
	Colors ThemeColors `yaml:"colors"`
}

// Config represents the application configuration.
type Config struct {
	Theme Theme `yaml:"theme"`
}

// DefaultConfig returns the default configuration with the built-in color scheme.
func DefaultConfig() Config {
	return Config{
		Theme: Theme{
			Colors: ThemeColors{
				// Primary colors - purple accent
				Primary:   AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"},
				OnPrimary: AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"},

				// Text colors
				Text:      AdaptiveColor{Light: "#333333", Dark: "#CCCCCC"},
				TextMuted: AdaptiveColor{Light: "#666666", Dark: "#888888"},

				// Border colors
				Border: AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"},

				// Success (green)
				Success:   AdaptiveColor{Light: "#2E7D32", Dark: "#4CAF50"},
				OnSuccess: AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"},

				// Error (red)
				Error:   AdaptiveColor{Light: "#C62828", Dark: "#EF5350"},
				OnError: AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"},

				// Info (blue)
				Info:   AdaptiveColor{Light: "#1565C0", Dark: "#42A5F5"},
				OnInfo: AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"},
			},
		},
	}
}

// DefaultConfigPath returns the default path for the application configuration file.
// Uses XDG Base Directory Specification (~/.config/grove/config.yaml).
func DefaultConfigPath() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		configDir = filepath.Join(homeDir, ".config")
	}
	return filepath.Join(configDir, "grove", "config.yaml")
}

// LoadConfig loads configuration from the specified path.
// If the file doesn't exist, returns default configuration with no error.
// If the file exists but is invalid, returns default configuration with an error.
func LoadConfig(path string) (Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist - use defaults silently
			return cfg, nil
		}
		return cfg, fmt.Errorf("reading config file: %w", err)
	}

	// Parse YAML into a temporary config to merge with defaults
	var fileCfg Config
	if err := yaml.Unmarshal(data, &fileCfg); err != nil {
		return cfg, fmt.Errorf("parsing config file: %w", err)
	}

	// Merge file config with defaults (file values override defaults)
	mergeConfig(&cfg, &fileCfg)

	return cfg, nil
}

// mergeConfig merges source config into dest, overriding only non-empty values.
func mergeConfig(dest, source *Config) {
	mergeTheme(&dest.Theme, &source.Theme)
}

func mergeTheme(dest, source *Theme) {
	mergeThemeColors(&dest.Colors, &source.Colors)
}

func mergeThemeColors(dest, source *ThemeColors) {
	mergeAdaptiveColor(&dest.Primary, &source.Primary)
	mergeAdaptiveColor(&dest.OnPrimary, &source.OnPrimary)
	mergeAdaptiveColor(&dest.Text, &source.Text)
	mergeAdaptiveColor(&dest.TextMuted, &source.TextMuted)
	mergeAdaptiveColor(&dest.Border, &source.Border)
	mergeAdaptiveColor(&dest.Success, &source.Success)
	mergeAdaptiveColor(&dest.Error, &source.Error)
	mergeAdaptiveColor(&dest.Info, &source.Info)
	mergeAdaptiveColor(&dest.OnSuccess, &source.OnSuccess)
	mergeAdaptiveColor(&dest.OnError, &source.OnError)
	mergeAdaptiveColor(&dest.OnInfo, &source.OnInfo)
}

func mergeAdaptiveColor(dest, source *AdaptiveColor) {
	if source.Light != "" {
		dest.Light = source.Light
	}
	if source.Dark != "" {
		dest.Dark = source.Dark
	}
}

// GenerateSampleConfig generates a sample configuration YAML string with comments.
func GenerateSampleConfig() string {
	return `# Grove Theme Configuration
# This file allows customization of the application's color scheme.
# Colors use hex format (#RRGGBB) and support light/dark terminal themes.
#
# Location: ~/.config/grove/config.yaml
# Changes require application restart to take effect.

theme:
  colors:
    # Primary accent color (used for selection, active states)
    primary:
      light: "#874BFD"  # Purple accent for light terminals
      dark: "#7D56F4"   # Purple accent for dark terminals

    # Text on primary color (for contrast)
    on_primary:
      light: "#FFFFFF"
      dark: "#FFFFFF"

    # Main text color
    text:
      light: "#333333"  # Dark text for light terminals
      dark: "#CCCCCC"   # Light text for dark terminals

    # Muted/secondary text color
    text_muted:
      light: "#666666"
      dark: "#888888"

    # Border color
    border:
      light: "#874BFD"
      dark: "#7D56F4"

    # Success (green) - for success messages
    success:
      light: "#2E7D32"
      dark: "#4CAF50"

    on_success:
      light: "#FFFFFF"
      dark: "#FFFFFF"

    # Error (red) - for error messages
    error:
      light: "#C62828"
      dark: "#EF5350"

    on_error:
      light: "#FFFFFF"
      dark: "#FFFFFF"

    # Info (blue) - for informational messages
    info:
      light: "#1565C0"
      dark: "#42A5F5"

    on_info:
      light: "#FFFFFF"
      dark: "#FFFFFF"
`
}

// WriteSampleConfig writes a sample configuration file to the specified path.
// Creates parent directories if they don't exist.
func WriteSampleConfig(path string) error {
	// Create parent directories
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	// Write the sample config
	sample := GenerateSampleConfig()
	if err := os.WriteFile(path, []byte(sample), 0644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}
