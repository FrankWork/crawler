package main

import (
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
	RedisPoolIdleTimeout duration
}

// Duration struct for Config struct
type duration struct {
	time.Duration //anonymous field
}

// UnmarshalText method of Duration for parse text
func (d *duration) UnmarshalText(text []byte) error {
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

// parseDoc download and parse a document
func parseDoc(url *URLWrapper) (links []string) {
	rawurl := url.RawURL
	doc := Request(rawurl)
	// fmt.Print(doc.Html())
	if doc == nil {
		return
	}
	log.Printf("%d %s\n", url.Depth, GetTitle(doc))

	links = GetAllLinks(doc)
	// if len(links) == 0 {
	// 	log.Printf("no links: %s\n", url)
	// }
	return
}

func main() {
	// f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE, 0666) // |os.O_APPEND
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()
	// log.SetOutput(f)

	configPath := "config.toml"
	cfg := NewConfig(configPath)

	redis := NewRedisClient(cfg.RedisHost, cfg.RedisAuth, cfg.RedisDb,
		cfg.RedisPoolCapacity, cfg.RedisPoolMaxCapacity,
		cfg.RedisPoolIdleTimeout.Duration)

	engine := NewEngine(cfg, redis)
	defer engine.Close()

	engine.Start(parseDoc)
}
