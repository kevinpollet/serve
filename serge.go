package serge

import (
	"fmt"
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/kevinpollet/serge/log"
)

type FileServer struct {
	host        string
	port        int
	dir         string
	middlewares []alice.Constructor
	server      *http.Server
}

func NewFileServer(options ...fileServerOption) *FileServer {
	fs := &FileServer{host: "127.0.0.1", port: 8080, dir: "."}

	for _, optionSetter := range options {
		optionSetter(fs)
	}

	fs.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", fs.host, fs.port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      alice.New(fs.middlewares...).Then(http.FileServer(http.Dir(fs.dir))),
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
