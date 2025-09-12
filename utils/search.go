package utils

import (
	"context"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/HubertasVin/chanseq"
	"github.com/HubertasVin/findstr/models"
)

func SearchMatchLines(ctx context.Context, flags models.ProgramFlags) (<-chan models.FileMatch, error) {
	paths, err := FilePathWalkDir(ctx,
		flags.Root,
		flags.ExcludeDir,
		flags.ExcludeFile,
		flags.ThreadCount,
		flags.SkipGit,
	)
	if err != nil {
		return nil, err
	}

	numWorkers := min(flags.ThreadCount, len(paths))
	out := runParallel(ctx, paths, flags.Pattern, flags.Root, numWorkers, flags.ContextSize)
	return out, nil
}

func runParallel(
	ctx context.Context,
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
	for range numWorkers {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case j, ok := <-jobs:
					if !ok {
						return
					}
					match := processFile(j.rel, root, contextSize, pattern)
					select {
					case <-ctx.Done():
						return
					case tmp <- chanseq.Seq[models.FileMatch]{Index: j.idx, Val: match}:
					}
				}
			}
		}()
	}

	go func() {
		for i, rel := range paths {
			select {
			case <-ctx.Done():
				break
			case jobs <- job{idx: i, rel: rel}:
			}
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
