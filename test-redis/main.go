package main

import (
	"fmt"
	"log"
	"time"

	// 注意包名
	"github.com/gomodule/redigo/redis"
)

func main() {
	addr := "192.168.1.208:6379"
	password := "huanshao"

	conn, err := redis.Dial("tcp", addr, redis.DialPassword(password))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	/*
			conn, err := conn.Do("SET", "key", "value","NX")
			NX key存在则不更新
		    EX 过期时间 单位 s
	*/
	conn.Do("SET", "time", time.Now().Unix(), "NX", "EX", 60)

	// GET key
	stime, err := redis.Int(conn.Do("GET", "time"))
	fmt.Println(stime, err)

	// args,类似于go操作MySQL查询结果有多个条件那个玩意,加一个条件就 "args=args.Add()" 一次
	// 底下这一堆相当于 conn.Do("SET", "time2", time.Now().Unix(), "NX", "EX", 60)
	args := redis.Args{}
	args = args.Add("time2")
	args = args.Add(time.Now().Unix())
	args = args.Add("NX")
	args = args.Add("EX")
	args = args.Add(60)
	conn.Do("SET", args...)
}
