package gogocontext

import (
	"golang.org/x/net/context"
)

type (
	Params map[string]string

	key int
)

const (
	paramsKey key = 0
)

func ParamsWithContext(c context.Context, params Params) context.Context {
	return context.WithValue(c, paramsKey, params)
}

func ParamsFromContext(c context.Context) Params {
	return c.Value(paramsKey).(Params)
}
