package utils

import (
	"log"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/HubertasVin/findstr/models"
	"github.com/remeh/sizedwaitgroup"
)

// SearchMatchLines walks files under root and returns matches via a channel.
// The returned channel is closed when all workers are done.
func SearchMatchLines(flags models.ProgramFlags) (<-chan models.FileMatch, error) {
	paths, err := FilePathWalkDir(
		flags.Root,
		flags.ExcludeDir,
		flags.ExcludeFile,
		flags.ThreadCount,
	)
	if err != nil {
		return nil, err
	}

	re, err := regexp.Compile(flags.Pattern)
	if err != nil {
		return nil, err
	}

	numWorkers := min(flags.ThreadCount, len(paths))

	out := runParallel(paths, re, flags.Root, numWorkers, flags.ContextSize)

	return out, nil
}

func runParallel(
	paths []string,
	re *regexp.Regexp,
	root string,
	numWorkers int,
	contextSize int,
) chan models.FileMatch {
	out := make(chan models.FileMatch, numWorkers*2)
	wg := sizedwaitgroup.New(numWorkers)

	go func() {
		for _, rel := range paths {
			wg.Add()
			go func(rel string) {
				defer wg.Done()
				processFile(rel, root, contextSize, re, out)
			}(rel)
		}
		wg.Wait()
		close(out)
	}()

	return out
}

func processFile(
	relPath, root string,
	contextSize int,
	re *regexp.Regexp,
	ch chan<- models.FileMatch,
) {
	full := filepath.Join(root, relPath)
	lines, err := ReadFileLines(full)
	if err != nil {
		log.Println("Failed to read file:", relPath)
		return
	}

	var ctxLines, matchLines []int
	for i, line := range lines {
		if re.MatchString(line) {
			ctxLines = append(ctxLines, GetMatchContextLines(i, lines, contextSize)...)
			matchLines = append(matchLines, i)
		}
	}

	if len(ctxLines) == 0 {
		return
	}

	ctxLines = RemoveDuplicate(ctxLines)
	sort.Ints(ctxLines)

	ch <- models.FileMatch{
		File:            full,
		ContextLineNums: ctxLines,
		MatchLineNums:   matchLines,
		FileContent:     lines,
	}
}
