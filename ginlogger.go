package zlogger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var gl ginLogger

type ginLogger struct {
	*zap.Logger
}

func NewGinLoggerConfig(loggerConfig loggerConfig, skipRoutes []string) gin.LoggerConfig {
	if skipRoutes != nil {
		skipRoutes = []string{}
	}
	_libLogger := generateZapLogger(&loggerConfig.config, "lib")
	loggerConfig.config.DisableCaller = true

	loggerConfig.config.EncoderConfig.MessageKey = "requestUrl"
	gl = ginLogger{generateZapLogger(&loggerConfig.config, loggerConfig.loggerName)}
	gin.DebugPrintRouteFunc = ginDebugLogger

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
	if gl.Level() > zapcore.DebugLevel {
		// PRODUCTION

		gl.Named("gin").Info(params.Path,
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
				gl.Named("gin").Sugar().Errorf("%-18s%-20s%s\t%s\t%s\t%s",
					formatedStatusCode,
					formatedRequestMethod,
					params.Path,
					formattedError,
					params.ClientIP,
					formatedLatency)
			} else {
				gl.Named("gin").Sugar().Infof("%-18s%-20s%s\t%s\t%s",
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
func ginDebugLogger(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	if gl.Level() > zapcore.DebugLevel {
		// PRODUCTION
		gl.Named("gin").Info(absolutePath, 
		zap.String("requestMethod", httpMethod),
	)
	}else {
		// DEBUG
		gl.Named("gin").Sugar().Infof("%-18s%s", colorifyRequestMethod(httpMethod), absolutePath)
	}
}