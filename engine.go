package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"time"

	"sync"

	"github.com/BurntSushi/toml"
)

type Config struct {
	StartURLs            []string
	Domains              []string
	MaxDepth             int
	Distributed          bool
	RedisHost            string
	RedisDb              int
	RedisAuth            string
	RedisPoolCapacity    int
	RedisPoolMaxCapacity int
	RedisPoolIdleTimeout Duration
}
type Duration struct {
	time.Duration //anonymous field
}

func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func getExecutablePath() string {
	ex, err := os.Executable()
	if err != nil {
		log.Println("getExecutablePath failed")
		panic(err)
	}
	return path.Dir(ex)
}

type Engine struct {
	cfg       *Config
	urlQueue  Queue
	dupFilter DuplicateURLFilter
	rc        *RedisClient
	wg        *sync.WaitGroup
}

func NewEngine() *Engine {
	// parse config and auth toml file
	cwd := getExecutablePath()

	var err error
	var cfg Config

	if _, err = toml.DecodeFile(path.Join(cwd, "config.toml"), &cfg); err != nil {
		log.Println("toml.DecodeFile failed")
		log.Printf("toml file path: %s\n", path.Join(cwd, "config.toml"))
		panic(err)
	}
	fmt.Println("decode config.toml")

	var dupFilter DuplicateURLFilter
	var urlQueue Queue
	var rc *RedisClient
	rc = nil

	if cfg.Distributed {
		fmt.Println("Distribute")
		rc := NewRedisClient(cfg.RedisHost, cfg.RedisAuth, cfg.RedisDb,
			cfg.RedisPoolCapacity, cfg.RedisPoolMaxCapacity,
			cfg.RedisPoolIdleTimeout.Duration)
		dupFilter = NewDuplicateURLFilterDistribute("defaultKey", rc)
		urlQueue = NewURLQueueDistributed("urlQueue", rc)
	} else {
		fmt.Println("Local")
		dupFilter = NewDuplicateURLFilterLocal()
		urlQueue = NewURLQueueLocal()
	}
	return &Engine{&cfg, urlQueue, dupFilter, rc, new(sync.WaitGroup)}
}

func (e *Engine) GetConfig() *Config {
	return e.cfg
}

func (e *Engine) Close() {
	if e.rc != nil {
		e.rc.Close()
	}
}

func (e *Engine) Start(parser func(*sync.WaitGroup, *URLWrapper, int)) {
	cfg := e.GetConfig()

	// start requests
	depth := 0
	for _, rawurl := range cfg.StartURLs {
		if e.dupFilter.isDuplicate(rawurl) {
			uw := NewURLWrapper(rawurl, depth)
			// wg is nil: block here until the parser is finished
			parser(nil, uw, cfg.MaxDepth)
		}
	}

	// run until no url in queue
	initNumGortn := runtime.NumGoroutine()
	log.Printf("init NumGoroutine: %d\n", initNumGortn)

	for runtime.NumGoroutine() > initNumGortn || !e.urlQueue.isEmpty() {
		if e.urlQueue.isEmpty() {
			time.Sleep(time.Second)
		}
		uw := e.urlQueue.dequeue()
		if uw != nil && !e.dupFilter.isDuplicate(uw.RawURL) {
			e.dupFilter.addURL(uw.RawURL)
			e.wg.Add(1)
			go parseDoc(e.wg, uw, e.cfg.MaxDepth)
		}
	}

	e.wg.Wait()
	// storage()
}
