package zlogger_test

import "github.com/Zbyteio/zlogger-lib"
var (
  ZBlocksAppLogger zlogger.AppLogger
)

func init(){
  ZBlocksAppLogger = zlogger.NewZlogger()
}