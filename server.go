package gogo

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pilwon/gogo/middleware"
	"golang.org/x/net/context"
)

type middlewareNode struct {
	handler middleware.Handler
	next    *middlewareNode
}

func (m middlewareNode) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.ServeHTTPWithContext(Context, w, r)
}

func (m middlewareNode) ServeHTTPWithContext(c context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	nextCalled := false
	ctx := nextMiddlewareWithContext(c, func(c context.Context) context.Context {
		if !nextCalled && m.next != nil {
			return m.next.ServeHTTPWithContext(c, w, r)
		}
		return c
	})
	ctx = m.handler.ServeHTTP(ctx, w, r)
	if !nextCalled && ctx != nil && m.next != nil {
		return m.next.ServeHTTPWithContext(c, w, r)
	}
	return ctx
}

type server struct {
	handlers []middleware.Handler
	root     middlewareNode
}

func newServer(config Config) *server {
	if registeredRouter == nil {
		fmt.Fprintf(os.Stderr, "error: Missing \"Router\". Register a router.\n")
		os.Exit(1)
	}
	return &server{
		handlers: []middleware.Handler{},
		root:     voidMiddlewareNode(),
	}
}

func (g *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.root.ServeHTTP(NewResponseWriter(w), r)
}

func (g *server) Use(h middleware.Handler) {
	g.handlers = append(g.handlers, h)
	g.root = build(g.handlers)
}

func (g *server) UseFunc(h middleware.HandlerFunc) {
	g.Use(middleware.HandlerFunc(h))
}

func (g *server) UseHandler(h http.Handler) {
	g.Use(middleware.WrapHandler(h))
}

func (g *server) UseHandlerFunc(h http.HandlerFunc) {
	g.UseHandler(http.HandlerFunc(h))
}

func (g *server) Get(path string, h middleware.HandlerFunc) {
	registeredRouter.AddRoute(Context, "GET", path, h)
}

func (g *server) Head(path string, h middleware.HandlerFunc) {
	registeredRouter.AddRoute(Context, "HEAD", path, h)
}

func (g *server) Options(path string, h middleware.HandlerFunc) {
	registeredRouter.AddRoute(Context, "OPTIONS", path, h)
}

func (g *server) Post(path string, h middleware.HandlerFunc) {
	registeredRouter.AddRoute(Context, "POST", path, h)
}

func (g *server) Put(path string, h middleware.HandlerFunc) {
	registeredRouter.AddRoute(Context, "PUT", path, h)
}

func (g *server) Patch(path string, h middleware.HandlerFunc) {
	registeredRouter.AddRoute(Context, "PATCH", path, h)
}

func (g *server) Delete(path string, h middleware.HandlerFunc) {
	registeredRouter.AddRoute(Context, "DELETE", path, h)
}

func (g *server) Run(addr string) {
	g.UseHandler(registeredRouter.Handler())

	l := log.New(os.Stdout, "[gogo] ", 0)
	l.Printf("Listening on %s", addr)
	l.Fatal(http.ListenAndServe(addr, g))
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
