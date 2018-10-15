package engine

import (
	"os"
)

var redisServer int

func init() {
	if os.Args[1] == "live" {
		redisServer = 1
	}
}
