package main

import (
	"context"
	"net/http"
	"os"

	interposeMiddleware "github.com/carbocation/interpose/middleware"
	"github.com/goincremental/negroni-sessions"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/pilwon/gogo"
	"github.com/pilwon/gogo/adapter/interpose"
	"github.com/pilwon/gogo/adapter/negroni"
	"github.com/pilwon/gogo/middleware"
	"github.com/pilwon/gogo/middleware/logger"
	"github.com/pilwon/gogo/middleware/recovery"
	"github.com/pilwon/gogo/middleware/static"
)

var (
	NewLoggerMiddleware   = logger.New
	NewRecoveryMiddleware = recovery.New
	NewStaticMiddleware   = static.New
)

func NewBasicAuth(username string, password string) middleware.Handler {
	return interpose.Middleware(interposeMiddleware.BasicAuth(username, password))
}

func NewGorillaLogger() middleware.Handler {
	return middleware.HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		gorillaHandlers.LoggingHandler(os.Stdout, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gogo.Next(c, w, r)
		})).ServeHTTP(w, r)
		return nil
	})
}

func NewGzipMiddleware(compressionLevel int) middleware.Handler {
	return negroni.Middleware(gzip.Gzip(compressionLevel))
}

func NewSessions(name string, store sessions.Store) middleware.Handler {
	return negroni.Middleware(sessions.Sessions(name, store))
}
