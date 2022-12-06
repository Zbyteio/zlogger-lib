package zlogger

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var _ginLogger GinLogger

type GinLogger interface {
	ginDebugLogger(httpMethod, absolutePath, handlerName string, nuHandlers int)
	GinRequestLoggerMiddleware() gin.HandlerFunc
}

type ginLogger struct {
	*zap.Logger
}

func NewGinLogger(loggerConfig loggerConfig) GinLogger {
	_libLogger := generateZapLogger(&loggerConfig.config, "lib")
	_ginLogger = &ginLogger{generateZapLogger(&loggerConfig.config, loggerConfig.loggerName)}

	if loggerConfig.loggerType == DEBUG_LOGGER {
		_libLogger.Info("created a [DEBUG-GIN-LOGGER] with logger-name :: " + loggerConfig.loggerName)
	} else if loggerConfig.loggerType == JSON_LOGGER {
		_libLogger.Info("created a [JSON-GIN-LOGGER] with logger-name :: " + loggerConfig.loggerName)
	}
	// set logger function to
	// print routes for this logger
	gin.DebugPrintRouteFunc = _ginLogger.ginDebugLogger
	return _ginLogger
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
	// TODO: remove color coding in release mode
	gl.Info(fmt.Sprintf("%s\t%s", colorifyRequestMethod(httpMethod), absolutePath))
}