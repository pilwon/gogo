package router

import (
	"net/http"

	"github.com/pilwon/gogo/middleware"
	"golang.org/x/net/context"
)

type Router interface {
	AddRoute(c context.Context, httpVerb string, path string, h middleware.Handler) error
	Handler() http.Handler
}
