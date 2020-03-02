package serge

import (
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/andybalholm/brotli"
)

const (
	encodingBrotli   = "br"
	encodingDeflate  = "deflate"
	encodingGzip     = "gzip"
	encodingIdentity = "identity"
)

type responseWriterEncoder struct {
	io.WriteCloser
	http.ResponseWriter
}

func newResponseWriterEncoder(encoding string, rw http.ResponseWriter) (*responseWriterEncoder, error) {
	switch encoding {
	case encodingBrotli:
		return &responseWriterEncoder{brotli.NewWriter(rw), rw}, nil

	case encodingGzip:
		return &responseWriterEncoder{gzip.NewWriter(rw), rw}, nil

	case encodingDeflate:
		writer, err := flate.NewWriter(rw, flate.DefaultCompression)
		if err != nil {
			return nil, err
		}

		return &responseWriterEncoder{writer, rw}, nil
	}

	return nil, fmt.Errorf("unsupported encoding: %s", encoding)
}

func (rw *responseWriterEncoder) Write(bytes []byte) (int, error) {
	return rw.WriteCloser.Write(bytes)
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
	acceptEncodingValue := req.Header.Get("Accept-Encoding")

	if len(acceptEncodingValue) == 0 {
		acceptedEncodings["*"] = float32(1.0)
		return acceptedEncodings, nil
	}

	for _, acceptedEncoding := range strings.Split(acceptEncodingValue, ",") {
		q := 1.0
		parts := strings.Split(strings.TrimSpace(acceptedEncoding), ";q=")

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
