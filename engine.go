package main

import (
	"sync/atomic"
)

// Engine struct contrl the crawler behavior
type Engine struct {
	cfg       *Config
	redis     *RedisClient
	urlQueue  Queue
	dupFilter DupFilter
	nActive   uint64 // number of active worker goroutines
	counter   chan uint64
	// wg        *sync.WaitGroup
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
	return &Engine{cfg, redis, urlQueue, dupFilter,
		0, make(chan uint64)}
}

// Close the engine and redis client
func (e *Engine) Close() {
	e.redis.Close()
}

// parse a document
func (e *Engine) parse(url *URLWrapper, parser func(*URLWrapper) []string) {
	if url.Depth > e.cfg.MaxDepth {
		return
	}

	if !e.dupFilter.IsDuplicate(url.RawURL) {
		atomic.AddUint64(&e.nActive, 1)
		defer func() {
			atomic.AddUint64(&e.nActive, ^uint64(1-1)) // nActive--
			e.counter <- e.nActive
		}()

		e.dupFilter.AddURL(url.RawURL)
		links := parser(url)

		for _, rawurl := range links {
			if !e.dupFilter.IsDuplicate(rawurl) {
				newurl := NewURLWrapper(rawurl, url.Depth+1)
				e.urlQueue.Enqueue(newurl)
			}
		}
	}
}

// Start the engine TODO:storage
func (e *Engine) Start(parser func(*URLWrapper) []string) {
	// start requests
	depth := 0
	for _, rawurl := range e.cfg.StartURLs {
		url := NewURLWrapper(rawurl, depth)
		go e.parse(url, parser)
	}
	// block here till the start urls are all parsed
	for i := 0; i < len(e.cfg.StartURLs); i++ {
		<-e.counter
	}

	// start a goroutine to recive message from channel e.counter
	go func() {
		for e.nActive != 0 {
			<-e.counter
		}
	}()

	for {
		if e.urlQueue.IsEmpty() {
			<-e.counter
			if e.nActive == 0 {
				break
			}
		} else {
			url := e.urlQueue.Dequeue()
			if url != nil && !e.dupFilter.IsDuplicate(url.RawURL) {
				// https://golang.org/doc/faq#goroutines
				// It is practical to create hundreds of thousands of goroutines
				go e.parse(url, parser)
			}
		}
	}

	// e.wg.Wait()
	// TODO: storage()
}
