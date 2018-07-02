package redisclient

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// InitConnection returns a new connection to redis
func InitConnection(uri string) redis.Conn {
	c, err := redis.DialURL(uri)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return c
}
