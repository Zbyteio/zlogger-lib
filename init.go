package zlogger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
	gormv2 "gorm.io/gorm"
)

var (
	_defaultAppLogger AppLogger
	_defaultGinConfig gin.LoggerConfig
)

func init() {
	SetupLoggerWithConfig("default", DEBUG_LOGGER, nil, nil)
}

func SetupLoggerWithConfig(serviceName string, loggerType LoggerType, db *gormv2.DB, skipRoutes []string) {
	var loggerConfig loggerConfig

	if loggerType == JSON_LOGGER {
		loggerConfig = NewLoggerConfig(
			serviceName,
			JSON_LOGGER,
			zapcore.InfoLevel)
	} else {
		loggerConfig = NewLoggerConfig(
			serviceName,
			DEBUG_LOGGER,
			zapcore.DebugLevel)
	}


  // init app logger
	_defaultAppLogger = NewAppLogger(loggerConfig)


  // init gorm logger
	SetupGormLogger(db, loggerConfig)
	
  // init gin logger
  _defaultGinConfig = NewGinLoggerConfig(loggerConfig, skipRoutes)
}


func GetAppLogger() AppLogger {
		return _defaultAppLogger
}

func GetGinConfig() gin.LoggerConfig {
	return _defaultGinConfig
}