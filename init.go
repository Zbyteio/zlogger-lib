package zlogger

import (
	"errors"

	"github.com/gin-gonic/gin"
	gormv1 "github.com/jinzhu/gorm"
	"go.uber.org/zap/zapcore"
	gormv2 "gorm.io/gorm"
)

var (
	_isInitialised    bool
	_defaultAppLogger AppLogger
	_defaultGinConfig gin.LoggerConfig
)

func SetupLoggerWithConfig(serviceName string, loggerType LoggerType, dbv1 *gormv1.DB, dbv2 *gormv2.DB, skipRoutes []string) {
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
  SetupGormLogger(dbv1, dbv2, loggerConfig)

  // init gin logger
  _defaultGinConfig = NewGinLoggerConfig(loggerConfig, skipRoutes)
	
	// set _isInitialised
	_isInitialised = true
}


func GetDefaultAppLogger() (*AppLogger, error) {
	if _isInitialised {
		return &_defaultAppLogger, nil
	}
	return nil, errors.New("logger not initialised")
}

func GetDefaultGinConfig() (*gin.LoggerConfig, error) {
	if _isInitialised {
		return &_defaultGinConfig, nil
	}
	return nil, errors.New("logger not initialised")
}