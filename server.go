package gogo

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pilwon/gogo/middleware"
	"github.com/pilwon/gogo/router"
)

type middlewareNode struct {
	handler middleware.Handler
	next    *middlewareNode
}

func (m middlewareNode) ServeHTTP(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	nextCalled := false
	ctx := nextMiddlewareWithContext(c, func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		if !nextCalled && m.next != nil {
			return m.next.ServeHTTP(c, w, r)
		}
		return c
	})
	ctx = m.handler.ServeHTTP(ctx, w, r)
	if !nextCalled && ctx != nil && m.next != nil {
		return m.next.ServeHTTP(c, w, r)
	}
	return ctx
}

type Server struct {
	Context context.Context

	handlers []middleware.Handler
	root     middlewareNode
}

func newServer(c context.Context, config Config) *Server {
	if registeredRouter == nil {
		fmt.Fprintf(os.Stderr, "error: Missing \"Router\". Register a router.\n")
		os.Exit(1)
	}
	return &Server{
		Context:  c,
		handlers: []middleware.Handler{},
		root:     voidMiddlewareNode(),
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.root.ServeHTTP(s.Context, NewResponseWriter(w), r)
}

func (s *Server) Use(h middleware.Handler) {
	s.handlers = append(s.handlers, h)
	s.root = build(s.handlers)
}

func (s *Server) UseFunc(h middleware.HandlerFunc) {
	s.Use(middleware.HandlerFunc(h))
}

func (s *Server) UseHandler(h http.Handler) {
	s.Use(middleware.WrapHandler(h))
}

func (s *Server) UseHandlerFunc(h http.HandlerFunc) {
	s.UseHandler(http.HandlerFunc(h))
}

func (s *Server) Get(path string, h router.HandlerFunc) {
	registeredRouter.AddRoute(s.Context, "GET", path, router.MiddlewareFromRouterHandler(h))
}
func (s *Server) GetAll(h router.HandlerFunc) {
	registeredRouter.AddRouteAll(s.Context, "GET", router.MiddlewareFromRouterHandler(h))
}

func (s *Server) Post(path string, h router.HandlerFunc) {
	registeredRouter.AddRoute(s.Context, "POST", path, router.MiddlewareFromRouterHandler(h))
}
func (s *Server) PostAll(h router.HandlerFunc) {
	registeredRouter.AddRouteAll(s.Context, "POST", router.MiddlewareFromRouterHandler(h))
}

func (s *Server) Put(path string, h router.HandlerFunc) {
	registeredRouter.AddRoute(s.Context, "PUT", path, router.MiddlewareFromRouterHandler(h))
}
func (s *Server) PutAll(h router.HandlerFunc) {
	registeredRouter.AddRouteAll(s.Context, "PUT", router.MiddlewareFromRouterHandler(h))
}

func (s *Server) Delete(path string, h router.HandlerFunc) {
	registeredRouter.AddRoute(s.Context, "DELETE", path, router.MiddlewareFromRouterHandler(h))
}
func (s *Server) DeleteAll(h router.HandlerFunc) {
	registeredRouter.AddRouteAll(s.Context, "DELETE", router.MiddlewareFromRouterHandler(h))
}

func (s *Server) Head(path string, h router.HandlerFunc) {
	registeredRouter.AddRoute(s.Context, "HEAD", path, router.MiddlewareFromRouterHandler(h))
}
func (s *Server) HeadAll(h router.HandlerFunc) {
	registeredRouter.AddRouteAll(s.Context, "HEAD", router.MiddlewareFromRouterHandler(h))
}

func (s *Server) Patch(path string, h router.HandlerFunc) {
	registeredRouter.AddRoute(s.Context, "PATCH", path, router.MiddlewareFromRouterHandler(h))
}
func (s *Server) PatchAll(h router.HandlerFunc) {
	registeredRouter.AddRouteAll(s.Context, "PATCH", router.MiddlewareFromRouterHandler(h))
}

func (s *Server) Options(path string, h router.HandlerFunc) {
	registeredRouter.AddRoute(s.Context, "OPTIONS", path, router.MiddlewareFromRouterHandler(h))
}
func (s *Server) OptionsAll(h router.HandlerFunc) {
	registeredRouter.AddRouteAll(s.Context, "OPTIONS", router.MiddlewareFromRouterHandler(h))
}

func (s *Server) Run(addr string) error {
	if addr == "" {
		return fmt.Errorf("Missing addr")
	}

	handler, err := registeredRouter.Handler()
	if err != nil {
		return err
	}
	s.UseHandler(handler)

	l := log.New(os.Stdout, "[gogo] ", 0)
	l.Printf("Listening on %s", addr)
	l.Fatal(http.ListenAndServe(addr, s))

	return nil
}

func build(handlers []middleware.Handler) middlewareNode {
	var next middlewareNode
	switch len(handlers) {
	case 0:
		return voidMiddlewareNode()
	case 1:
		next = voidMiddlewareNode()
	default:
		next = build(handlers[1:])
	}
	return middlewareNode{
		handler: handlers[0],
		next:    &next,
	}
}

func voidMiddlewareNode() middlewareNode {
	return middlewareNode{
		handler: middleware.HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
			return nil
		}),
		next: nil,
	}
}
