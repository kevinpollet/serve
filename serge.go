package serge

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kevinpollet/serge/log"
)

const (
	DefaultDir  = "."
	DefaultHost = "127.0.0.1"
	DefaultPort = 8080
)

type dir string

func (d dir) Open(name string) (http.File, error) {
	fullName := filepath.Join(string(d), filepath.FromSlash(name))

	if strings.HasPrefix(filepath.Base(fullName), ".") {
		return nil, os.ErrNotExist
	}

	file, err := os.Open(fullName)
	if err != nil {
		return nil, err
	}

	return file, nil
}

type FileServer struct {
	host   string
	port   int
	dir    string
	server *http.Server
}

func NewFileServer(options ...fileServerOption) *FileServer {
	fs := &FileServer{host: DefaultHost, port: DefaultPort, dir: DefaultDir}

	for _, optionSetter := range options {
		optionSetter(fs)
	}

	fs.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", fs.host, fs.port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      http.FileServer(dir(fs.dir)),
	}

	return fs
}

func (fs *FileServer) ListenAndServe() error {
	log.Logger().Infof("server is listening on: %s", fs.server.Addr)
	return fs.server.ListenAndServe()
}

func (fs *FileServer) ListenAndServeTLS(certFile, keyFile string) error {
	log.Logger().Infof("TLS server is listening on: %s", fs.server.Addr)
	return fs.server.ListenAndServeTLS(certFile, keyFile)
}
