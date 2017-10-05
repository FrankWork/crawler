package main

import (
	"bytes"
	"container/list"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"sync"
)

// URLWrapper store the parsed url and it's depth
type URLWrapper struct {
	RawURL string
	Depth  int
}

// NewURLWrapper new a URLWrapper object
func NewURLWrapper(rawurl string, depth int) *URLWrapper {
	return &URLWrapper{rawurl, depth}
}

// String return depth and url of URLWrapper as a string
func (r *URLWrapper) String() string {
	return fmt.Sprintf("%d: %s", r.Depth, r.RawURL)
}

// serialize URLWrapper to store into a database or transport across internet
func serialize(uw *URLWrapper) string {
	var b = new(bytes.Buffer)
	e := gob.NewEncoder(b)
	// Encoding the data
	err := e.Encode(uw)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b.Bytes())
}

// deserialize an URLWrapper object from hex string
func deserialize(hexStr string) *URLWrapper {
	byteArr, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(err)
	}
	var uw *URLWrapper
	d := gob.NewDecoder(bytes.NewReader(byteArr))

	// Decoding the serialized data
	err = d.Decode(&uw)
	if err != nil {
		panic(err)
	}
	return uw
}

// Queue interface to store URLWrapper
type Queue interface {
	Enqueue(uw *URLWrapper)
	Dequeue() *URLWrapper
	IsEmpty() bool
}

// URLQueue : Local URL Messaging Queue
// list.List is not concurrent safe, sync.Mutex
type URLQueue struct {
	List *list.List
	lock *sync.Mutex // make the list concurrent safe
}

// NewURLQueue return an URLQueue object
func NewURLQueue() *URLQueue {
	return &URLQueue{list.New(), new(sync.Mutex)}
}

// Enqueue an URLWrapper object
func (q *URLQueue) Enqueue(r *URLWrapper) {
	defer q.lock.Unlock()
	q.lock.Lock()
	q.List.PushBack(r)
}

// Dequeue an URLWrapper object
func (q *URLQueue) Dequeue() *URLWrapper {
	defer q.lock.Unlock()
	q.lock.Lock()
	if q.List.Len() > 0 {
		r := q.List.Front()
		q.List.Remove(r)
		return r.Value.(*URLWrapper)
	}
	return nil
}

// IsEmpty return a bool
func (q *URLQueue) IsEmpty() bool {
	defer q.lock.Unlock()
	q.lock.Lock()
	return q.List.Len() == 0
}

// URLQueueRedis implement an Queue interface based on redis
type URLQueueRedis struct {
	name string // name of the queue, used as a key in redis
	rc   *RedisClient
}

// NewURLQueueRedis new an queue object, name is used as a key in redis
func NewURLQueueRedis(name string, rc *RedisClient) *URLQueueRedis {
	return &URLQueueRedis{name, rc}
}

// Enqueue an URLWrapper
func (q *URLQueueRedis) Enqueue(uw *URLWrapper) {
	uwStr := serialize(uw)
	q.rc.LPush(q.name, uwStr)
}

// Dequeue an URLWrapper
func (q *URLQueueRedis) Dequeue() *URLWrapper {
	uwStr := q.rc.RPop(q.name)
	if uwStr != "" {
		return deserialize(uwStr)
	}
	return nil
}

// IsEmpty returns a bool
func (q *URLQueueRedis) IsEmpty() bool {
	n := q.rc.LLen(q.name)
	return n == 0
}
