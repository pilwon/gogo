package negroni

import (
	"context"
	"net/http"

	"github.com/pilwon/gogo"
	"github.com/pilwon/gogo/middleware"
)

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h(w, r, next)
}

func Middleware(h Handler) middleware.Handler {
	return middleware.HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		h.ServeHTTP(w, r, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gogo.Next(c, w, r)
		}))
		return nil
	})
}

func MiddlewareFunc(h func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)) middleware.Handler {
	return Middleware(HandlerFunc(h))
}
