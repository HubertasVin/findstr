package utils

import (
	"archive/zip"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

var ArchFileTypes = [2]string{"zip", "tar"}

func IsCompatibleArchive(fileName string) bool {
	ext := strings.TrimPrefix(filepath.Ext(fileName), ".")
	return slices.Contains(ArchFileTypes[:], ext)
}

func GetArchiveFiles(
	archPath, excludeDir, excludeFile string,
	skipGit, searchArchives bool,
) ([]string, error) {
	return getZipFiles(archPath, excludeDir, excludeFile, skipGit)
}

func getZipFiles(
	path, excludeDir, excludeFile string,
	skipGit bool,
) ([]string, error) {
	var files []string

	zf, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer zf.Close()

	for _, file := range zf.File {
		files = append(files, combineArchivePath(path, file.Name))
	}

	return files, nil
}

// CombineArchivePath combines an archive file path with an internal archive path
func combineArchivePath(archivePath, internalPath string) string {
    archivePath = filepath.Clean(archivePath)
    
    internalPath = path.Clean(internalPath)
    internalPath = strings.TrimPrefix(internalPath, "/")
    
    if internalPath == "" || internalPath == "." {
        return archivePath
    }
    
    return archivePath + "#" + internalPath
}
