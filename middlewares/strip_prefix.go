package middlewares

import (
	"net/http"
	"strings"

	"github.com/justinas/alice"
)

type stripPrefixHandler struct {
	prefix string
	next   http.Handler
}

func NewStripPrefixHandler(prefix string) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return &stripPrefixHandler{prefix, next}
	}
}

func (h *stripPrefixHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	req.URL.Path = strings.TrimPrefix(req.URL.Path, h.prefix)
	h.next.ServeHTTP(rw, req)
}
