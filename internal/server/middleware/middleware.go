package middleware

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler
type ContextKey string

func Chain(h http.HandlerFunc, mws ...Middleware) http.Handler {
	return ChainMux(http.HandlerFunc(h), mws...)
}

func ChainMux(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
