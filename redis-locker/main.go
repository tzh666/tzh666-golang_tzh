package main

import (
	"fmt"
	"log"
	"time"

	// 请注意包名
	"github.com/garyburd/redigo/redis"
	"github.com/google/uuid"
)

type Locker struct {
	conn     redis.Conn
	name     string
	interval int
	tag      string
}

func NewLocker(conn redis.Conn, name string, interval int) *Locker {
	locker := &Locker{
		conn:     conn,
		name:     fmt.Sprintf("locker:redis:%s", name),
		interval: interval,
		tag:      uuid.New().String(),
	}
	go locker.boostrap()
	return locker
}

// 续时
func (l *Locker) boostrap() {
	for range time.Tick(time.Second) {
		// 在例程里面执行 SET key
		l.conn.Do("SET", l.name, l.tag, "EX", l.interval, "NX")
		locker, err := redis.String(l.conn.Do("GET", l.name))
		if err == nil && locker == l.tag {
			// 给原来的key续时
			l.conn.Do("EXPIRE", l.name, l.interval)
		}
	}
}

// 判断是否存在,存在就上锁,从而达到选主的过程
func (l *Locker) Lock() bool {
	l.conn.Do("SET", l.name, l.tag, "EX", l.interval, "NX")
	locker, err := redis.String(l.conn.Do("GET", l.name))
	if err != nil {
		return false
	}
	return locker == l.tag
}

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

	locker := NewLocker(conn, "test", 10)
	for {
		if !locker.Lock() {
			log.Println("not locker")
			time.Sleep(3 * time.Second)
			continue
		}
		log.Println("exec")
		time.Sleep(time.Second * 2)
	}
}
