package serge

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type basicAuthHandler struct {
	credentials map[string][]byte
	next        http.Handler
}

func NewBasicAuthHandler(reader io.Reader, next http.Handler) (http.Handler, error) {
	credentials, err := parseCredentials(reader)
	if err != nil {
		return nil, err
	}

	return &basicAuthHandler{credentials, next}, nil
}

func (h *basicAuthHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	user, pass, hasAuth := req.BasicAuth()

	if !hasAuth {
		rw.Header().Add("WWW-Authenticate", "Basic realm=\"serge\"")
		rw.WriteHeader(http.StatusUnauthorized)

		return
	}

	if !h.authenticate(user, pass) {
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	h.next.ServeHTTP(rw, req)
}

func (h *basicAuthHandler) authenticate(user, pass string) bool {
	passwordHash, userExists := h.credentials[user]
	return userExists && bcrypt.CompareHashAndPassword(passwordHash, []byte(pass)) == nil
}

func parseCredentials(reader io.Reader) (map[string][]byte, error) {
	scanner := bufio.NewScanner(reader)
	credentialPrefixRegexp := regexp.MustCompile(`\$2[axy]*\$`)
	credentials := make(map[string][]byte)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")

		if len(parts) != 2 || !credentialPrefixRegexp.MatchString(parts[1]) {
			return nil, fmt.Errorf("unsupported password encoding")
		}

		credentials[parts[0]] = []byte(parts[1])
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return credentials, nil
}
