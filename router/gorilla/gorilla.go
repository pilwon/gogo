package gorilla

import (
	"fmt"
	// "net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	router *mux.Router
}

func New() *Router {
	return &Router{}
}

func (r *Router) AddRoute(httpVerb string, path string, handler Handler) error {
	fmt.Printf("Add route: %v\n", path)
  return nil
}
