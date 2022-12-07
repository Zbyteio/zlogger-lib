package zlogger_test

import (
	"testing"
)

func TestAppLogger(t *testing.T)  {
  t.Run("Test App logger", func(t *testing.T) {
    //ZBlocksAppDebugLogger.Debugf("%s", "success print debug via applogger[DEBUG]")
    ZBlocksAppReleaseLogger.Debugf("%s", "success print debug via applogger[RELEASE]")
  })
}