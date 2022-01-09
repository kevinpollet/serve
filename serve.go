package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/justinas/alice"
	"github.com/kevinpollet/serve/log"
	"github.com/kevinpollet/serve/middlewares"
)

var (
	flagAddr     = flag.String("addr", "127.0.0.1:8080", "")
	flagAuth     = flag.String("auth", "", "")
	flagAuthFile = flag.String("auth-file", "", "")
	flagDir      = flag.String("dir", ".", "")
	flagCert     = flag.String("cert", "", "")
	flagKey      = flag.String("key", "", "")
)

const usage = `Usage: serve [options]

Options:
-addr       Sets the server address. Default is "127.0.0.1:8080".
-auth       Sets the basic auth credentials (password must be hashed with bcrypt and escaped with '').
-auth-file  Sets the basic auth credentials following the ".htpasswd" format.
-dir        Sets the directory containing the files to serve. Default is ".".
-cert       Sets the TLS certificate.
-key        Sets the TLS private key.
-help       Prints this text.
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}

	flag.Parse()

	var handlers []alice.Constructor

	switch {
	case len(*flagAuth) > 0:
		reader := strings.NewReader(*flagAuth)

		basicAuthHandler, err := middlewares.NewBasicAuthHandler("serve", reader)
		if err != nil {
			errExit(err)
		}

		handlers = append(handlers, basicAuthHandler)

	case len(*flagAuthFile) > 0:
		file, err := os.Open(*flagAuthFile)
		if err != nil {
			errExit(err)
		}

		defer func() { _ = file.Close() }()

		basicAuthHandler, err := middlewares.NewBasicAuthHandler("serve", file)
		if err != nil {
			errExit(err)
		}

		handlers = append(handlers, basicAuthHandler)
	}

	handlers = append(handlers, middlewares.NewNegotiateEncodingHandler(
		middlewares.EncodingBrotli,
		middlewares.EncodingGzip,
		middlewares.EncodingDeflate,
	))

	server := http.Server{
		Addr:         *flagAddr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      NewFileServer(*flagDir, WithMiddlewares(handlers...)),
	}

	log.Logger().Printf("server is listening on: %s", server.Addr)

	if len(*flagCert) > 0 && len(*flagKey) > 0 {
		log.Logger().Fatal(server.ListenAndServeTLS(*flagCert, *flagKey))
	}

	log.Logger().Fatal(server.ListenAndServe())
}

func errExit(err error) {
	log.Logger().Error(err)
	os.Exit(1)
}
