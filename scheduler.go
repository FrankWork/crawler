package main

import (
	"bytes"
	"container/list"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"sync"
)

type URLWrapper struct {
	RawURL string
	Depth  int
}

func NewURLWrapper(rawurl string, depth int) *URLWrapper {
	// request, err := http.NewRequest("GET", rawurl, nil)
	// if err != nil {
	// 	log.Printf("new request failed!, %s\n", rawurl)
	// 	return nil
	// }
	// request.Header.Add("If-None-Match", `W/"wyzzy"`)
	// request.Header.Add("depth", "1")
	return &URLWrapper{rawurl, depth}
}

func (r *URLWrapper) print() {
	fmt.Println(r.RawURL, r.Depth)
}

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

// RequestQueue : Request Messaging Queue
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

func enqueue(uw *URLWrapper) {

}
func dequeue() *URLWrapper {
	return nil
}
func queueIsEmpty() bool {
	return false
}
