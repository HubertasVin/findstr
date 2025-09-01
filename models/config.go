package models

type PartsJSON struct {
	Parts []string `json:"parts"`
}

type LayoutJSON struct {
	Align     string    `json:"align,omitempty"`     // "left" or "right"
	AutoWidth *bool     `json:"autoWidth,omitempty"` // nil => default true
	Header    PartsJSON `json:"header"`
	Match     PartsJSON `json:"match"`
	Context   PartsJSON `json:"context"`
}

type StyleJson struct {
	Fg   *string `json:"fg"`
	Bg   *string `json:"bg"`
	Bold *bool   `json:"bold"`
}

type ThemeJSON struct {
	Styles map[string]StyleJson `json:"styles"`
}

type ConfigJSON struct {
	Theme  ThemeJSON  `json:"theme"`
	Layout LayoutJSON `json:"layout"`
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
