// Package log provides context-aware and structured logging capabilities.
package log

import (
	"database/sql/driver"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// Foreground colors.
const (
	FgBlack color = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Background colors.
const (
	BgBlack color = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

type color uint8

const (
	Production  environment = "production"
	Development environment = "development"
	Test        environment = "test"
	Local       environment = "local"
)

type environment string

// Logger is a logger that supports log levels, context and structured logging.
type Logger interface {
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

type zapLogger struct {
	*zap.Logger
}

type testGormLogger struct {
	gorm.LogWriter
}

func (tgl testGormLogger) Print(values ...interface{}) {
	var ( 
		sqlRegexp                = regexp.MustCompile(`\?`)
		numericPlaceHolderRegexp = regexp.MustCompile(`\$\d+`)

		gormValuesLen = len(values)
		formattedSource string
		formattedDuration string
		formattedValues []string
		sql string

		messages []interface{}
	)

	if gormValuesLen > 1 {
		// colorise file path
		formattedSource = values[1].(string)
	}

	if values[0].(string) == "sql" {
		// duration
		messages = append(messages, FgCyan.coloredString(
			fmt.Sprintf("%.2f", float64(values[2].(time.Duration).Nanoseconds()/1e4)/100.0)))
		// sql

		for _, value := range values[4].([]interface{}) {
			indirectValue := reflect.Indirect(reflect.ValueOf(value))
			if indirectValue.IsValid() {
				value = indirectValue.Interface()
				if t, ok := value.(time.Time); ok {
					if t.IsZero() {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", "0000-00-00 00:00:00"))
					} else {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format("2006-01-02 15:04:05")))
					}
				} else if b, ok := value.([]byte); ok {
					if str := string(b); isPrintable(str) {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
					} else {
						formattedValues = append(formattedValues, "'<binary>'")
					}
				} else if r, ok := value.(driver.Valuer); ok {
					if value, err := r.Value(); err == nil && value != nil {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
					} else {
						formattedValues = append(formattedValues, "NULL")
					}
				} else {
					switch value.(type) {
					case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
						formattedValues = append(formattedValues, fmt.Sprintf("%v", value))
					default:
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
					}
				}
			} else {
				formattedValues = append(formattedValues, "NULL")
			}
		}

		// differentiate between $n placeholders or else treat like ?
		if numericPlaceHolderRegexp.MatchString(values[3].(string)) {
			sql = values[3].(string)
			for index, value := range formattedValues {
				placeholder := fmt.Sprintf(`\$%d([^\d]|$)`, index+1)
				sql = regexp.MustCompile(placeholder).ReplaceAllString(sql, value+"$1")
			}
		} else {
			formattedValuesLength := len(formattedValues)
			for index, value := range sqlRegexp.Split(values[3].(string), -1) {
				sql += value
				if index < formattedValuesLength {
					sql += formattedValues[index]
				}
			}
		}

		messages = append(messages, sql)
		messages = append(messages, "\n", FgCyan.coloredString(fmt.Sprintf("%v", strconv.FormatInt(values[5].(int64), 10)+" rows affected or returned")))
	} else {
		messages = append(messages, values[2:]...)
	}
	var args []interface{}
	args = append(args, formattedSource, formattedDuration, sql)
	args = append(args, messages...)
	tgl.printErrorUsingProperType(values[0].(string), args)
}

func (tgl testGormLogger) printErrorUsingProperType(logType string, args []interface{}) {
	switch (logType) {
	case "error":
		zbyteGormLogger.Errorf(fmt.Sprintln(args...))
	case "warn":
		zbyteGormLogger.Warnf(fmt.Sprintln(args...))
	case "info":
		zbyteGormLogger.Infof(fmt.Sprintln(args...))
	case "sql":
		zbyteGormLogger.Debugf(fmt.Sprintln(args...))
	default:
		zbyteGormLogger.Debugf(fmt.Sprintln(args...))
	}
}

var (
	ZbyteLogger            Logger
	systemEnv              environment
	ZbyteGromDefaultLogger testGormLogger
	zbyteGormLogger Logger

)

// New creates a new logger using the default configuration.
func init() {
	godotenv.Load()
	var err error
	var config zap.Config
	var localLogger *zap.Logger
	var gormLogger *zap.Logger

	systemEnv.setEnv(os.Getenv("ENVIRONMENT"))

	if systemEnv == Local {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig = zap.NewDevelopmentEncoderConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.Level.SetLevel(zap.DebugLevel)
	} else {
		config = zap.NewProductionConfig()
		config.EncoderConfig = zap.NewProductionEncoderConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		config.Level.SetLevel(zap.InfoLevel)
	}
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
	ZbyteLogger = &zapLogger{localLogger}

	if systemEnv == Local {
		ZbyteLogger.Infof("creating debug-logger for environment: %s", systemEnv)
	} else {
		ZbyteLogger.Infof("creating production-logger created for environment: %s", systemEnv)
	}
	gin.DebugPrintRouteFunc = ginDebugLogger

	// create a gorm logger by tweaking config
	// to disable file path printing
	config.DisableCaller = true
	gormLogger, _ = config.Build()
	zbyteGormLogger = &zapLogger{gormLogger}
}

func RequestLoggerZap() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// Before calling handler
		c.Next()
		stop := time.Now()
		// After calling handler
		// create color coding for status codes
		var statusCode int = c.Writer.Status()
		var formatedStatusCode string = colorifySatusCode(statusCode)
		var formatedRequestMethod string = colorifyRequestMethod(c.Request.Method)
		var formatedLatency string = colorifyRequestLatency(stop.Sub(start))

		if systemEnv == Local {
			ZbyteLogger.Infof("%s\t%s\t%s%s\t%s\t%s",
				formatedStatusCode,
				formatedRequestMethod,
				c.Request.Host,
				c.Request.RequestURI,
				c.Request.Proto,
				formatedLatency)
		} else {
			ZbyteLogger.Infof("%d\t[%s]\t%s%s", statusCode, c.Request.Method, c.Request.Host, c.Request.RequestURI)
		}
	}
}

func ginDebugLogger(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	ZbyteLogger.Infof("%s\t%s", colorifyRequestMethod(httpMethod), absolutePath)
}

// NewForTest returns a new logger and the corresponding observed logs which can be used in unit tests to verify log entries.
func NewZbyteLoggerForTest() (Logger, *observer.ObservedLogs) {
	core, recorded := observer.New(zapcore.InfoLevel)
	return newWithZap(zap.New(core)), recorded
}

func (l *zapLogger) Debugf(template string, args ...interface{}) {
	errorString := fmt.Sprintf(template, args...)
	l.Debug(errorString)
}

func (l *zapLogger) Infof(template string, args ...interface{}) {
	errorString := fmt.Sprintf(template, args...)
	l.Info(errorString)
}

func (l *zapLogger) Warnf(template string, args ...interface{}) {
	errorString := fmt.Sprintf(template, args...)
	l.Warn(errorString)
}

func (l *zapLogger) Errorf(template string, args ...interface{}) {
	errorString := fmt.Sprintf(template, args...)
	l.Error(errorString)
}

func (l *zapLogger) GormLogger(args ...interface{}) {

}

func RequestLoggerGin(param gin.LogFormatterParams) string {
	// your custom format
	return fmt.Sprintf("%s %s [%s] %s %s %s\n",
		colorifyRequestMethod(param.Method),
		param.Path,
		param.Request.Proto,
		colorifySatusCode(param.StatusCode),
		colorifyRequestLatency(param.Latency),
		FgRed.coloredString(param.ErrorMessage),
	)
}

// ---------------------- INTERNAL --------------------------------------

func (env environment) setEnv(sysenv string) {
	switch sysenv {
	case "production":
		systemEnv = Production
	case "development":
		systemEnv = Development
	case "test":
		systemEnv = Test
	default:
		systemEnv = Local
	}
}

// newWithZap creates a new logger using the preconfigured zap logger.
func newWithZap(l *zap.Logger) Logger {
	l.Sync()
	return &zapLogger{l}
}

func (c color) coloredString(value string) string {
	return fmt.Sprintf("\x1b[%d;%dm %s \x1b[0m", c, 1, value)
}

func colorifySatusCode(statusCode int) string {
	var statusCodeString string = strconv.Itoa(statusCode)
	if statusCode >= 500 {
		return FgRed.coloredString(statusCodeString)
	} else if statusCode >= 400 {
		return FgYellow.coloredString(statusCodeString)
	} else if statusCode >= 300 {
		return FgCyan.coloredString(statusCodeString)
	} else if statusCode >= 200 {
		return FgGreen.coloredString(statusCodeString)
	}
	// default value
	return FgWhite.coloredString(statusCodeString)
}

func colorifyRequestMethod(methodName string) string {
	switch methodName {
	case "GET":
		return BgGreen.coloredString(methodName)
	case "POST":
		return BgYellow.coloredString(methodName)
	case "PUT":
		return BgBlue.coloredString(methodName)
	case "DELETE":
		return BgRed.coloredString(methodName)
	case "OPTION":
		return BgCyan.coloredString(methodName)
	case "PATCH":
		return BgMagenta.coloredString(methodName)
	default:
		return BgWhite.coloredString(methodName)
	}
}

func colorifyRequestLatency(latency time.Duration) string {
	if latency < time.Second {
		return FgGreen.coloredString(latency.String())
	} else if latency < time.Second*2 {
		return FgYellow.coloredString(latency.String())
	}
	return FgRed.coloredString(latency.String())
}



func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
