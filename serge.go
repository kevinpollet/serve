package serge

import (
	"compress/flate"
	"compress/gzip"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/kevinpollet/serge/log"
)

const headerLocation = "Location"

type fileServer struct {
	root http.FileSystem
}

func NewFileServer(dir string) http.Handler {
	return &fileServer{root: http.Dir(dir)}
}

func (fs *fileServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	urlPath := path.Clean(req.URL.Path)

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
			rw.Header().Add(headerLocation, urlPath+"/")
			rw.WriteHeader(http.StatusMovedPermanently)
		} else {
			req.URL.Path = urlPath + "/index.html"
			fs.ServeHTTP(rw, req)
		}

		return
	}

	if strings.HasPrefix(fileInfo.Name(), ".") {
		toHTTPResponse(rw, os.ErrNotExist)
		return
	}

	contentEncoding, err := negotiateContentEncoding(
		req,
		encodingGzip, encodingDeflate, encodingIdentity,
	)
	if err != nil {
		toHTTPResponse(rw, err)
		return
	}

	rw.Header().Add(headerContentEncoding, contentEncoding)

	switch contentEncoding {
	case encodingGzip:
		gzipWriter := gzip.NewWriter(rw)
		defer gzipWriter.Close()

		rw = &encodedResponseWriter{gzipWriter, rw}

	case encodingDeflate:
		flateWriter, _ := flate.NewWriter(rw, flate.DefaultCompression)
		defer flateWriter.Close()

		rw = &encodedResponseWriter{flateWriter, rw}
	}

	http.ServeContent(rw, req, fileInfo.Name(), fileInfo.ModTime(), file)
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
