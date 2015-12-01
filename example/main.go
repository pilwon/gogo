package main

import (
	"fmt"
	"net/http"

	"github.com/pilwon/gogo"
	_ "github.com/pilwon/gogo/router/httprouter"
	"golang.org/x/net/context"
)

func main() {
	app := gogo.New()

	app.Use(NewRecoveryMiddleware())
	app.Use(NewLoggerMiddleware())
	app.Use(NewStaticMiddleware(http.Dir("static")))

	app.Get("/", func(c context.Context, w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "gogo server\n")
	})

	app.Get("/hello/:name", func(c context.Context, w http.ResponseWriter, r *http.Request) {
		name, _ := gogo.Param(c, "name")
		fmt.Fprintf(w, "Hello, %s!\n", name)
	})

	app.Run(":8080")
}
