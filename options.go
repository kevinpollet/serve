package serge

import "github.com/justinas/alice"

type fileServerOption func(*FileServer)

func Address(address string) fileServerOption {
	return func(fs *FileServer) {
		fs.addr = address
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
