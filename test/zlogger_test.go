package zlogger_test

import (
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLogger(t *testing.T)  {

  
  t.Run("Test App logger", func(t *testing.T) {
    ZBlocksAppDebugLogger.Debugf("%s", "success print debug via applogger")
  })

  t.Run("Test Gin logger", func(t *testing.T) {
    ginEng, ginSrv := createServer()
    gin.SetMode(gin.DebugMode)
    ginEng.Use(
      ZBlocksGinDebugLogger.GinRequestLoggerMiddleware(),
    )
    ginEng.GET("/abc", func(c *gin.Context) {
      c.String(http.StatusOK, "Welcome Gin Server")
    })

    go runServerAndClose(ginSrv)
    resp, err := http.Get("http://localhost:8080/abc")
    if err != nil {
      log.Panicln(err)
    }
    respByte, _ := io.ReadAll(resp.Body)
    ZBlocksAppDebugLogger.Debug(string(respByte))
  })
}