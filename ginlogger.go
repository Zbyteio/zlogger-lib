package zlogger

import (
	"time"

	"github.com/gin-gonic/gin"
)

type GinLogger interface {
	ginDebugLogger(httpMethod, absolutePath, handlerName string, nuHandlers int)
	GinRequestLoggerMiddleware() gin.HandlerFunc
}
type ginLogger struct {
	applogger AppLogger
}

func NewGinLogger(appLogger AppLogger) GinLogger {
	var newGinLogger GinLogger = ginLogger{applogger: appLogger}
	gin.DebugPrintRouteFunc = newGinLogger.ginDebugLogger
	return newGinLogger
}

func (gl ginLogger) GinRequestLoggerMiddleware() gin.HandlerFunc {
	if gin.Mode() == gin.DebugMode {
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

			gl.applogger.Infof("%s\t%s\t%s%s\t%s\t%s",
				formatedStatusCode,
				formatedRequestMethod,
				c.Request.Host,
				c.Request.RequestURI,
				c.Request.Proto,
				formatedLatency)

		}
	} else {
		return func(c *gin.Context) {
			start := time.Now()
			// Before calling handler
			c.Next()
			stop := time.Now()
			// After calling handler
			// create color coding for status codes

			gl.applogger.Infof("%s\t%s\t%s%s\t%s\t%s",
				c.Writer.Status(),
				c.Request.Method,
				c.Request.Host,
				c.Request.RequestURI,
				c.Request.Proto,
				stop.Sub(start))
		}
	}
}

// for printing all the routes defined in gin
func (gl ginLogger) ginDebugLogger(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	gl.applogger.Infof("%s\t%s", colorifyRequestMethod(httpMethod), absolutePath)
}
