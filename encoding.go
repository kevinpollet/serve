package serge

import (
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	encodingGzip     = "gzip"
	encodingDeflate  = "deflate"
	encodingIdentity = "identity"

	headerAcceptEncoding  = "Accept-Encoding"
	headerContentEncoding = "Content-Encoding"
)

type encodedResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (rw *encodedResponseWriter) Write(bytes []byte) (int, error) {
	return rw.Writer.Write(bytes)
}

func negotiateContentEncoding(req *http.Request, offers ...string) (string, error) {
	bestEncoding := ""
	bestEncodingWeight := float32(0.0)

	acceptedEncodings, err := parseAcceptEncodingHeader(req)
	if err != nil {
		return "", err
	}

	for _, offer := range offers {
		weight, exists := acceptedEncodings[offer]
		if !exists {
			weight, exists = acceptedEncodings["*"]
		}

		if exists && weight > bestEncodingWeight {
			bestEncoding = offer
			bestEncodingWeight = weight
		}
	}

	return bestEncoding, nil
}

func parseAcceptEncodingHeader(req *http.Request) (map[string]float32, error) {
	acceptedEncodings := make(map[string]float32)
	acceptEncodingValue := req.Header.Get(headerAcceptEncoding)

	if len(acceptEncodingValue) == 0 {
		acceptedEncodings["*"] = float32(1.0)
		return acceptedEncodings, nil
	}

	for _, acceptedEncoding := range strings.Split(acceptEncodingValue, ",") {
		q := 1.0
		parts := strings.Split(acceptedEncoding, ";q=")

		if len(parts) == 2 {
			parsedQ, err := strconv.ParseFloat(parts[1], 32)
			if err != nil {
				return acceptedEncodings, err
			}

			q = parsedQ
		}

		acceptedEncodings[parts[0]] = float32(q)
	}

	return acceptedEncodings, nil
}
