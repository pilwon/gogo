package gogo

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/pilwon/gogo/gogocontext"
	"github.com/pilwon/gogo/middleware"
	"github.com/pilwon/gogo/router"
	"golang.org/x/net/context"
)

type Config map[string]interface{}

var (
	registeredRouter router.Router
)

func New() *Server {
	return newServer(context.Background(), Config{})
}

func NewWithConfig(config Config) *Server {
	return newServer(context.Background(), config)
}

func NewWithContext(c context.Context) *Server {
	return newServer(c, Config{})
}

func NewWithContextConfig(c context.Context, config Config) *Server {
	return newServer(c, config)
}

func Next(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	if next := nextMiddlewareFromContext(c); next != nil {
		return next(c, w, r)
	}
	return c
}

func Params(c context.Context) gogocontext.Params {
	return gogocontext.ParamsFromContext(c)
}

func Param(c context.Context, key string) (string, error) {
	val, ok := gogocontext.ParamsFromContext(c)[key]
	if ok {
		return val, nil
	} else {
		return val, errors.New(fmt.Sprintf("Invalid param: %s", key))
	}
}

func RegisterRouter(r router.Router) {
	if registeredRouter != nil {
		panic("Router already registered.")
	}
	registeredRouter = r
}

func MiddlewareFromHandler(h http.Handler) middleware.Handler {
	return middleware.WrapHandler(h)
}

func MiddlewareFromRouteHandler(h router.Handler) middleware.Handler {
	return router.MiddlewareFromRouterHandler(h)
}

func Wrap(h http.Handler) router.HandlerFunc {
	return router.WrapHandler(h).(router.HandlerFunc)
}
