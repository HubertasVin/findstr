package utils

import (
	"archive/tar"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ulikunitz/xz"
)

type TarHandler struct{}

var tarExtensions = []string{
	".tar", ".tar.gz", ".tgz", ".tar.bz2", ".tbz2", ".tar.xz", ".txz",
}

func (t *TarHandler) CanHandle(fileName string) bool {
	lower := strings.ToLower(fileName)
	for _, ext := range tarExtensions {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

func (t *TarHandler) Iterate(archPath string, callback func(name string, isDir bool) error) error {
	file, err := os.Open(archPath)
	if err != nil {
		return fmt.Errorf("failed to open tar: %w", err)
	}
	defer file.Close()

	tr, closer, err := createTarReader(file, archPath)
	if err != nil {
		return err
	}
	if closer != file {
		defer closer.Close()
	}

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		isDir := header.Typeflag == tar.TypeDir
		if err := callback(header.Name, isDir); err != nil {
			return err
		}
	}
	return nil
}

func (t *TarHandler) ReadFile(archPath, targetPath string) (io.ReadCloser, error) {
	file, err := os.Open(archPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open tar: %w", err)
	}

	tr, closer, err := createTarReader(file, archPath)
	if err != nil {
		file.Close()
		return nil, err
	}

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			if closer != file {
				closer.Close()
			}
			file.Close()
			return nil, err
		}

		if filepath.ToSlash(header.Name) == targetPath {
			if header.FileInfo().IsDir() {
				if closer != file {
					closer.Close()
				}
				file.Close()
				return nil, fmt.Errorf("path is a directory")
			}
			return &tarFileReader{
				Reader: io.LimitReader(tr, header.Size),
				file:   file,
				closer: closer,
			}, nil
		}
	}

	if closer != file {
		closer.Close()
	}
	file.Close()
	return nil, fmt.Errorf("file %s not found in archive", targetPath)
}

type tarFileReader struct {
	io.Reader
	file   *os.File
	closer io.Closer
}

func (t *tarFileReader) Close() error {
	if t.closer != nil && t.closer != t.file {
		t.closer.Close()
	}
	return t.file.Close()
}

func createTarReader(file *os.File, archPath string) (*tar.Reader, io.Closer, error) {
	var reader io.Reader = file
	var closer io.Closer = file

	switch {
	case strings.HasSuffix(archPath, ".tar.gz") || strings.HasSuffix(archPath, ".tgz"):
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return nil, nil, err
		}
		reader = gzReader
		closer = gzReader
	case strings.HasSuffix(archPath, ".tar.bz2") || strings.HasSuffix(archPath, ".tbz2"):
		reader = bzip2.NewReader(file)
	case strings.HasSuffix(archPath, ".tar.xz") || strings.HasSuffix(archPath, ".txz"):
		xzReader, err := xz.NewReader(file)
		if err != nil {
			return nil, nil, err
		}
		reader = xzReader
	}

	return tar.NewReader(reader), closer, nil
}
