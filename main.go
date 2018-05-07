package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	addr = flag.String("addr", ":8080", "http service address")
)

type TrimmedOrder struct {
	Id     uint64 `json:"id,int,omitempty"`
	Symbol string `json:"symbol,string,omitempty"`
	Price  uint32 `json:"price,int,omitempty"`
	Amount uint32 `json:"amount,int,omitempty"`
}

var done = make(chan bool)
var BTC_USD = NewPair("BTC", "USDT")

func main() {
	flag.Parse()
	engine := NewTradingEngine()
	engine.CreateNewOrderBook(BTC_USD, done)

	server := NewServer(engine)
	go server.start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		registerClient(server, w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// serveWs handles websocket requests from the peer.
func registerClient(server *Server, w http.ResponseWriter, r *http.Request) {
	var order Order
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer connection.Close()

	client := &Client{server: server, connection: connection}
	client.server.register <- client

	for {

		if err := client.connection.ReadJSON(&order); err != nil {
			log.Println("Read", err)
			break
		}

		if err = server.engine.AddOrder(&order); err != nil {
			log.Println("Failed processing order: %v", &order)
		}

		server.engine.CloseOrderBookChannel(BTC_USD)
		<-done
		server.engine.PrintLogs()

		if err = client.connection.WriteJSON(&order); err != nil {
			log.Println("Write:", err)
			break
		}
	}
}
