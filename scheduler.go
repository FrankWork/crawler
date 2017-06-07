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

// URLQueue : Local URL Messaging Queue
// list.List is not concurrent safe
type URLQueue struct {
	List *list.List
}

func (q *URLQueue) enqueue(r *URLWrapper) {
	defer lock.Unlock()
	lock.Lock()
	q.List.PushBack(r)
}
func (q *URLQueue) dequeue() *URLWrapper {
	defer lock.Unlock()
	lock.Lock()
	if q.List.Len() > 0 {
		r := q.List.Front()
		q.List.Remove(r)
		return r.Value.(*URLWrapper)
	}
	return nil
}
func (q *URLQueue) isEmpty() bool {
	defer lock.Unlock()
	lock.Lock()
	return q.List.Len() == 0
}

var urlQueue *URLQueue
var lock sync.Mutex

func init() {
	urlQueue = &URLQueue{list.New()}

}

// URL Messaging Queue across internet
func enqueue(uw *URLWrapper) {

}
func dequeue() *URLWrapper {
	return nil
}
func queueIsEmpty() bool {
	return false
}
