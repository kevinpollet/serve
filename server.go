package srv

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/justinas/alice"
	"github.com/kevinpollet/srv/log"
)

const autoIndexTemplate = `
<!DOCTYPE html>
<html>
<body>
<ul style="list-style: none;">
{{range $file := .}}
<li><a href="{{$file.Name}}">{{$file.Name}}</a></li>
{{end}}
</ul>
</body>
</html>
`

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
	path := path.Clean(req.URL.Path)

	file, fileInfo, err := fs.fileSystem.OpenWithStat(path)
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

			t, err := template.New("autoIndex").Parse(autoIndexTemplate)
			if err != nil {
				fs.handleError(rw, err)
				return
			}

			t.Execute(rw, files) // nolint

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
