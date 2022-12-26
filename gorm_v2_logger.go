package zlogger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

func setupGormLogger(db *gorm.DB, loggerConfig loggerConfig) GormLogger {
	loggerConfig.config.DisableCaller = true
	loggerConfig.config.DisableStacktrace = true
	
	_libLogger := generateZapLogger(&loggerConfig.config, "lib")
	_gormLogger := generateZapLogger(&loggerConfig.config, loggerConfig.loggerName)


	gormLogger := GormLogger{
		ZapLogger:                 _gormLogger,
		LoggerMode:                gin.DebugMode,
		LogLevel:                  gormlogger.Info,
		SlowThreshold:             100 * time.Millisecond,
		SkipCallerLookup:          false,
		IgnoreRecordNotFoundError: false,
	}

	if loggerConfig.loggerType == DEBUG_LOGGER {
		_libLogger.Info("created a [DEBUG-GORM-LOGGER] with logger-name :: " + loggerConfig.loggerName)
	} else if loggerConfig.loggerType == JSON_LOGGER {
		gormLogger.LoggerMode = gin.ReleaseMode
		_libLogger.Info("created a [JSON-GORM-LOGGER] with logger-name :: " + loggerConfig.loggerName)
	}
	gormlogger.Default = gormLogger
	if db != nil {
		db.Logger = gormLogger
	}
	return gormLogger
}

// try to accomodate this in NewGormLogger func
func (l GormLogger) SetAsDefault() {
	gormlogger.Default = l
}

func (l GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return GormLogger{
		ZapLogger:                 l.ZapLogger,
		SlowThreshold:             l.SlowThreshold,
		LogLevel:                  level,
		SkipCallerLookup:          l.SkipCallerLookup,
		IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
	}
}

func (l GormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Info {
		return
	}
	l.ZapLogger.Named("gorm").Sugar().Debugf(str, args...)
}

func (l GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Warn {
		return
	}
	l.ZapLogger.Named("gorm").Sugar().Warnf(str, args...)
}

func (l GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Error {
		return
	}
	l.ZapLogger.Named("gorm").Sugar().Errorf(str, args...)
}

func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		if l.LoggerMode == gin.DebugMode {
			formattedError := colorPallet.colorfgRed(err.Error())
			formattedElapsed := colorifySqlLatency(elapsed, l.SlowThreshold)
			formattedSql := colorPallet.colorfgMagenta(sql)
			l.ZapLogger.Named("gorm").Error(fmt.Sprintf("error=%stime=%v\trows= %d\tsql=%s", formattedError, formattedElapsed, rows, formattedSql))
		} else {
			l.ZapLogger.Named("gorm").Error("trace",
				zap.Error(err),
				zap.Duration("elapsed", elapsed),
				zap.Int64("rows", rows),
				zap.String("sql", sql))
		}
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormlogger.Warn:
		sql, rows := fc()
		if l.LoggerMode == gin.DebugMode {
			formattedElapsed := colorifySqlLatency(elapsed, l.SlowThreshold)
			formattedSql := colorPallet.colorfgMagenta(sql)
			l.ZapLogger.Named("gorm").Warn(fmt.Sprintf("time=%v\trows=%d\tsql=%s", formattedElapsed, rows, formattedSql))
			
		} else {
			l.ZapLogger.Named("gorm").Debug("trace",
				zap.Duration("elapsed", elapsed),
				zap.Int64("rows", rows),
				zap.String("sql", sql))
		}
	case l.LogLevel >= gormlogger.Info:
		sql, rows := fc()
		if l.LoggerMode  == gin.DebugMode {
			formattedElapsed := colorifySqlLatency(elapsed, l.SlowThreshold)
			formattedSql := colorPallet.colorfgMagenta(sql)
			l.ZapLogger.Named("gorm").Debug(fmt.Sprintf("time=%v\trows=%d\tsql=%s", formattedElapsed, rows, formattedSql))
		} else {
			l.ZapLogger.Named("gorm").Debug("trace",
				zap.Duration("elapsed", elapsed),
				zap.Int64("rows", rows),
				zap.String("sql", sql))
		}
	}
}


func SetupGormLogger(db *gorm.DB, loggerConfig loggerConfig) {
	setupGormLogger(db, loggerConfig)
}