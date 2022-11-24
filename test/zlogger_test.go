package zlogger_test

import (
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func TestGormLogger(t *testing.T)  {
  t.Run("test gorm logger", func(t *testing.T) {
    dsn := "host=localhost\nuser=postgres\npassword=foxbat\ndbname=postgres\nport=5432\nsslmode=disable\nTimeZone=Asia/Shanghai"
    ZBlocksGormDebugLogger.SetAsDefault()
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
      Logger: ZBlocksGormDebugLogger,
    })
    if err != nil {
      log.Println(err)
    }
    type User struct {
      gorm.Model
      Email string `json:"email"`
      Name string `json:"name"`
    }
    err = db.AutoMigrate(User{})
    if err != nil {
      log.Println(err.Error())
    }
    user := &User{
      Email: "19mandal97@gmail.com",
      Name: "Sourabh Mandal",
    }
    tx := db.Create(user)
    if tx.Error != nil {
      log.Println(tx.Error)
    }
    tx = db.Delete(user)
    if tx.Error != nil {
      log.Println(tx.Error)
    }
  })
}