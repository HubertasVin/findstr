package utils

import (
	"path/filepath"
	"strings"
)

// RemoveDuplicate returns a new slice that preserves order
// while eliminating duplicates.
func RemoveDuplicate[T comparable](in []T) []T {
	seen := make(map[T]struct{})
	var out []T

	for _, v := range in {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}

func SplitStringToArray(input, separator string) []string {
	var result []string
	if input != "" {
		for part := range strings.SplitSeq(input, separator) {
			cleaned := filepath.Clean(strings.TrimSpace(part))
			if cleaned != "" {
				result = append(result, cleaned)
			}
		}
	}

	return result
}
