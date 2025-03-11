package main

import (
	"fmt"
	"os"

	"github.com/guchey/cursorruleshub/pkg/config"
	"github.com/guchey/cursorruleshub/pkg/downloader"
	"github.com/spf13/cobra"
)

var (
	configFile string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "cursorruleshub",
		Short: "Cursor Rules Hub - Cursorルールをダウンロードするツール",
		Long: `Cursor Rules Hub はYAMLファイルで定義されたCursorルールを
.cursor/rules ディレクトリにダウンロードするツールです。`,
	}

	var pullCmd = &cobra.Command{
		Use:   "pull",
		Short: "設定ファイルに記載されたルールをダウンロードします",
		RunE: func(cmd *cobra.Command, args []string) error {
			// 設定ファイルを読み込む
			cfg, err := config.LoadConfig(configFile)
			if err != nil {
				return err
			}

			// すべてのルールをダウンロード
			if err := downloader.DownloadAllRules(cfg); err != nil {
				return err
			}

			return nil
		},
	}

	// フラグを設定
	pullCmd.Flags().StringVarP(&configFile, "config", "c", "cursorrules.yaml", "設定ファイルのパス")

	// コマンドを追加
	rootCmd.AddCommand(pullCmd)

	// コマンドを実行
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
