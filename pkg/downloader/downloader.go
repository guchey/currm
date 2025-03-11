package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/guchey/currm/pkg/config"
)

// getURLWithRevision returns the URL with the revision if specified
func getURLWithRevision(rule config.Rule) string {
	// If revision is not specified or is "latest", use the original URL
	if rule.Revision == "" || rule.Revision == "latest" {
		return rule.URL
	}

	// For GitHub URLs, apply the revision
	if strings.Contains(rule.URL, "github.com") {
		// For raw.githubusercontent.com URLs
		if strings.Contains(rule.URL, "raw.githubusercontent.com") {
			// URL format: https://raw.githubusercontent.com/owner/repo/branch/path/to/file
			parts := strings.Split(rule.URL, "/")
			if len(parts) >= 6 {
				// Replace the branch part with the revision
				parts[5] = rule.Revision
				return strings.Join(parts, "/")
			}
		}

		// For github.com URLs
		if strings.Contains(rule.URL, "github.com") && !strings.Contains(rule.URL, "raw.githubusercontent.com") {
			// URL format: https://github.com/owner/repo/blob/branch/path/to/file
			parts := strings.Split(rule.URL, "/")
			if len(parts) >= 7 && parts[5] == "blob" {
				// Replace the branch part with the revision
				parts[6] = rule.Revision
				return strings.Join(parts, "/")
			}
		}
	}

	// For other URLs, use the original URL
	return rule.URL
}

// getShortRevision returns a shortened version of the revision if it's a commit hash
func getShortRevision(revision string) string {
	// If it looks like a commit hash (40 hexadecimal characters), shorten it
	if len(revision) >= 40 && isHexString(revision) {
		return revision[:8] // Use the first 8 characters
	}
	return revision
}

// isHexString checks if a string consists only of hexadecimal characters
func isHexString(s string) bool {
	for _, r := range s {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
			return false
		}
	}
	return true
}

// DownloadRule downloads the specified rule from the given URL and saves it to the rules directory
func DownloadRule(rule config.Rule, rulesDir string) error {
	// Get URL with revision consideration
	url := getURLWithRevision(rule)

	// Use rule name as the base filename instead of extracting from URL
	fileName := rule.Name + ".mdc"

	// If revision is specified, add it to the filename
	if rule.Revision != "" && rule.Revision != "latest" {
		fileExt := filepath.Ext(fileName)
		fileBase := strings.TrimSuffix(fileName, fileExt)
		shortRev := getShortRevision(rule.Revision)
		fileName = fmt.Sprintf("%s-%s%s", fileBase, shortRev, fileExt)
	}

	// Create the full path where the file will be saved
	filePath := filepath.Join(rulesDir, fileName)

	// Create and execute HTTP request to download the rule
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download rule '%s': %w", rule.Name, err)
	}
	defer resp.Body.Close()

	// Verify the HTTP response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download rule '%s': HTTP status code %d", rule.Name, resp.StatusCode)
	}

	// Read the content from the response
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read content for rule '%s': %w", rule.Name, err)
	}

	// If the URL ends with .cursorrules, convert it to .mdc format
	isCursorRules := strings.HasSuffix(url, ".cursorrules")
	if isCursorRules {
		// Get description from rule or use name as fallback
		description := rule.Description
		if description == "" {
			description = rule.Name
		}

		// Get globs from rule or use "*" as default
		globs := rule.Globs
		if globs == "" {
			globs = "*"
		}

		// Get alwaysApply from rule
		alwaysApply := rule.AlwaysApply

		// Create the .mdc format with YAML front matter
		mdcContent := fmt.Sprintf("---\ndescription: %s\nglobs: %s\nalwaysApply: %t\n---\n\n%s",
			description,
			globs,
			alwaysApply,
			string(content))
		content = []byte(mdcContent)
	}

	// Create the destination file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file '%s': %w", filePath, err)
	}
	defer file.Close()

	// Write the content to the file
	_, err = file.Write(content)
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

// RuleStatus represents the status of a rule
type RuleStatus struct {
	Name           string
	LocalPath      string
	HasLocalFile   bool
	NeedsUpdate    bool
	LastModified   time.Time
	RemoteModified time.Time
	Revision       string
}

// CheckRuleUpdates checks if any rules need to be updated
func CheckRuleUpdates(cfg *config.Config) ([]RuleStatus, error) {
	// Get the directory where rules are stored
	rulesDir, err := config.GetRulesDir()
	if err != nil {
		return nil, err
	}

	var statuses []RuleStatus

	// Check each rule defined in the configuration
	for _, rule := range cfg.Rules {
		// Get URL with revision consideration
		url := getURLWithRevision(rule)

		// Use rule name as the base filename
		fileName := rule.Name + ".mdc"

		// If revision is specified, add it to the filename
		if rule.Revision != "" && rule.Revision != "latest" {
			fileExt := filepath.Ext(fileName)
			fileBase := strings.TrimSuffix(fileName, fileExt)
			shortRev := getShortRevision(rule.Revision)
			fileName = fmt.Sprintf("%s-%s%s", fileBase, shortRev, fileExt)
		}

		// Create the full path where the file should be
		filePath := filepath.Join(rulesDir, fileName)

		status := RuleStatus{
			Name:      rule.Name,
			LocalPath: filePath,
			Revision:  rule.Revision,
		}

		// Check if the file exists locally
		fileInfo, err := os.Stat(filePath)
		if err == nil {
			status.HasLocalFile = true
			status.LastModified = fileInfo.ModTime()
		} else if os.IsNotExist(err) {
			status.HasLocalFile = false
		} else {
			return nil, fmt.Errorf("failed to check file '%s': %w", filePath, err)
		}

		// If a specific revision is specified and the file exists, no update is needed
		if rule.Revision != "" && rule.Revision != "latest" && status.HasLocalFile {
			status.NeedsUpdate = false
			statuses = append(statuses, status)
			continue
		}

		// Check remote file
		req, err := http.NewRequest("HEAD", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request for rule '%s': %w", rule.Name, err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to check rule '%s': %w", rule.Name, err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to check rule '%s': HTTP status code %d", rule.Name, resp.StatusCode)
		}

		// Get last modified time from header if available
		lastModHeader := resp.Header.Get("Last-Modified")
		if lastModHeader != "" {
			remoteTime, err := time.Parse(time.RFC1123, lastModHeader)
			if err == nil {
				status.RemoteModified = remoteTime
				// Check if remote file is newer than local file
				if !status.HasLocalFile || remoteTime.After(status.LastModified) {
					status.NeedsUpdate = true
				}
			}
		} else {
			// If no Last-Modified header, assume update is needed if file doesn't exist
			status.NeedsUpdate = !status.HasLocalFile
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}
