package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/guchey/currm/pkg/config"
)

// DownloadRule downloads the specified rule
func DownloadRule(rule config.Rule, rulesDir string) error {
	// Get filename (use the last part of the URL)
	urlParts := strings.Split(rule.URL, "/")
	fileName := urlParts[len(urlParts)-1]

	// Add .mdc extension if not present
	if !strings.HasSuffix(fileName, ".mdc") {
		fileName = fileName + ".mdc"
	}

	// Create path to save the file
	filePath := filepath.Join(rulesDir, fileName)

	// Create HTTP request
	resp, err := http.Get(rule.URL)
	if err != nil {
		return fmt.Errorf("failed to download rule '%s': %w", rule.Name, err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download rule '%s': HTTP status code %d", rule.Name, resp.StatusCode)
	}

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file '%s': %w", filePath, err)
	}
	defer file.Close()

	// Write response body to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write to file '%s': %w", filePath, err)
	}

	fmt.Printf("Downloaded rule '%s' to '%s'\n", rule.Name, filePath)
	return nil
}

// DownloadAllRules downloads all rules specified in the configuration file
func DownloadAllRules(cfg *config.Config) error {
	rulesDir, err := config.GetRulesDir()
	if err != nil {
		return err
	}

	fmt.Printf("Downloading rules to '%s'\n", rulesDir)

	for _, rule := range cfg.Rules {
		if err := DownloadRule(rule, rulesDir); err != nil {
			fmt.Printf("Warning: %v\n", err)
			// Continue even if there are errors
			continue
		}
	}

	return nil
}
