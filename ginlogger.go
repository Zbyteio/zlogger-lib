package zlogger

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type GinLogger interface {
	ginDebugLogger(httpMethod, absolutePath, handlerName string, nuHandlers int)
	GinRequestLoggerMiddleware() gin.HandlerFunc
}
type ginLogger struct {
	*zap.Logger
}

func NewGinLogger(ginMode string) GinLogger {
	// create a new zap logger
	var err error
	var config zap.Config
	var _ginLogger *zap.Logger
	gin.SetMode(ginMode)

	if gin.Mode() == gin.DebugMode {
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
	config.DisableCaller = true
	_ginLogger, err = config.Build()
	defer _ginLogger.Sync()
	if err != nil {
		// zap logger unable to initialize
		// use default logger to log this
		log.Printf("ERROR :: %s", err.Error())
	}

	libraryLogger := _ginLogger.Named("lib.gin")
	if gin.Mode() == gin.DebugMode {
		libraryLogger.Info("creating a [DEBUG-LOGGER] for :: " + gin.Mode())
	} else {
		libraryLogger.Info("creating a [JSON-LOGGER] for :: " + gin.Mode())
	}

	var newGinLogger GinLogger = ginLogger{_ginLogger.Named("gin")}

	// set logger function to
	// print routes for this logger
	gin.DebugPrintRouteFunc = newGinLogger.ginDebugLogger
	return newGinLogger
}

func (gl ginLogger) GinRequestLoggerMiddleware() gin.HandlerFunc {

	if gin.Mode() == gin.DebugMode {
		return func(c *gin.Context) {
			reqUrl := fmt.Sprintf("%s%s", c.Request.Host, c.Request.URL.String())
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

			gl.Named(c.Request.Proto).Info(fmt.Sprintf("%s\t%s\t%s\t%s",
				formatedStatusCode,
				formatedRequestMethod,
				reqUrl,
				formatedLatency))
		}
	} else {
		return func(c *gin.Context) {
			reqUrl := fmt.Sprintf("%s%s", c.Request.Host, c.Request.URL.String())	
			start := time.Now()
			// Before calling handler
			c.Next()
			stop := time.Now()
			// After calling handler
			// create color coding for status codes
			gl.Named(c.Request.Proto).Info("",
				zap.Int("statusCode", c.Writer.Status()),
				zap.String("requestMethod", c.Request.Method),
				zap.String("requestUrl", reqUrl),
				zap.Duration("latency", stop.Sub(start)),
			)
		}
	}
}

// for printing all the routes defined in gin
func (gl ginLogger) ginDebugLogger(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	gl.Info(fmt.Sprintf("%s\t%s", colorifyRequestMethod(httpMethod), absolutePath))
}
