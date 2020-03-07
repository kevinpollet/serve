# serge <!-- omit in toc -->

[![Build Status](https://github.com/kevinpollet/serge/workflows/build/badge.svg)](https://github.com/kevinpollet/serge/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/kevinpollet/serge)](https://goreportcard.com/report/github.com/kevinpollet/serge)
[![GoDoc](https://godoc.org/github.com/kevinpollet/serge?status.svg)](https://pkg.go.dev/github.com/kevinpollet/serge)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org)
[![License](https://img.shields.io/github/license/kevinpollet/serge)](./LICENSE.md)

**serge** is a simple, secure and modern HTTP server, written in [Go](https://go.dev/), that can be used to serve static sites, files or single-page applications. You can use it through its command-line interface, [Docker](https://www.docker.com/) image or programmatically. As it's an `http.Handler` implementation, it should be easy to integrate it into your application. The following key features make **serge** unique and differentiable from the existing solutions and the `http.FileServer` implementation.

- Programmatic API.
- Custom error handler and pages.
- Basic HTTP authentication.
- Hide dot files by default.
- Directory listing is disabled by default. 
- Content encoding negotiation with support of [gzip](https://www.gzip.org/), [Deflate](https://en.wikipedia.org/wiki/DEFLATE) and [Brotli](https://en.wikipedia.org/wiki/Brotli) compression algorithms.

## Table of Contents <!-- omit in toc -->

- [Install](#install)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

## Install

> This project is actively developed and at its early stages. Expect some breaking changes before its first release.

```
go get -v github.com/kevinpollet/serge
```

## Examples

The [examples](./examples) directory contains the following examples:

- [hello](./examples/hello) — A simple static site that can be serve from the command-line.
- [docker](./examples/docker) — A simple static site that can be serve from a docker container.

## Contributing

Contributions are welcome!

Want to file a bug, request a feature or contribute some code?

1. Check out the [Code of Conduct](./CODE_OF_CONDUCT.md).
2. Check for an existing issue corresponding to your bug or feature request.
3. Open an issue to describe your bug or feature request.

## License

[MIT](./LICENSE.md)
