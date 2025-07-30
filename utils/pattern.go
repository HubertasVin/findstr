package utils

import "regexp"

// CheckPattern is a thin wrapper around regexp.MatchString.
func CheckPattern(input, pattern string) bool {
	return regexp.MustCompile(pattern).MatchString(input)
}
