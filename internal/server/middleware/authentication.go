package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
)

var contextKeyBearerToken ContextKey = "BEARER_TOKEN"

func BearerToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := parseBearerToken(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			slog.Error("bearer token error", "error", err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyBearerToken, *token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

var (
	errMissingAuthHeader              = errors.New("missing authorization header")
	errUnsupportedAuthorizationSchema = errors.New("unsupported authorization scheme")
	errMissingBearerToken             = errors.New("missing bearer token")
)

func parseBearerToken(r *http.Request) (*string, error) {
	authHeader := r.Header.Get("Authorization")
	authHeader = strings.TrimSpace(authHeader)
	if authHeader == "" {
		return nil, errMissingAuthHeader
	}

	const bearerPrefix string = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return nil, errUnsupportedAuthorizationSchema
	}

	token := strings.TrimPrefix(authHeader, bearerPrefix)
	token = strings.TrimSpace(token)

	if token == "" {
		return nil, errMissingBearerToken
	}

	return &token, nil
}

func GetBearerToken(ctx context.Context) string {
	return ctx.Value(contextKeyBearerToken).(string)
}
