package gogo

import "golang.org/x/net/context"

type (
	key                    int
	nextMiddlewareCallback func(context.Context) context.Context
)

const (
	nextMiddlewareKey key = 0
)

func nextMiddlewareWithContext(c context.Context, cb nextMiddlewareCallback) context.Context {
	return context.WithValue(c, nextMiddlewareKey, cb)
}

func nextMiddlewareFromContext(c context.Context) nextMiddlewareCallback {
	return c.Value(nextMiddlewareKey).(nextMiddlewareCallback)
}
