package utils

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// GetMatchContextLines returns the numeric range around a match.
func GetMatchContextLines(lineNum int, fileContent []string) []int {
	left, right := getLinesRange(lineNum, fileContent)

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

	var files []string
	err = filepath.Walk(absRoot, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(absRoot, path)
		if err != nil {
			return err
		}

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

		files = append(files, rel)
		return nil
	})
	return files, err
}

func getLinesRange(lineNum int, fileContent []string) (int, int) {
	left := subtractTo0(lineNum, 2)

	right := addToBound(lineNum, 2, len(fileContent)-1)

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
