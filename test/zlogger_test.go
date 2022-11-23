package zlogger_test

import "testing"

func TestMain(t *testing.T)  {

  
  t.Run("Test App logger", func(t *testing.T) {
    ZBlocksAppLogger.Debugf("%s", "success")
  })
}