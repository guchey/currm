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

// DownloadRule downloads the specified rule from the given URL and saves it to the rules directory
func DownloadRule(rule config.Rule, rulesDir string) error {
	// Extract filename from the last part of the URL
	urlParts := strings.Split(rule.URL, "/")
	fileName := urlParts[len(urlParts)-1]

	// Add .mdc extension if not already present
	if !strings.HasSuffix(fileName, ".mdc") {
		fileName = fileName + ".mdc"
	}

	// Create the full path where the file will be saved
	filePath := filepath.Join(rulesDir, fileName)

	// Create and execute HTTP request to download the rule
	resp, err := http.Get(rule.URL)
	if err != nil {
		return fmt.Errorf("failed to download rule '%s': %w", rule.Name, err)
	}
	defer resp.Body.Close()

	// Verify the HTTP response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download rule '%s': HTTP status code %d", rule.Name, resp.StatusCode)
	}

	// Create the destination file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file '%s': %w", filePath, err)
	}
	defer file.Close()

	// Copy the downloaded content to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write to file '%s': %w", filePath, err)
	}

	fmt.Printf("Downloaded rule '%s' to '%s'\n", rule.Name, filePath)
	return nil
}

// DownloadAllRules downloads all rules specified in the configuration file
// It continues downloading even if some rules fail to download
func DownloadAllRules(cfg *config.Config) error {
	// Get the directory where rules should be stored
	rulesDir, err := config.GetRulesDir()
	if err != nil {
		return err
	}

	fmt.Printf("Downloading rules to '%s'\n", rulesDir)

	// Download each rule defined in the configuration
	for _, rule := range cfg.Rules {
		if err := DownloadRule(rule, rulesDir); err != nil {
			fmt.Printf("Warning: %v\n", err)
			// Continue with the next rule even if this one failed
			continue
		}
	}

	return nil
}
