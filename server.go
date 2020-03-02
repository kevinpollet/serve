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

type fileSystem string

func (fs fileSystem) Open(name string) (http.File, error) {
	rootDir := string(fs)
	fullName := filepath.Join(rootDir, filepath.FromSlash(name))

	if strings.HasPrefix(filepath.Base(fullName), ".") {
		return nil, os.ErrNotExist
	}

	file, err := os.Open(fullName)
	if err != nil {
		return nil, err
	}

	return file, nil
}

type fileServer struct {
	fileSystem http.FileSystem
}

func NewFileServer(dir string) http.Handler {
	return &fileServer{fileSystem(dir)}
}

func (server *fileServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	indexPageName := "index.html"
	urlPath := path.Clean(req.URL.Path)

	file, err := server.fileSystem.Open(urlPath)
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
			server.ServeHTTP(rw, req)
		}

		return
	}

	contentEncoding, err := negotiateContentEncoding(
		req,
		encodingBrotli, encodingGzip, encodingDeflate, encodingIdentity,
	)
	if err != nil {
		toResponse(rw, err)
		return
	}

	rw.Header().Add("Content-Encoding", contentEncoding)

	if contentEncoding != "" && contentEncoding != encodingIdentity {
		rwEncoder, err := newResponseWriterEncoder(contentEncoding, rw)
		if err != nil {
			toResponse(rw, err)
			return
		}

		rw = rwEncoder
		defer rwEncoder.Close()
	}

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
