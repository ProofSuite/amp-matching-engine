package redis

import (
	"github.com/alicebob/miniredis"
	"github.com/gomodule/redigo/redis"
)

type RedisConnection struct {
	redis.Conn
}

// InitConnection returns a new connection to redis
func InitConnection(uri string) redis.Conn {
	c, err := redis.DialURL(uri)
	if err != nil {
		panic(err)
	}
	return c
}

func NewRedisConnection(uri string) *RedisConnection {
	c, err := redis.DialURL(uri)
	if err != nil {
		panic(err)
	}

	return &RedisConnection{c}
}

func NewMiniRedisConnection() *RedisConnection {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	c, err := redis.Dial("tcp", s.Addr())
	if err != nil {
		panic(err)
	}

	return &RedisConnection{c}
}

func (c *RedisConnection) FlushAll() {
	c.Do("FLUSHALL")
}

func (c *RedisConnection) GetValue(key string) (string, error) {
	value, err := redis.String(c.Do("GET", key))
	if err != nil {
		return "", err
	}

	return value, nil
}

func (c *RedisConnection) Exists(key string) (bool, error) {
	exists, err := redis.Bool(c.Do("EXISTS", key))
	if err != nil {
		return false, err
	}

	return exists, nil
}
