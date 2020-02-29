package serge

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/kevinpollet/serge/log"
)

type fileServer struct {
	root http.FileSystem
}

func NewFileServer(dir string) http.Handler {
	return &fileServer{root: http.Dir(dir)}
}

func (fs *fileServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	urlPath := path.Clean(req.URL.Path)
	fmt.Println(urlPath)

	file, err := fs.root.Open(urlPath)
	if err != nil {
		toHTTPResponse(rw, err)
		return
	}

	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		toHTTPResponse(rw, err)
		return
	}

	if fileInfo.IsDir() {
		if !strings.HasSuffix(req.URL.Path, "/") {
			rw.Header().Add("Location", urlPath+"/")
			rw.WriteHeader(http.StatusMovedPermanently)
		} else {
			req.URL.Path = urlPath + "/index.html"
			fs.ServeHTTP(rw, req)
		}

		return
	}

	http.ServeContent(rw, req, filepath.Base(urlPath), fileInfo.ModTime(), file)
}

func toHTTPResponse(rw http.ResponseWriter, err error) {
	switch {
	case os.IsNotExist(err):
		rw.WriteHeader(http.StatusNotFound)

	case os.IsPermission(err):
		rw.WriteHeader(http.StatusForbidden)

	default:
		log.Logger().Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
