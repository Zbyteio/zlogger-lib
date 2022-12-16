package zlogger

import (
	"time"

	gormv1 "github.com/jinzhu/gorm"
	"go.uber.org/zap"
	gormv2 "gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type GormLogger struct {
	LoggerMode                string
	ZapLogger                 *zap.Logger
	LogLevel                  gormlogger.LogLevel
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
}

func SetupGormLogger(dbv1 *gormv1.DB, dbv2 *gormv2.DB, loggerConfig loggerConfig) {
	if dbv1 != nil {
		setupGormLoggerV1(dbv1, loggerConfig)
	} else if dbv2 != nil {
		setupGormLoggerV2(dbv2, loggerConfig)
	}
}

func setupGormLoggerV1(dbv1 *gormv1.DB, loggerConfig loggerConfig) GormLogger {
	loggerConfig.config.DisableCaller = true
	loggerConfig.config.DisableStacktrace = true
	loggerConfig.config.EncoderConfig.MessageKey = "query"
	_gormLogger := generateZapLogger(&loggerConfig.config, loggerConfig.loggerName)
	dbv1.SetLogger(GormLogger{ZapLogger: _gormLogger})
	dbv1.LogMode(true)
	return GormLogger{ZapLogger: _gormLogger}
}

func (l GormLogger) Print(values ...interface{}) {
	if len(values) < 2 {
		return
	}

	switch values[0] {
	case "sql":
		if(l.ZapLogger.Level() > zap.DebugLevel) {
			l.ZapLogger.Named("gorm").Debug(values[3].(string),
				zap.Any("values", values[4]),
				zap.Float64("duration", float64(values[2].(time.Duration))/float64(time.Millisecond)),
				zap.Int64("affected-rows", values[5].(int64)),
			)
		}else {
			latency := colorifySqlLatency(
				values[2].(time.Duration), 
				l.SlowThreshold)

			l.ZapLogger.Named("gorm").Sugar().Debugf("duration=%s rows=%d sql=%s", 
			latency, values[5].(int64), colorPallet.colorfgMagenta(values[3].(string)))
		}
	default:
		l.ZapLogger.Debug("",
			zap.Any("values", values[2:]),
			zap.String("filePath", values[1].(string)),
		)
	}
}
