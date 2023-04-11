// Package httputil implements HTTP utility functions for working with HTTP requests and responses.
package httputil

import (
	"net/http"
	"strconv"
	"strings"
)

// IsBodySizeWithinLimit checks if the Content-Length of the response is within
// the specified limit. If the limit is zero or negative, it always returns
// true.
func IsBodySizeWithinLimit(header http.Header, limit int64) bool {
	if limit <= 0 {
		return true
	}

	value := header.Get("Content-Length")

	contentLength, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return true // Assume the size is within limits if we can't parse it.
	}

	return contentLength <= limit
}

// MaxAge returns the max-age value in seconds if found in the Cache-Control
// header. It returns -1 if the header is not present or if the value is not a
// valid number.
func MaxAge(header http.Header) int {
	if header == nil {
		return -1
	}

	value := header.Get("Cache-Control")

	if strings.Contains(value, "max-age") {
		parts := strings.Split(value, ",")

		for _, part := range parts {
			part = strings.TrimSpace(part)

			if strings.HasPrefix(part, "max-age") {
				ageStr := strings.TrimPrefix(part, "max-age=")

				maxAge, err := strconv.Atoi(ageStr)
				if err == nil {
					return maxAge
				}
			}
		}
	}

	return -1
}
