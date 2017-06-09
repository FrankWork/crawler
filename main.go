package main

import (
	"fmt"
	"log"
	"sync"
)

func parseDoc(wg *sync.WaitGroup, uw *URLWrapper, maxDepth int) {
	if wg != nil {
		defer wg.Done()
	}

	url := uw.RawURL
	doc := request(uw)
	// fmt.Print(doc.Html())
	if doc == nil {
		return
	}
	// $ ./crawler > log.txt
	// fmt.Printf("%d %s %s\n", depth, getTitle(doc), url)
	fmt.Printf("%d %s\n", uw.Depth, getTitle(doc))

	urlCount := getLinks(doc)
	if len(urlCount) == 0 {
		log.Printf("no links: %s\n", url)
	}

	if uw.Depth+1 > maxDepth {
		return
	}
	for newurl := range urlCount {
		// fmt.Println(newurl)
		newuw := NewURLWrapper(newurl, uw.Depth+1)
		urlQueue.enqueue(newuw)
	}
}

func main() {
	engine := NewEngine()
	defer engine.Close()

	engine.Start(parseDoc)
}
