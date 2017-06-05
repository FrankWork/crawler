package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func processDoc(wg *sync.WaitGroup, uw *URLWrapper, maxDepth int) {
	if wg != nil {
		defer wg.Done()
	}

	depth := uw.Depth
	if depth > maxDepth {
		// log.Println("reach maxDepth")
		return
	}

	conn, resource := redisPoolConnect()
	defer RedisResourcePool.Put(resource)

	url := uw.RawURL
	if isDuplicateSet(conn, url) {
		// log.Printf("URL: %v is duplicate\n", url)
		return
	}

	doc := request(uw)
	// fmt.Print(doc.Html())

	if doc == nil {
		return
	}
	// maskDupURL(conn, url)
	maskDupURLSet(conn, url)
	// maskDupURLSet(conn, url)
	// $ ./crawler > log.txt
	// fmt.Printf("%d %s %s\n", depth, getTitle(doc), url)
	fmt.Printf("%d %s\n", depth, getTitle(doc))

	urlCount := getLinks(doc)
	if len(urlCount) == 0 {
		log.Printf("no links: %s\n", url)
	}
	for newurl := range urlCount {
		// fmt.Println(newurl)
		newuw := NewURLWrapper(newurl, depth+1)
		urlQueue.enqueue(newuw)
	}
}

func downloadOnePage() {
	// url := "http://mtj.163.com/?from=nietop"
	// // url := "http://www.163.com"

	// rw := NewRequestWrapper(url, 0)
	// doc := request(rw)

	// if doc == nil {

	// }
	// html, err := doc.Html()
	// if err != nil {

	// }
	// log.Println(html)
}

var domains = []string{"163.com"}

// "163.com", "netease.com","kaola.com", "bobo.com", "126.com", "youdao.com", "lofter.com", "126.net"

func main() {
	// benchmarkMain()
	// downloadOnePage()

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

	// http://xf.house.163.com/bj/0RCG.html#lpk-lpxxzc-ss // No encoding!
	rawurl := "http://music.163.com" // html codes
	// rawurl := "http://g.qq.com?ADTAG=pcqq.home.sidenav" // connection reset by peer

	// rawurl := "http://www.163.com"
	// rawurl := "http://www.qq.com"
	// rawurl := "http://www.baidu.com"

	depth := 0
	maxDepth := 0

	var wg sync.WaitGroup
	defer RedisResourcePool.Close()

	// init url
	uw := NewURLWrapper(rawurl, depth)
	processDoc(nil, uw, maxDepth)

	// emptyQueueCount := 0
	for { //emptyQueueCount < 100 {
		if urlQueue.isEmpty() {
			// emptyQueueCount++
			// log.Println("empty queue")
			time.Sleep(3 * time.Microsecond)
		} else {
			uw := urlQueue.dequeue()
			wg.Add(1)
			go processDoc(&wg, uw, maxDepth)
		}
	}
	// for !urlQueue.isEmpty() {
	// 	uw := urlQueue.dequeue()
	// 	wg.Add(1)
	// 	go processDoc(&wg, uw, maxDepth)
	// }

	wg.Wait() // unreachable code
	// storage()
}
