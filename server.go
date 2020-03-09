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
	autoIndex    bool
	errorHandler func(http.FileSystem, http.ResponseWriter, error)
	fileSystem   dotFileHidingFileSystem
	middlewares  []alice.Constructor
}

// NewFileServer returns a new handler instance that serves HTTP requests
// with the contents of the given directory.
func NewFileServer(dir string, options ...Option) http.Handler {
	fs := &fileServer{fileSystem: dotFileHidingFileSystem{http.Dir(dir)}}

	for _, option := range options {
		option(fs)
	}

	return alice.New(fs.middlewares...).Then(fs)
}

func (fs *fileServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	contentEncodings := []string{encodingBrotli, encodingGzip, encodingDeflate}

	if !strings.HasPrefix(req.URL.Path, "/") {
		req.URL.Path = "/" + req.URL.Path
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
	file, fileInfo, err := fs.fileSystem.OpenWithStat(path.Clean(req.URL.Path))
	if err != nil {
		fs.handleError(rw, err)
		return
	}

	defer file.Close()

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
	file http.File,
	fileInfo os.FileInfo,
) {
	if !fileInfo.IsDir() {
		http.ServeContent(rw, req, fileInfo.Name(), fileInfo.ModTime(), file)
		return
	}

	// enforce trailing slash
	if !strings.HasSuffix(req.URL.Path, "/") {
		redirectTo(rw, req, fmt.Sprint(req.URL.Path, "/"))
		return
	}

	indexFilePath := path.Clean(req.URL.Path + "/index.html")
	indexFile, indexFileInfo, err := fs.fileSystem.OpenWithStat(indexFilePath)

	if err != nil {
		switch {
		case fs.autoIndex && os.IsNotExist(err):
			files, err := file.Readdir(-1)
			if err != nil {
				fs.handleError(rw, err)
				return
			}

			rw.Header().Add("Content-Type", "text/html")
			fmt.Fprintln(rw, "<!DOCTYPE html>", "<html>", "<body>", "<ul style=\"list-style: none;\">")

			for _, file := range files {
				fmt.Fprintf(rw, "<li><a href=\"%s\">%s</a></li>", file.Name(), file.Name())
			}

			fmt.Fprintln(rw, "</ul>", "</html>", "</body>")

		default:
			fs.handleError(rw, err)
			return
		}

		return
	}

	defer indexFile.Close()

	http.ServeContent(rw, req, indexFileInfo.Name(), indexFileInfo.ModTime(), indexFile)
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
