package serge

import (
	"net/http"

	"github.com/justinas/alice"
)

type Option func(*fileServer)

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
