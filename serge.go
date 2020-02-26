package serge

import (
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/kevinpollet/serge/log"
)

type FileServer struct {
	addr        string
	dir         string
	middlewares []alice.Constructor
	server      *http.Server
}

func NewFileServer(options ...fileServerOption) *FileServer {
	fs := &FileServer{addr: "127.0.0.1:8080", dir: "."}

	for _, optionSetter := range options {
		optionSetter(fs)
	}

	fs.server = &http.Server{
		Addr:         fs.addr,
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
