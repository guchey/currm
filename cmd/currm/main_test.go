package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/guchey/currm/pkg/config"
	"github.com/spf13/cobra"
)

func TestRootCommand(t *testing.T) {
	// Prepare to capture standard output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute the command
	os.Args = []string{"currm"}
	main()

	// Restore standard output and get captured output
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check if help message is included
	if output == "" {
		t.Error("Root command output is empty")
	}
}

func TestPullCommand(t *testing.T) {
	// Save the current working directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "cmd-test-*")
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

	// Create a test configuration file
	configContent := `rules:
  - name: "test-rule"
    url: "https://example.com/rule"
`
	configPath := filepath.Join(tempDir, "test-config.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create configuration file: %v", err)
	}

	// Prepare to capture standard output and standard error
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = wOut
	os.Stderr = wErr

	// Execute the pull command (should fail because the URL doesn't exist)
	os.Args = []string{"currm", "pull", "--config", configPath}
	main()

	// Restore standard output and standard error, and get captured output
	wOut.Close()
	wErr.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	var bufOut, bufErr bytes.Buffer
	io.Copy(&bufOut, rOut)
	io.Copy(&bufErr, rErr)

	stdoutOutput := bufOut.String()
	stderrOutput := bufErr.String()

	// Check the output
	// Note: In a real test, it would be better to use a mock server
	// or specify a URL that actually exists in the configuration file
	if stdoutOutput == "" && stderrOutput == "" {
		t.Error("Pull command output is empty")
	}
}

func TestInvalidConfigFile(t *testing.T) {
	// This test cannot directly call main() because os.Exit(1) would terminate the test process.
	// Instead, we simulate the command execution by checking the command behavior.

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "invalid-config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Path to a non-existent configuration file
	nonExistentConfig := filepath.Join(tempDir, "non-existent-config.yaml")

	// Prepare to capture standard error
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Create rootCmd directly and execute it
	var rootCmd = &cobra.Command{
		Use:   "currm",
		Short: "Currm - A tool for downloading Cursor rules",
		Long: `Currm is a tool for downloading Cursor rules defined in YAML files 
to the .cursor/rules directory in your current directory.`,
	}

	var pullCmd = &cobra.Command{
		Use:   "pull",
		Short: "Download rules specified in the configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load configuration file
			_, err := config.LoadConfig(nonExistentConfig)
			if err != nil {
				return err
			}
			return nil
		},
	}

	// Set flags
	pullCmd.Flags().StringVarP(&configFile, "config", "c", nonExistentConfig, "Path to configuration file")

	// Add commands
	rootCmd.AddCommand(pullCmd)

	// Execute command
	rootCmd.SetArgs([]string{"pull"})
	err = rootCmd.Execute()

	// Restore standard error and get captured output
	w.Close()
	os.Stderr = oldStderr
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Check that an error occurred
	if err == nil {
		t.Error("No error occurred for a non-existent configuration file")
	}

	// Check if the error message includes content related to the "configuration file"
	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "configuration file") && !strings.Contains(errorMsg, "config") {
		t.Errorf("Error message does not include content related to the configuration file: %s", errorMsg)
	}
}
