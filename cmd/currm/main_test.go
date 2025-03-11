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
	// 標準出力をキャプチャするための準備
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// コマンドを実行
	os.Args = []string{"currm"}
	main()

	// 標準出力の復元とキャプチャした出力の取得
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// ヘルプメッセージが含まれているか確認
	if output == "" {
		t.Error("ルートコマンドの出力が空です")
	}
}

func TestPullCommand(t *testing.T) {
	// 現在の作業ディレクトリを保存
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("現在のディレクトリの取得に失敗しました: %v", err)
	}

	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "cmd-test-*")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 一時ディレクトリに移動
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("一時ディレクトリへの移動に失敗しました: %v", err)
	}
	// テスト終了後に元のディレクトリに戻る
	defer os.Chdir(originalDir)

	// テスト用の設定ファイルを作成
	configContent := `rules:
  - name: "test-rule"
    url: "https://example.com/rule"
`
	configPath := filepath.Join(tempDir, "test-config.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("設定ファイルの作成に失敗しました: %v", err)
	}

	// 標準出力と標準エラー出力をキャプチャするための準備
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = wOut
	os.Stderr = wErr

	// pullコマンドを実行（存在しないURLなのでエラーになるはず）
	os.Args = []string{"currm", "pull", "--config", configPath}
	main()

	// 標準出力と標準エラー出力の復元とキャプチャした出力の取得
	wOut.Close()
	wErr.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	var bufOut, bufErr bytes.Buffer
	io.Copy(&bufOut, rOut)
	io.Copy(&bufErr, rErr)

	stdoutOutput := bufOut.String()
	stderrOutput := bufErr.String()

	// 出力の確認
	// 注意: 実際のテストでは、モックサーバーを使用するか、
	// または設定ファイルに実際に存在するURLを指定する方が良いでしょう
	if stdoutOutput == "" && stderrOutput == "" {
		t.Error("pullコマンドの出力が空です")
	}
}

func TestInvalidConfigFile(t *testing.T) {
	// このテストはmain()関数を直接呼び出すと、os.Exit(1)によりテストプロセスが終了してしまうため、
	// 実際のコマンド実行をシミュレートする代わりに、コマンドの動作を確認します

	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "invalid-config-test-*")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 存在しない設定ファイルのパス
	nonExistentConfig := filepath.Join(tempDir, "non-existent-config.yaml")

	// 標準エラー出力をキャプチャするための準備
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// rootCmdを直接作成して実行
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

	// 標準エラー出力の復元とキャプチャした出力の取得
	w.Close()
	os.Stderr = oldStderr
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// エラーが発生したことを確認
	if err == nil {
		t.Error("存在しない設定ファイルに対してエラーが発生しませんでした")
	}

	// エラーメッセージに「設定ファイル」に関する内容が含まれているか確認
	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "configuration file") && !strings.Contains(errorMsg, "config") {
		t.Errorf("エラーメッセージに設定ファイルに関する内容が含まれていません: %s", errorMsg)
	}
}
