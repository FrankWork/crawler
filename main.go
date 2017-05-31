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
		// log.Printf("URL: %v is duplicate\n", url)
		return
	}

	doc := request(url)
	if doc == nil {
		return
	}
	// maskDupURL(conn, url)
	maskDupURLDebug(conn, url)
	log.Printf("%d %s\n", depth, getTitle(doc))

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

	// ========== encoding error ===========
	// rawurl := "http://x3.163.com/"

	// ========== read error ===============
	// rawurl := "http://v.163.com/open/"
	// http://open.163.com/movie/2017/5/U/1/MCK194LGV_MCK196RU1.html
	// http://open.163.com/movie/2017/5/H/T/MCKH42S7I_MCKH5A2HT.html
	// http://open.163.com/#f=topnav

	// ========== server shutdown ===============
	// http://xdw.zhidao.163.com?from=index
	// http://g.163.com/a?CID=49141&Values=132441315&Redirect=http://www.elianhong.com/zhuanti/ucsd/index.html

	// ========== image ===============
	// http://img2.cache.netease.com/f2e/www/index2014/images/cert.png // image

	rawurl := "http://www.163.com"
	depth := 0
	maxDepth := 3

	wg.Add(1)
	go processLink(&wg, rawurl, depth, maxDepth)

	wg.Wait()
	// storage()

}
