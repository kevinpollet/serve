package serge

import (
	"net/http"

	"github.com/justinas/alice"
)

func BasicAuth(user, pass string) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			reqUser, reqPass, ok := req.BasicAuth()

			switch {
			case !ok:
				rw.Header().Add("WWW-Authenticate", "Basic realm=\"serge\"")
				rw.WriteHeader(http.StatusUnauthorized)

			case reqUser != user || reqPass != pass:
				rw.WriteHeader(http.StatusForbidden)

			default:
				next.ServeHTTP(rw, req)
			}
		})
	}
}
