package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

func redisExample() {
	cmd := flag.String("cmd", "set", "set or get")
	key := flag.String("key", "foo", "key for redis")
	value := flag.String("value", "110", "value for set")
	flag.Parse()

	if *cmd == "set" {
		redisSET(*key, *value)
	} else {
		redisGET(*key)
	}

}

func fpExample() {
	fp := fingerPrint("hello world")
	log.Println(fp)
}

func goQuery() {
	url := "http://www.163.com"
	doc := request(url)

	title := getTitle(doc)
	fmt.Println(title)

	urlCount := getLinks(doc)
	fmt.Println("\"\":", urlCount[""])
	c := 0
	for key, value := range urlCount {
		if idx := strings.Index(key, "http"); idx != 0 {
			fmt.Println(key, ":", value)
		}
		if value > 1 {
			// fmt.Println(key)
			c++
		}
	}
	fmt.Println("c: ", c)
}

func main() {
	rawurl := "http://www.163.com"

	if isDuplicate(rawurl) {
		log.Printf("URL: %v is duplicate\n", rawurl)
	} else {
		log.Printf("URL: %v is NOT duplicate\n", rawurl)

		doc := request(rawurl)
		log.Printf("request: %v", getTitle(doc))

		urlCount := getLinks(doc)
		for url, count := range urlCount {
			if count == 0 {
				log.Printf("url: %v count == 0\n", url)
			}
			if len(url) == 0 {
				log.Printf("len(url)==0, count: %v\n", count)
			}

		}
		// storage()
	}

}
