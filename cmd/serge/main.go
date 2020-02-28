package main

import (
	"net/http"

	"github.com/kevinpollet/serge"
	"github.com/kevinpollet/serge/log"
)

func main() {
	server := serge.NewFileServer(
		serge.Host("0.0.0.0"),
		serge.Port(8080),
		serge.Dir("examples/hello"),
		serge.Middlewares(serge.BasicAuth("user", "pass")),
	)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Logger().Fatal(err)
	}
}
