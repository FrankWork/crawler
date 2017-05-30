package main

import (
	"context"
	"flag"
	"log"
	"strconv"
	"sync"

	"github.com/garyburd/redigo/redis"
)

// Vitess pooling
func redisPoolingBenchmarkAsync(wg *sync.WaitGroup, n int) {

	op := func(conn ResourceConn, key string, value string) {
		cmd := "set"
		reply, err := conn.Do(cmd, key, value)
		if err != nil {
			log.Println("redis set failed: ", err)
		} else {
			log.Printf("redis set key: %v value: %v \n", key, value)
			log.Println(reply)
		}
	}

	ctx := context.TODO()
	wrapper := func(wg *sync.WaitGroup, key string, value string) {
		defer wg.Done()
		resource, err := RedisResourcePool.Get(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer RedisResourcePool.Put(resource)
		conn := resource.(ResourceConn)

		op(conn, key, value)
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go wrapper(wg, "foo", strconv.Itoa(i))
	}

	wg.Wait()
}

// generate one active connection
func redisPoolingBenchmarkSync(wg *sync.WaitGroup, n int) {
	conn := redisConnect()
	defer conn.Close()

	op := func(conn redis.Conn, key string, value string) {
		cmd := "set"
		reply, err := conn.Do(cmd, key, value)
		if err != nil {
			log.Println("redis set failed: ", err)
		} else {
			log.Printf("redis set key: %v value: %v \n", key, value)
			log.Println(reply)
		}
	}

	wrapper := func(wg *sync.WaitGroup, conn redis.Conn) {
		defer wg.Done()
		for i := 0; i < n; i++ {
			op(conn, "foo", strconv.Itoa(i))
		}
	}

	wg.Add(1)
	go wrapper(wg, conn)
	wg.Wait()
}

// main func to be used in main()
func benchmarkMain() {
	n := flag.Int("n", 10, "num of connections")
	mode := flag.String("m", "sync", "mode: async or sync")
	flag.Parse()

	log.Printf("n   : %d\n", *n)
	log.Printf("mode: %s\n", *mode)

	var wg sync.WaitGroup

	defer RedisResourcePool.Close()

	if *mode == "sync" {
		redisPoolingBenchmarkSync(&wg, *n)
	} else if *mode == "async" {
		redisPoolingBenchmarkAsync(&wg, *n)
	} else {
		log.Println("Oops!")
	}

	wg.Wait()
}
