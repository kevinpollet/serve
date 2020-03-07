package serge

import (
	"net/http"

	"github.com/justinas/alice"
)

func NewStripPrefixHandler(prefix string) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.StripPrefix(prefix, next)
	}
}
