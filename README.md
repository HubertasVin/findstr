# findstr

A simple command‑line utility to search for occurrences of a string within files under a directory, displaying matching lines with context and highlighting.

## Features

* **Recursive search** through subdirectories
* **Context lines**: shows two lines before and after each match
* **Configurable root search directory** and (soon) **exclude paths**

## Installation

Install the latest release binary via Go:

```bash
sudo ln -sf "$(go env GOPATH)/bin/findstr" /usr/local/bin/findstr
go install github.com/HubertasVin/findstr@latest
```

## Update local binary

```bash
go install github.com/HubertasVin/findstr@latest
```

## Uninstall

Remove the installed binary:
```bash
sudo rm -f /usr/local/bin/findstr
rm -f "$(go env GOPATH)/bin/findstr"
```

## Usage

Search for file content matching `<pattern>` under the specified root:
```plain
findstr [flags] <pattern>
```

### Flags

- `-r, --root` <dir> root directory (default ./)
- `-e, --exclude-dir` <paths> comma-separated relative directories to ignore
- `-x, --exclude-file` <glob> comma-separated bash-style globs to ignore; special pattern noext matches files with no extension
- `-t, --thread` <num> worker count (default 1)
- `-c, --context` <num> context lines around a matched line (default 2)
- `--json print` results as JSON and exit
- `--create-config` write default config to ~/.config/findstr.toml and exit
- `-v, --version` print version info

### Examples

Search for `TODO` in the current directory:
```bash
findstr TODO
```

Search in a specific folder:
```bash
findstr -r ./src "func main"
```

## First-time config

Generate a default config file:
```bash
findstr --create-config
```

This creates `~/.config/findstr.toml`. Edit it to change colors or layout. Example:
```toml
[layout]
align = "right"
autoWidth = true

[layout.header]
parts = ["---", " ", "{filepath}", ":"]

[layout.match]
parts = ["{ln}", " | ", "{text}"]

[layout.context]
parts = ["{ln}", " | ", "{text}"]

[theme.styles.header]
fg = "#ffffff"
bold = true

[theme.styles.match]
fg = "#ffffff"
bold = true

[theme.styles.context]
fg = "#cccccc"
bold = false

[theme.styles.highlight]
fg = "#ff0000"
bold = true
```

### Layout tokens
- {filepath} {dir} {base} {clean}
- {ln} line number
- {text} the line’s text

## Contributing

Contributions welcome! Please open issues or pull requests on [GitHub](https://github.com/HubertasVin/findstr).

## License

MIT © 2025
