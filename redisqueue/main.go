package main

import (
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
)

func main() {
	addr := "192.168.1.208:6379"
	password := "huanshao"
	dsn := fmt.Sprintf("redis://tzh:%s@%s/0", password, addr)

	// 请注意包名
	conn, err := redis.DialURL(dsn, redis.DialPassword(password))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

}
