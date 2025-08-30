# findstr

A simple command‑line utility to search for occurrences of a regex pattern within files under a directory, displaying matching lines with context and highlighting.

## Features

* **Recursive search** through subdirectories
* **Context lines**: shows two lines before and after each match
* **Configurable root search directory** and (soon) **exclude paths**

## Installation

Install the latest release binary via Go:

```bash
go install github.com/HubertasVin/findstr
sudo mv "$(go env GOPATH)/bin/findstr" /usr/local/bin/
```

Ensure `$GOBIN` (or `$GOPATH/bin`) is on your `PATH`:

```bash
export PATH="$(go env GOBIN):$PATH"
```

## Usage

```plain
findstr [flags] <pattern>
```

Search for file content matching `<pattern>` under the specified root.

### Flags

* `-r, --root <dir>`
  Root directory to search (default `./`)
* `-e, --exclude-dir <paths>`
  Comma-separated relative paths to ignore
* `-x, --exclude-file <glob>`
  Comma-separated bash-style glob patterns of files to ignore.
  Pattern `noext` can be used for files with no extension.
* `-t, --thread-count <num>`
  Thread count to use for file parsing. (default `1`)
* `-c, --context-size <num>` 
  Number of context lines to show around a matched line. (default `2`)
* `--style <json>` 
  Custom style in valid json format for highlighting.
  Available keys:
  - `matchFg`: `<hex>`,
  - `matchBg`: `<hex>`,
  - `matchBold`: `<bool>`.
* `--json`
  Print result in json format.
### Examples

Search for `TODO` in the current directory:

```bash
findstr TODO
```

Search in a specific folder:

```bash
findstr -r ./src "func main"
```

## Uninstall

Remove the installed binary:

```bash
sudo rm /usr/local/bin/findstr
```

## Contributing

Contributions welcome! Please open issues or pull requests on [GitHub](https://github.com/HubertasVin/findstr).

## License

MIT © 2025
