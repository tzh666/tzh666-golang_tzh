package redis

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/go-redis/redis"
)

// 声明一个全局的rdb变量
var rdb *redis.Client

// 初始化连接
func Init() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			viper.GetString("redis.host"),
			viper.GetInt("redis.port"),
		),
		Password: "huanshao",
		DB:       viper.GetInt("redis.db"),
	})
	_, err = rdb.Ping().Result()
	return
}

// 对外暴露一个关闭的方法
func Close() {
	_ = rdb.Close()
}
