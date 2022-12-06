package zlogger

import (
	"fmt"
	"log"

	"go.uber.org/zap"
)

func CreateLoggerName(serviceName string, packageName string, rest ...string) string {
  var loggerName string = fmt.Sprintf("%s.%s", serviceName, packageName)
  var restLen int = len(rest)
  
  for i := 0; i < restLen; i++ {
    loggerName += "." + rest[i]
  }
  return loggerName
}

func generateZapLogger(zapconfig *zap.Config,loggerName string)(*zap.Logger) {
  var _logger *zap.Logger
  var err error
	_logger, err = zapconfig.Build(zap.AddCallerSkip(1))
	defer _logger.Sync()
	if err != nil {
		// zap logger unable to initialize
		// use default logger to log this
		log.Printf("ERROR :: %s", err.Error())
		return nil
	}
	_logger = _logger.Named(loggerName)
	return _logger
}