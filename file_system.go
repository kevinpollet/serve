package serge

import (
	"net/http"
	"os"
	"strings"
)

type dotFileHiddingFileSystem struct {
	http.FileSystem
}

func (df dotFileHiddingFileSystem) Open(name string) (http.File, error) {
	if containsDotFile(name) {
		return nil, os.ErrNotExist
	}

	return df.FileSystem.Open(name)
}

func containsDotFile(name string) bool {
	for _, file := range strings.Split(name, "/") {
		if strings.HasPrefix(file, ".") {
			return true
		}
	}

	return false
}
