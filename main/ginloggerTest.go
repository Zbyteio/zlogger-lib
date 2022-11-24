package main

import (
	"log"
	"net/http"

	"github.com/Zbyteio/zlogger-lib"
	"github.com/gin-gonic/gin"
)


func main() {
  ginEng := gin.New()



  ZBlocksGinDebugLogger := zlogger.NewGinLogger(gin.Mode())
  ginEng.Use(
    ZBlocksGinDebugLogger.GinRequestLoggerMiddleware(),
  )

  ginEng.GET("/", func(c *gin.Context) {
    c.String(http.StatusOK, "Welcome Gin Server")
  })
   // service connections
  if err := ginEng.Run(); err != nil {
    log.Printf("listen: %s\n", err)
  }
}