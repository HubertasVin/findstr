package utils

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"text/tabwriter"

	"github.com/HubertasVin/findstr/models"
	"github.com/fatih/color"
)

// PrintMatches pretty‑prints all matches with context.
func PrintMatches(matches <-chan models.FileMatch) {
	header := color.New(color.Bold, color.FgWhite).SprintFunc()
	high := color.RGB(0, 135, 0).SprintfFunc()

    first := true
	for m := range matches {
		if !first {
			fmt.Println()
		}
        first = false
		printMatchLines(m, header, high)
	}
}

func printMatchLines(
	m models.FileMatch,
	headerFn func(...any) string,
	highFn func(string, ...any) string,
) {
	const reset = "\x1b[0m"
	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 4, ' ', tabwriter.TabIndent)

	last := m.ContextLineNums[len(m.ContextLineNums)-1]
	lnWidth := len(strconv.Itoa(last + 1))

	fmt.Fprintln(tw, headerFn("--- "+m.File))

	prev := -1
	for _, ln := range m.ContextLineNums {
		if prev != -1 && ln-prev >= 2 {
			fmt.Fprintln(tw, headerFn("..."))
		}

		numFmt := fmt.Sprintf("%-*d", lnWidth, ln+1)
		text := fmt.Sprintf("%s | %s", numFmt, m.FileContent[ln])

		if slices.Contains(m.MatchLineNums, ln) {
			fmt.Fprintln(tw, highFn("%s", text)+reset)
		} else {
			fmt.Fprintln(tw, text)
		}
		tw.Flush()
		prev = ln
	}
}
