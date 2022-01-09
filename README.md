# serve

[![Build Status](https://github.com/kevinpollet/serve/workflows/build/badge.svg)](https://github.com/kevinpollet/serve/actions)

Simple and secure [Go](https://go.dev/) `HTTP` server to serve static sites or files from the command-line.

- HTTP/2 and TLS support.
- Custom Error pages.
- Basic HTTP authentication.
- Hide dot files by default.
- Directory listing is disabled by default.
- Encoding negotiation with support of [gzip](https://www.gzip.org/), [Deflate](https://en.wikipedia.org/wiki/DEFLATE)
  and [Brotli](https://en.wikipedia.org/wiki/Brotli) compression algorithms.

## Installation

```shell
go install github.com/kevinpollet/serve
```

## Usage

```shell
Usage: serve [options]

Options:
-addr       Sets the server address. Default is "127.0.0.1:8080".
-auth       Sets the basic auth credentials (password must be hashed with bcrypt and escaped with '').
-auth-file  Sets the basic auth credentials following the ".htpasswd" format.
-dir        Sets the directory containing the files to serve. Default is ".".
-cert       Sets the TLS certificate.
-key        Sets the TLS private key.
-help       Prints this text.
```

### Docker

A Docker [image](https://hub.docker.com/r/kevinpollet/serve) is available to serve static files from a running Docker
container. By default, all files located in the `/var/www/` directory will be made available through TCP connections on
port `8080`. For more details, check out the Docker [example](./examples/docker).

## Examples

- [hello](./examples/hello) — Simple static site that can be served from the command-line.
- [docker](./examples/docker) — Simple static site that can be served from a Docker container.

## Contributing

PRs welcome!

Want to file a bug or request a feature?

1. Check out the [Code of Conduct](./CODE_OF_CONDUCT.md).
2. Check for an existing issue corresponding to your bug or feature request.
3. Open an issue to describe your bug or feature request.

## License

[MIT](./LICENSE)
