package httprouter

import (
	"context"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pilwon/gogo"
	"github.com/pilwon/gogo/gogocontext"
	"github.com/pilwon/gogo/middleware"
)

func init() {
	gogo.RegisterRouter(New())
}

func wrapHandler(c context.Context, h middleware.Handler) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		params := gogocontext.Params{}
		for _, param := range p {
			params[param.Key] = param.Value
		}
		c = gogocontext.ParamsWithContext(c, params)
		h.ServeHTTP(c, w, r)
	})
}

type Router struct {
	router *httprouter.Router
}

func New() *Router {
	return &Router{
		router: httprouter.New(),
	}
}

func (r *Router) AddRoute(c context.Context, httpVerb string, path string, h middleware.Handler) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("Failed to add route")
		}
	}()
	r.router.Handle(httpVerb, path, wrapHandler(c, h))
	return
}

func (r *Router) AddRouteAll(c context.Context, httpVerb string, h middleware.Handler) (err error) {
	return r.AddRoute(c, httpVerb, "/*all", h)
}

func (r *Router) Handler() (http.Handler, error) {
	return r.router, nil
}
