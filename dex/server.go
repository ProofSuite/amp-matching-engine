package dex

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Server type handles a mapping of socket structs and the trading engine
type Server struct {
	clients    map[*Socket]bool
	engine     *TradingEngine
	operator   *Operator
	actionLogs chan *Action
	txLogs     chan *types.Transaction
}

// NewServer returns a a new empty Server instance. The operator is deployed
// according to the given configuration. There are currently 4 different choices of
// configurations:
// - NewConfiguration() new configuration
func NewServer() *Server {
	return &Server{
		clients: make(map[*Socket]bool),
		engine:  nil,
	}
}

func (s *Server) SetupTradingEngine(config *OperatorConfig, quotes Tokens, pairs TokenPairs, done chan bool) error {
	s.engine = NewTradingEngine()

	for _, val := range quotes {
		log.Printf("val is equal to %x: %v\n", val.Address, val.Symbol)
		if err := s.engine.RegisterNewQuoteToken(val); err != nil {
			return err
		}
	}

	for _, p := range pairs {
		if err := s.engine.RegisterNewPair(p, done); err != nil {
			return err
		}
	}

	if err := s.engine.RegisterOperator(config); err != nil {
		return err
	}

	return nil
}

// Setup registers a list of quote tokens and token pairs
func (s *Server) SetupCurrencies(quotes Tokens, pairs TokenPairs, done chan bool) {
	s.engine = NewTradingEngine()

	for _, val := range quotes {
		s.engine.RegisterNewQuoteToken(val)
	}

	for _, p := range pairs {
		err := s.engine.RegisterNewPair(p, done)
		if err != nil {
			fmt.Printf("\nError registering token pair %v: %v\n", p, err)
		}

	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api" {
		w.WriteHeader(http.StatusNotFound)
	}

	s.OpenWebsocketConnection(w, r)
}

func (s *Server) Start() {
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		s.OpenWebsocketConnection(w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// OpenWebsocketConnection opens a new websocket connection
func (s *Server) OpenWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	out := make(chan *Message)
	in := make(chan *Message)
	events := make(chan *Event)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// defer connection.Close()
	socket := &Socket{server: s, connection: conn, messagesOut: out, messagesIn: in, events: events}

	go socket.handleMessagesOut() //Handle messages from server socket to client
	go socket.handleMessagesIn()  //Handle messages from client to server socket
	socket.listenToMessagesIn()

}
