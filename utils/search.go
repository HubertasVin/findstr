package utils

import (
	"path/filepath"
	"slices"

	"github.com/HubertasVin/findstr/models"
)

// SearchMatchLines walks files under root and returns matches.
func SearchMatchLines(root, pattern string) ([]models.FileMatch, error) {
	paths, err := FilePathWalkDir(root)
	if err != nil {
		return nil, err
	}

	var matches []models.FileMatch
	for _, rel := range paths {
		full := filepath.Join(root, rel)
		lines, err := ReadFileLines(full)
		if err != nil {
			return nil, err
		}

		var ctxLines, highLines []int
		for i, line := range lines {
			if CheckPattern(line, pattern) {
				ctxLines = append(ctxLines, GetMatchContextLines(i, lines)...)
				highLines = append(highLines, i)
			}
		}

		ctxLines = RemoveDuplicate(ctxLines)
		slices.Sort(ctxLines)

		if len(ctxLines) > 0 {
			matches = append(matches, models.FileMatch{
				File:            full,
				ContextLineNums: ctxLines,
				HighLineNums:    highLines,
				FileContent:     lines,
			})
		}
	}
	return matches, nil
}
