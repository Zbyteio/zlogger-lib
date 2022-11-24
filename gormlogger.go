package zlogger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

func NewGormLogger(ginMode string) GormLogger {

	// create a zap logger
	var err error
	var _config zap.Config
	var _gormLogger *zap.Logger

	gin.SetMode(ginMode)

	if gin.Mode() == gin.DebugMode {
		_config = zap.NewDevelopmentConfig()
		_config.EncoderConfig = zap.NewDevelopmentEncoderConfig()
		_config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		_config.Level.SetLevel(zap.DebugLevel)
	} else {
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
	_config.DisableCaller = true
	_config.DisableStacktrace = true

	_gormLogger, err = _config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	defer _gormLogger.Sync()

	return GormLogger{
		ZapLogger:                 _gormLogger.Named("gorm"),
		LoggerMode:                gin.Mode(),
		LogLevel:                  gormlogger.Info,
		SlowThreshold:             100 * time.Millisecond,
		SkipCallerLookup:          false,
		IgnoreRecordNotFoundError: false,
	}
}

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
	l.ZapLogger.Sugar().Debugf(str, args...)
}

func (l GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Warn {
		return
	}
	l.ZapLogger.Sugar().Warnf(str, args...)
}

func (l GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Error {
		return
	}
	l.ZapLogger.Sugar().Errorf(str, args...)
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
			l.ZapLogger.Error(fmt.Sprintf("error=%stime=%v\trows= %d\tsql=%s", formattedError, formattedElapsed, rows, formattedSql))
		} else {
			l.ZapLogger.Error("trace",
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
			l.ZapLogger.Warn(fmt.Sprintf("time=%v\trows=%d\tsql=%s", formattedElapsed, rows, formattedSql))
			
		} else {
			l.ZapLogger.Debug("trace",
				zap.Duration("elapsed", elapsed),
				zap.Int64("rows", rows),
				zap.String("sql", sql))
		}
	case l.LogLevel >= gormlogger.Info:
		sql, rows := fc()
		if l.LoggerMode  == gin.DebugMode {
			formattedElapsed := colorifySqlLatency(elapsed, l.SlowThreshold)
			formattedSql := colorPallet.colorfgMagenta(sql)
			l.ZapLogger.Debug(fmt.Sprintf("time=%v\trows=%d\tsql=%s", formattedElapsed, rows, formattedSql))
		} else {
			l.ZapLogger.Debug("trace",
				zap.Duration("elapsed", elapsed),
				zap.Int64("rows", rows),
				zap.String("sql", sql))
		}
	}
}