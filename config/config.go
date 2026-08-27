package config

import (
	"log"
	"os"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db *gorm.DB
)

func ConnectDB() error {
	dbUserName := os.Getenv("DB_USERNAME")
	dbUserPassword := os.Getenv("DB_PASSWORD")
	dbProtocol := os.Getenv("DB_PROTOCOL")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_DATABASE")

	metricsLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)

	dsn := dbUserName + ":" + dbUserPassword + "@" + dbProtocol + "(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: metricsLogger})

	if err != nil {
		return err
	}

	sqlDB, err := d.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(25)
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetConnMaxLifetime(6 * time.Minute)
    sqlDB.SetConnMaxIdleTime(1 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return err
	}

	db = d
	return nil
}

func GetDB() *gorm.DB {
	return db
}
