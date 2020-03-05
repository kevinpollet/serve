package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/kevinpollet/serge"
	"github.com/kevinpollet/serge/log"
)

const (
	defaultDir  = "."
	defaultAddr = "127.0.0.1:8080"
)

var (
	addr      string
	dir       string
	cert, key string
)

func init() {
	flag.StringVar(&addr, "addr", defaultAddr, "the server listening address")
	flag.StringVar(&dir, "dir", defaultDir, "the directory to serve")
	flag.StringVar(&cert, "cert", "", "the TLS certificate")
	flag.StringVar(&key, "key", "", "the TLS key")
}

func main() {
	flag.Parse()

	server := http.Server{
		Addr:         addr,
		Handler:      serge.NewFileServer(dir),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Logger().Printf("server is listening on: %s", server.Addr)

	if len(cert) > 0 && len(key) > 0 {
		log.Logger().Fatal(server.ListenAndServeTLS(cert, key))
	}

	log.Logger().Fatal(server.ListenAndServe())
}
