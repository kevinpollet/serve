package serge

import "github.com/justinas/alice"

type fileServerOption func(*FileServer)

// Host is a functional that defines the host for the server to listen on.
// nolint
func Host(host string) fileServerOption {
	return func(fs *FileServer) {
		fs.host = host
	}
}

// Port is a functional that defines the port for the server to listen on.
// nolint
func Port(port int) fileServerOption {
	return func(fs *FileServer) {
		fs.port = port
	}
}

// Dir is a functional that defines the directory path to serve.
// nolint
func Dir(dir string) fileServerOption {
	return func(fs *FileServer) {
		fs.dir = dir
	}
}

// Dir is a functional that defines the middlewares to run before serving the request.
// nolint
func Middlewares(middlewares ...alice.Constructor) fileServerOption {
	return func(fs *FileServer) {
		fs.middlewares = middlewares
	}
}
