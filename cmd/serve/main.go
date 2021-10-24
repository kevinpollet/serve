package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/justinas/alice"
	"github.com/kevinpollet/serve"
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

const helpText = `Usage: serve [options]

Options:
-addr       The server address, "127.0.0.1:8080" by default.
-auth       The basic auth credentials (password must be hashed with bcrypt and escaped with '').
-auth-file  The basic auth credentials following the ".htpasswd" format.
-dir        The directory containing the files to serve, "." by default.
-cert       The TLS certificate.
-key        The TLS private key.
-help       Prints this text.
`

func init() {
	flag.Usage = help
	flag.Parse()
}

func main() {
	var handlers []alice.Constructor

	switch {
	case len(*flagAuth) > 0:
		reader := strings.NewReader(*flagAuth)

		basicAuthHandler, err := middlewares.NewBasicAuthHandler("serve", reader)
		if err != nil {
			exitWithError(err)
		}

		handlers = append(handlers, basicAuthHandler)

	case len(*flagAuthFile) > 0:
		file, err := os.Open(*flagAuthFile)
		if err != nil {
			exitWithError(err)
		}

		defer func() { _ = file.Close() }()

		basicAuthHandler, err := middlewares.NewBasicAuthHandler("serve", file)
		if err != nil {
			exitWithError(err)
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
		Handler:      serve.NewFileServer(*flagDir, serve.WithMiddlewares(handlers...)),
	}

	log.Logger().Printf("server is listening on: %s", server.Addr)

	if len(*flagCert) > 0 && len(*flagKey) > 0 {
		log.Logger().Fatal(server.ListenAndServeTLS(*flagCert, *flagKey))
	}

	log.Logger().Fatal(server.ListenAndServe())
}

func exitWithError(err error) {
	log.Logger().Error(err)
	os.Exit(1)
}

func help() {
	fmt.Print(helpText) // nolint
	os.Exit(2)
}
