package zlogger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)


type loggerConfig struct {
	loggerName string
	loggerType LoggerType
	loggerLevel zapcore.Level
  config zap.Config
}

func (lc *loggerConfig) GetLoggerName() string {
	return lc.loggerName
}

func (lc *loggerConfig) GetLoggerLevel() zapcore.Level {
	return lc.loggerLevel
}

func (lc *loggerConfig) GetLoggerType() LoggerType {
	return lc.loggerType
}

func (lc *loggerConfig) GetZapConfig() zap.Config {
	return lc.config
}

func (lc *loggerConfig) SetLoggerName(loggerName string) string {
	lc.loggerName = loggerName
	return lc.loggerName
}

func (lc *loggerConfig) SetLoggerLevel(loggerLevel zapcore.Level) zapcore.Level {
	lc.loggerLevel = loggerLevel
	return lc.loggerLevel
}

func (lc *loggerConfig) SetLoggerType(loggerType LoggerType) LoggerType {
	lc.loggerType = loggerType
	return lc.loggerType
}


func NewLoggerConfig(loggerName string, loggerType LoggerType, loggerLevel zapcore.Level) (loggerConfig) {
	if loggerType != DEBUG_LOGGER && loggerType != JSON_LOGGER {
		loggerType = DEBUG_LOGGER
	}
	if loggerName == "" {
		loggerName = "app"
	}

	_loggerConfig := loggerConfig{
		loggerName: loggerName,
		loggerType: loggerType,
		loggerLevel: loggerLevel,
		config:  zap.Config{
			Level:            zap.NewAtomicLevelAt(loggerLevel),
			Development:      false,
			Encoding:         "json",
			EncoderConfig:    zapcore.EncoderConfig{
				// Keys can be anything except the empty string.
				TimeKey:        "timestamp",
				LevelKey:       "logLevel",
				NameKey:        "loggerName",
				CallerKey:      "filePath",
				FunctionKey:    zapcore.OmitKey,
				MessageKey:     "message",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.CapitalLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		},
	}
	if loggerType == DEBUG_LOGGER {
		_loggerConfig.config.Encoding = "console"
		_loggerConfig.config.Development = true
		_loggerConfig.config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		_loggerConfig.config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	}
	return _loggerConfig
}