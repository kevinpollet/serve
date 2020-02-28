package serge

import "github.com/justinas/alice"

type fileServerOption func(*FileServer)

func Host(host string) fileServerOption { // nolint
	return func(fs *FileServer) {
		fs.host = host
	}
}

func Port(port int) fileServerOption { // nolint
	return func(fs *FileServer) {
		fs.port = port
	}
}

func Dir(dir string) fileServerOption { // nolint
	return func(fs *FileServer) {
		fs.dir = dir
	}
}

func Middlewares(middlewares ...alice.Constructor) fileServerOption { // nolint
	return func(fs *FileServer) {
		fs.middlewares = middlewares
	}
}
