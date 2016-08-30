package router

import (
	"context"
	"net/http"

	"github.com/pilwon/gogo/middleware"
)

type Router interface {
	AddRoute(c context.Context, httpVerb string, path string, h middleware.Handler) error
	AddRouteAll(c context.Context, httpVerb string, h middleware.Handler) error
	Handler() (http.Handler, error)
}

type Handler interface {
	ServeHTTP(c context.Context, w http.ResponseWriter, r *http.Request)
}

type HandlerFunc func(c context.Context, w http.ResponseWriter, r *http.Request)

func (h HandlerFunc) ServeHTTP(c context.Context, w http.ResponseWriter, r *http.Request) {
	h(c, w, r)
}

func WrapHandler(h http.Handler) Handler {
	return HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

func MiddlewareFromRouterHandler(h Handler) middleware.Handler {
	return middleware.HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		h.ServeHTTP(c, w, r)
		return nil
	})
}
