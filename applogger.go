package zlogger

import (
	"fmt"
	"log"

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
	var config zap.Config
	var localLogger *zap.Logger


	config = zap.NewProductionConfig()
	config.EncoderConfig = zap.NewProductionEncoderConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.Level.SetLevel(zap.InfoLevel)

	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.CallerKey = "filePath"
	config.EncoderConfig.LevelKey = "logLevel"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	localLogger, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		// zap logger unable to initialize
		// use default logger to log this
		log.Printf("ERROR :: %s", err.Error())
	}


	localLogger.Info("creating production-logger created")
	return &appLogger{localLogger}
}

func NewDebugZlogger() (AppLogger){
	var err error
	var config zap.Config
	var localLogger *zap.Logger


	config = zap.NewDevelopmentConfig()
	config.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.Level.SetLevel(zap.DebugLevel)
	
	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.CallerKey = "filePath"
	config.EncoderConfig.LevelKey = "logLevel"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	localLogger, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		// zap logger unable to initialize
		// use default logger to log this
		log.Printf("ERROR :: %s", err.Error())
	}

	localLogger.Info("creating a debug-logger")
	return &appLogger{localLogger}
}