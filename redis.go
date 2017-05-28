package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/garyburd/redigo/redis"
)

func redisConnect() redis.Conn {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
		os.Exit(1)
	}

	return conn
}
func redisSet(key string, value string) {
	conn := redisConnect()
	defer conn.Close()

	ok, err := conn.Do("set", key, value)
	if err != nil {
		fmt.Println("redis set failed: ", err)
	} else {
		fmt.Printf("redis set key: %v value: %v \n", key, value)
		fmt.Println(ok)
	}
}

func redisGet(key string) {
	conn := redisConnect()
	defer conn.Close()

	value, err := redis.String(conn.Do("Get", key))
	if err != nil {
		fmt.Println("redis get failed: ", err)
	} else {
		fmt.Printf("redis get key: %v value: %v \n", key, value)
	}

}

func redisExample() {
	cmd := flag.String("cmd", "set", "set or get")
	key := flag.String("key", "foo", "key for redis")
	value := flag.String("value", "110", "value for set")
	flag.Parse()

	if *cmd == "set" {
		redisSet(*key, *value)
	} else {
		redisGet(*key)
	}

}
