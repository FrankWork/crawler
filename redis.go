package main

import (
	"context"
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
	// RedisResourcePool is vitess pools wrapper of redigo conn
	RedisResourcePool *pools.ResourcePool
)

func init() {
	RedisHost = "localhost:6379"
	RedisDb = 0

	// Vitess pooling
	factory := func() (pools.Resource, error) {
		conn, err := redis.Dial("tcp", ":6379")
		return ResourceConn{conn}, err
	}
	capacity := 1
	maxCap := 10
	idleTimeout := time.Minute

	RedisResourcePool = pools.NewResourcePool(factory, capacity, maxCap, idleTimeout)

}

func redisPoolConnect() (ResourceConn, pools.Resource) {
	ctx := context.TODO()
	resource, err := RedisResourcePool.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// defer RedisResourcePool.Put(resource)
	conn := resource.(ResourceConn)

	return conn, resource
}

func redisConnect() redis.Conn {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal("Connect to redis error", err)
	}
	return conn
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

func redisSISMember(conn ResourceConn, key string, member string) bool {
	// conn := RedisPool.Get()
	// defer conn.Close()

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
