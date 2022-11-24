package zlogger

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

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

type appLogger struct {
	*zap.Logger
}

func (l *appLogger) Debugf(template string, args ...interface{}) {
	errorString := fmt.Sprintf(template, args...)
	l.Debug(errorString)
}

func (l *appLogger) Infof(template string, args ...interface{}) {
	errorString := fmt.Sprintf(template, args...)
	l.Info(errorString)
}

func (l *appLogger) Warnf(template string, args ...interface{}) {
	errorString := fmt.Sprintf(template, args...)
	l.Warn(errorString)
}

func (l *appLogger) Errorf(template string, args ...interface{}) {
	errorString := fmt.Sprintf(template, args...)
	l.Error(errorString)
}

// NewZloggerForTest returns a new logger and the corresponding observed logs which can be used in unit tests to verify log entries.
func NewZloggerForTest() (AppLogger, *observer.ObservedLogs) {
	var testLogger *zap.Logger
	var testCore zapcore.Core
	var recorded *observer.ObservedLogs

	testCore, recorded = observer.New(zapcore.InfoLevel)

	testLogger = zap.New(testCore)
	return &appLogger{testLogger}, recorded
}

func NewZlogger() (AppLogger){
	var err error
	var _config zap.Config
	var _appLogger *zap.Logger

	if gin.Mode() == gin.DebugMode {
		_config = zap.NewDevelopmentConfig()
		_config.EncoderConfig = zap.NewDevelopmentEncoderConfig()
		_config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		_config.Level.SetLevel(zap.DebugLevel)
	}else {
		_config = zap.NewProductionConfig()
		_config.EncoderConfig = zap.NewProductionEncoderConfig()
		_config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		_config.Level.SetLevel(zap.InfoLevel)
	}
	_config.EncoderConfig.TimeKey = "time"
	_config.EncoderConfig.CallerKey = "filePath"
	_config.EncoderConfig.LevelKey = "logLevel"
	_config.EncoderConfig.MessageKey = "message"
	_config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	_config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	_appLogger, err = _config.Build(zap.AddCallerSkip(1))
	defer _appLogger.Sync()


	libraryLogger := _appLogger.Named("library.zlogger.applogger")
	if err != nil {
		// zap logger unable to initialize
		// use default logger to log this
		log.Printf("ERROR :: %s", err.Error())
	}


	if gin.Mode() == gin.DebugMode{
		libraryLogger.Info("creating a [DEBUG-LOGGER] for :: " + gin.Mode())
	} else {
		libraryLogger.Info("creating a [JSON-LOGGER] for :: " + gin.Mode())
	}
	return &appLogger{_appLogger}
}