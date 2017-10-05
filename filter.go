package main

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"net/url"
	"sync"
)

// fingerPrint encode the src string using sha1 and hex
func fingerPrint(src string) string {
	hash := sha1.New()
	hash.Write([]byte(src))
	return hex.EncodeToString(hash.Sum(nil))
}

// hostAndFingerPrint get hostname and finger print from the raw url
func hostAndFingerPrint(rawURL string) (string, string) {
	urlStrct, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal(err)
	}

	fingerPrint := fingerPrint(rawURL)

	return urlStrct.Host, fingerPrint
}

// DupFilter interface to filter duplicate url
type DupFilter interface {
	IsDuplicate(rawurl string) bool
	AddURL(rawurl string)
	RemoveURL(rawurl string)
}

// DupURLFilter implement DupFilter interface based on map
// map is not concurrent safe, sync.Mutex
type DupURLFilter struct {
	urlSet map[string]int
	lock   *sync.Mutex // make the map concurrent safe
}

// NewDupURLFilter is constructor of DupURLFilter struct
func NewDupURLFilter() *DupURLFilter {
	return &DupURLFilter{make(map[string]int), new(sync.Mutex)}
}

// IsDuplicate interface of DupURLFilter
func (d *DupURLFilter) IsDuplicate(rawurl string) bool {
	urlfp := fingerPrint(rawurl)

	d.lock.Lock()
	defer d.lock.Unlock()
	if d.urlSet[urlfp] == 0 {
		return false
	}
	return true
}

// AddURL interface of DupURLFilter
func (d *DupURLFilter) AddURL(rawurl string) {
	urlfp := fingerPrint(rawurl)
	d.lock.Lock()
	defer d.lock.Unlock()
	d.urlSet[urlfp]++
}

// RemoveURL interface of DupURLFilter
func (d *DupURLFilter) RemoveURL(rawurl string) {
	urlfp := fingerPrint(rawurl)
	d.lock.Lock()
	defer d.lock.Unlock()
	d.urlSet[urlfp] = 0
}

// DupURLFilterRedis implement DupFilter interface based on redis
type DupURLFilterRedis struct {
	name  string // name of the redis set
	redis *RedisClient
}

// NewDupURLFilterRedis is constrctor of DupURLFilterRedis struct
func NewDupURLFilterRedis(name string, redis *RedisClient) *DupURLFilterRedis {
	return &DupURLFilterRedis{name, redis}
}

// IsDuplicate interface of DupURLFilter
func (d *DupURLFilterRedis) IsDuplicate(rawurl string) bool {
	urlfp := fingerPrint(rawurl)
	return d.redis.SIsMember(d.name, urlfp)
}

// AddURL interface of DupURLFilter
func (d *DupURLFilterRedis) AddURL(rawurl string) {
	urlfp := fingerPrint(rawurl)
	d.redis.SAdd(d.name, urlfp)
}

// RemoveURL interface of DupURLFilter
func (d *DupURLFilterRedis) RemoveURL(rawurl string) {
	urlfp := fingerPrint(rawurl)
	d.redis.SRem(d.name, urlfp)
}
