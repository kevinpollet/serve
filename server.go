package serve

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/justinas/alice"
)

const defaultAutoIndexTemplate = `
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

// FileServer is a http.Handler implementation that serves HTTP requests with the contents of the corresponding file.
type FileServer struct {
	autoIndex     bool
	autoIndexTmpl *template.Template
	fileSystem    dotFileHidingFileSystem
	middlewares   []alice.Constructor
}

// NewFileServer returns a new handler instance that serves HTTP requests with the contents of the given directory.
func NewFileServer(dir string, options ...Option) http.Handler {
	autoIndexTmpl := template.New("autoIndex")

	fs := &FileServer{
		fileSystem:    dotFileHidingFileSystem{http.Dir(dir)},
		autoIndexTmpl: template.Must(autoIndexTmpl.Parse(defaultAutoIndexTemplate)),
	}

	for _, option := range options {
		option(fs)
	}

	return alice.New(fs.middlewares...).Then(fs)
}

// ServeHTTP responds to an HTTP request.
func (fs *FileServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	p := path.Clean(req.URL.Path)

	file, fileInfo, err := fs.fileSystem.OpenWithStat(p)
	if err != nil {
		fs.handleError(rw, err)
		return
	}

	defer func() { _ = file.Close() }()

	fs.serveContent(rw, req, file, fileInfo)
}

func (fs *FileServer) handleError(rw http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError

	if os.IsNotExist(err) || os.IsPermission(err) {
		statusCode = http.StatusNotFound
	}

	rw.WriteHeader(statusCode)

	errorPageName := fmt.Sprintf("%d.html", statusCode)
	if file, err := fs.fileSystem.Open(errorPageName); err == nil {
		defer func() { _ = file.Close() }()

		_, _ = io.Copy(rw, file)
	}
}

func (fs *FileServer) serveContent(rw http.ResponseWriter, req *http.Request, file http.File, fileInfo os.FileInfo) {
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
	if err == nil {
		defer func() { _ = indexFile.Close() }()
		http.ServeContent(rw, req, indexFileInfo.Name(), indexFileInfo.ModTime(), indexFile)

		return
	}

	if !fs.autoIndex || !os.IsNotExist(err) {
		fs.handleError(rw, err)
		return
	}

	files, err := file.Readdir(-1)
	if err != nil {
		fs.handleError(rw, err)
		return
	}

	rw.Header().Add("Content-Type", "text/html")

	_ = fs.autoIndexTmpl.Execute(rw, files)
}

func redirectTo(rw http.ResponseWriter, req *http.Request, path string) {
	if query := req.URL.RawQuery; query != "" {
		path += fmt.Sprint("?", query)
	}

	rw.Header().Add("Location", path)
	rw.WriteHeader(http.StatusMovedPermanently)
}
