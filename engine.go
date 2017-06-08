package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/BurntSushi/toml"
)

type Auth struct {
	RedisHost string
	RedisDb   int
	RedisAuth string
}

type Config struct {
	StartURLs            []string
	Domains              []string
	MaxDepth             int
	Distributed          bool
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
	cfg       Config
	auth      Auth
	urlQueue  Queue
	dupFilter DuplicateURLFilter
	rc        *RedisClient
}

func NewEngine() *Engine {
	// parse config and auth toml file
	cwd := getExecutablePath()

	var err error
	var cfg Config
	var auth Auth

	if _, err = toml.DecodeFile(path.Join(cwd, "config.toml"), &cfg); err != nil {
		log.Println("toml.DecodeFile failed")
		log.Printf("toml file path: %s\n", path.Join(cwd, "config.toml"))
		panic(err)
	}
	fmt.Println("decode config.toml")

	if _, err = toml.DecodeFile(path.Join(cwd, "auth.toml"), &auth); err != nil {
		log.Println("toml.DecodeFile failed")
		log.Printf("toml file path: %s\n", path.Join(cwd, "auth.toml"))
		panic(err)
	}

	rc := NewRedisClient(auth.RedisHost, auth.RedisAuth, auth.RedisDb,
		cfg.RedisPoolCapacity, cfg.RedisPoolMaxCapacity,
		cfg.RedisPoolIdleTimeout.Duration)

	var dupFilter DuplicateURLFilter
	var urlQueue Queue

	if cfg.Distributed {
		fmt.Println("Distribute")
		dupFilter = NewDuplicateURLFilterDistribute("defaultKey", rc)
		urlQueue = NewURLQueueDistributed("urlQueue", rc)
	} else {
		fmt.Println("Local")
		dupFilter = NewDuplicateURLFilterLocal()
		urlQueue = NewURLQueueLocal()
	}
	return &Engine{cfg, auth, urlQueue, dupFilter, rc}
}
