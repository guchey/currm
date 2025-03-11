package downloader

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/guchey/currm/pkg/config"
)

func TestDownloadRule(t *testing.T) {
	// Create HTTP test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return different responses based on path
		switch r.URL.Path {
		case "/success":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Test rule content"))
		case "/not-found":
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "downloader-test-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test successful case
	successRule := config.Rule{
		Name: "Success Rule",
		URL:  server.URL + "/success",
	}

	err = DownloadRule(successRule, tempDir)
	if err != nil {
		t.Errorf("Error occurred in success case: %v", err)
	}

	// Check if file was created
	expectedFilePath := filepath.Join(tempDir, successRule.Name+".mdc")
	if _, err := os.Stat(expectedFilePath); os.IsNotExist(err) {
		t.Errorf("File was not created: %s", expectedFilePath)
	}

	// Check file content
	content, err := os.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if string(content) != "Test rule content" {
		t.Errorf("File content differs from expected. Expected: %s, Actual: %s", "Test rule content", string(content))
	}

	// Test 404 error
	notFoundRule := config.Rule{
		Name: "Non-existent Rule",
		URL:  server.URL + "/not-found",
	}

	err = DownloadRule(notFoundRule, tempDir)
	if err == nil {
		t.Error("No error was returned for 404 error case")
	}

	// Test invalid URL
	invalidRule := config.Rule{
		Name: "Invalid URL",
		URL:  "http://invalid-url-that-does-not-exist.example",
	}

	err = DownloadRule(invalidRule, tempDir)
	if err == nil {
		t.Error("No error was returned for invalid URL case")
	}
}

func TestDownloadAllRules(t *testing.T) {
	// Create HTTP test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return different responses based on path
		switch r.URL.Path {
		case "/rule1":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Rule 1 content"))
		case "/rule2":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Rule 2 content"))
		case "/error":
			w.WriteHeader(http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Save current working directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "all-rules-test-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to temporary directory
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temporary directory: %v", err)
	}
	// Return to original directory after test
	defer os.Chdir(originalDir)

	// Create test configuration
	cfg := &config.Config{
		Rules: []config.Rule{
			{Name: "Rule 1", URL: server.URL + "/rule1"},
			{Name: "Rule 2", URL: server.URL + "/rule2"},
			{Name: "Error Rule", URL: server.URL + "/error"},
		},
	}

	// Execute the function under test
	err = DownloadAllRules(cfg)
	if err != nil {
		t.Fatalf("DownloadAllRules function returned an error: %v", err)
	}

	// Get rules directory path
	rulesDir, err := config.GetRulesDir()
	if err != nil {
		t.Fatalf("Failed to get rules directory: %v", err)
	}

	// Check if files for successful rules were created
	rule1Path := filepath.Join(rulesDir, "Rule 1.mdc")
	if _, err := os.Stat(rule1Path); os.IsNotExist(err) {
		t.Errorf("Rule 1 file was not created: %s", rule1Path)
	}

	rule2Path := filepath.Join(rulesDir, "Rule 2.mdc")
	if _, err := os.Stat(rule2Path); os.IsNotExist(err) {
		t.Errorf("Rule 2 file was not created: %s", rule2Path)
	}

	// Check file contents
	content1, err := os.ReadFile(rule1Path)
	if err != nil {
		t.Fatalf("Failed to read Rule 1 file: %v", err)
	}
	if string(content1) != "Rule 1 content" {
		t.Errorf("Rule 1 file content differs from expected. Expected: %s, Actual: %s", "Rule 1 content", string(content1))
	}

	content2, err := os.ReadFile(rule2Path)
	if err != nil {
		t.Fatalf("Failed to read Rule 2 file: %v", err)
	}
	if string(content2) != "Rule 2 content" {
		t.Errorf("Rule 2 file content differs from expected. Expected: %s, Actual: %s", "Rule 2 content", string(content2))
	}
}

func TestDownloadRuleWithCursorRules(t *testing.T) {
	// Create HTTP test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return different responses based on path
		switch r.URL.Path {
		case "/test.cursorrules":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("# Test Cursor Rules\n\nThis is a test rule content for .cursorrules format."))
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "cursorrules-test-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test cases
	testCases := []struct {
		name         string
		rule         config.Rule
		expectedYAML string
	}{
		{
			name: "Default values",
			rule: config.Rule{
				Name: "Cursor Rules Test",
				URL:  server.URL + "/test.cursorrules",
			},
			expectedYAML: "---\ndescription: Cursor Rules Test\nglobs: *\nalwaysApply: false\n---\n\n",
		},
		{
			name: "Custom description",
			rule: config.Rule{
				Name:        "Cursor Rules Test",
				URL:         server.URL + "/test.cursorrules",
				Description: "Custom description for test",
			},
			expectedYAML: "---\ndescription: Custom description for test\nglobs: *\nalwaysApply: false\n---\n\n",
		},
		{
			name: "Custom globs",
			rule: config.Rule{
				Name:  "Cursor Rules Test",
				URL:   server.URL + "/test.cursorrules",
				Globs: "*.go,*.md",
			},
			expectedYAML: "---\ndescription: Cursor Rules Test\nglobs: *.go,*.md\nalwaysApply: false\n---\n\n",
		},
		{
			name: "AlwaysApply true",
			rule: config.Rule{
				Name:        "Cursor Rules Test",
				URL:         server.URL + "/test.cursorrules",
				AlwaysApply: true,
			},
			expectedYAML: "---\ndescription: Cursor Rules Test\nglobs: *\nalwaysApply: true\n---\n\n",
		},
		{
			name: "All custom values",
			rule: config.Rule{
				Name:        "Cursor Rules Test",
				URL:         server.URL + "/test.cursorrules",
				Description: "Full custom description",
				Globs:       "src/**/*.ts",
				AlwaysApply: true,
			},
			expectedYAML: "---\ndescription: Full custom description\nglobs: src/**/*.ts\nalwaysApply: true\n---\n\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Download the rule
			err = DownloadRule(tc.rule, tempDir)
			if err != nil {
				t.Errorf("Error occurred when downloading .cursorrules file: %v", err)
				return
			}

			// Check if file was created with .mdc extension
			expectedFilePath := filepath.Join(tempDir, tc.rule.Name+".mdc")
			if _, err := os.Stat(expectedFilePath); os.IsNotExist(err) {
				t.Errorf("File was not created: %s", expectedFilePath)
				return
			}

			// Check file content
			content, err := os.ReadFile(expectedFilePath)
			if err != nil {
				t.Fatalf("Failed to read file: %v", err)
			}

			// Check if the content has the expected YAML front matter
			if !strings.HasPrefix(string(content), tc.expectedYAML) {
				t.Errorf("File content does not have expected YAML front matter.\nExpected prefix:\n%s\n\nActual content:\n%s",
					tc.expectedYAML, string(content))
			}

			// Check if the content has the expected rule content
			expectedContent := "# Test Cursor Rules\n\nThis is a test rule content for .cursorrules format."
			if !strings.Contains(string(content), expectedContent) {
				t.Errorf("File content does not contain expected rule content.\nExpected content:\n%s\n\nActual content:\n%s",
					expectedContent, string(content))
			}
		})
	}
}
