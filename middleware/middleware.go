package middleware

import (
	"net/http"

	"golang.org/x/net/context"
)

type Handler interface {
	ServeHTTP(c context.Context, w http.ResponseWriter, r *http.Request) context.Context
}

type HandlerFunc func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context

func (h HandlerFunc) ServeHTTP(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	return h(c, w, r)
}

func WrapHandler(h http.Handler) Handler {
	return HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		h.ServeHTTP(w, r)
		return nil
	})
}
