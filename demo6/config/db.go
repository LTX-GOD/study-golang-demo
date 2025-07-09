package config

import (
	"log"
	"time"

	"github.com/zsm/CurrencyExchangeApp/global"
	"github.com/zsm/CurrencyExchangeApp/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDB() {
	dsn := AppConfig.Database.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to initialize database,got err: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to configure database, got error: %v", err)
	}

	sqlDB.SetMaxIdleConns(AppConfig.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(AppConfig.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = db.AutoMigrate(&models.ExchangeRate{})
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	global.Db = db
}
