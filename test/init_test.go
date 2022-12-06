package zlogger_test

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Zbyteio/zlogger-lib"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
)
var (
  ZBlocksAppDebugLogger zlogger.AppLogger
  ZBlocksAppReleaseLogger zlogger.AppLogger

  ZBlocksGinDebugLogger zlogger.GinLogger
  ZBlocksGinReleaseLogger zlogger.GinLogger

  ZBlocksGormDebugLogger zlogger.GormLogger
  ZBlocksGormReleaseLogger zlogger.GormLogger
)

func init(){
  debugConf := zlogger.NewLoggerConfig("applogger", zlogger.DEBUG_LOGGER, zapcore.DebugLevel)
  prodConf := zlogger.NewLoggerConfig("applogger", zlogger.JSON_LOGGER, zapcore.InfoLevel)
  
  ZBlocksAppDebugLogger = zlogger.NewAppLogger(debugConf)
  ZBlocksAppReleaseLogger = zlogger.NewAppLogger(prodConf)
  
  ZBlocksGinDebugLogger = zlogger.NewGinLogger(debugConf)
  ZBlocksGinReleaseLogger = zlogger.NewGinLogger(prodConf)

  zlogger.SetupGormLogger(debugConf)
}

func createServer() (*gin.Engine, *http.Server){
  ginEng := gin.New()

  ginSrv := &http.Server{
    Addr:    ":8080",
    Handler: ginEng,
  }

  return ginEng, ginSrv
}

func runServerAndClose(ginSrv *http.Server) {
  // service connections
  if err := ginSrv.ListenAndServe(); err != nil {
    log.Printf("listen: %s\n", err)
  }
  // Wait for a timeout of 3 seconds gracefully shutdown the server with
  ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
  defer ginSrv.Shutdown(ctx)
  defer cancel()
}