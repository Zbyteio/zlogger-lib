package zlogger_test

import (
	"log"
	"testing"

	"github.com/Zbyteio/zlogger-lib"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	gormv2 "gorm.io/gorm"

	gormv1 "github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

func TestGormLogger(t *testing.T) {
	t.Run("test gorm_v2 logger", func(t *testing.T) {
		gormdebugConf := zlogger.NewLoggerConfig("gormlogger_v2", zlogger.JSON_LOGGER, zapcore.DebugLevel)
		gormLogger := zlogger.SetupGormLoggerV2(gormdebugConf)

		dsn := "host=localhost\nuser=postgres\npassword=foxbat\ndbname=postgres\nport=5432\nsslmode=disable\nTimeZone=Asia/Shanghai"
		db, err := gormv2.Open(postgres.Open(dsn), &gormv2.Config{
			Logger: gormLogger,
		})
		if err != nil {
			log.Println(err)
		}
    
		type User struct {
			gormv2.Model
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

  t.Run("test gorm_v1 logger", func(t *testing.T) {
		gormdebugConf := zlogger.NewLoggerConfig("gormlogger_v1", zlogger.JSON_LOGGER, zapcore.DebugLevel)

		dsn := "host=localhost\nuser=postgres\npassword=foxbat\ndbname=postgres\nport=5432\nsslmode=disable\nTimeZone=Asia/Shanghai"
		db, err := gormv1.Open("postgres", dsn)
		if err != nil {
			log.Println(err)
		}

    db.SetLogger(zlogger.SetupGormLoggerV1(gormdebugConf))
    db.LogMode(true)
    
		type User struct {
			gormv1.Model
			Email string `json:"email"`
			Name  string `json:"name"`
		}
		db.AutoMigrate(User{})
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
