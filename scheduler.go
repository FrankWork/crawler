package main

import (
	"bytes"
	"container/list"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"sync"
)

// URLWrapper: store the parsed url and it's depth
type URLWrapper struct {
	RawURL string
	Depth  int
}

func NewURLWrapper(rawurl string, depth int) *URLWrapper {
	return &URLWrapper{rawurl, depth}
}

func (r *URLWrapper) print() {
	fmt.Println(r.RawURL, r.Depth)
}

// Serialize URLWrapper to store into a database or transport across internet
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

type Queue interface {
	enqueue(uw *URLWrapper)
	dequeue() *URLWrapper
	isEmpty() bool
}

// URLQueueLocal : Local URL Messaging Queue
// list.List is not concurrent safe
type URLQueueLocal struct {
	List *list.List
	lock *sync.Mutex
}

func (q *URLQueueLocal) enqueue(r *URLWrapper) {
	defer q.lock.Unlock()
	q.lock.Lock()
	q.List.PushBack(r)
}
func (q *URLQueueLocal) dequeue() *URLWrapper {
	defer q.lock.Unlock()
	q.lock.Lock()
	if q.List.Len() > 0 {
		r := q.List.Front()
		q.List.Remove(r)
		return r.Value.(*URLWrapper)
	}
	return nil
}
func (q *URLQueueLocal) isEmpty() bool {
	defer q.lock.Unlock()
	q.lock.Lock()
	return q.List.Len() == 0
}

type URLQueueDistributed struct {
	name string
}

// URL Messaging Queue across internet
func (q *URLQueueDistributed) enqueue(uw *URLWrapper) {
	conn, resource := redisPoolConnect()
	defer RedisResourcePool.Put(resource)

	uwStr := serialize(uw)
	redisLPUSH(conn, q.name, uwStr)
}
func (q *URLQueueDistributed) dequeue() *URLWrapper {
	conn, resource := redisPoolConnect()
	defer RedisResourcePool.Put(resource)

	uwStr := redisRPOP(conn, q.name)
	if uwStr != "" {
		return deserialize(uwStr)
	}
	return nil
}
func (q *URLQueueDistributed) isEmpty() bool {
	conn, resource := redisPoolConnect()
	defer RedisResourcePool.Put(resource)

	n := redisLLen(conn, q.name)
	return n == 0
}
