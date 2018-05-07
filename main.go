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

func main() {
	flag.Parse()
	server := NewServer()
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
	// var v interface{}
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer connection.Close()

	client := &Client{server: server, connection: connection}
	client.server.register <- client

	for {
		err := client.connection.ReadJSON(&order)
		if err != nil {
			log.Println("Read", err)
			break
		}

		log.Printf("Received: %v", order)
		err = client.connection.WriteJSON(&order)
		if err != nil {
			log.Println("Write:", err)
			break
		}
	}
}
