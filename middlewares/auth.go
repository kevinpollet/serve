package middlewares

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/justinas/alice"
	"golang.org/x/crypto/bcrypt"
)

type hashedCredentials map[string][]byte

func parseHashedCredentials(reader io.Reader) (hashedCredentials, error) {
	scanner := bufio.NewScanner(reader)
	bcryptPrefixRegexp := regexp.MustCompile(`\$2[aby]*\$`)
	credentials := make(hashedCredentials)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")

		if len(parts) != 2 || !bcryptPrefixRegexp.MatchString(parts[1]) {
			return nil, errors.New("unsupported password hash: only bcrypt is supported")
		}

		credentials[parts[0]] = []byte(parts[1])
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return credentials, nil
}

func (c hashedCredentials) match(user, password string) bool {
	hash, exists := c[user]
	return exists && bcrypt.CompareHashAndPassword(hash, []byte(password)) == nil
}

type basicAuthHandler struct {
	credentials hashedCredentials
	next        http.Handler
	realm       string
}

func NewBasicAuthHandler(realm string, reader io.Reader) (alice.Constructor, error) {
	credentials, err := parseHashedCredentials(reader)
	if err != nil {
		return nil, err
	}

	return func(next http.Handler) http.Handler {
		return &basicAuthHandler{credentials, next, realm}
	}, nil
}

func (h *basicAuthHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	user, password, hasAuth := req.BasicAuth()
	if !hasAuth {
		rw.Header().Add("WWW-Authenticate", fmt.Sprintf("Basic realm=\"%s\"", h.realm))
		rw.WriteHeader(http.StatusUnauthorized)

		return
	}

	if !h.credentials.match(user, password) {
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	h.next.ServeHTTP(rw, req)
}
