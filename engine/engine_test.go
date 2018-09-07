package engine

import (
	"sync"

	"github.com/Proofsuite/amp-matching-engine/rabbitmq"

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

	amqp := rabbitmq.InitConnection("amqp://guest:guest@localhost:5672/")

	if redisServer == 0 {
		c := redis.NewRedisConnection("redis://localhost:6379")
		// Clear redis before starting tests
		c.FlushAll()
		return &Engine{c, amqp, &sync.Mutex{}}
	}
	return &Engine{redis.NewMiniRedisConnection(), amqp, &sync.Mutex{}}
}
