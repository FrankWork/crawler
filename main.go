package main

import (
	"flag"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/garyburd/redigo/redis"
)

func processLink(wg *sync.WaitGroup, url string, depth int, maxDepth int) {
	defer wg.Done()

	if depth > maxDepth {
		log.Println("reach maxDepth")
		return
	}
	if isDuplicate(url) {
		log.Printf("URL: %v is duplicate\n", url)
		return
	}

	doc := request(url)
	log.Printf("request: %v depth: %v", getTitle(doc), depth)

	urlCount := getLinks(doc)
	log.Printf("total urls: %v\n", len(urlCount))

	for url2, count := range urlCount {
		if len(url2) == 0 {
			log.Printf("len(url)==0, count: %v\n", count)
		} else if idx := strings.Index(url2, "javascript"); idx == 0 {
			log.Printf("url: %v, count: %v\n", url2, count)
		} else {
			// doc := request(url)
			wg.Add(1)
			go processLink(wg, url2, depth+1, maxDepth)
		}
	}
}

func redisPoolingConcurrentBenchmark(wg *sync.WaitGroup) {
	n := flag.Int("n", 10, "num of connection")
	flag.Parse()

	log.Println(*n)
	wrapper := func(wg *sync.WaitGroup, key string, value string) {
		conn := RedisClient.Get()
		// conn := redisConnect()
		defer conn.Close()

		defer wg.Done()
		redisSET(conn, key, value)
	}

	for i := 0; i < *n; i++ {
		wg.Add(1)
		go wrapper(wg, "foo", strconv.Itoa(i))
	}

	wg.Wait()
}

func redisPoolingSequentialBenchmark(wg *sync.WaitGroup) {
	n := flag.Int("n", 10, "num of connection")
	flag.Parse()

	// conn := RedisClient.Get()
	conn := redisConnect()
	defer conn.Close()

	log.Println(*n)
	wrapper := func(wg *sync.WaitGroup, conn redis.Conn) {
		defer wg.Done()
		for i := 0; i < *n; i++ {
			redisSET(conn, "foo", strconv.Itoa(i))
		}
	}

	wg.Add(1)
	go wrapper(wg, conn)
	wg.Wait()
}

func main() {
	var wg sync.WaitGroup

	// redisPoolingSequentialBenchmark(&wg)
	redisPoolingConcurrentBenchmark(&wg)
	wg.Wait()

	// reset := flag.String("reset", "false", "whther to reset")
	// flag.Parse()

	// if *reset == "true" {
	// 	if redisDEL("www.163.com") {
	// 		log.Println("reset success")
	// 	} else {
	// 		log.Println("reset failed")
	// 	}
	// 	return
	// }

	// rawurl := "http://www.163.com"
	// depth := 0
	// maxDepth := 3

	// wg.Add(1)
	// go processLink(&wg, rawurl, depth, maxDepth)

	// wg.Wait()
	// storage()

}
