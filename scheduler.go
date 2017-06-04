package main

import (
	"bytes"
	"container/list"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type RequestWrapper struct {
	request *http.Request
	depth   int
}

func NewRequestWrapper(rawurl string, depth int) *RequestWrapper {
	request, err := http.NewRequest("GET", rawurl, nil)
	if err != nil {
		log.Printf("new request failed!, %s\n", rawurl)
		return nil
	}
	// request.Header.Add("If-None-Match", `W/"wyzzy"`)
	// request.Header.Add("depth", "1")
	return &RequestWrapper{request, depth}
}

func (r *RequestWrapper) print() {
	fmt.Println(r.request.URL.String(), r.depth)
}

func serialize(m map[string]int) string {
	var b = new(bytes.Buffer)
	e := gob.NewEncoder(b)
	// Encoding the map
	err := e.Encode(m)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b.Bytes())
}

func deserialize(hexStr string) map[string]int {
	byteArr, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(err)
	}
	var decodedMap map[string]int
	d := gob.NewDecoder(bytes.NewReader(byteArr))

	// Decoding the serialized data
	err = d.Decode(&decodedMap)
	if err != nil {
		panic(err)
	}
	return decodedMap
}

// RequestQueue : Request Messaging Queue
type RequestQueue struct {
	List *list.List
}

func (q *RequestQueue) enqueue(r *RequestWrapper) {
	defer lock.Unlock()
	lock.Lock()
	q.List.PushBack(r)
}
func (q *RequestQueue) dequeue() *RequestWrapper {
	defer lock.Unlock()
	lock.Lock()
	r := q.List.Front()
	q.List.Remove(r)
	return r.Value.(*RequestWrapper)
}
func (q *RequestQueue) isEmpty() bool {
	defer lock.Unlock()
	lock.Lock()
	return q.List.Len() == 0
}

var requestQueue *RequestQueue
var lock sync.Mutex

func init() {
	requestQueue = &RequestQueue{list.New()}

}
