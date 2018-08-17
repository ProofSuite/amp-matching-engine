package engine

import (
	"sync"

	"github.com/alicebob/miniredis"
	"github.com/gomodule/redigo/redis"
)

func init() {

}
func getResource() (*Resource, *miniredis.Miniredis) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	c, err := redis.Dial("tcp", s.Addr())
	if err != nil {
		panic(err)
	}
	return &Resource{c, &sync.Mutex{}}, s
}
