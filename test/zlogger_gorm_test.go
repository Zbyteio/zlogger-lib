package zlogger_test

import (
	"log"
	"testing"

	"github.com/Zbyteio/zlogger-lib"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

func TestGormLogger(t *testing.T) {
	t.Run("test gorm_v2 logger", func(t *testing.T) {
		dsn := "host=localhost\nuser=postgres\npassword=foxbat\ndbname=postgres\nport=5432\nsslmode=disable\nTimeZone=Asia/Shanghai"
		db, err := gorm.Open(postgres.Open(dsn))
		if err != nil {
			log.Println(err)
		}
		
		gormdebugConf := zlogger.NewLoggerConfig("gormlogger_v2", zlogger.DEBUG_LOGGER, zapcore.DebugLevel)
		zlogger.SetupGormLogger(db, gormdebugConf)    
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
