package srv

import (
	"net/http"

	"github.com/justinas/alice"
)

type Option func(*fileServer)

// WithAutoIndex is a functional option that enables auto index.
func WithAutoIndex() Option {
	return func(fs *fileServer) {
		fs.autoIndex = true
	}
}

// WithMiddlewares is a functional option that sets the middlewares to apply before serving the request.
func WithMiddlewares(middlewares ...alice.Constructor) Option {
	return func(fs *fileServer) {
		fs.middlewares = middlewares
	}
}

// WithErrorHandler is a functional option that sets the Error handler.
func WithErrorHandler(errorHandler func(http.FileSystem, http.ResponseWriter, error)) Option {
	return func(fs *fileServer) {
		fs.errorHandler = errorHandler
	}
}
