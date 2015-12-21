package denco

import (
	"errors"
	"net/http"

	"github.com/naoina/denco"
	"github.com/pilwon/gogo"
	"github.com/pilwon/gogo/gogocontext"
	"github.com/pilwon/gogo/middleware"
	"golang.org/x/net/context"
)

func init() {
	gogo.RegisterRouter(New())
}

func wrapHandler(c context.Context, h middleware.Handler) denco.HandlerFunc {
	return denco.HandlerFunc(func(w http.ResponseWriter, r *http.Request, p denco.Params) {
		params := gogocontext.Params{}
		for _, param := range p {
			params[param.Name] = param.Value
		}
		c = gogocontext.ParamsWithContext(c, params)
		h.ServeHTTP(c, w, r)
	})
}

type Router struct {
	handlers []denco.Handler
	mux      *denco.Mux
}

func New() *Router {
	return &Router{
		handlers: []denco.Handler{},
		mux:      denco.NewMux(),
	}
}

func (r *Router) AddRoute(c context.Context, httpVerb string, path string, h middleware.Handler) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("Failed to add route")
		}
	}()
	handler := r.mux.Handler(httpVerb, path, wrapHandler(c, h))
	r.handlers = append(r.handlers, handler)
	return
}

func (r *Router) AddRouteAll(c context.Context, httpVerb string, h middleware.Handler) (err error) {
	return r.AddRoute(c, httpVerb, "*", h)
}

func (r *Router) Handler() (http.Handler, error) {
	return r.mux.Build(r.handlers)
}
