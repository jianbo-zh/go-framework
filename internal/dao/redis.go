package dao

import (
	"github.com/gomodule/redigo/redis"
)

type RedisConfig struct {
	Network  string // 网络协议： tcp
	Address  string // 地址格式： host:port
	Username string // 用户名
	Password string // 密码
	Database int    // 数据库
}

func Redis(config RedisConfig) (*redis.Pool, error) {
	pool := &redis.Pool{
		MaxIdle:     8,
		MaxActive:   0,
		IdleTimeout: 100,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				config.Network,
				config.Address,
				redis.DialUsername(config.Username),
				redis.DialPassword(config.Password),
				redis.DialDatabase(config.Database),
			)
		},
	}

	return pool, nil
}
