package middlewares

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/justinas/alice"
)

const (
	EncodingBrotli  = "br"
	EncodingDeflate = "deflate"
	EncodingGzip    = "gzip"

	encodingIdentity = "identity"
)

type acceptEncoding map[string]float64

func (m acceptEncoding) qvalue(encoding string) (float64, bool) {
	qvalue, exists := m[encoding]
	if !exists {
		qvalue, exists = m["*"]
	}

	return qvalue, exists
}

type responseWriterEncoder struct {
	io.WriteCloser
	http.ResponseWriter
}

func newResponseWriterEncoder(encoding string, rw http.ResponseWriter) *responseWriterEncoder {
	switch encoding {
	case EncodingBrotli:
		return &responseWriterEncoder{brotli.NewWriter(rw), rw}

	case EncodingGzip:
		return &responseWriterEncoder{gzip.NewWriter(rw), rw}

	case EncodingDeflate:
		writer, _ := flate.NewWriter(rw, flate.DefaultCompression)
		return &responseWriterEncoder{writer, rw}
	}

	return nil
}

func (rw *responseWriterEncoder) Write(bytes []byte) (int, error) {
	return rw.WriteCloser.Write(bytes)
}

type encodingHandler struct {
	encodings []string
	next      http.Handler
}

func NewNegotiateEncodingHandler(encodings ...string) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return &encodingHandler{encodings, next}
	}
}

func (h *encodingHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	bestQvalue := 0.0
	bestEncoding := ""

	acceptEncoding, err := parseAcceptEncoding(req)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, encoding := range h.encodings {
		if qvalue, exists := acceptEncoding.qvalue(encoding); exists && qvalue > bestQvalue {
			bestEncoding = encoding
			bestQvalue = qvalue
		}
	}

	if bestEncoding != "" {
		rwEncoder := newResponseWriterEncoder(bestEncoding, rw)
		defer rwEncoder.Close()

		rw.Header().Add("Content-Encoding", bestEncoding)
		h.next.ServeHTTP(rwEncoder, req)

		return
	}

	qvalue, exists := acceptEncoding.qvalue(encodingIdentity)
	if exists && qvalue == 0.0 {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	rw.Header().Add("Content-Encoding", encodingIdentity)
	h.next.ServeHTTP(rw, req)
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
