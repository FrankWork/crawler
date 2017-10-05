package main

import (
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/toml"
)

// Config struct to store configuration info
type Config struct {
	Name                 string
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

// Duration struct for Config struct
type Duration struct {
	time.Duration //anonymous field
}

// UnmarshalText method of Duration for parse text
func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

// NewConfig returns a Config object
func NewConfig(configPath string) *Config {
	var cfg Config
	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		panic(err)
	}
	return &cfg
}

// FIXME
func parseDoc(url *URLWrapper) {
	rawurl := url.RawURL
	doc := request(rawurl)
	// fmt.Print(doc.Html())
	if doc == nil {
		return
	}
	fmt.Printf("%d %s\n", url.Depth, getTitle(doc))

	urlCount := getLinks(doc)
	if len(urlCount) == 0 {
		log.Printf("no links: %s\n", url)
	}

	for newurl := range urlCount {
		// fmt.Println(newurl)
		newuw := NewURLWrapper(newurl, uw.Depth+1)
		// FIXME: need an urlQueue across the function
		urlQueue.enqueue(newuw)
	}
}

func main() {
	configPath := "config.toml"
	cfg := NewConfig(configPath)

	redis := NewRedisClient(cfg.RedisHost, cfg.RedisAuth, cfg.RedisDb,
		cfg.RedisPoolCapacity, cfg.RedisPoolMaxCapacity,
		cfg.RedisPoolIdleTimeout.Duration)

	engine := NewEngine(cfg, redis)
	defer engine.Close()

	engine.Start(parseDoc)
}
