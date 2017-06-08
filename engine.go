package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/youtube/vitess/go/pools"
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
	redisPool *pools.ResourcePool
}

func (e *Engine) Init() {
	// parse config and auth toml file
	cwd := getExecutablePath()
	var err error

	cfg := &e.cfg
	auth := &e.auth

	if _, err = toml.DecodeFile(path.Join(cwd, "config.toml"), cfg); err != nil {
		log.Println("toml.DecodeFile failed")
		log.Printf("toml file path: %s\n", path.Join(cwd, "config.toml"))
		panic(err)
	}
	fmt.Println("decode config.toml")

	if _, err = toml.DecodeFile(path.Join(cwd, "auth.toml"), auth); err != nil {
		log.Println("toml.DecodeFile failed")
		log.Printf("toml file path: %s\n", path.Join(cwd, "auth.toml"))
		panic(err)
	}

	if cfg.Distributed {
		fmt.Println("Distribute")
		e.dupFilter = NewDuplicateURLFilterDistribute("defaultKey")
		e.urlQueue = NewURLQueueDistributed("urlQueue")
	} else {
		fmt.Println("Local")
		e.dupFilter = NewDuplicateURLFilterLocal()
		e.urlQueue = NewURLQueueLocal()
	}
	e.redisPool = NewRedisPool(auth.RedisHost, auth.RedisAuth, auth.RedisDb,
		cfg.RedisPoolCapacity, cfg.RedisPoolMaxCapacity,
		cfg.RedisPoolIdleTimeout.Duration)

}
