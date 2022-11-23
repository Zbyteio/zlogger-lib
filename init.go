package zlogger

import (
	"github.com/gin-gonic/gin"
)

type ZLoggers struct {
	LoggerEnv LogEnvironment
	ZAppLogger  AppLogger
	ZGinLogger GinLogger
	ZGormLogger GormLogger
}

// New creates a new logger using the default configuration.
func InitLoggers(systemEnv LogEnvironment) ZLoggers {
	var zlogger ZLoggers = ZLoggers{}

	zlogger.LoggerEnv = systemEnv

	if zlogger.LoggerEnv == Local {
		zlogger.ZAppLogger = newDebugZlogger()
	} else {
		zlogger.ZAppLogger = newZlogger()
	}

	zlogger.ZGinLogger = newGinLogger(zlogger.ZAppLogger)
	gin.DebugPrintRouteFunc = zlogger.ZGinLogger.ginDebugLogger
	zlogger.ZGormLogger = newGormlogger(zlogger.ZAppLogger)

	return zlogger
}