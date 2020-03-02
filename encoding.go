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

type acceptEncoding map[string]float64

func (m acceptEncoding) qvalue(encoding string) (float64, bool) {
	qvalue, exists := m[encoding]
	if !exists {
		qvalue, exists = m["*"]
	}

	return qvalue, exists
}

func negotiateContentEncoding(req *http.Request, contentEncodings ...string) (string, error) {
	bestEncoding := ""
	bestQvalue := 0.0

	acceptEncoding, err := parseAcceptEncoding(req)
	if err != nil {
		return "", err
	}

	for _, contentEncoding := range contentEncodings {
		qvalue, exists := acceptEncoding.qvalue(contentEncoding)

		if exists && qvalue > bestQvalue {
			bestEncoding = contentEncoding
			bestQvalue = qvalue
		}
	}

	if bestEncoding == "" {
		qvalue, exists := acceptEncoding.qvalue(encodingIdentity)
		if qvalue != 0.0 || !exists {
			bestEncoding = encodingIdentity
		}
	}

	return bestEncoding, nil
}

func parseAcceptEncoding(req *http.Request) (acceptEncoding, error) {
	encodings := make(map[string]float64)
	acceptEncoding := req.Header.Get("Accept-Encoding")

	if acceptEncoding == "" {
		return encodings, nil
	}

	for _, encoding := range strings.Split(acceptEncoding, ",") {
		qValue := 1.0
		trimedEncoding := strings.TrimSpace(encoding)
		parts := strings.Split(trimedEncoding, ";q=")

		if len(parts) == 2 {
			parsedQValue, err := strconv.ParseFloat(parts[1], 32)
			if err != nil {
				return nil, err
			}

			qValue = parsedQValue
		}

		encodings[parts[0]] = qValue
	}

	return encodings, nil
}
