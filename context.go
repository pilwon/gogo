package gogo

import (
	"context"
	"net/http"
)

type (
	key                    int
	nextMiddlewareCallback func(c context.Context, w http.ResponseWriter, r *http.Request) context.Context
)

const (
	nextMiddlewareKey key = 0
)

func nextMiddlewareWithContext(c context.Context, cb nextMiddlewareCallback) context.Context {
	return context.WithValue(c, nextMiddlewareKey, cb)
}

func nextMiddlewareFromContext(c context.Context) (nextMiddlewareCallback, bool) {
	result, ok := c.Value(nextMiddlewareKey).(nextMiddlewareCallback)
	return result, ok
}
