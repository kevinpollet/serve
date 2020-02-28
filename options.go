package serge

import "github.com/justinas/alice"

type fileServerOption func(*FileServer)

func Host(host string) fileServerOption {
	return func(fs *FileServer) {
		fs.host = host
	}
}

func Port(port int) fileServerOption {
	return func(fs *FileServer) {
		fs.port = port
	}
}

func Dir(dir string) fileServerOption {
	return func(fs *FileServer) {
		fs.dir = dir
	}
}

func Middlewares(middlewares ...alice.Constructor) fileServerOption {
	return func(fs *FileServer) {
		fs.middlewares = middlewares
	}
}
