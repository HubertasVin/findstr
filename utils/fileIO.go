package utils

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// GetMatchContextLines returns the numeric range around a match.
func GetMatchContextLines(lineNum int, fileContent []string, contextSize int) []int {
	left, right := getLinesRange(lineNum, fileContent, contextSize)
	return makeRange(left, right)
}

// ReadFileLines reads an entire file into a slice of lines.
func ReadFileLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*2048), 2048*2048)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

// FilePathWalkDir returns a slice of relative file paths under root.
func FilePathWalkDir(root, excludeDir, excludeFile string, threadCount int) ([]string, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	excludeDirs := SplitStringToArray(excludeDir, ",")
	excludeFiles := SplitStringToArray(excludeFile, ",")

	var files []string

	skipSpecial := map[string]struct{}{
		"proc":       {},
		"sys":        {},
		"dev":        {},
		"run":        {},
		"lost+found": {},
	}

	err = filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			if os.IsNotExist(walkErr) {
				return nil
			}
			return walkErr
		}

		rel, err := filepath.Rel(absRoot, path)
		if err != nil {
			return err
		}

		// skip special virtual/volatile directories
		if d.IsDir() {
			for skip := range skipSpecial {
				if rel == skip || strings.HasPrefix(rel, skip+string(os.PathSeparator)) {
					return fs.SkipDir
				}
			}
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}
		if !info.Mode().IsRegular() {
			return nil
		}

		// apply exclude-dir logic
		dir := filepath.Dir(rel)
		for _, ex := range excludeDirs {
			if ex == "." {
				return nil
			}
			relToEx, err := filepath.Rel(ex, dir)
			if err == nil && !strings.HasPrefix(relToEx, "..") {
				return nil
			}
		}

		// apply exclude-file logic
		if fileExcludedByPattern(rel, excludeFiles) {
			return nil
		}

		files = append(files, rel)
		return nil
	})

	return files, err
}

func getLinesRange(lineNum int, fileContent []string, contextSize int) (int, int) {
	left := subtractTo0(lineNum, contextSize)
	right := addToBound(lineNum, contextSize, len(fileContent)-1)
	return left, right
}

func subtractTo0(x, y int) int {
	if x-y < 0 {
		return 0
	}
	return x - y
}

func addToBound(x, y, bound int) int {
	if x+y > bound {
		return bound
	}
	return x + y
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func fileExcludedByPattern(rel string, patterns []string) bool {
	if len(patterns) == 0 {
		return false
	}

	fileBase := filepath.Base(rel)
	for _, pat := range patterns {
		if pat == "" {
			continue
		}
		if pat == "noext" {
			if filepath.Ext(fileBase) == "" {
				return true
			}
			continue
		}
		if ok, _ := filepath.Match(pat, rel); ok {
			return true
		}
		if ok, _ := filepath.Match(pat, filepath.Base(rel)); ok {
			return true
		}
	}

	return false
}
