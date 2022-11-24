package zlogger_test

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Zbyteio/zlogger-lib"
	"github.com/gin-gonic/gin"
)
var (
  ZBlocksAppDebugLogger zlogger.AppLogger
  ZBlocksAppReleaseLogger zlogger.AppLogger

  ZBlocksGinDebugLogger zlogger.GinLogger
  ZBlocksGinReleaseLogger zlogger.GinLogger
)

func init(){
  ZBlocksAppDebugLogger = zlogger.NewZlogger(gin.DebugMode)
  ZBlocksAppReleaseLogger = zlogger.NewZlogger(gin.ReleaseMode)
  ZBlocksGinDebugLogger = zlogger.NewGinLogger(gin.DebugMode)
  ZBlocksGinReleaseLogger = zlogger.NewGinLogger(gin.ReleaseMode)

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