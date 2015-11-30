package gogo

import (
	"log"

	"github.com/pilwon/gogo/gogocontext"
	"github.com/pilwon/gogo/middleware"
	"github.com/pilwon/gogo/router"
	"golang.org/x/net/context"
)

type Config map[string]interface{}

var (
	Context          = context.Background()
	registeredRouter router.Router
)

func New() *server {
	return newServer(Config{})
}

func NewWithConfig(config Config) *server {
	return newServer(config)
}

func Next(c context.Context) context.Context {
	if next := nextMiddlewareFromContext(c); next != nil {
		return next(c)
	}
	return c
}

func Params(c context.Context) gogocontext.Params {
	return gogocontext.ParamsFromContext(c)
}

func Param(c context.Context, key string) string {
	return gogocontext.ParamsFromContext(c)[key]
}

func RegisterRouter(r router.Router) {
	if registeredRouter != nil {
		log.Panicln("Router already registered.")
	}
	registeredRouter = r
}

var Wrap = middleware.WrapHandler
