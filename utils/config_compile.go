package utils

import (
	"encoding/json"
	"errors"
	"image/color"

	"github.com/HubertasVin/findstr/models"
	"github.com/icza/gox/imagex/colorx"
)

// Parses the config JSON (layout+theme or layout-only),
// applies defaults, then applies an optional legacy style override.
func ParseConfig(rawConfig string) (models.CompiledLayout, models.Style, error) {
	var layout models.LayoutJSON
	var theme models.Theme

	if rawConfig == "" {
		layout = defaultLayoutJSON()
		theme = defaultThemeResolved()
	} else {
		var cfg models.ConfigJSON
		if err := json.Unmarshal([]byte(rawConfig), &cfg); err != nil {
			if err2 := json.Unmarshal([]byte(rawConfig), &layout); err2 != nil {
				return models.CompiledLayout{}, models.Style{}, err
			}
			layout = fillLayoutDefaults(layout)
			theme = defaultThemeResolved()
		} else {
			layout = fillLayoutDefaults(cfg.Layout)
			theme = resolveThemeWithDefaults(cfg.Theme)
		}
	}

	match := theme.Styles["match"]

	return CompileLayout(layout), match, nil
}

func CompileLayout(l models.LayoutJSON) models.CompiledLayout {
	autoWidth := true
	if l.AutoWidth != nil {
		autoWidth = *l.AutoWidth
	}
	return models.CompiledLayout{
		Header:     compileParts(l.Header.Parts),
		Match:      compileParts(l.Match.Parts),
		Context:    compileParts(l.Context.Parts),
		AlignRight: l.Align != "left",
		AutoWidth:  autoWidth,
	}
}

func compileParts(parts []string) []models.Token {
	vars := map[string]models.VarKind{
		"{filepath}": models.VarFilepath,
		"{dir}":      models.VarDir,
		"{base}":     models.VarBase,
		"{clean}":    models.VarClean,
		"{ln}":       models.VarLn,
		"{text}":     models.VarText,
	}
	out := make([]models.Token, 0, len(parts))
	for _, p := range parts {
		if vk, ok := vars[p]; ok {
			out = append(out, models.Token{IsVar: true, Var: vk})
		} else {
			out = append(out, models.Token{Lit: p})
		}
	}
	return out
}

func defaultLayoutJSON() models.LayoutJSON {
	return models.LayoutJSON{
		Align:     "right",
		AutoWidth: ptrBool(true),
		Header:    models.PartsJSON{Parts: []string{"---", " ", "{filepath}", ":"}},
		Match:     models.PartsJSON{Parts: []string{"{ln}", " | ", "{text}"}},
		Context:   models.PartsJSON{Parts: []string{"{ln}", " | ", "{text}"}},
	}
}

func fillLayoutDefaults(in models.LayoutJSON) models.LayoutJSON {
	d := defaultLayoutJSON()
	if in.Align != "" {
		d.Align = in.Align
	}
	if in.AutoWidth != nil {
		d.AutoWidth = in.AutoWidth
	}
	if len(in.Header.Parts) != 0 {
		d.Header = in.Header
	}
	if len(in.Match.Parts) != 0 {
		d.Match = in.Match
	}
	if len(in.Context.Parts) != 0 {
		d.Context = in.Context
	}
	return d
}

func defaultThemeResolved() models.Theme {
	return models.Theme{
		Styles: map[string]models.Style{
			"header":  {MatchFg: color.RGBA{255, 255, 255, 255}, MatchBg: color.RGBA{0, 0, 0, 0}, MatchBold: true},
			"match":   {MatchFg: color.RGBA{255, 255, 255, 255}, MatchBg: color.RGBA{0, 138, 0, 255}, MatchBold: true},
			"context": {MatchFg: color.RGBA{204, 204, 204, 255}, MatchBg: color.RGBA{0, 0, 0, 0}, MatchBold: false},
		},
	}
}

func resolveThemeWithDefaults(t models.ThemeJSON) models.Theme {
	out := defaultThemeResolved()
	if len(t.Styles) == 0 {
		return out
	}
	for name, sj := range t.Styles {
		def, ok := out.Styles[name]
		if !ok {
			def = models.Style{MatchFg: color.RGBA{255, 255, 255, 255}}
		}
		out.Styles[name] = mergeStyleJSON(sj, def)
	}
	return out
}

func mergeStyleJSON(sj models.StyleJson, def models.Style) models.Style {
	res := def
	if sj.Fg != nil {
		if c, err := parseColor(sj.Fg, nil); err == nil {
			res.MatchFg = c
		}
	}
	if sj.Bg != nil {
		if c, err := parseColor(sj.Bg, nil); err == nil {
			res.MatchBg = c
		}
	}
	if sj.Bold != nil {
		res.MatchBold = *sj.Bold
	}
	return res
}

func parseLegacyStyle(s string) (models.Style, error) {
	if s == "" {
		s = "{}"
	}
	var dto models.StyleJson
	if err := json.Unmarshal([]byte(s), &dto); err != nil {
		return models.Style{}, err
	}
	fg, err := parseColor(dto.Fg, nil)
	if err != nil {
		return models.Style{}, err
	}
	bg, err := parseColor(dto.Bg, &color.RGBA{R: 0, G: 138, B: 0, A: 255})
	if err != nil {
		return models.Style{}, err
	}
	res := models.Style{MatchFg: fg, MatchBg: bg}
	if dto.Bold != nil {
		res.MatchBold = *dto.Bold
	}
	return res, nil
}

func ptrBool(b bool) *bool { return &b }

func parseColor(c *string, def *color.RGBA) (color.RGBA, error) {
	if c != nil {
		rgba, err := colorx.ParseHexColor(*c)
		if err != nil {
			return color.RGBA{}, errors.New("Error: Parsing hex color got \"" + err.Error() + "\"")
		}
		return rgba, nil
	}
	if def != nil {
		return *def, nil
	}
	return color.RGBA{255, 255, 255, 255}, nil
}
