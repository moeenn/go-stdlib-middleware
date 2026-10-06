package middleware

import (
	"context"
	"net/http"
	"strings"
)

var contextKeyRequestId ContextKey = "REQUEST_ID"

const requestIdHeader string = "X-Request-Id"

type RequestIdFactory func() string

type RequestIdArgs struct {
	Factory        RequestIdFactory
	ReadFromHeader bool
}

func RequestId(args RequestIdArgs) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqId, ok := readRequestIdFromHeader(r)
			if !ok || !args.ReadFromHeader {
				reqId = args.Factory()
			}

			ctx := context.WithValue(r.Context(), contextKeyRequestId, reqId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func readRequestIdFromHeader(r *http.Request) (string, bool) {
	value := r.Header.Get(requestIdHeader)
	value = strings.TrimSpace(value)
	return value, value != ""
}

func GetRequestId(ctx context.Context) (string, bool) {
	value := ctx.Value(contextKeyRequestId)
	if value == nil {
		return "", false
	}
	return value.(string), true
}
