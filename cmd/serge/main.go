package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/kevinpollet/serge"
	"github.com/kevinpollet/serge/log"
)

const (
	defaultDir  = "."
	defaultHost = "127.0.0.1"
	defaultPort = 8080
)

var (
	port                 int
	host, dir, cert, key string
)

func init() {
	flag.IntVar(&port, "port", defaultPort, "the port to serve")
	flag.StringVar(&host, "host", defaultHost, "the server host")
	flag.StringVar(&dir, "dir", defaultDir, "the directory to serve")
	flag.StringVar(&cert, "cert", "", "the TLS certificate")
	flag.StringVar(&key, "key", "", "the TLS key")
}

func main() {
	flag.Parse()

	server := http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, port),
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
