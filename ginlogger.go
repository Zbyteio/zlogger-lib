package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type GinLogger interface {
	ginDebugLogger(httpMethod, absolutePath, handlerName string, nuHandlers int) 
	GinDebugRequestLoggerMiddleware() gin.HandlerFunc
	GinRequestLoggerMiddleware() gin.HandlerFunc
	RequestLoggerGin(param gin.LogFormatterParams) string
}
type ginLogger struct {
	applogger AppLogger
}

func newGinLogger(appLogger AppLogger) GinLogger {
	return &ginLogger{applogger: appLogger}
}

func (gl ginLogger)GinDebugRequestLoggerMiddleware() gin.HandlerFunc {
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
}

func (gl ginLogger)GinRequestLoggerMiddleware() gin.HandlerFunc {
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


func (gl ginLogger)ginDebugLogger(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	gl.applogger.Infof("%s\t%s", colorifyRequestMethod(httpMethod), absolutePath)
}


func (gl ginLogger)RequestLoggerGin(param gin.LogFormatterParams) string {
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