package serge

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/justinas/alice"
	"github.com/kevinpollet/serge/log"
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

type fileServer struct {
	fileSystem http.FileSystem
}

func NewFileServer(dir string, middlewares ...alice.Constructor) http.Handler {
	fs := &fileServer{dotFileHiddingFileSystem{http.Dir(dir)}}

	return alice.New(middlewares...).Then(fs)
}

func (fs *fileServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	indexPageName := "index.html"
	urlPath := path.Clean(req.URL.Path)
	contentEncodings := []string{encodingBrotli, encodingGzip, encodingDeflate}

	file, err := fs.fileSystem.Open(urlPath)
	if err != nil {
		toResponse(rw, err)
		return
	}

	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		toResponse(rw, err)
		return
	}

	if fileInfo.IsDir() {
		if !strings.HasSuffix(req.URL.Path, "/") {
			relRedirect(rw, req, fmt.Sprint(req.URL.Path, "/"))
		} else {
			req.URL.Path = fmt.Sprintf("%s/%s", urlPath, indexPageName)
			fs.ServeHTTP(rw, req)
		}

		return
	}

	contentEncoding, err := negotiateContentEncoding(req, contentEncodings...)
	if err != nil {
		toResponse(rw, err)
		return
	}

	if contentEncoding == "" {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if contentEncoding != encodingIdentity {
		rwEncoder, err := newResponseWriterEncoder(contentEncoding, rw)
		if err != nil {
			toResponse(rw, err)
			return
		}

		rw = rwEncoder
		defer rwEncoder.Close()
	}

	rw.Header().Add("Content-Encoding", contentEncoding)

	http.ServeContent(rw, req, fileInfo.Name(), fileInfo.ModTime(), file)
}

func relRedirect(rw http.ResponseWriter, req *http.Request, relPath string) {
	query := req.URL.RawQuery
	if query != "" {
		relPath += fmt.Sprint("?", query)
	}

	rw.Header().Add("Location", relPath)
	rw.WriteHeader(http.StatusMovedPermanently)
}

func toResponse(rw http.ResponseWriter, err error) {
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
