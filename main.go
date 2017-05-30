package main

import (
	"flag"
	"log"
	"sync"
)

func processLink(wg *sync.WaitGroup, url string, depth int, maxDepth int) {
	defer wg.Done()

	if depth > maxDepth {
		log.Println("reach maxDepth")
		return
	}

	conn, resource := redisPoolConnect()
	defer RedisResourcePool.Put(resource)

	if isDuplicateDebug(conn, url) {
		log.Printf("URL: %v is duplicate\n", url)
		return
	}

	doc := request(url)
	if doc == nil {
		return
	}
	// maskDupURL(conn, url)
	maskDupURLDebug(conn, url)
	log.Printf("%d %s %s\n", depth, getTitle(doc), url)

	urlCount := getLinks(doc)
	// log.Printf("total urls: %v\n", len(urlCount))

	for url2 := range urlCount {
		wg.Add(1)
		go processLink(wg, url2, depth+1, maxDepth)
	}
}

func main() {
	// benchmarkMain()
	reset := flag.String("reset", "false", "whther to reset")
	flag.Parse()

	var wg sync.WaitGroup

	defer RedisResourcePool.Close()

	if *reset == "true" {
		if redisDEL("www.163.com") {
			log.Println("reset success")
		} else {
			log.Println("reset failed")
		}
		return
	}

	//rawurl := "http://www.163.com/newsapp"
	//rawurl := "http://gb.corp.163.com/gb/about/overview.html"
	rawurl := "http://www.163.com"
	depth := 0
	maxDepth := 3

	wg.Add(1)
	go processLink(&wg, rawurl, depth, maxDepth)

	wg.Wait()
	// storage()

}
