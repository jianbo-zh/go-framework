package dao

import (
	"time"

	"goframework/internal/env"

	"github.com/jianbo-zh/go-config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger "gorm.io/gorm/logger"
)

func Mysql(mysqlDSN string) (*gorm.DB, error) {

	defaultLevel := logger.Info
	if config.GetString("app.env") == env.Production {
		defaultLevel = logger.Silent
	}

	gormDB, err := gorm.Open(mysql.New(
		mysql.Config{
			DSN: mysqlDSN,
		}),
		&gorm.Config{
			Logger: logger.Default.LogMode(defaultLevel),
		})
	if err != nil {
		return nil, err
	}

	gormDB.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 auto_increment=1")

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return gormDB, nil
}
