package main

import (
	"log"
	"strings"
	"sync"
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

func main() {
	benchmarkMain()

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
