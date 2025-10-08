package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

type ZipHandler struct{}

func (z *ZipHandler) CanHandle(fileName string) bool {
	return strings.HasSuffix(strings.ToLower(fileName), ".zip")
}

func (z *ZipHandler) Iterate(archPath string, callback func(name string, isDir bool) error) error {
	reader, err := zip.OpenReader(archPath)
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		if err := callback(file.Name, file.FileInfo().IsDir()); err != nil {
			return err
		}
	}
	return nil
}

func (z *ZipHandler) ReadFile(archPath, targetPath string) (io.ReadCloser, error) {
	reader, err := zip.OpenReader(archPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open zip: %w", err)
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
			return &zipFileReader{rc: rc, zr: reader}, nil
		}
	}

	reader.Close()
	return nil, fmt.Errorf("file %s not found in archive", targetPath)
}

type zipFileReader struct {
	rc io.ReadCloser
	zr *zip.ReadCloser
}

func (z *zipFileReader) Read(p []byte) (n int, err error) {
	return z.rc.Read(p)
}

func (z *zipFileReader) Close() error {
	z.rc.Close()
	return z.zr.Close()
}
