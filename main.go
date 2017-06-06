package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"sync"

	"time"

	"github.com/BurntSushi/toml"
)

type Auth struct {
	RedisAuth string
}

type Config struct {
	StartURLs []string
	Domains   []string
	MaxDepth  int
}

var (
	cfg  Config
	auth Auth
)

func getExecutablePath() string {
	ex, err := os.Executable()
	if err != nil {
		log.Println("getExecutablePath failed")
		panic(err)
	}
	return path.Dir(ex)
}

func init() {
	// parse config and auth toml file
	cwd := getExecutablePath()

	var err error
	if _, err = toml.DecodeFile(path.Join(cwd, "config.toml"), &cfg); err != nil {
		log.Println("toml.DecodeFile failed")
		log.Printf("toml file path: %s\n", path.Join(cwd, "config.toml"))
		panic(err)
	}

	if _, err = toml.DecodeFile(path.Join(cwd, "auth.toml"), &auth); err != nil {
		log.Println("toml.DecodeFile failed")
		log.Printf("toml file path: %s\n", path.Join(cwd, "auth.toml"))
		panic(err)
	}
}

func parseDoc(wg *sync.WaitGroup, uw *URLWrapper, maxDepth int) {
	if wg != nil {
		defer wg.Done()
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
	// goroutine wait group and redis connection pool
	var wg sync.WaitGroup
	defer RedisResourcePool.Close()

	// start requests
	depth := 0
	for _, rawurl := range cfg.StartURLs {
		uw := NewURLWrapper(rawurl, depth)
		// do NOT dispatch new goroutine in order to block here until the func is finished
		parseDoc(nil, uw, cfg.MaxDepth)
	}

	initNumGo := runtime.NumGoroutine()
	log.Printf("init NumGoroutine: %d\n", initNumGo)

	for runtime.NumGoroutine() > initNumGo || !urlQueue.isEmpty() {
		if urlQueue.isEmpty() {
			time.Sleep(time.Second)
		}
		uw := urlQueue.dequeue()
		if uw != nil {
			wg.Add(1)
			go parseDoc(&wg, uw, cfg.MaxDepth)
		}
	}

	wg.Wait()
	// storage()
}
