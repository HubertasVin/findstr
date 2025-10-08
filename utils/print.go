package utils

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/HubertasVin/findstr/models"
	"github.com/fatih/color"
)

type fileVars struct {
	filepath string
	dir      string
	base     string
	clean    string
}

func PrintMatches(
	ctx context.Context,
	matches <-chan models.FileMatch,
	layout models.CompiledLayout,
	theme models.Theme,
	contextSize int,
) {
	w := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer w.Flush()

	headerStyleFn := buildStyleFn(theme.Styles["header"])
	matchStyleFn := buildStyleFn(theme.Styles["match"])
	contextStyleFn := buildStyleFn(theme.Styles["context"])
	const resetClear = "\x1b[0m\x1b[K"
	const tabWidth = 4

	first := true
	for {
		select {
		case <-ctx.Done():
			return
		case fm, ok := <-matches:
			if !ok {
				return
			}

			if !first {
				fmt.Fprintln(w)
			}
			first = false

			fv := fileVars{
				filepath: fm.File,
				dir:      filepath.Dir(fm.File),
				base:     filepath.Base(fm.File),
				clean:    filepath.Clean(fm.File),
			}

			leftWidth := 0
			if layout.AutoWidth && len(fm.ContextLineNums) > 0 {
				leftWidth = numDigits(fm.ContextLineNums[len(fm.ContextLineNums)-1] + 1)
			}

			if len(layout.Header) > 0 {
				line := renderTokens(layout.Header, fv, 0, "", leftWidth, layout.AlignRight, tabWidth)
				fmt.Fprint(w, headerStyleFn("%s", line))
				fmt.Fprint(w, resetClear)
				fmt.Fprintln(w)
			}

			matchSet := make(map[int]struct{}, len(fm.MatchLineNums))
			for _, ln := range fm.MatchLineNums {
				matchSet[ln] = struct{}{}
			}

			prev := -1
			for _, ln := range fm.ContextLineNums {
				if prev != -1 && ln-prev > contextSize {
					fmt.Fprint(w, headerStyleFn("%s", "..."))
					fmt.Fprint(w, resetClear)
					fmt.Fprintln(w)
				}

				text := fm.FileContent[ln]
				var tokens []models.Token
				if _, ok := matchSet[ln]; ok {
					tokens = layout.Match
				} else {
					tokens = layout.Context
				}

				line := renderTokens(tokens, fv, ln+1, text, leftWidth, layout.AlignRight, tabWidth)

				if _, ok := matchSet[ln]; ok {
					fmt.Fprint(w, matchStyleFn("%s", line))
				} else {
					fmt.Fprint(w, contextStyleFn("%s", line))
				}
				fmt.Fprint(w, resetClear)
				fmt.Fprintln(w)
				prev = ln
			}
		}
	}
}

func buildStyleFn(s models.Style) func(format string, a ...any) string {
	c := color.RGB(int(s.Fg.R), int(s.Fg.G), int(s.Fg.B))
	if s.Bg.A != 0 {
		c = c.AddBgRGB(int(s.Bg.R), int(s.Bg.G), int(s.Bg.B))
	}
	if s.Bold {
		c = c.Add(color.Bold)
	}
	return c.SprintfFunc()
}

func renderTokens(
	toks []models.Token,
	fv fileVars,
	ln int,
	text string,
	lnWidth int,
	alignRight bool,
	tabWidth int,
) string {
	var buf bytes.Buffer
	for _, t := range toks {
		if !t.IsVar {
			buf.WriteString(t.Lit)
			continue
		}
		switch t.Var {
		case models.VarFilepath:
			buf.WriteString(fv.filepath)
		case models.VarDir:
			buf.WriteString(fv.dir)
		case models.VarBase:
			buf.WriteString(fv.base)
		case models.VarClean:
			buf.WriteString(fv.clean)
		case models.VarLn:
			buf.WriteString(renderLineNum(ln, lnWidth, alignRight))
		case models.VarText:
			buf.WriteString(expandTabs(text, 0, tabWidth))
		}
	}
	return buf.String()
}

func expandTabs(s string, startCol, tabWidth int) string {
	if tabWidth <= 0 {
		return s
	}
	var b strings.Builder
	b.Grow(len(s) + 8)
	col := startCol
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\t' {
			n := tabWidth - (col % tabWidth)
			if n == 0 {
				n = tabWidth
			}
			for j := 0; j < n; j++ {
				b.WriteByte(' ')
				col++
			}
			continue
		}
		b.WriteByte(c)
		if c == '\n' || c == '\r' {
			col = 0
		} else {
			col++
		}
	}
	return b.String()
}

func renderLineNum(ln int, lnWidth int, alignRight bool) string {
	if alignRight && lnWidth > 0 {
		s := strconv.Itoa(ln)
		if pad := lnWidth - len(s); pad > 0 {
			return fmt.Sprintf("%*s%s", pad, "", s)
		}
		return s
	}
	return strconv.Itoa(ln)
}

func numDigits(n int) int {
	if n < 10 {
		return 1
	}
	if n < 100 {
		return 2
	}
	if n < 1000 {
		return 3
	}
	if n < 10000 {
		return 4
	}
	return len(strconv.Itoa(n))
}
