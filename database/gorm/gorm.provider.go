package gorm

import (
	"errors"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDbProvider(config *viper.Viper) (*gorm.DB, error) {
	driver := config.GetString("db.driver")
	connectionString := config.GetString("db.connection_string")

	if len(driver) == 0 {
		return nil, errors.New("driver is empty")
	}

	if len(connectionString) == 0 {
		return nil, errors.New("connection_string is empty")
	}
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.Debug()

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(3 * time.Minute)
	return db, nil
}
