package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")
var done = make(chan bool)
var BTC_USD = Pair{BaseToken: "ETH", QuoteToken: "EOS"}
var pairs = []Pair{BTC_USD}

func main() {
	flag.Parse()
	server := NewServer()
	go server.start(pairs, done)

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		server.openWebsocketConnection(w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
