# Currm (Cursor Rule Manager)

Currm is a tool for downloading Cursor rules defined in YAML files to the `.cursor/rules` directory in your current directory.

## Installation

```bash
go install github.com/guchey/currm/cmd/currm@latest
```

## Usage

1. Create a `currm.yaml` file and define the rules you want to download:

```yaml
rules:
  - name: "Rule Name 1"
    url: "https://example.com/path/to/rule1.mdc"
  - name: "Rule Name 2"
    url: "https://example.com/path/to/rule2.mdc"
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

- Loads rule information (name and URL) from a YAML file
- Downloads rule files from specified URLs
- Saves downloaded files to the `.cursor/rules` directory in your current directory
- Filenames are automatically retrieved from URLs and the `.mdc` extension is added if necessary

## License

This project is released under the MIT License.
