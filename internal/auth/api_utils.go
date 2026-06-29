package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {

	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", errors.New("Missing header")
	}

	return strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer ")), nil
}
