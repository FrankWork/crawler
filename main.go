package main

import (
	"flag"
	"log"
	"sync"
)

func processLink(wg *sync.WaitGroup, url string, depth int, maxDepth int) {
	defer wg.Done()

	if depth > maxDepth {
		// log.Println("reach maxDepth")
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
	log.Printf("%d %s %s\n", depth, getTitle(doc), url)
	// log.Printf("%d %s\n", depth, getTitle(doc))

	urlCount := getLinks(doc)
	if len(urlCount) == 0 {
		log.Printf("no links: %s\n", url)
	}

	for url2 := range urlCount {
		wg.Add(1)
		go processLink(wg, url2, depth+1, maxDepth)
	}
}

func downloadOnePage() {
	url := "http://mtj.163.com/?from=nietop"
	// url := "http://www.163.com"

	doc := request(url)

	if doc == nil {

	}
	html, err := doc.Html()
	if err != nil {

	}
	log.Println(html)
}
func main_() {
	// benchmarkMain()
	downloadOnePage()
}
func main() {
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
	// rawurl := "http://mtj.163.com/?from=nietop" // gb2312
	// http://x3.163.com/2015/cloud/

	// ========== read error ===============
	// rawurl := "http://v.163.com/open/"
	rawurl := "http://open.163.com"
	// rawurl := "http://open.163.com/movie/2017/5/U/1/MCK194LGV_MCK196RU1.html"
	// http://open.163.com/movie/2017/5/H/T/MCKH42S7I_MCKH5A2HT.html
	// http://open.163.com/movie/2017/5/M/A/MCK3SQQP8_MCK3T2RMA.html
	// http://open.163.com/#f=topnav
	// http://mhws.163.com/?from=nietop
	// http://bz.163.com/?from=nietop

	// ========== server shutdown ===============
	// rawurl := "http://xdw.zhidao.163.com?from=index"
	// http://g.163.com/a?CID=49141&Values=132441315&Redirect=http://www.elianhong.com/zhuanti/ucsd/index.html

	// ========== image ===============
	// http://img2.cache.netease.com/f2e/www/index2014/images/cert.png // image

	// rawurl := "http://www.163.com"
	depth := 0
	maxDepth := 1

	wg.Add(1)
	go processLink(&wg, rawurl, depth, maxDepth)

	wg.Wait()
	// storage()
}
