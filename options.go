package serge

import (
	"net/http"

	"github.com/justinas/alice"
)

type fileServerOption func(*fileServer)

func Middlewares(middlewares ...alice.Constructor) fileServerOption { // nolint
	return func(fs *fileServer) {
		fs.middlewares = middlewares
	}
}

func ErrorHandler(errorHandler func(http.FileSystem, http.ResponseWriter, error)) fileServerOption { //nolint
	return func(fs *fileServer) {
		fs.errorHandler = errorHandler
	}
}
