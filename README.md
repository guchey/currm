# Cursor Rules Hub

Cursor Rules Hub は、YAMLファイルで定義されたCursorルールを現在のディレクトリの `.cursor/rules` ディレクトリにダウンロードするためのツールです。

## インストール

```bash
go install github.com/guchey/cursorruleshub/cmd/cursorruleshub@latest
```

## 使い方

1. `cursorrules.yaml` ファイルを作成し、ダウンロードしたいルールを定義します：

```yaml
rules:
  - name: "ルール名1"
    url: "https://example.com/path/to/rule1.mdc"
  - name: "ルール名2"
    url: "https://example.com/path/to/rule2.mdc"
```

2. 以下のコマンドを実行してルールをダウンロードします：

```bash
cursorruleshub pull
```

別の設定ファイルを使用する場合は、`--config` または `-c` フラグを使用します：

```bash
cursorruleshub pull -c 別の設定ファイル.yaml
```

## 機能

- YAMLファイルからルール情報（名前とURL）を読み込みます
- 指定されたURLからルールファイルをダウンロードします
- ダウンロードしたファイルを現在のディレクトリの `.cursor/rules` ディレクトリに保存します
- ファイル名は自動的にURLから取得され、必要に応じて `.mdc` 拡張子が追加されます

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。