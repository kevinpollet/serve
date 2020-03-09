package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/justinas/alice"
	"github.com/kevinpollet/serge"
	"github.com/kevinpollet/serge/log"
	mddlw "github.com/kevinpollet/serge/middlewares"
)

var (
	flagAddr = flag.String("addr", "127.0.0.1:8080", "")
	flagAuth = flag.String("auth", "", "")
	flagDir  = flag.String("dir", ".", "")
	flagCert = flag.String("cert", "", "")
	flagKey  = flag.String("key", "", "")
)

const helpText = `serge [options]

Options:
-addr    The server address, "127.0.0.1:8080" by default.
-auth    The basic auth credentials (password must be hashed with bcrypt and escaped with ''), e.g. hello:'$2y$12$...'
-dir     The directory containing the files to serve, "." by default.
-cert    The TLS certificate.
-key     The TLS private key.
-help    Prints this text.
`

func main() {
	middlewares := make([]alice.Constructor, 0)

	flag.Usage = help
	flag.Parse()

	if len(*flagAuth) > 0 {
		reader := strings.NewReader(*flagAuth)
		basicAuthHandler, err := mddlw.NewBasicAuthHandler(reader)

		if err != nil {
			log.Logger().Fatal(err)
		}

		middlewares = append(middlewares, basicAuthHandler)
	}

	middlewares = append(middlewares, mddlw.NewNegotiateEncodingHandler(
		mddlw.EncodingBrotli,
		mddlw.EncodingGzip,
		mddlw.EncodingDeflate,
	))

	server := http.Server{
		Addr:         *flagAddr,
		Handler:      serge.NewFileServer(*flagDir, serge.WithMiddlewares(middlewares...)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Logger().Printf("server is listening on: %s", server.Addr)

	if len(*flagCert) > 0 && len(*flagKey) > 0 {
		log.Logger().Fatal(server.ListenAndServeTLS(*flagCert, *flagKey))
	}

	log.Logger().Fatal(server.ListenAndServe())
}

func help() {
	fmt.Println(helpText)
}
