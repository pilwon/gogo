package gorillalogger

import (
	"net/http"
	"os"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/pilwon/gogo"
	"golang.org/x/net/context"
)

// GorillaLogger is a middleware handler that logs using Gorilla logger
type GorillaLogger struct {
}

// New returns a new GorillaLogger instance
func New() *GorillaLogger {
	return &GorillaLogger{}
}

func (l *GorillaLogger) ServeHTTP(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	gorillaHandlers.LoggingHandler(os.Stdout, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gogo.Next(c, w, r)
	})).ServeHTTP(w, r)
	return nil
}
