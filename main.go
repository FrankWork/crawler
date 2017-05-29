package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

func processLink(url string, depth int, maxDepth int) {
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
	for url2, count := range urlCount {
		if len(url2) == 0 {
			log.Printf("len(url)==0, count: %v\n", count)
		} else if idx := strings.Index(url2, "javascript"); idx == 0 {
			log.Printf("url: %v, count: %v\n", url2, count)
		} else {
			// doc := request(url)
			go processLink(url2, depth+1, maxDepth)
		}
	}
}

func redisPoolingBenchmark() {
	n := flag.Int("n", 10, "num of connection")
	flag.Parse()

	log.Println(*n)
	for i := 0; i < *n; i++ {
		// redisGET("foo")
		go redisSISMember("www.163.com", "123456789")
	}

}

func f(n int) {
	for i := 0; i < 10; i++ {
		fmt.Println(n, ":", i)
	}
}

func main() {
	go f(0)

	// redisPoolingBenchmark()

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
	// go processLink(rawurl, depth, maxDepth)

	// storage()

}
