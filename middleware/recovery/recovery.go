package recovery

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/pilwon/gogo"
	"golang.org/x/net/context"
)

// Recovery is a middleware that recovers from any panics and writes a 500 if there was one.
type Recovery struct {
	Logger     *log.Logger
	PrintStack bool
	StackAll   bool
	StackSize  int
}

// New returns a new instance of Recovery
func New() *Recovery {
	return &Recovery{
		Logger:     log.New(os.Stdout, "[gogo] ", 0),
		PrintStack: false,
		StackAll:   false,
		StackSize:  1024 * 8,
	}
}

func (rec *Recovery) ServeHTTP(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			stack := make([]byte, rec.StackSize)
			stack = stack[:runtime.Stack(stack, rec.StackAll)]

			f := "PANIC: %s\n%s"
			rec.Logger.Printf(f, err, stack)

			if rec.PrintStack {
				fmt.Fprintf(w, f, err, stack)
			} else {
				fmt.Fprintf(w, "Internal Server Error")
			}
		}
	}()
	gogo.Next(c, w, r)
	return nil
}
