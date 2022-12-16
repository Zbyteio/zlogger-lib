package zlogger_test

import (
	"net/http"
	"testing"

	"github.com/Zbyteio/zlogger-lib"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
)

func TestGinLogger(t *testing.T) {
	var ZBlocksGinDebugLogger gin.LoggerConfig = zlogger.NewGinLoggerConfig(
		zlogger.NewLoggerConfig(
			"ginlogger",
			zlogger.DEBUG_LOGGER,
			zapcore.DebugLevel),
		[]string{},
	)

	t.Run("Test Gin logger", func(t *testing.T) {
		ginEng, ginSrv := createServer()
		ginEng.Use(
			gin.LoggerWithConfig(ZBlocksGinDebugLogger),
		)

		ginEng.GET("/debug-api", func(c *gin.Context) {
			c.String(http.StatusOK, "Welcome Gin Server")
		})

		ginEng.GET("/release-api", func(c *gin.Context) {
			c.String(http.StatusOK, "Welcome Gin Server")
		})

		go runServerAndClose(ginSrv)
		http.Get("http://localhost:8080/debug-api")
		http.Get("http://localhost:8080/release-api")
	})
}
