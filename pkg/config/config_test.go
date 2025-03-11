package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// テスト用の一時設定ファイルを作成
	tempFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("一時ファイルの作成に失敗しました: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// テスト用の設定内容
	configContent := `rules:
  - name: "test-rule-1"
    url: "https://example.com/rule1"
  - name: "test-rule-2"
    url: "https://example.com/rule2"
`
	if _, err := tempFile.Write([]byte(configContent)); err != nil {
		t.Fatalf("一時ファイルへの書き込みに失敗しました: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("一時ファイルのクローズに失敗しました: %v", err)
	}

	// テスト対象の関数を実行
	cfg, err := LoadConfig(tempFile.Name())
	if err != nil {
		t.Fatalf("LoadConfig関数がエラーを返しました: %v", err)
	}

	// 結果の検証
	if len(cfg.Rules) != 2 {
		t.Errorf("期待されるルール数: 2, 実際: %d", len(cfg.Rules))
	}

	// 各ルールの内容を検証
	expectedRules := []Rule{
		{Name: "test-rule-1", URL: "https://example.com/rule1"},
		{Name: "test-rule-2", URL: "https://example.com/rule2"},
	}

	for i, expected := range expectedRules {
		if cfg.Rules[i].Name != expected.Name {
			t.Errorf("ルール %d の名前が一致しません。期待: %s, 実際: %s", i, expected.Name, cfg.Rules[i].Name)
		}
		if cfg.Rules[i].URL != expected.URL {
			t.Errorf("ルール %d のURLが一致しません。期待: %s, 実際: %s", i, expected.URL, cfg.Rules[i].URL)
		}
	}

	// 存在しないファイルのテスト
	_, err = LoadConfig("non-existent-file.yaml")
	if err == nil {
		t.Error("存在しないファイルに対してエラーが発生しませんでした")
	}

	// 不正なYAMLのテスト
	invalidYamlFile, err := os.CreateTemp("", "invalid-*.yaml")
	if err != nil {
		t.Fatalf("不正なYAMLファイルの作成に失敗しました: %v", err)
	}
	defer os.Remove(invalidYamlFile.Name())

	if _, err := invalidYamlFile.Write([]byte("invalid: yaml: content:")); err != nil {
		t.Fatalf("不正なYAMLファイルへの書き込みに失敗しました: %v", err)
	}
	if err := invalidYamlFile.Close(); err != nil {
		t.Fatalf("不正なYAMLファイルのクローズに失敗しました: %v", err)
	}

	_, err = LoadConfig(invalidYamlFile.Name())
	if err == nil {
		t.Error("不正なYAMLファイルに対してエラーが発生しませんでした")
	}
}

func TestGetRulesDir(t *testing.T) {
	// 現在の作業ディレクトリを保存
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("現在のディレクトリの取得に失敗しました: %v", err)
	}

	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "rules-test-*")
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

	// テスト対象の関数を実行
	rulesDir, err := GetRulesDir()
	if err != nil {
		t.Fatalf("GetRulesDir関数がエラーを返しました: %v", err)
	}

	// MacOSでは/var/foldersと/private/var/foldersが同じディレクトリを指すため
	// パスの比較ではなく、ディレクトリが存在するかどうかを確認する
	if _, err := os.Stat(rulesDir); os.IsNotExist(err) {
		t.Errorf("ルールディレクトリが作成されていません: %s", rulesDir)
	}

	// 期待されるディレクトリ構造（.cursor/rules）が作成されているか確認
	cursorDir := filepath.Join(tempDir, ".cursor")
	if _, err := os.Stat(cursorDir); os.IsNotExist(err) {
		t.Errorf(".cursorディレクトリが作成されていません: %s", cursorDir)
	}

	expectedRulesDir := filepath.Join(cursorDir, "rules")
	if _, err := os.Stat(expectedRulesDir); os.IsNotExist(err) {
		t.Errorf("rulesディレクトリが作成されていません: %s", expectedRulesDir)
	}
}
