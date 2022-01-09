# serve

[![Build Status](https://github.com/kevinpollet/serve/workflows/build/badge.svg)](https://github.com/kevinpollet/serve/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/kevinpollet/serve)](https://goreportcard.com/report/github.com/kevinpollet/serve)
[![GoDoc](https://godoc.org/github.com/kevinpollet/serve?status.svg)](https://pkg.go.dev/github.com/kevinpollet/serve)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org)
[![License](https://img.shields.io/github/license/kevinpollet/serve)](./LICENSE)

Simple, secure [Go](https://go.dev/) `HTTP` server to serve static sites, single-page applications or files. It can be
used through its command-line interface, [Docker](https://www.docker.com/) image or programmatically. As it's
an `http.Handler` implementation, it can be used as a drop-in replacement of the `http.FileServer` from the standard
library and can be integrated easily.

- Programmatic API.
- HTTP/2 and TLS support.
- Custom error pages.
- Basic HTTP authentication.
- Hide dot files by default.
- Directory listing is disabled by default.
- Encoding negotiation with support of [gzip](https://www.gzip.org/), [Deflate](https://en.wikipedia.org/wiki/DEFLATE)
  and [Brotli](https://en.wikipedia.org/wiki/Brotli) compression algorithms.

## Installation

```shell
go get github.com/kevinpollet/serve                 # get dependency
go install github.com/kevinpollet/serve/cmd/serve   # build and install command-line interface bin
```

## Usage

### Command-line

The following text is the output of the `serve -help` command.

```shell
Usage: serve [options]

Options:
-addr       The server address, "127.0.0.1:8080" by default.
-auth       The basic auth credentials (password must be hashed with bcrypt and escaped with '').
-auth-file  The basic auth credentials following the ".htpasswd" format.
-dir        The directory containing the files to serve, "." by default.
-cert       The TLS certificate.
-key        The TLS private key.
-help       Prints this text.
```

### Library

```go
package main

import (
	"log"
	"net/http"

	"github.com/kevinpollet/serve"
	"github.com/kevinpollet/serve/middlewares"
)

func main() {
	http.Handle("/static", serve.NewFileServer("examples/hello",
		serve.WithAutoIndex(),
		serve.WithMiddlewares(middlewares.NewStripPrefixHandler("/static")),
	))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Docker

An official docker image is available on [Docker Hub](https://hub.docker.com/r/kevinpollet/serve). The
following `Dockerfile` shows how to use the provided base image to serve your static sites or files through a running
Docker container. By default, the base image will serve all files available in the `/var/www/` directory and listen for
TCP connections on `8080`.

```
FROM kevinpollet/serve:latest
COPY . /var/www/
```

Then, you can build and run your Docker image with the following commands. Your static site or files will be available
on http://localhost:8080.

```shell
docker build . -t moby:latest
docker run -d -p 8080:8080 moby:latest
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

[MIT](./LICENSE)
