package logger

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pilwon/gogo"
	"golang.org/x/net/context"
)

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type Logger struct {
	// Logger inherits from log.Logger used to log messages with the Logger middleware
	*log.Logger
}

// New returns a new Logger instance
func New() *Logger {
	return &Logger{log.New(os.Stdout, "[gogo] ", 0)}
}

func (l *Logger) ServeHTTP(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	start := time.Now()
	l.Printf("Started %s %s", r.Method, r.URL.Path)

	gogo.Next(c, w, r)

	res := w.(gogo.ResponseWriter)
	l.Printf("Completed %v %s in %v", res.Status(), http.StatusText(res.Status()), time.Since(start))

	return nil
}
