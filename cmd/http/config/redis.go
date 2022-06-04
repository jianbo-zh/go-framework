package config

import (
	"github.com/jianbo-zh/go-config"
)

func redisConfig() {

	config.Add("redis", config.StrMap{
		"host":     config.Env("REDIS_HOST", "127.0.0.1"),
		"port":     config.Env("REDIS_PORT", 6379),
		"username": config.Env("REDIS_USERNAME", ""),
		"password": config.Env("REDIS_PASSWORD", ""),
	})

}
