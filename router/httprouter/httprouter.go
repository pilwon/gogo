package httprouter

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pilwon/gogo"
	"github.com/pilwon/gogo/middleware"
	"github.com/pilwon/gogo/gogocontext"
	"golang.org/x/net/context"
)

func init() {
	gogo.RegisterRouter(New())
}

type Router struct {
	router *httprouter.Router
}

func New() *Router {
	return &Router{
		router: httprouter.New(),
	}
}

func (r *Router) AddRoute(c context.Context, httpVerb string, path string, h middleware.Handler) error {
	switch httpVerb {
	case "GET":
		r.router.GET(path, wrapHandler(c, h))
	case "HEAD":
		r.router.HEAD(path, wrapHandler(c, h))
	case "OPTIONS":
		r.router.OPTIONS(path, wrapHandler(c, h))
	case "POST":
		r.router.POST(path, wrapHandler(c, h))
	case "PUT":
		r.router.PUT(path, wrapHandler(c, h))
	case "PATCH":
		r.router.PATCH(path, wrapHandler(c, h))
	case "DELETE":
		r.router.DELETE(path, wrapHandler(c, h))
	default:
		panic(fmt.Sprintf("Unsupported HTTP verb: (%s, %s)", httpVerb, path))
	}
	return nil
}

func (r *Router) Handler() http.Handler {
	return r.router
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
