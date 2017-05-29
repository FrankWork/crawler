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

func redisPoolingBenchmark(wg *sync.WaitGroup, conn redis.Conn) {
	n := flag.Int("n", 10, "num of connection")
	flag.Parse()

	log.Println(*n)
	for i := 0; i < *n; i++ {
		wg.Add(1)
		// go redisSISMember("www.163.com", "123456789")
		wrapper := func(wg *sync.WaitGroup, conn redis.Conn, key string, value string) {
			defer wg.Done()
			ok, err := conn.Do("set", key, value)
			if err != nil {
				log.Println("redis set failed: ", err)
			} else {
				log.Printf("redis set key: %v value: %v \n", key, value)
				log.Println(ok)
			}
		}
		go wrapper(wg, conn, "foo", strconv.Itoa(i))
	}
}

func close() {
	log.Println("close now")
}
func main() {
	var wg sync.WaitGroup
	conn := RedisClient.Get()
	defer conn.Close()
	defer close()

	redisPoolingBenchmark(&wg, conn)

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

	wg.Wait()
	// storage()

}
