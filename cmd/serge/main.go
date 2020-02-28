package main

import (
	"flag"

	"github.com/kevinpollet/serge"
	"github.com/kevinpollet/serge/log"
)

var port int
var host, dir, cert, key string

func init() {
	flag.IntVar(&port, "port", serge.DefaultPort, "the port to serve")
	flag.StringVar(&host, "host", serge.DefaultHost, "the server host")
	flag.StringVar(&dir, "dir", serge.DefaultDir, "the directory to serve")
	flag.StringVar(&cert, "cert", "", "the TLS certificate")
	flag.StringVar(&key, "key", "", "the TLS key")
}

func main() {
	flag.Parse()

	server := serge.NewFileServer(
		serge.Host(host),
		serge.Port(port),
		serge.Dir(dir),
	)

	if len(cert) > 0 && len(key) > 0 {
		log.Logger().Fatal(server.ListenAndServeTLS(cert, key))
	}

	log.Logger().Fatal(server.ListenAndServe())
}
