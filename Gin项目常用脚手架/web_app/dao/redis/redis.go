package redis

import (
	"fmt"
	"web_app/settings"

	"github.com/go-redis/redis"
)

// 声明一个全局的rdb变量
var rdb *redis.Client

// 初始化连接
func Init(cfg *settings.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port,
		),
		Password: cfg.Password,
		DB:       cfg.DBName,
	})
	_, err = rdb.Ping().Result()
	return
}

// 对外暴露一个关闭的方法
func Close() {
	_ = rdb.Close()
}
