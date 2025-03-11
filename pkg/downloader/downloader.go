package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/guchey/cursorruleshub/pkg/config"
)

// DownloadRule は指定されたルールをダウンロードします
func DownloadRule(rule config.Rule, rulesDir string) error {
	// ファイル名を取得（URLの最後の部分を使用）
	urlParts := strings.Split(rule.URL, "/")
	fileName := urlParts[len(urlParts)-1]

	// 拡張子がない場合は.mdcを追加
	if !strings.HasSuffix(fileName, ".mdc") {
		fileName = fileName + ".mdc"
	}

	// 保存先のパスを作成
	filePath := filepath.Join(rulesDir, fileName)

	// HTTPリクエストを作成
	resp, err := http.Get(rule.URL)
	if err != nil {
		return fmt.Errorf("ルール '%s' のダウンロードに失敗しました: %w", rule.Name, err)
	}
	defer resp.Body.Close()

	// レスポンスコードをチェック
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ルール '%s' のダウンロードに失敗しました: HTTPステータスコード %d", rule.Name, resp.StatusCode)
	}

	// ファイルを作成
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("ファイル '%s' の作成に失敗しました: %w", filePath, err)
	}
	defer file.Close()

	// レスポンスボディをファイルに書き込み
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("ファイル '%s' への書き込みに失敗しました: %w", filePath, err)
	}

	fmt.Printf("ルール '%s' を '%s' にダウンロードしました\n", rule.Name, filePath)
	return nil
}

// DownloadAllRules は設定ファイルに記載されたすべてのルールをダウンロードします
func DownloadAllRules(cfg *config.Config) error {
	rulesDir, err := config.GetRulesDir()
	if err != nil {
		return err
	}

	fmt.Printf("ルールを '%s' にダウンロードします\n", rulesDir)

	for _, rule := range cfg.Rules {
		if err := DownloadRule(rule, rulesDir); err != nil {
			fmt.Printf("警告: %v\n", err)
			// エラーがあっても続行
			continue
		}
	}

	return nil
}
