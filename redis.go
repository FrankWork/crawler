package main

import (
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/youtube/vitess/go/pools"
)

// ResourceConn adapts a Redigo connection to a Vitess Resource.
type ResourceConn struct {
	redis.Conn
}

// Close methods of struct ResourceConn
func (r ResourceConn) Close() {
	r.Conn.Close()
}

var (
	//RedisHost comment
	RedisHost string
	//RedisDb comment
	RedisDb int
	//RedisPool redigo default pool
	RedisPool *redis.Pool
	// RedisResourcePool is vitess pools wrapper of redigo conn
	RedisResourcePool *pools.ResourcePool
)

func init() {
	RedisHost = "localhost:6379"
	RedisDb = 0

	MaxIdle := 1
	MaxActive := 10

	// pooling
	RedisPool = &redis.Pool{
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

	factory := func() (pools.Resource, error) {
		conn, err := redis.Dial("tcp", ":6379")
		return ResourceConn{conn}, err
	}
	capacity := MaxIdle
	maxCap := MaxActive
	idleTimeout := time.Minute

	RedisResourcePool = pools.NewResourcePool(factory, capacity, maxCap, idleTimeout)

}

func redisConnect() redis.Conn {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal("Connect to redis error", err)
	}
	return conn
}

func redisSET2(conn ResourceConn, key string, value string) {
	// conn := RedisPool.Get()
	// defer conn.Close()

	ok, err := conn.Do("set", key, value)
	if err != nil {
		log.Println("redis set failed: ", err)
	} else {
		log.Printf("redis set key: %v value: %v \n", key, value)
		log.Println(ok)
	}
}

func redisSET(conn redis.Conn, key string, value string) {
	// conn := RedisPool.Get()
	// defer conn.Close()

	ok, err := conn.Do("set", key, value)
	if err != nil {
		log.Println("redis set failed: ", err)
	} else {
		log.Printf("redis set key: %v value: %v \n", key, value)
		log.Println(ok)
	}
}

func redisGET(key string) {
	conn := RedisPool.Get()
	defer conn.Close()

	value, err := redis.String(conn.Do("Get", key))
	if err != nil {
		log.Println("redis get failed: ", err)
	} else {
		log.Printf("redis get key: %v value: %v \n", key, value)
	}

}

func redisSISMember(key string, member string) bool {
	conn := RedisPool.Get()
	defer conn.Close()

	value, err := redis.Int(conn.Do("SISMEMBER", key, member))
	if err != nil {
		log.Printf("key: %v member: %v", key, member)
		log.Fatal("redis SISMEMBER failed: ", err)
	}

	if value == 1 {
		// `member` is a member of `key`
		return true
	}

	return false
}

func redisSADD(key string, member string) bool {
	conn := RedisPool.Get()
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
	conn := RedisPool.Get()
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
	conn := RedisPool.Get()
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
