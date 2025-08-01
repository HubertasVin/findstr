# findstr

A simple command‑line utility to search for occurrences of a regex pattern within files under a directory, displaying matching lines with context and highlighting.

## Features

* **Recursive search** through subdirectories
* **Context lines**: shows two lines before and after each match
* **Configurable root search directory** and (soon) **exclude paths**

## Installation

Install the latest release binary via Go:

```bash
go install github.com/HubertasVin/findstr@latest
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
* `-e, --exclude-dir <paths>` *(coming soon)*
  Comma‑separated relative paths to ignore

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
eval "rm \"$(go env GOBIN)/findstr\" 2>/dev/null || rm \"$(go env GOPATH)/bin/findstr\""
```

## Contributing

Contributions welcome! Please open issues or pull requests on [GitHub](https://github.com/HubertasVin/findstr).

## License

MIT © 2025
