package zlogger

import (
	"time"

	"go.uber.org/zap"
)

type Logger struct {
	zap *zap.Logger
}

func SetupGormLoggerV1(loggerConfig loggerConfig) Logger {
	loggerConfig.config.DisableCaller = true
	loggerConfig.config.DisableStacktrace = true
  _gormLogger := generateZapLogger(&loggerConfig.config, loggerConfig.loggerName)
	return Logger{zap: _gormLogger}
}

func (l Logger) Print(values ...interface{}) {
	if len(values) < 2 {
		return
	}

	switch values[0] {
	case "sql":
		l.zap.Debug("gorm.v1.debug.sql",
			zap.String("query", values[3].(string)),
			zap.Any("values", values[4]),
			zap.Float64("duration", float64(values[2].(time.Duration))/float64(time.Millisecond)),
			zap.Int64("affected-rows", values[5].(int64)),
			zap.String("filePath", values[1].(string)), // if AddCallerSkip(6) is well defined, we can safely remove this field
		)
	default:
		l.zap.Debug("gorm.v1.debug.other",
			zap.Any("values", values[2:]),
			zap.String("source", values[1].(string)), // if AddCallerSkip(6) is well defined, we can safely remove this field
		)
	}
}