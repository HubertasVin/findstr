package utils

import (
	"path/filepath"
	"slices"
	"strings"
)

var ArchFileTypes = [2]string{"zip", "tar"}

func IsCompatibleArchive(fileName string) bool {
	ext := strings.TrimPrefix(filepath.Ext(fileName), ".")
	return slices.Contains(ArchFileTypes[:], ext)
}
