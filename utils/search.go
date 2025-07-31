package utils

import (
	"path/filepath"
	"slices"

	"github.com/HubertasVin/findstr/models"
)

// SearchMatchLines walks files under root and returns matches.
func SearchMatchLines(flags models.ProgramFlags) ([]models.FileMatch, error) {
	paths, err := FilePathWalkDir(
		flags.Root,
		flags.ExcludeDir,
		flags.ExcludeFile,
		flags.ThreadCount,
	)
	if err != nil {
		return nil, err
	}

	var matches []models.FileMatch
	for _, rel := range paths {
		full := filepath.Join(flags.Root, rel)
		lines, err := ReadFileLines(full)
		if err != nil {
			return nil, err
		}

		var ctxLines, matchLines []int
		for i, line := range lines {
			if CheckPattern(line, flags.Pattern) {
				ctxLines = append(ctxLines, GetMatchContextLines(i, lines)...)
				matchLines = append(matchLines, i)
			}
		}

		ctxLines = RemoveDuplicate(ctxLines)
		slices.Sort(ctxLines)

		if len(ctxLines) > 0 {
			matches = append(matches, models.FileMatch{
				File:            full,
				ContextLineNums: ctxLines,
				MatchLineNums:   matchLines,
				FileContent:     lines,
			})
		}
	}
	return matches, nil
}
