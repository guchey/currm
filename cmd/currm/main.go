package main

import (
	"fmt"
	"os"

	"github.com/guchey/currm/pkg/config"
	"github.com/guchey/currm/pkg/downloader"
	"github.com/spf13/cobra"
)

var (
	configFile string
)

// isLikelyCommitHash checks if a string looks like a commit hash
func isLikelyCommitHash(s string) bool {
	// Consider 40-character hexadecimal as SHA hash
	if len(s) >= 40 {
		for _, r := range s {
			if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
				return false
			}
		}
		return true
	}
	return false
}

// formatRevision formats the revision for display
func formatRevision(revision string) string {
	if revision == "latest" {
		return "Latest version"
	} else if isLikelyCommitHash(revision) {
		// For commit hashes, display shortened version
		return fmt.Sprintf("Commit: %s", revision[:8])
	}
	return revision
}

func main() {
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
			cfg, err := config.LoadConfig(configFile)
			if err != nil {
				return err
			}

			// Download all rules
			if err := downloader.DownloadAllRules(cfg); err != nil {
				return err
			}

			return nil
		},
	}

	var checkCmd = &cobra.Command{
		Use:   "check",
		Short: "Check for updates to rules specified in the configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load configuration file
			cfg, err := config.LoadConfig(configFile)
			if err != nil {
				return err
			}

			// Check for updates
			statuses, err := downloader.CheckRuleUpdates(cfg)
			if err != nil {
				return err
			}

			// Display results
			updatesAvailable := false
			fmt.Println("Checking for updates...")

			for _, status := range statuses {
				revInfo := ""
				if status.Revision != "" {
					revInfo = fmt.Sprintf(" (%s)", formatRevision(status.Revision))
				}

				if !status.HasLocalFile {
					fmt.Printf("- %s%s: Rule is not installed\n", status.Name, revInfo)
					updatesAvailable = true
				} else if status.NeedsUpdate {
					fmt.Printf("- %s%s: Update available\n", status.Name, revInfo)
					updatesAvailable = true
				} else {
					fmt.Printf("- %s%s: Up to date\n", status.Name, revInfo)
				}
			}

			if updatesAvailable {
				fmt.Println("\nRun 'currm pull' to install updates")
			} else {
				fmt.Println("\nAll rules are up to date")
			}

			return nil
		},
	}

	// Set flags
	pullCmd.Flags().StringVarP(&configFile, "config", "c", "currm.yaml", "Path to configuration file")
	checkCmd.Flags().StringVarP(&configFile, "config", "c", "currm.yaml", "Path to configuration file")

	// Add commands
	rootCmd.AddCommand(pullCmd)
	rootCmd.AddCommand(checkCmd)

	// Execute command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
