package interpose

import (
	"net/http"

	"github.com/pilwon/gogo"
	"github.com/pilwon/gogo/middleware"
	"golang.org/x/net/context"
)

func Middleware(h func(http.Handler) http.Handler) middleware.Handler {
	return middleware.HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		h(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gogo.Next(c, w, r)
		})).ServeHTTP(w, r)
		return nil
	})
}
