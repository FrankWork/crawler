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

// RedisClient is a Vitess connection pool
type RedisClient struct {
	pool *pools.ResourcePool
}

// NewRedisClient return a RedisClient object
func NewRedisClient(host, auth string, db, cap, maxCap int,
	timeout time.Duration) *RedisClient {
	// Vitess pooling
	factory := func() (pools.Resource, error) {
		var conn redis.Conn
		var err error

		dbOpt := redis.DialDatabase(db)
		if auth == "" {
			conn, err = redis.Dial("tcp", host, dbOpt)
		} else {
			authOpt := redis.DialPassword(auth)
			conn, err = redis.Dial("tcp", host, authOpt, dbOpt)
		}
		return ResourceConn{conn}, err
	} // factory
	pool := pools.NewResourcePool(factory, cap, maxCap, timeout)
	return &RedisClient{pool}
}

// Close the connection pools of RedisClient
func (rc *RedisClient) Close() {
	rc.pool.Close()
}

// Connect get a connection from the pool
func (rc *RedisClient) Connect() (ResourceConn, pools.Resource) {
	ctx := context.TODO()
	resource, err := rc.pool.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
	conn := resource.(ResourceConn)
	return conn, resource
}

func (rc *RedisClient) Set(key string, value string) {
	conn, resource := rc.Connect()
	defer rc.pool.Put(resource)

	ok, err := conn.Do("SET", key, value)
	if err != nil {
		log.Println("redis SET failed: ", err)
	} else {
		log.Printf("redis SET key: %v value: %v \n", key, value)
		log.Println(ok)
	}
}

func (rc *RedisClient) Get(key string) {
	conn, resource := rc.Connect()
	defer rc.pool.Put(resource)

	value, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Println("redis GET failed: ", err)
	} else {
		log.Printf("redis GET key: %v value: %v \n", key, value)
	}
}

func (rc *RedisClient) SIsMember(key string, member string) bool {
	conn, resource := rc.Connect()
	defer rc.pool.Put(resource)

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

func (rc *RedisClient) SAdd(key string, member string) bool {
	conn, resource := rc.Connect()
	defer rc.pool.Put(resource)

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

func (rc *RedisClient) SRem(key string, member string) bool {
	conn, resource := rc.Connect()
	defer rc.pool.Put(resource)

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

func (rc *RedisClient) Del(key string) bool {
	conn, resource := rc.Connect()
	defer rc.pool.Put(resource)

	value, err := redis.Int(conn.Do("DEL", key))
	if err != nil {
		log.Fatal("redis DEL failed: ", err)
	}

	if value == 1 {
		// del `key` successfully
		return true
	}

	return false
}

func (rc *RedisClient) LPush(key string, value string) bool {
	conn, resource := rc.Connect()
	defer rc.pool.Put(resource)

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

func (rc *RedisClient) RPop(key string) string {
	conn, resource := rc.Connect()
	defer rc.pool.Put(resource)

	ret, err := redis.String(conn.Do("RPOP", key))
	if err != nil {
		log.Fatal("redis RPOP failed: ", err)
		return ""
	}

	return ret
}

func (rc *RedisClient) LLen(key string) int {
	conn, resource := rc.Connect()
	defer rc.pool.Put(resource)

	ret, err := redis.Int(conn.Do("LLen", key))
	if err != nil {
		log.Fatal("redis LLen failed: ", err)
		return 0
	}

	return ret
}
