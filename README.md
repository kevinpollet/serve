# srv <!-- omit in toc -->

[![Build Status](https://github.com/kevinpollet/srv/workflows/build/badge.svg)](https://github.com/kevinpollet/srv/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/kevinpollet/srv)](https://goreportcard.com/report/github.com/kevinpollet/srv)
[![GoDoc](https://godoc.org/github.com/kevinpollet/srv?status.svg)](https://pkg.go.dev/github.com/kevinpollet/srv)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org)
[![License](https://img.shields.io/github/license/kevinpollet/srv)](./LICENSE.md)

**srv** is a simple, secure and modern `HTTP` server, written in [Go](https://go.dev/), to serve static sites,
single-page applications or a file with ease. You can use it through its command-line
interface, [Docker](https://www.docker.com/) image or programmatically. As it's an `http.Handler` implementation, it
should be easy to integrate it into your application. The following key features make **srv** unique and differentiable
from the existing solutions and the `http.FileServer` implementation.

- Programmatic API.
- HTTP/2 and TLS support.
- Custom error handler and pages.
- Basic HTTP authentication.
- Hide dot files by default.
- Directory listing is disabled by default.
- Encoding negotiation with support of [gzip](https://www.gzip.org/), [Deflate](https://en.wikipedia.org/wiki/DEFLATE)
  and [Brotli](https://en.wikipedia.org/wiki/Brotli) compression algorithms.

## Install

```shell
go get github.com/kevinpollet/srv                  # get dependency
go install -i github.com/kevinpollet/srv/cmd/srv   # build and install command-line interface bin
```

## Usage

### Command-line <!-- omit in toc -->

**srv** can be use through its provided command-line. The following text is the output of the `srv -help` command.

```shell
srv [options]

Options:
-addr      The server address, "127.0.0.1:8080" by default.
-auth      The basic auth credentials (password must be hashed with bcrypt and escaped with '').
-authfile  The basic auth credentials file following the ".htpasswd" format.
-dir       The directory containing the files to serve, "." by default.
-cert      The TLS certificate.
-key       The TLS private key.
-help      Prints this text.
```

### Docker <!-- omit in toc -->

An official docker image is available on [Docker Hub](https://hub.docker.com/r/kevinpollet/srv). The
following `Dockerfile` shows how to use the provided base image to serve your static sites or files through a running
Docker container. By default, the base image will serve all files available in the `/var/www/` directory and listen for
TCP connections on `8080`.

```
FROM kevinpollet/srv:latest
COPY . /var/www/
```

Then, you can build and run your Docker image with the following commands. Your static site or files will be available
on http://localhost:8080.

```shell
docker build . -t moby:latest
docker run -d -p 8080:8080 moby:latest
```

### API <!-- omit in toc -->

```go
package main

import (
	"log"
	"net/http"

	"github.com/kevinpollet/srv"
	"github.com/kevinpollet/srv/middlewares"
)

func main() {
	customErrorHandler := func(fs http.FileSystem, rw http.ResponseWriter, err error) {
			log.Print(err)
			rw.WriteHeader(http.StatusInternalServerError)
	}

	http.Handle("/static", srv.NewFileServer("examples/hello",
		srv.WithAutoIndex(),
		srv.WithMiddlewares(middlewares.NewStripPrefixHandler("/static")),
		srv.WithErrorHandler(customErrorHandler),
	))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Examples

The [examples](./examples) directory contains the following examples:

- [hello](./examples/hello) — A simple static site that can be served from the command-line.
- [docker](./examples/docker) — A simple static site that can be served from a docker container.

## Contributing

Contributions are welcome!

Want to file a bug, request a feature or contribute some code?

1. Check out the [Code of Conduct](./CODE_OF_CONDUCT.md).
2. Check for an existing issue corresponding to your bug or feature request.
3. Open an issue to describe your bug or feature request.

## License

[MIT](./LICENSE.md)
