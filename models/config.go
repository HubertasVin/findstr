package models

type PartsJSON struct {
	Parts []string `toml:"parts"`
}

type LayoutJSON struct {
	Align     string    `toml:"align,omitempty"`
	AutoWidth *bool     `toml:"autoWidth,omitempty"`
	Header    PartsJSON `toml:"header"`
	Match     PartsJSON `toml:"match"`
	Context   PartsJSON `toml:"context"`
}

type StyleJson struct {
	Fg   *string `toml:"fg"`
	Bg   *string `toml:"bg"`
	Bold *bool   `toml:"bold"`
}

type ThemeJSON struct {
	Styles map[string]StyleJson `toml:"styles"`
}

type ConfigJSON struct {
	Theme  ThemeJSON  `toml:"theme"`
	Layout LayoutJSON `toml:"layout"`
}

type VarKind uint8

const (
	VarFilepath VarKind = iota
	VarDir
	VarBase
	VarClean
	VarLn
	VarText
)

type Token struct {
	Lit   string
	Var   VarKind
	IsVar bool
}

type CompiledLayout struct {
	Header     []Token
	Match      []Token
	Context    []Token
	AlignRight bool
	AutoWidth  bool
}

type Theme struct {
	Styles map[string]Style
}
