package config

import (
	"github.com/jianbo-zh/go-config"
)

func mysqlConfig() {

	config.Add("mysql-main", config.StrMap{
		"host":     config.Env("DB_HOST", "127.0.0.1"),
		"port":     config.Env("DB_PORT", "3306"),
		"database": config.Env("DB_DATABASE", "cloud_new"),
		"username": config.Env("DB_USERNAME", "username"),
		"password": config.Env("DB_PASSWORD", "password"),
		"charset":  config.Env("DB_CHARSET", "utf8mb4"),
	})

	config.Add("mysql-main-readonly", config.StrMap{
		"host":     config.Env("RDB_HOST", "127.0.0.1"),
		"port":     config.Env("RDB_PORT", "3306"),
		"database": config.Env("RDB_DATABASE", "cloud_new"),
		"username": config.Env("RDB_USERNAME", "username"),
		"password": config.Env("RDB_PASSWORD", "password"),
		"charset":  config.Env("RDB_CHARSET", "utf8mb4"),
	})

	config.Add("mysql-log", config.StrMap{
		"host":     config.Env("LOGDB_HOST", "127.0.0.1"),
		"port":     config.Env("LOGDB_PORT", "3306"),
		"database": config.Env("LOGDB_DATABASE", "cloud_new"),
		"username": config.Env("LOGDB_USERNAME", "username"),
		"password": config.Env("LOGDB_PASSWORD", "password"),
		"charset":  config.Env("LOGDB_CHARSET", "utf8mb4"),
	})
}
