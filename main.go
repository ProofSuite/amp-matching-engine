package main

// import (
// 	"flag"

// 	"github.com/Proofsuite/amp-matching-engine/dex"
// )

// var addr = flag.String("addr", "localhost:8080", "http service address")
// var done = make(chan bool)
// var config = dex.NewDefaultConfiguration()

// func main() {
// 	flag.Parse()

// 	quoteTokens := config.QuoteTokens
// 	pairs := config.TokenPairs

// 	server := dex.NewServer()
// 	server.Setup(quoteTokens, pairs, done)

// 	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
// 		server.OpenWebsocketConnection(w, r)
// 	})

// 	err := http.ListenAndServe(*addr, nil)
// 	if err != nil {
// 		log.Fatal("ListenAndServe: ", err)
// 	}
// }
