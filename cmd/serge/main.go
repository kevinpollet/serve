package main

import (
	"flag"
	"net/http"
	"strings"
	"time"

	"github.com/justinas/alice"
	"github.com/kevinpollet/serge"
	"github.com/kevinpollet/serge/log"
)

const (
	defaultDir  = "."
	defaultAddr = "127.0.0.1:8080"
)

func main() {
	middlewares := make([]alice.Constructor, 0)
	addr := flag.String("addr", defaultAddr, "the server listening address")
	auth := flag.String("auth", "", "the basic auth credentials")
	dir := flag.String("dir", defaultDir, "the directory to serve")
	cert := flag.String("cert", "", "the TLS certificate")
	key := flag.String("key", "", "the TLS key")

	flag.Parse()

	if len(*auth) > 0 {
		reader := strings.NewReader(*auth)

		basicAuthHandler, err := serge.NewBasicAuthHandler(reader)
		if err != nil {
			log.Logger().Fatal(err)
		}

		middlewares = append(middlewares, basicAuthHandler)
	}

	server := http.Server{
		Addr:         *addr,
		Handler:      serge.NewFileServer(*dir, serge.Middlewares(middlewares...)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Logger().Printf("server is listening on: %s", server.Addr)

	if len(*cert) > 0 && len(*key) > 0 {
		log.Logger().Fatal(server.ListenAndServeTLS(*cert, *key))
	}

	log.Logger().Fatal(server.ListenAndServe())
}
