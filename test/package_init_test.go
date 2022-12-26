package zlogger_test

import (
	"testing"

	"github.com/Zbyteio/zlogger-lib"
	"github.com/go-playground/assert/v2"
)


func TestPackageInit(t *testing.T) {


  t.Run("Test for init function", func(t *testing.T) {
    zlogger.SetupLoggerWithConfig("default", zlogger.DEBUG_LOGGER, nil, nil)
    assert.NotEqual(t, zlogger.GetDefaultAppLogger(), nil)
  })
}