package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nwaples/rardecode"
)

type RarHandler struct{}

func (r *RarHandler) CanHandle(fileName string) bool {
	return strings.HasSuffix(strings.ToLower(fileName), ".rar")
}

func (r *RarHandler) Iterate(archPath string, callback func(name string, isDir bool) error) error {
	file, err := os.Open(archPath)
	if err != nil {
		return fmt.Errorf("failed to open rar: %w", err)
	}
	defer file.Close()

	reader, err := rardecode.NewReader(file, "")
	if err != nil {
		return fmt.Errorf("failed to read rar: %w", err)
	}

	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := callback(header.Name, header.IsDir); err != nil {
			return err
		}
	}
	return nil
}

func (r *RarHandler) ReadFile(archPath, targetPath string) (io.ReadCloser, error) {
	file, err := os.Open(archPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open rar: %w", err)
	}

	reader, err := rardecode.NewReader(file, "")
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to read rar: %w", err)
	}

	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			file.Close()
			return nil, err
		}

		if filepath.ToSlash(header.Name) == targetPath {
			if header.IsDir {
				file.Close()
				return nil, fmt.Errorf("path is a directory")
			}
			return &rarFileReader{Reader: reader, file: file}, nil
		}
	}

	file.Close()
	return nil, fmt.Errorf("file %s not found in archive", targetPath)
}

type rarFileReader struct {
	io.Reader
	file *os.File
}

func (r *rarFileReader) Close() error {
	return r.file.Close()
}
