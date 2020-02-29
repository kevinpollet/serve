package serge

import (
	"net/http"
	"strconv"
	"strings"
)

func parseAcceptEncoding(req *http.Request) (map[string]float32, error) {
	acceptEncodingValue := req.Header.Get("Accept-Encoding")
	acceptedEncodings := map[string]float32{"*": float32(1.0)}

	if len(acceptEncodingValue) == 0 {
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

func negotiateContentEncoding(req *http.Request, offers ...string) (string, error) {
	bestEncoding := ""
	bestEncodingWeight := float32(0.0)

	acceptedEncodings, err := parseAcceptEncoding(req)
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
