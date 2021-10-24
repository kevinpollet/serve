package srv

import (
	"net/http"

	"github.com/justinas/alice"
)

// Option is the functional option type.
type Option func(*fileServer)

// WithAutoIndex enables auto index.
func WithAutoIndex() Option {
	return func(fs *fileServer) {
		fs.autoIndex = true
	}
}

// WithMiddlewares sets the middlewares to apply before serving the request.
func WithMiddlewares(middlewares ...alice.Constructor) Option {
	return func(fs *fileServer) {
		fs.middlewares = middlewares
	}
}

// WithErrorHandler sets the Error handler.
func WithErrorHandler(errorHandler func(http.FileSystem, http.ResponseWriter, error)) Option {
	return func(fs *fileServer) {
		fs.errorHandler = errorHandler
	}
}
