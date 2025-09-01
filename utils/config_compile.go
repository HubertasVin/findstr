package utils

import (
	"errors"
	"image/color"

	"github.com/HubertasVin/findstr/models"
	"github.com/icza/gox/imagex/colorx"
)

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

func defaultThemeResolved() models.Theme {
	return models.Theme{
		Styles: map[string]models.Style{
			"header":  {Fg: color.RGBA{255, 255, 255, 255}, Bg: color.RGBA{0, 0, 0, 0}, Bold: true},
			"match":   {Fg: color.RGBA{255, 255, 255, 255}, Bg: color.RGBA{0, 0, 0, 0}, Bold: true},
			"context": {Fg: color.RGBA{255, 255, 255, 255}, Bg: color.RGBA{0, 0, 0, 0}, Bold: false},
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
			def = models.Style{Fg: color.RGBA{255, 255, 255, 255}}
		}
		out.Styles[name] = mergeStyleJSON(sj, def)
	}
	return out
}

func mergeStyleJSON(sj models.StyleJson, def models.Style) models.Style {
	res := def
	if sj.Fg != nil {
		if c, err := parseColor(sj.Fg, nil); err == nil {
			res.Fg = c
		}
	}
	if sj.Bg != nil {
		if c, err := parseColor(sj.Bg, nil); err == nil {
			res.Bg = c
		}
	}
	if sj.Bold != nil {
		res.Bold = *sj.Bold
	}
	return res
}

func fillLayoutDefaults(in models.LayoutJSON) models.LayoutJSON {
	d := models.LayoutJSON{
		Align:     "right",
		AutoWidth: ptrBool(true),
		Header:    models.PartsJSON{Parts: []string{"---", " ", "{filepath}", ":"}},
		Match:     models.PartsJSON{Parts: []string{"{ln}", " | ", "{text}"}},
		Context:   models.PartsJSON{Parts: []string{"{ln}", " | ", "{text}"}},
	}
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
