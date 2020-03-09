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

type fileServer struct {
	fileSystem   http.FileSystem
	middlewares  []alice.Constructor
	errorHandler func(http.FileSystem, http.ResponseWriter, error)
}

// NewFileServer returns a new handler instance that serves HTTP requests
// with the contents of the given directory.
func NewFileServer(dir string, options ...Option) http.Handler {
	fs := &fileServer{fileSystem: dotFileHiddingFileSystem{http.Dir(dir)}}

	for _, option := range options {
		option(fs)
	}

	return alice.New(fs.middlewares...).Then(fs)
}

func (fs *fileServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	urlPath := path.Clean(req.URL.Path)
	contentEncodings := []string{encodingBrotli, encodingGzip, encodingDeflate}

	if !strings.HasPrefix(req.URL.Path, "/") {
		req.URL.Path = "/" + req.URL.Path
		urlPath = path.Clean(req.URL.Path)
	}

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

	fs.serveContent(rw, req, file, fileInfo)
}

func (fs *fileServer) handleError(rw http.ResponseWriter, err error) {
	if fs.errorHandler == nil {
		defaultErrorHandler(fs.fileSystem, rw, err)
		return
	}

	fs.errorHandler(fs.fileSystem, rw, err)
}

func (fs *fileServer) serveContent(
	rw http.ResponseWriter,
	req *http.Request,
	file io.ReadSeeker,
	fileInfo os.FileInfo,
) {
	indexPageName := "index.html"

	if !fileInfo.IsDir() {
		http.ServeContent(rw, req, fileInfo.Name(), fileInfo.ModTime(), file)
		return
	}

	if !strings.HasSuffix(req.URL.Path, "/") {
		redirectTo(rw, req, fmt.Sprint(req.URL.Path, "/"))
		return
	}

	req.URL.Path = fmt.Sprintf("%s/%s", req.URL.Path, indexPageName)
	fs.ServeHTTP(rw, req)
}

func defaultErrorHandler(fs http.FileSystem, rw http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError

	if os.IsNotExist(err) || os.IsPermission(err) {
		statusCode = http.StatusNotFound
	}

	if statusCode == http.StatusInternalServerError {
		log.Logger().Error(err)
	}

	rw.WriteHeader(statusCode)

	errorPageName := fmt.Sprintf("%d.html", statusCode)
	if file, err := fs.Open(errorPageName); err == nil {
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
