package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary configuration file for testing
	tempFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Test configuration content
	configContent := `rules:
  - name: "test-rule-1"
    url: "https://example.com/rule1"
  - name: "test-rule-2"
    url: "https://example.com/rule2"
`
	if _, err := tempFile.Write([]byte(configContent)); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	// Execute the function under test
	cfg, err := LoadConfig(tempFile.Name())
	if err != nil {
		t.Fatalf("LoadConfig function returned an error: %v", err)
	}

	// Verify the results
	if len(cfg.Rules) != 2 {
		t.Errorf("Expected number of rules: 2, Actual: %d", len(cfg.Rules))
	}

	// Verify the content of each rule
	expectedRules := []Rule{
		{Name: "test-rule-1", URL: "https://example.com/rule1"},
		{Name: "test-rule-2", URL: "https://example.com/rule2"},
	}

	for i, expected := range expectedRules {
		if cfg.Rules[i].Name != expected.Name {
			t.Errorf("Rule %d name does not match. Expected: %s, Actual: %s", i, expected.Name, cfg.Rules[i].Name)
		}
		if cfg.Rules[i].URL != expected.URL {
			t.Errorf("Rule %d URL does not match. Expected: %s, Actual: %s", i, expected.URL, cfg.Rules[i].URL)
		}
	}

	// Test with non-existent file
	_, err = LoadConfig("non-existent-file.yaml")
	if err == nil {
		t.Error("No error occurred for a non-existent file")
	}

	// Test with invalid YAML
	invalidYamlFile, err := os.CreateTemp("", "invalid-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create invalid YAML file: %v", err)
	}
	defer os.Remove(invalidYamlFile.Name())

	if _, err := invalidYamlFile.Write([]byte("invalid: yaml: content:")); err != nil {
		t.Fatalf("Failed to write to invalid YAML file: %v", err)
	}
	if err := invalidYamlFile.Close(); err != nil {
		t.Fatalf("Failed to close invalid YAML file: %v", err)
	}

	_, err = LoadConfig(invalidYamlFile.Name())
	if err == nil {
		t.Error("No error occurred for an invalid YAML file")
	}
}

func TestGetRulesDir(t *testing.T) {
	// Save the current working directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "rules-test-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to the temporary directory
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temporary directory: %v", err)
	}
	// Return to the original directory after the test
	defer os.Chdir(originalDir)

	// Execute the function under test
	rulesDir, err := GetRulesDir()
	if err != nil {
		t.Fatalf("GetRulesDir function returned an error: %v", err)
	}

	// On macOS, /var/folders and /private/var/folders point to the same directory
	// So instead of comparing paths, check if the directory exists
	if _, err := os.Stat(rulesDir); os.IsNotExist(err) {
		t.Errorf("Rules directory was not created: %s", rulesDir)
	}

	// Check if the expected directory structure (.cursor/rules) was created
	cursorDir := filepath.Join(tempDir, ".cursor")
	if _, err := os.Stat(cursorDir); os.IsNotExist(err) {
		t.Errorf(".cursor directory was not created: %s", cursorDir)
	}

	expectedRulesDir := filepath.Join(cursorDir, "rules")
	if _, err := os.Stat(expectedRulesDir); os.IsNotExist(err) {
		t.Errorf("rules directory was not created: %s", expectedRulesDir)
	}
}
