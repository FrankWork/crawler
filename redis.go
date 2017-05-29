package main

import (
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	//RedisClient comment
	RedisClient *redis.Pool
	//RedisHost comment
	RedisHost string
	//RedisDb comment
	RedisDb int
	//MaxIdle comment
	MaxIdle int
	//MaxActive comment
	MaxActive int
)

func init() {
	RedisHost = "localhost:6379"
	RedisDb = 0
	MaxIdle = 1
	MaxActive = 1000

	// pooling
	RedisClient = &redis.Pool{
		MaxIdle:     MaxIdle,
		MaxActive:   MaxActive,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", RedisHost)
			if err != nil {
				return nil, err
			}
			c.Do("SELECT", RedisDb)
			return c, nil
		},
	}
}

func redisConnect() redis.Conn {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal("Connect to redis error", err)
	}

	return conn
}
func redisSET(key string, value string) {
	conn := RedisClient.Get() //redisConnect()
	defer conn.Close()

	ok, err := conn.Do("set", key, value)
	if err != nil {
		log.Println("redis set failed: ", err)
	} else {
		log.Printf("redis set key: %v value: %v \n", key, value)
		log.Println(ok)
	}
}

func redisGET(key string) {
	conn := RedisClient.Get() //redisConnect()
	defer conn.Close()

	value, err := redis.String(conn.Do("Get", key))
	if err != nil {
		log.Println("redis get failed: ", err)
	} else {
		log.Printf("redis get key: %v value: %v \n", key, value)
	}

}

func redisSISMember(key string, member string) bool {
	conn := RedisClient.Get() //redisConnect()
	defer conn.Close()

	value, err := redis.Int(conn.Do("SISMEMBER", key, member))
	if err != nil {
		log.Fatal("redis SISMEMBER failed: ", err)
	}

	if value == 1 {
		// `member` is a member of `key`
		return true
	}

	return false
}

func redisSADD(key string, member string) bool {
	conn := RedisClient.Get() //redisConnect()
	defer conn.Close()

	value, err := redis.Int(conn.Do("SADD", key, member))
	if err != nil {
		log.Fatal("redis SADD failed: ", err)
	}

	if value == 1 {
		// add `member` to `key` successfully
		return true
	}

	return false
}

func redisSREM(key string, member string) bool {
	conn := RedisClient.Get() //redisConnect()
	defer conn.Close()

	value, err := redis.Int(conn.Do("SREM", key, member))
	if err != nil {
		log.Fatal("redis SREM failed: ", err)
	}

	if value == 1 {
		// remove `member` from `key` successfully
		return true
	}

	return false
}

func redisDEL(key string) bool {
	conn := RedisClient.Get() //redisConnect()
	defer conn.Close()

	value, err := redis.Int(conn.Do("Del", key))
	if err != nil {
		log.Fatal("redis SREM failed: ", err)
	}

	if value == 1 {
		// del `key` successfully
		return true
	}

	return false
}
