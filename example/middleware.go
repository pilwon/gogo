package main

import (
	"github.com/pilwon/gogo/middleware/logger"
	"github.com/pilwon/gogo/middleware/recovery"
	"github.com/pilwon/gogo/middleware/static"
)

var (
	NewRecoveryMiddleware = recovery.New
	NewLoggerMiddleware   = logger.New
	NewStaticMiddleware   = static.New
)
