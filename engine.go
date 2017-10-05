package main

import (
	"log"
	"runtime"
	"sync"
	"time"
)

// Engine struct contrl the crawler behavior
type Engine struct {
	cfg       *Config
	redis     *RedisClient
	urlQueue  Queue
	dupFilter DupFilter
	wg        *sync.WaitGroup
}

// NewEngine is constrctor of Engine struct
func NewEngine(cfg *Config, redis *RedisClient) *Engine {
	var dupFilter DupFilter
	var urlQueue Queue

	if cfg.Distributed {
		urlQueue = NewURLQueueRedis(cfg.Name+"URLQueue", redis)
		dupFilter = NewDupURLFilterRedis(cfg.Name+"URLSet", redis)
	} else {
		urlQueue = NewURLQueue()
		dupFilter = NewDupURLFilter()
	}
	return &Engine{cfg, redis, urlQueue, dupFilter, new(sync.WaitGroup)}
}

// Close the engine and redis client
func (e *Engine) Close() {
	e.redis.Close()
}

// parserWrapper wait goroutine to stop
// and decides whether to continue based on depth
// FIXME
func parserWrapper(wg *sync.WaitGroup, url *URLWrapper, maxDepth int) {
	if wg != nil {
		defer wg.Done()
	}

	if url.Depth+1 > maxDepth {
		return
	}
}

// Start the engine FIXME: parser TODO:storage
func (e *Engine) Start(parser func(*sync.WaitGroup, *URLWrapper, int)) {
	// start requests
	depth := 0
	for _, rawurl := range e.cfg.StartURLs {
		if e.dupFilter.IsDuplicate(rawurl) {
			url := NewURLWrapper(rawurl, depth)
			// wg is nil: block here until the parser is finished
			parser(nil, url, e.cfg.MaxDepth)
		}
	}

	// run until no url in queue
	initNumGortn := runtime.NumGoroutine()
	log.Printf("init NumGoroutine: %d\n", initNumGortn)

	for runtime.NumGoroutine() > initNumGortn || !e.urlQueue.IsEmpty() {
		if e.urlQueue.IsEmpty() {
			time.Sleep(time.Second)
		} else {
			url := e.urlQueue.Dequeue()
			if url != nil && !e.dupFilter.IsDuplicate(url.RawURL) {
				e.dupFilter.AddURL(url.RawURL)
				e.wg.Add(1)
				go parseDoc(e.wg, url, e.cfg.MaxDepth)
			}
		}
	}
	e.wg.Wait()
	// TODO: storage()
}
