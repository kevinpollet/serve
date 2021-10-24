package srv

import (
	"net/http"
	"os"
	"strings"
)

type dotFileHidingFile struct {
	http.File
}

func (f dotFileHidingFile) Readdir(n int) ([]os.FileInfo, error) {
	files, err := f.File.Readdir(n)
	if err != nil {
		return nil, err
	}

	var filteredFiles []os.FileInfo

	for _, file := range files {
		if !strings.HasPrefix(file.Name(), ".") {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return filteredFiles, nil
}

type dotFileHidingFileSystem struct {
	http.FileSystem
}

func (fs dotFileHidingFileSystem) Open(name string) (http.File, error) {
	if containsDotFile(name) {
		return nil, os.ErrNotExist
	}

	file, err := fs.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}

	return dotFileHidingFile{file}, err
}

func (fs dotFileHidingFileSystem) OpenWithStat(name string) (http.File, os.FileInfo, error) {
	file, err := fs.Open(name)
	if err != nil {
		return nil, nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}

	return dotFileHidingFile{file}, fileInfo, nil
}

func containsDotFile(name string) bool {
	for _, file := range strings.Split(name, "/") {
		if strings.HasPrefix(file, ".") {
			return true
		}
	}

	return false
}
