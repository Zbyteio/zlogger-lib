package zlogger_test

import (
	"testing"

	"github.com/Zbyteio/zlogger-lib"
	"github.com/go-playground/assert/v2"
)

func TestAppLogger(t *testing.T)  {
  t.Run("Test App logger", func(t *testing.T) {
    ZBlocksAppDebugLogger.Debugf("%s", "success print debug via applogger[DEBUG]")
    ZBlocksAppReleaseLogger.Debugf("%s", "success print debug via applogger[RELEASE]")
  })

  t.Run("Test Create logger name", func(t *testing.T) {
    loggername := zlogger.CreateLoggerName("svc_name", "pkg_name", "file_name", "function_name")
    assert.Equal(t, "svc_name.pkg_name.file_name.function_name", loggername)
  })
}