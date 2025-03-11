package downloader

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/guchey/currm/pkg/config"
)

func TestDownloadRule(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// パスに基づいて異なるレスポンスを返す
		switch r.URL.Path {
		case "/success":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("テストルールの内容"))
		case "/not-found":
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "downloader-test-*")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 成功ケースのテスト
	successRule := config.Rule{
		Name: "成功ルール",
		URL:  server.URL + "/success",
	}

	err = DownloadRule(successRule, tempDir)
	if err != nil {
		t.Errorf("成功ケースでエラーが発生しました: %v", err)
	}

	// ファイルが作成されたか確認
	expectedFilePath := filepath.Join(tempDir, "success.mdc")
	if _, err := os.Stat(expectedFilePath); os.IsNotExist(err) {
		t.Errorf("ファイルが作成されていません: %s", expectedFilePath)
	}

	// ファイルの内容を確認
	content, err := os.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("ファイルの読み込みに失敗しました: %v", err)
	}
	if string(content) != "テストルールの内容" {
		t.Errorf("ファイルの内容が期待と異なります。期待: %s, 実際: %s", "テストルールの内容", string(content))
	}

	// 404エラーのテスト
	notFoundRule := config.Rule{
		Name: "存在しないルール",
		URL:  server.URL + "/not-found",
	}

	err = DownloadRule(notFoundRule, tempDir)
	if err == nil {
		t.Error("404エラーの場合にエラーが返されませんでした")
	}

	// 無効なURLのテスト
	invalidRule := config.Rule{
		Name: "無効なURL",
		URL:  "http://invalid-url-that-does-not-exist.example",
	}

	err = DownloadRule(invalidRule, tempDir)
	if err == nil {
		t.Error("無効なURLの場合にエラーが返されませんでした")
	}
}

func TestDownloadAllRules(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// パスに基づいて異なるレスポンスを返す
		switch r.URL.Path {
		case "/rule1":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ルール1の内容"))
		case "/rule2":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ルール2の内容"))
		case "/error":
			w.WriteHeader(http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// 現在の作業ディレクトリを保存
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("現在のディレクトリの取得に失敗しました: %v", err)
	}

	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "all-rules-test-*")
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

	// テスト用の設定を作成
	cfg := &config.Config{
		Rules: []config.Rule{
			{Name: "ルール1", URL: server.URL + "/rule1"},
			{Name: "ルール2", URL: server.URL + "/rule2"},
			{Name: "エラールール", URL: server.URL + "/error"},
		},
	}

	// テスト対象の関数を実行
	err = DownloadAllRules(cfg)
	if err != nil {
		t.Fatalf("DownloadAllRules関数がエラーを返しました: %v", err)
	}

	// ルールディレクトリのパスを取得
	rulesDir, err := config.GetRulesDir()
	if err != nil {
		t.Fatalf("ルールディレクトリの取得に失敗しました: %v", err)
	}

	// 成功したルールのファイルが作成されたか確認
	rule1Path := filepath.Join(rulesDir, "rule1.mdc")
	if _, err := os.Stat(rule1Path); os.IsNotExist(err) {
		t.Errorf("ルール1のファイルが作成されていません: %s", rule1Path)
	}

	rule2Path := filepath.Join(rulesDir, "rule2.mdc")
	if _, err := os.Stat(rule2Path); os.IsNotExist(err) {
		t.Errorf("ルール2のファイルが作成されていません: %s", rule2Path)
	}

	// ファイルの内容を確認
	content1, err := os.ReadFile(rule1Path)
	if err != nil {
		t.Fatalf("ルール1のファイルの読み込みに失敗しました: %v", err)
	}
	if string(content1) != "ルール1の内容" {
		t.Errorf("ルール1のファイルの内容が期待と異なります。期待: %s, 実際: %s", "ルール1の内容", string(content1))
	}

	content2, err := os.ReadFile(rule2Path)
	if err != nil {
		t.Fatalf("ルール2のファイルの読み込みに失敗しました: %v", err)
	}
	if string(content2) != "ルール2の内容" {
		t.Errorf("ルール2のファイルの内容が期待と異なります。期待: %s, 実際: %s", "ルール2の内容", string(content2))
	}
}
