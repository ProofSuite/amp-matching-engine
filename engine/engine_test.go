package engine

import (
	"testing"

	"sync"

	"github.com/alicebob/miniredis"
	"github.com/gomodule/redigo/redis"
)

func init() {
	// mockConn, err := amqptest.Dial("amqp://localhost:5672/%2f") // will fail,
	// if err == nil {
	// 	fmt.Println("This shall fail, because no fake amqp server is running...")
	// }

	// fakeServer := server.NewServer("amqp://localhost:5672/%2f")
	// fakeServer.Start()

	// mockConn, err = amqptest.Dial("amqp://localhost:5672/%2f") // now it works =D

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// //Now you can use mockConn as a real amqp connection.
	// channel, err := mockConn.Channel()
	// fmt.Print(channel)
}
func TestAMQP(t *testing.T) {

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
