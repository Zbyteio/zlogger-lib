package zlogger_test

import (
	"testing"

	"github.com/Zbyteio/zlogger-lib"
	"github.com/go-playground/assert/v2"
	"go.uber.org/zap/zapcore"
)

func TestConfigLogger(t *testing.T)  {
  t.Run("Test Debug logger setup", func(t *testing.T) {
    appdebugConf := zlogger.NewLoggerConfig("applogger", zlogger.DEBUG_LOGGER, zapcore.DebugLevel)
    logLevel := appdebugConf.GetLoggerLevel()
    loggerName:= appdebugConf.GetLoggerName()
    logType := appdebugConf.GetLoggerType()

    assert.Equal(t, logLevel, zapcore.DebugLevel)
    assert.Equal(t, loggerName, "applogger")
    assert.Equal(t, logType, zlogger.DEBUG_LOGGER)

    // check zap config
    zapConfig := appdebugConf.GetZapConfig()
    encodingType := zapConfig.Encoding
    devMode := zapConfig.Development

    assert.Equal(t, devMode, true)
    assert.Equal(t, encodingType, "console")
    assert.Equal(t, logLevel, zapConfig.Level.Level())

  })

  t.Run("Test Json logger setup", func(t *testing.T) {
    appprodConf := zlogger.NewLoggerConfig("applogger", zlogger.JSON_LOGGER, zapcore.InfoLevel)
    logLevel := appprodConf.GetLoggerLevel()
    loggerName:= appprodConf.GetLoggerName()
    logType := appprodConf.GetLoggerType()
    assert.Equal(t, logLevel, zapcore.InfoLevel)
    assert.Equal(t, loggerName, "applogger")
    assert.Equal(t, logType, zlogger.JSON_LOGGER)

    // check zap config
    zapConfig := appprodConf.GetZapConfig()
    encodingType := zapConfig.Encoding
    devMode := zapConfig.Development

    assert.Equal(t, devMode, false)
    assert.Equal(t, encodingType, "json")
    assert.Equal(t, logLevel, zapConfig.Level.Level())
  })
}