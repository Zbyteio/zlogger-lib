package zlogger_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
)

func TestLogger(t *testing.T)  {

  
  t.Run("Test App logger", func(t *testing.T) {
    ZBlocksAppDebugLogger.Debugf("%s", "success print debug via applogger")
  })

  t.Run("Test Gin logger", func(t *testing.T) {
    ginEng, ginSrv := createServer()
    ginEng.Use(
      ZBlocksGinDebugLogger.GinRequestLoggerMiddleware(),
    )
    go runServerAndClose(ginSrv)
    resp, err := http.Get("http://localhost:8080")
    if err != nil {
      log.Panicln(err)
    }
    respByte, _ := io.ReadAll(resp.Body)
    fmt.Println(string(respByte))
  })
}