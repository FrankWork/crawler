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

func NewURLQueueLocal() *URLQueueLocal {
	return &URLQueueLocal{list.New(), new(sync.Mutex)}
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
	rc   *RedisClient
}

func NewURLQueueDistributed(name string, rc *RedisClient) *URLQueueDistributed {
	return &URLQueueDistributed{name, rc}
}

// URL Messaging Queue across internet
func (q *URLQueueDistributed) enqueue(uw *URLWrapper) {
	uwStr := serialize(uw)
	q.rc.LPush(q.name, uwStr)
}

// FIXME
func (q *URLQueueDistributed) dequeue() *URLWrapper {
	uwStr := q.rc.RPop(q.name)
	if uwStr != "" {
		return deserialize(uwStr)
	}
	return nil
}

// FIXME
func (q *URLQueueDistributed) isEmpty() bool {
	n := q.rc.LLen(q.name)
	return n == 0
}
