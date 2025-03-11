package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Rule はCursorルールの情報を表す構造体です
type Rule struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

// Config は設定ファイルの構造を表す構造体です
type Config struct {
	Rules []Rule `yaml:"rules"`
}

// LoadConfig は指定されたパスから設定ファイルを読み込みます
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("設定ファイルの読み込みに失敗しました: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("YAMLのパースに失敗しました: %w", err)
	}

	return &config, nil
}

// GetRulesDir はルールファイルを保存するディレクトリのパスを返します
func GetRulesDir() (string, error) {
	// 現在のカレントディレクトリを取得
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("カレントディレクトリの取得に失敗しました: %w", err)
	}

	// カレントディレクトリを基準として .cursor/rules ディレクトリのパスを作成
	rulesDir := filepath.Join(currentDir, ".cursor", "rules")
	if err := os.MkdirAll(rulesDir, 0755); err != nil {
		return "", fmt.Errorf("ルールディレクトリの作成に失敗しました: %w", err)
	}

	return rulesDir, nil
}
