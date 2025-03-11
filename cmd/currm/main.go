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

	// Set flags
	pullCmd.Flags().StringVarP(&configFile, "config", "c", "currm.yaml", "Path to configuration file")

	// Add commands
	rootCmd.AddCommand(pullCmd)

	// Execute command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
