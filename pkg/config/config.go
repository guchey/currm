package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Rule represents information about a Cursor rule
type Rule struct {
	Name        string `yaml:"name"`
	URL         string `yaml:"url"`
	Revision    string `yaml:"revision,omitempty"`    // Specific revision or "latest"
	Description string `yaml:"description,omitempty"` // Description for the rule
	Globs       string `yaml:"globs,omitempty"`       // Glob patterns for file matching
	AlwaysApply bool   `yaml:"alwaysApply,omitempty"` // Whether to always apply this rule
}

// Config represents the structure of the configuration file
type Config struct {
	Rules []Rule `yaml:"rules"`
}

// LoadConfig loads the configuration file from the specified path
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &config, nil
}

// GetRulesDir returns the path to the directory where rule files will be saved
func GetRulesDir() (string, error) {
	// Get current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	// Create path to .cursor/rules directory based on current directory
	rulesDir := filepath.Join(currentDir, ".cursor", "rules")
	if err := os.MkdirAll(rulesDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create rules directory: %w", err)
	}

	return rulesDir, nil
}
