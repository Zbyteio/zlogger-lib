package zlogger_test

import (
	"testing"

	"github.com/Zbyteio/zlogger-lib"
	"go.uber.org/zap/zapcore"
)

func TestAppLogger(t *testing.T)  {
  // var ZBlocksAppDebugLogger zlogger.AppLogger = zlogger.NewAppLogger(zlogger.NewLoggerConfig("applogger", zlogger.DEBUG_LOGGER, zapcore.DebugLevel))
  var ZBlocksAppReleaseLogger zlogger.AppLogger = zlogger.NewAppLogger(zlogger.NewLoggerConfig("applogger", zlogger.JSON_LOGGER, zapcore.InfoLevel))


  t.Run("Test App logger", func(t *testing.T) {
    //ZBlocksAppDebugLogger.Debugf("%s", "success print debug via applogger[DEBUG]")
    ZBlocksAppReleaseLogger.Debugf("%s", "success print debug via applogger[RELEASE]")
  })
}