package serve

import "github.com/justinas/alice"

// Option is the functional option type.
type Option func(*FileServer)

// WithAutoIndex enables auto index.
func WithAutoIndex() Option {
	return func(fs *FileServer) {
		fs.autoIndex = true
	}
}

// WithMiddlewares sets the middlewares to apply before serving the request.
func WithMiddlewares(middlewares ...alice.Constructor) Option {
	return func(fs *FileServer) {
		fs.middlewares = middlewares
	}
}
