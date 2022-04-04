// package main

// import (
// 	"fmt"
// 	"log"
// 	"time"

// 	// 请注意包名
// 	"github.com/garyburd/redigo/redis"
// 	"github.com/google/uuid"
// )

// func main() {
// 	addr := "192.168.1.208:6379"
// 	password := "huanshao"
// 	dsn := fmt.Sprintf("redis://tzh:%s@%s/0", password, addr)

// 	// 请注意包名
// 	conn, err := redis.DialURL(dsn, redis.DialPassword(password))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer conn.Close()

// 	// key
// 	key := fmt.Sprintf("locker:redis:%s", "test")

// 	// 唯一标识UUID,可以用UUID那个包
// 	tarid := uuid.New().String()

// 	// 超时时间得大于执行时间
// 	interval := 10

// 	/*
// 		127.0.0.1:6379> keys locker*
// 		1) "locker:redis:test"
// 		127.0.0.1:6379>
// 		127.0.0.1:6379> ttl "locker:redis:test"
// 		(integer) 10
// 	*/
// 	go func() {
// 		// 一秒执行一次
// 		for range time.Tick(time.Second) {
// 			// 在例程里面执行 SET key
// 			conn.Do("SET", key, tarid, "EX", interval, "NX")
// 			locker, err := redis.String(conn.Do("GET", key))
// 			if err == nil && locker == tarid {
// 				// 给原来的key续时
// 				conn.Do("EXPIRE", key, interval)
// 			}
// 		}
// 	}()

// 	for {
// 		conn.Do("SET", key, tarid, "EX", interval, "NX")
// 		locker, err := redis.String(conn.Do("GET", key))
// 		if err != nil {
// 			continue
// 		}
// 		if locker != tarid {
// 			time.Sleep(time.Second * 1)
// 			log.Printf("locker: %s", locker)

// 			continue
// 		}
// 		log.Println("exec")
// 		time.Sleep(time.Second * 2)
// 	}
// }
// //;