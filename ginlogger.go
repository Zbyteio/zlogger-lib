package zlogger

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var gl ginLogger

type GinLogger interface {
	ginDebugLogger(httpMethod, absolutePath, handlerName string, nuHandlers int)
	GinRequestLoggerMiddleware(params gin.LogFormatterParams) string 
}

type ginLogger struct {
	*zap.Logger
}

func NewGinLogger(loggerConfig loggerConfig, skipRoutes []string) gin.LoggerConfig {
	_libLogger := generateZapLogger(&loggerConfig.config, "lib")
	loggerConfig.config.DisableCaller = true

	loggerConfig.config.EncoderConfig.MessageKey = "requestUrl"
	
	gl = ginLogger{generateZapLogger(&loggerConfig.config, loggerConfig.loggerName)}
	
	if loggerConfig.loggerType == DEBUG_LOGGER {
		_libLogger.Info("created a [DEBUG-GIN-LOGGER] with logger-name :: " + loggerConfig.loggerName)
	} else if loggerConfig.loggerType == JSON_LOGGER {
		_libLogger.Info("created a [JSON-GIN-LOGGER] with logger-name :: " + loggerConfig.loggerName)
	}
	// set logger function to
	// print routes for this logger
	//gin.DebugPrintRouteFunc = _ginLogger.ginDebugLogger
	return gin.LoggerConfig{
		SkipPaths: skipRoutes,
		Formatter: gin.LogFormatter(ginRequestLoggerMiddleware),
	}
}

func ginRequestLoggerMiddleware(params gin.LogFormatterParams) string {
	if gl.Level().CapitalString() > zapcore.DebugLevel.CapitalString() {
		// PRODUCTION

		gl.Info(params.Path,
			zap.Int("statusCode", params.StatusCode),
			zap.String("requestMethod", params.Method),
			zap.String("error", params.ErrorMessage),
			zap.String("clientIP", params.ClientIP),
			zap.Duration("latency", params.Latency),
		)
	} else {
			// DEBUG
			var formatedStatusCode string = colorifySatusCode(params.StatusCode)
			var formatedRequestMethod string = colorifyRequestMethod(params.Method)
			var formatedLatency string = colorifyRequestLatency(params.Latency)

			if(params.ErrorMessage != "") {
				var formattedError string = colorifyRequestError(params.ErrorMessage)
				gl.Sugar().Errorf("%-18s%-20s%s\t%s\t%s\t%s",
					formatedStatusCode,
					formatedRequestMethod,
					params.Path,
					formattedError,
					params.ClientIP,
					formatedLatency)
			} else {
				gl.Sugar().Infof("%-18s%-20s%s\t%s\t%s",
					formatedStatusCode,
					formatedRequestMethod,
					params.Path,
					params.ClientIP,
					formatedLatency)
			}
			
	}
	return ""
}

// for printing all the routes defined in gin
func GinDebugLogger(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	if gl.Level().CapitalString() > zapcore.DebugLevel.CapitalString() {
		// PRODUCTION
		gl.Info(absolutePath, 
		zap.String("requestMethod", httpMethod),
	)
	}else {
		// DEBUG
		gl.Info(fmt.Sprintf("%-8s%s", colorifyRequestMethod(httpMethod), absolutePath))
	}
}