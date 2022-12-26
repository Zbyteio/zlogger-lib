package zlogger_test

import (
	"testing"

	"github.com/Zbyteio/zlogger-lib"
)


func TestPackageInit(t *testing.T) {
  t.Run("Test for init function", func(t *testing.T) {

    zlogger.GetAppLogger().Info("Logged via fefault app logger")
    zlogger.SetupLoggerWithConfig("zlogger", zlogger.DEBUG_LOGGER, nil, nil)
    zlogger.GetAppLogger().Info("Logged via new app logger")
  })
}