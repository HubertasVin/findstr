package utils

import (
	"fmt"
	"io"
	"path"
	"path/filepath"
	"strings"
)

type ArchiveHandler interface {
	CanHandle(fileName string) bool
	Iterate(archPath string, callback func(name string, isDir bool) error) error
	ReadFile(archPath, targetPath string) (io.ReadCloser, error)
}

var archiveHandlers = []ArchiveHandler{
	&ZipHandler{},
	&TarHandler{},
	&RarHandler{},
	&SevenZipHandler{},
}

// IsCompatibleArchive checks if a file is a supported archive
func IsCompatibleArchive(fileName string) bool {
	for _, handler := range archiveHandlers {
		if handler.CanHandle(fileName) {
			return true
		}
	}
	return false
}

// IsPathInArchive checks if a path contains archive separator
func IsPathInArchive(path string) bool {
	return strings.Contains(path, "#")
}

// GetArchiveFiles lists all files in an archive
func GetArchiveFiles(archPath, excludeDir, excludeFile string, skipGit, searchArch bool) ([]string, error) {
	handler := getHandler(archPath)
	if handler == nil {
		return nil, fmt.Errorf("unsupported archive format: %s", archPath)
	}

	var files []string
	err := handler.Iterate(archPath, func(name string, isDir bool) error {
		if !isDir {
			files = append(files, combineArchivePath(archPath, name))
		}
		return nil
	})

	return files, err
}

// ReadArchiveFileLines reads lines from a file within an archive
func ReadArchiveFileLines(path string) ([]string, error) {
	paths := strings.Split(path, "#")
	if len(paths) != 2 {
		return nil, fmt.Errorf("invalid archive path format: %s", path)
	}

	handler := getHandler(paths[0])
	if handler == nil {
		return nil, fmt.Errorf("unsupported archive format: %s", paths[0])
	}

	targetPath := filepath.ToSlash(paths[1])
	rc, err := handler.ReadFile(paths[0], targetPath)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	return readLines(rc)
}

// Helper functions
func combineArchivePath(archivePath, internalPath string) string {
	archivePath = filepath.Clean(archivePath)
	internalPath = path.Clean(internalPath)
	internalPath = strings.TrimPrefix(internalPath, "/")

	if internalPath == "" || internalPath == "." {
		return archivePath
	}

	return archivePath + "#" + internalPath
}

func getHandler(fileName string) ArchiveHandler {
	for _, handler := range archiveHandlers {
		if handler.CanHandle(fileName) {
			return handler
		}
	}
	return nil
}
