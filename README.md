# code structure

- redis.go
- scheduler.go
- filter.go
- query.go
- download.go

## redis.go

- ResourceConn
- RedisClient
- NewRedisClient


## scheduler.go

- URLWrapper
- NewURLWrapper
- Queue
- URLQueue
- NewURLQueue
- URLQueueRedis
- NewURLQueueRedis

## filter.go

- DupFilter
- DupURLFilter
- NewDupURLFilter
- DupURLFilterRedis
- NewDupURLFilterRedis

## query.go

- GetTitle
- GetAllLinks

## download.go

- Request



# benchmark

url :"http://open.163.com"
maxDepth : 1 

## dup links by reids

real	1m11.614s
user	0m3.564s
sys	0m0.352s

mutex

real	2m33.347s
user	0m3.680s
sys	0m0.368s


## dup links by map

real	0m30.619s
user	0m3.256s
sys	0m0.236s

mutex

real	1m7.696s
user	0m3.520s
sys	0m0.384s
