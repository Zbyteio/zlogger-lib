package zlogger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

var _appLogger AppLogger


type appLogger struct {
	*zap.Logger
}

// Logger is a logger that supports log levels, context and structured logging.
type AppLogger interface {
	// Debug uses fmt.Sprint to construct and log a message at DEBUG level
	Debug(msg string, fields ...zapcore.Field)
	// Info uses fmt.Sprint to construct and log a message at INFO level
	Info(msg string, fields ...zapcore.Field)
	// Warn uses fmt.Sprint to construct and log a message at INFO level
	Warn(msg string, fields ...zapcore.Field)
	// Error uses fmt.Sprint to construct and log a message at ERROR level
	Error(msg string, fields ...zapcore.Field)

	// Debug uses fmt.Sprint to construct and log a message at DEBUG level
	Debugf(template string, args ...interface{})
	// Info uses fmt.Sprint to construct and log a message at INFO level
	Infof(template string, args ...interface{})
	// Warn uses fmt.Sprint to construct and log a message at INFO level
	Warnf(template string, args ...interface{})
	// Error uses fmt.Sprint to construct and log a message at ERROR level
	Errorf(template string, args ...interface{})
}

func (l *appLogger) Debugf(template string, args ...interface{}) {
	debugString := fmt.Sprintf(template, args...)
	l.Named("app").Debug(debugString)
}

func (l *appLogger) Infof(template string, args ...interface{}) {
	infoString := fmt.Sprintf(template, args...)
	l.Named("app").Info(infoString)
}

func (l *appLogger) Warnf(template string, args ...interface{}) {
	warnString := fmt.Sprintf(template, args...)
	l.Named("app").Warn(warnString)
}

func (l *appLogger) Errorf(template string, args ...interface{}) {
	errorString := fmt.Sprintf(template, args...)
	l.Named("app").Error(errorString)
}

// NewZloggerForTest returns a new logger and the corresponding observed logs which can be used in unit tests to verify log entries.
func NewAppLoggerForTest() (AppLogger, *observer.ObservedLogs) {
	var testLogger *zap.Logger
	var testCore zapcore.Core
	var recorded *observer.ObservedLogs

	testCore, recorded = observer.New(zapcore.InfoLevel)

	testLogger = zap.New(testCore)
	return &appLogger{testLogger}, recorded
}

/*
* loggerConfig.loggerType - debug / json
* loggerName - name of the logger ("app" :default) 
*/
func NewAppLogger(loggerConfig loggerConfig) (AppLogger){
	_libLogger := generateZapLogger(&loggerConfig.config, "lib")
	_appLogger = &appLogger{generateZapLogger(&loggerConfig.config, loggerConfig.loggerName)}

	if loggerConfig.loggerType == DEBUG_LOGGER {
		_libLogger.Info("created a [DEBUG-APP-LOGGER] with logger-name :: " + loggerConfig.loggerName)
	} else if loggerConfig.loggerType == JSON_LOGGER {
		_libLogger.Info("created a [JSON-APP-LOGGER] with logger-name :: " + loggerConfig.loggerName)
	}
	return _appLogger
}