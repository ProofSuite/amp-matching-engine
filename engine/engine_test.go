package engine

import (
	"sync"

	"os"

	"github.com/Proofsuite/amp-matching-engine/redis"
)

var redisServer int

func init() {
	if os.Args[1] == "live" {
		redisServer = 1
	}
}

func getResource() *Engine {
	if redisServer == 0 {
		c := redis.NewRedisConnection("redis://localhost:6379")
		// Clear redis before starting tests
		c.FlushAll()
		return &Engine{c, &sync.Mutex{}}
	}
	return &Engine{redis.NewMiniRedisConnection(), &sync.Mutex{}}
}