package serge

import (
	"net/http"
	"os"
	"strings"
)

type dotFileHiddingFileSystem struct {
	http.FileSystem
}

func (fs dotFileHiddingFileSystem) Open(name string) (http.File, error) {
	if containsDotFile(name) {
		return nil, os.ErrNotExist
	}

	return fs.FileSystem.Open(name)
}

func (fs dotFileHiddingFileSystem) OpenWithStat(name string) (http.File, os.FileInfo, error) {
	file, err := fs.Open(name)
	if err != nil {
		return nil, nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}

	return file, fileInfo, nil
}

func containsDotFile(name string) bool {
	for _, file := range strings.Split(name, "/") {
		if strings.HasPrefix(file, ".") {
			return true
		}
	}

	return false
}
