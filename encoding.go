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

func negotiateContentEncoding(req *http.Request, contentEncodings ...string) (string, error) {
	bestEncoding := ""
	bestEncodingQvalue := 0.0

	acceptEncodings, err := parseAcceptEncoding(req)
	if err != nil {
		return "", err
	}

	for _, contentEncoding := range contentEncodings {
		qvalue, exists := acceptEncodings[contentEncoding]
		if !exists {
			qvalue = acceptEncodings["*"]
		}

		if qvalue > bestEncodingQvalue {
			bestEncoding = contentEncoding
			bestEncodingQvalue = qvalue
		}
	}

	return bestEncoding, nil
}

func parseAcceptEncoding(req *http.Request) (map[string]float64, error) {
	encodings := make(map[string]float64)
	acceptEncoding := req.Header.Get("Accept-Encoding")

	if acceptEncoding == "" {
		encodings[encodingIdentity] = 1.0
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
