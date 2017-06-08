package main

import (
	"context"
	"log"

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

var RedisResourcePool *pools.ResourcePool

func init() {
	// Vitess pooling
	factory := func() (pools.Resource, error) {
		var conn redis.Conn
		var err error

		dbOpt := redis.DialDatabase(auth.RedisDb)
		if auth.RedisAuth == "" {
			conn, err = redis.Dial("tcp", auth.RedisHost, dbOpt)
		} else {
			authOpt := redis.DialPassword(auth.RedisAuth)
			conn, err = redis.Dial("tcp", auth.RedisHost, authOpt, dbOpt)
		}

		return ResourceConn{conn}, err
	}

	RedisResourcePool = pools.NewResourcePool(
		factory,
		cfg.RedisPoolCapacity,
		cfg.RedisPoolMaxCapacity,
		cfg.RedisPoolIdleTimeout.Duration)

}

func redisPoolConnect() (ResourceConn, pools.Resource) {
	ctx := context.TODO()
	resource, err := RedisResourcePool.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
	conn := resource.(ResourceConn)
	return conn, resource
}

func redisSET(conn ResourceConn, key string, value string) {
	ok, err := conn.Do("set", key, value)
	if err != nil {
		log.Println("redis set failed: ", err)
	} else {
		log.Printf("redis set key: %v value: %v \n", key, value)
		log.Println(ok)
	}
}

func redisGET(conn ResourceConn, key string) {
	value, err := redis.String(conn.Do("Get", key))
	if err != nil {
		log.Println("redis get failed: ", err)
	} else {
		log.Printf("redis get key: %v value: %v \n", key, value)
	}

}

func redisSISMember(conn ResourceConn, key string, member string) bool {
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

func redisSADD(conn ResourceConn, key string, member string) bool {
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

func redisSREM(conn ResourceConn, key string, member string) bool {
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

func redisDEL(conn ResourceConn, key string) bool {
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

func redisLPUSH(conn ResourceConn, key string, value string) bool {
	ret, err := redis.Int(conn.Do("LPUSH", key, value))
	if err != nil {
		log.Fatal("redis LPUSH failed: ", err)
	}

	if ret == 1 {
		// successfully
		return true
	}

	return false
}

func redisRPOP(conn ResourceConn, key string) string {
	ret, err := redis.String(conn.Do("RPOP", key))
	if err != nil {
		log.Fatal("redis RPOP failed: ", err)
		return ""
	}

	return ret
}

func redisLLen(conn ResourceConn, key string) int {
	ret, err := redis.Int(conn.Do("LLen", key))
	if err != nil {
		log.Fatal("redis LLen failed: ", err)
		return 0
	}

	return ret
}
