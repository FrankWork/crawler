package main

import (
	"fmt"
	"log"
	"sync"
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
	if isDuplicateDebug(conn, url) {
		// log.Printf("URL: %v is duplicate\n", url)
		return
	}

	doc := request(uw)
	if doc == nil {
		return
	}
	// maskDupURL(conn, url)
	maskDupURLDebug(conn, url)
	// maskDupURLSet(conn, url)
	// $ ./crawler > log.txt
	fmt.Printf("%d %s %s\n", depth, getTitle(doc), url)
	// log.Printf("%d %s\n", depth, getTitle(doc))

	urlCount := getLinks(doc)
	if len(urlCount) == 0 {
		log.Printf("no links: %s\n", url)
	}
	for newurl := range urlCount {
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

	// rawurl := "http://mtj.163.com/?from=nietop" // gb2312
	rawurl := "http://open.163.com" // read error
	// rawurl := "http://open.163.com/movie/2017/5/U/1/MCK194LGV_MCK196RU1.html"
	// rawurl := "http://xdw.zhidao.163.com?from=index"//server shutdown
	// http://img2.cache.netease.com/f2e/www/index2014/images/cert.png // image
	// rawurl := "http://www.163.com"

	depth := 0
	maxDepth := 1

	var wg sync.WaitGroup
	defer RedisResourcePool.Close()

	// init url
	uw := NewURLWrapper(rawurl, depth)
	if uw != nil {
		urlQueue.enqueue(uw)
	}
	processDoc(nil, uw, maxDepth)

	for !urlQueue.isEmpty() {
		uw := urlQueue.dequeue()
		wg.Add(1)
		go processDoc(&wg, uw, maxDepth)
	}

	wg.Wait()
	// storage()
}
