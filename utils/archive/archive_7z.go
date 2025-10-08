package utils

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/bodgit/sevenzip"
)

type SevenZipHandler struct{}

func (s *SevenZipHandler) CanHandle(fileName string) bool {
	return strings.HasSuffix(strings.ToLower(fileName), ".7z")
}

func (s *SevenZipHandler) Iterate(archPath string, callback func(name string, isDir bool) error) error {
	reader, err := sevenzip.OpenReader(archPath)
	if err != nil {
		return fmt.Errorf("failed to open 7z: %w", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		isDir := file.FileHeader.FileInfo().IsDir()
		if err := callback(file.Name, isDir); err != nil {
			return err
		}
	}
	return nil
}

func (s *SevenZipHandler) ReadFile(archPath, targetPath string) (io.ReadCloser, error) {
	reader, err := sevenzip.OpenReader(archPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open 7z: %w", err)
	}

	for _, file := range reader.File {
		if filepath.ToSlash(file.Name) == targetPath {
			if file.FileInfo().IsDir() {
				reader.Close()
				return nil, fmt.Errorf("path is a directory")
			}
			rc, err := file.Open()
			if err != nil {
				reader.Close()
				return nil, err
			}
			return &sevenZipFileReader{rc: rc, zr: reader}, nil
		}
	}

	reader.Close()
	return nil, fmt.Errorf("file %s not found in archive", targetPath)
}

type sevenZipFileReader struct {
	rc io.ReadCloser
	zr *sevenzip.ReadCloser
}

func (s *sevenZipFileReader) Read(p []byte) (n int, err error) {
	return s.rc.Read(p)
}

func (s *sevenZipFileReader) Close() error {
	s.rc.Close()
	return s.zr.Close()
}
