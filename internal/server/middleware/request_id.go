package middleware

import (
	"context"
	"net/http"
)

var contextKeyRequestId ContextKey = "REQUEST_ID"

type RequestIdFactory func() string

func RequestId(f RequestIdFactory) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), contextKeyRequestId, f())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetRequestId(ctx context.Context) (string, bool) {
	value := ctx.Value(contextKeyRequestId)
	if value == nil {
		return "", false
	}
	return value.(string), true
}
