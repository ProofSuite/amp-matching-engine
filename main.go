package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	addr = flag.String("addr", ":8080", "http service address")
)

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
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer connection.Close()

	client := &Client{server: server, connection: connection}
	client.server.register <- client

	for {
		messageType, message, err := client.connection.ReadMessage()
		if err != nil {
			log.Println("Read", err)
			break
		}
		log.Printf("Received: %s", message)
		err = client.connection.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Write:", err)
			break
		}
	}
}
