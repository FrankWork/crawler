package main

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"flag"
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
)

type RequestWrapper struct {
	URL   string
	Depth int
}

func (r *RequestWrapper) print() {
	fmt.Println(r.URL, r.Depth)
}

func NewRequestWrapper(rawurl string, depth int) *RequestWrapper {
	return &RequestWrapper{rawurl, depth}
}

func serialize(m *RequestWrapper) string {
	var b = new(bytes.Buffer)
	e := gob.NewEncoder(b)
	// Encoding the map
	err := e.Encode(m)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b.Bytes())
}

func deserialize(hexStr string) *RequestWrapper {
	byteArr, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(err)
	}
	var decodedMap *RequestWrapper
	d := gob.NewDecoder(bytes.NewReader(byteArr))

	// Decoding the serialized data
	err = d.Decode(&decodedMap)
	if err != nil {
		panic(err)
	}
	return decodedMap
}

func redisConnect() redis.Conn {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal("Connect to redis error", err)
	}
	return conn
}

func redisSET(key string, value string) {
	conn := redisConnect()
	defer conn.Close()

	ok, err := conn.Do("set", key, value)
	if err != nil {
		log.Println("redis set failed: ", err)
	} else {
		log.Printf("redis set key: %v value: %v \n", key, value)
		log.Println(ok)
	}
}

func redisGET(key string) string {
	conn := redisConnect()
	defer conn.Close()

	value, err := redis.String(conn.Do("Get", key))
	if err != nil {
		log.Println("redis get failed: ", err)
	} else {
		log.Printf("redis get key: %v value: %v \n", key, value)
	}
	return value
}

func main2() {
	reset := flag.String("encode", "true", "true to encode")
	flag.Parse()
	if *reset == "true" {
		rw := NewRequestWrapper("http://www.163.com", 101)
		rw.print()

		hexStr := serialize(rw)
		redisSET("163", hexStr)
		return
	}

	hexStr := redisGET("163")
	d := deserialize(hexStr)
	d.print()
}
