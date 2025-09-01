package utils

import (
	"log"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/HubertasVin/chanseq"
	"github.com/HubertasVin/findstr/models"
)

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

	numWorkers := min(flags.ThreadCount, len(paths))
	out := runParallel(paths, flags.Pattern, flags.Root, numWorkers, flags.ContextSize)
	return out, nil
}

func runParallel(
	paths []string,
	pattern string,
	root string,
	numWorkers int,
	contextSize int,
) <-chan models.FileMatch {
	type job struct {
		idx int
		rel string
	}

	jobs := make(chan job, numWorkers*2)
	tmp := make(chan chanseq.Seq[models.FileMatch], numWorkers*2)

	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for w := 0; w < numWorkers; w++ {
		go func() {
			defer wg.Done()
			for j := range jobs {
				match := processFile(j.rel, root, contextSize, pattern)
				tmp <- chanseq.Seq[models.FileMatch]{Index: j.idx, Val: match}
			}
		}()
	}

	go func() {
		for i, rel := range paths {
			jobs <- job{idx: i, rel: rel}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(tmp)
	}()

	out := chanseq.ReorderByIndex(tmp)
	return out
}

func processFile(
	relPath, root string,
	contextSize int,
	pattern string,
) *models.FileMatch {
	full := filepath.Join(root, relPath)

	// Cheap binary check to avoid scanning non-text files.
	if IsLikelyBinary(full) {
		return nil
	}

	lines, err := ReadFileLines(full)
	if err != nil {
		log.Println("Failed to read file:", relPath)
		return nil
	}

	var ctxLines, matchLines []int
	for i, line := range lines {
		if strings.Contains(line, pattern) {
			ctxLines = append(ctxLines, GetMatchContextLines(i, lines, contextSize)...)
			matchLines = append(matchLines, i)
		}
	}
	if len(ctxLines) == 0 {
		return nil
	}

	ctxLines = RemoveDuplicate(ctxLines)
	sort.Ints(ctxLines)

	return &models.FileMatch{
		File:            full,
		ContextLineNums: ctxLines,
		MatchLineNums:   matchLines,
		FileContent:     lines,
	}
}
