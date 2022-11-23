package main

import (
	"github.com/gin-gonic/gin"
)



var (
	LoggerEnv LogEnvironment
	ZAppLogger  AppLogger
	ZGinLogger GinLogger
	ZGormLogger GormLogger
)
// New creates a new logger using the default configuration.
func InitLoggers(systemEnv LogEnvironment) {
	LoggerEnv = systemEnv

	if LoggerEnv == Local {
		ZAppLogger = newDebugZlogger()
	} else {
		ZAppLogger = newZlogger()
	}

	ZGinLogger = newGinLogger(ZAppLogger)
	gin.DebugPrintRouteFunc = ZGinLogger.ginDebugLogger
	ZGormLogger = newGormlogger(ZAppLogger)
}