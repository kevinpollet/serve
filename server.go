package serge

import (
	"fmt"
	"io"
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

	// negotiate content encoding
	contentEncoding, err := negotiateContentEncoding(req, contentEncodings...)
	if err != nil {
		fs.handleError(rw, err)
		return
	}

	if contentEncoding == "" {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if contentEncoding != encodingIdentity {
		rwEncoder, err := newResponseWriterEncoder(contentEncoding, rw)
		if err != nil {
			fs.handleError(rw, err)
			return
		}

		rw = rwEncoder
		defer rwEncoder.Close()
	}

	rw.Header().Add("Content-Encoding", contentEncoding)

	// serve file
	file, err := fs.fileSystem.Open(urlPath)
	if err != nil {
		fs.handleError(rw, err)
		return
	}

	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fs.handleError(rw, err)
		return
	}

	if fileInfo.IsDir() {
		if !strings.HasSuffix(req.URL.Path, "/") {
			redirectTo(rw, req, fmt.Sprint(req.URL.Path, "/"))
		} else {
			req.URL.Path = fmt.Sprintf("%s/%s", urlPath, indexPageName)
			fs.ServeHTTP(rw, req)
		}

		return
	}

	http.ServeContent(rw, req, fileInfo.Name(), fileInfo.ModTime(), file)
}

func (fs *fileServer) handleError(rw http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError

	if os.IsNotExist(err) || os.IsPermission(err) {
		statusCode = http.StatusNotFound
	}

	if statusCode == http.StatusInternalServerError {
		log.Logger().Error(err)
	}

	rw.WriteHeader(statusCode)

	errorPageName := fmt.Sprintf("%d.html", statusCode)
	if file, err := fs.fileSystem.Open(errorPageName); err == nil {
		defer file.Close()
		io.Copy(rw, file) //nolint
	}
}

func redirectTo(rw http.ResponseWriter, req *http.Request, path string) {
	query := req.URL.RawQuery
	if query != "" {
		path += fmt.Sprint("?", query)
	}

	rw.Header().Add("Location", path)
	rw.WriteHeader(http.StatusMovedPermanently)
}
