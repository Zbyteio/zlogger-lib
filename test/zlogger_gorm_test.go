package zlogger_test

import (
	"log"
	"testing"

	"github.com/Zbyteio/zlogger-lib"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGormLogger(t *testing.T) {
	t.Run("test gorm logger", func(t *testing.T) {
		gindebugConf := zlogger.NewLoggerConfig("ginlogger", zlogger.DEBUG_LOGGER, zapcore.DebugLevel)
		zlogger.SetupGormLogger(gindebugConf)

		dsn := "host=localhost\nuser=postgres\npassword=foxbat\ndbname=postgres\nport=5432\nsslmode=disable\nTimeZone=Asia/Shanghai"
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: ZBlocksGormDebugLogger,
		})
		if err != nil {
			log.Println(err)
		}
		type User struct {
			gorm.Model
			Email string `json:"email"`
			Name  string `json:"name"`
		}
		err = db.AutoMigrate(User{})
		if err != nil {
			log.Println(err.Error())
		}
		user := &User{
			Email: "19mandal97@gmail.com",
			Name:  "Sourabh Mandal",
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
