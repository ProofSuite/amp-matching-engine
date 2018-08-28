package engine

import (
	"sync"

	"os"
	"strconv"

	"github.com/alicebob/miniredis"
	"github.com/gomodule/redigo/redis"
)

var redisServer int

func init() {
	if os.Args[1] == "live" {
		redisServer = 1
	}
}

func getResource() *Engine {
	if redisServer == 0 {
		c, err := redis.DialURL("redis://localhost:6379")
		if err != nil {
			panic(err)
		}
		// Clear redis before starting tests
		flushData(c)
		return &Engine{c, &sync.Mutex{}}
	}

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	c, err := redis.Dial("tcp", s.Addr())
	if err != nil {
		panic(err)
	}

	return &Engine{c, &sync.Mutex{}}
}

func getSortedSet(c redis.Conn, key string) (map[string]float64, error) {
	resMap := make(map[string]float64)
	res, err := redis.Strings(c.Do("ZRANGE", key, "0", "-1", "WITHSCORES"))
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(res); i = i + 2 {
		resMap[res[i]], _ = strconv.ParseFloat(res[i+1], 64)
	}

	return resMap, nil
}

func getValue(c redis.Conn, key string) (string, error) {
	return redis.String(c.Do("GET", key))
}

func exists(c redis.Conn, key string) bool {
	exists, err := redis.Bool(c.Do("EXISTS", key))
	if err != nil {
		panic(err)
	}
	return exists
}

func flushData(c redis.Conn) {
	c.Do("FLUSHALL")
}
