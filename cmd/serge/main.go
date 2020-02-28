package main

import (
	"flag"
	"net/http"

	"github.com/kevinpollet/serge"
	"github.com/kevinpollet/serge/log"
)

func main() {
	hostPtr := flag.String("host", serge.DefaultHost, "the server host")
	portPtr := flag.Int("port", serge.DefaultPort, "the server port")
	dirPtr := flag.String("dir", serge.DefaultDir, "the directory to serve")

	flag.Parse()

	server := serge.NewFileServer(
		serge.Host(*hostPtr),
		serge.Port(*portPtr),
		serge.Dir(*dirPtr),
	)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Logger().Fatal(err)
	}
}
