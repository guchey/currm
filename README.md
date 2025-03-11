# Currm

Currm is a tool for downloading Cursor rules defined in YAML files to the `.cursor/rules` directory in your current directory.

## Installation

```bash
go install github.com/guchey/currm/cmd/currm@latest
```

## Usage

1. Create a `currm.yaml` file and define the rules you want to download:

```yaml
rules:
  - name: go
    url: "https://example.com/path/to/go.mdc"
    revision: "latest"
  - name: language
    url: "https://example.com/path/to/language.mdc"
    revision: "latest"
  - name: dry-solid-principles
    url: "https://example.com/path/to/dry-solid-principles.cursorrules"
    revision: "latest"
    description: "DRY and SOLID principles for code organization"
    globs: "*.go,*.js,*.ts"
    alwaysApply: false
```

2. Run the following command to download the rules:

```bash
currm pull
```

To use a different configuration file, use the `--config` or `-c` flag:

```bash
currm pull -c another-config-file.yaml
```

## Features

- Loads rule information (name, URL, revision, description, globs, alwaysApply) from a YAML file
- Downloads rule files from specified URLs
- Saves downloaded files to the `.cursor/rules` directory in your current directory
- Filenames are generated from the rule's `name` field with the `.mdc` extension
- Automatically converts `.cursorrules` format to `.mdc` format with YAML front matter
- Supports specifying a specific revision (e.g., commit hash) for GitHub URLs
- Checks for updates to rules with the `check` command

## License

This project is released under the MIT License.