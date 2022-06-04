package dao

import (
	"goframework/internal/env"

	"github.com/jianbo-zh/go-config"
	"gorm.io/driver/sqlite" // Sqlite driver based on GGO
	logger "gorm.io/gorm/logger"

	// "github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"gorm.io/gorm"
)

func Sqlite(sqliteDSN string) (*gorm.DB, error) {

	defaultLevel := logger.Info
	if config.GetString("app.env") == env.Production {
		defaultLevel = logger.Silent
	}

	// github.com/mattn/go-sqlite3
	gormDB, err := gorm.Open(sqlite.Open(sqliteDSN), &gorm.Config{
		Logger: logger.Default.LogMode(defaultLevel),
	})
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}
