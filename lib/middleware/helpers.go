package middleware

import (
	"errors"
	"net/http"
	"strings"
)

var (
	errInvalidSchema     = errors.New("unexpected authorization scheme provided")
	errMissingTokenValue = errors.New("missing bearer token value in header")
)

func getBearerToken(r *http.Request) (*string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("authorization header missing")
	}

	prefix := "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return nil, errInvalidSchema
	}

	authHeader = strings.TrimPrefix(authHeader, prefix)
	if len(authHeader) == 0 {
		return nil, errMissingTokenValue
	}

	authHeader = strings.TrimSpace(authHeader)
	return &authHeader, nil
}

func getRemoteIp(r *http.Request) string {
	ipAddress := r.Header.Get("X-Real-Ip")
	if ipAddress == "" {
		ipAddress = r.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}
	return ipAddress
}
