package dex

import (
	"encoding/json"
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

// NewServer returns a a new empty Server instance.
// Note: Currently note sure whether actionLogs and txLogs are still used
func NewServer() *Server {
	e := NewTradingEngine()

	return &Server{
		clients:    make(map[*Socket]bool),
		engine:     e,
		operator:   &Operator{},
		actionLogs: make(chan *Action),
		txLogs:     make(chan *types.Transaction),
	}
}

func (s *Server) SetupTradingEngine(config *OperatorConfig, quotes Tokens, pairs TokenPairs, done chan bool) error {
	s.engine = NewTradingEngine()

	for _, val := range quotes {
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

	http.HandleFunc("/pair/new", s.RegisterNewPair)
	http.HandleFunc("/quote/new", s.RegisterNewQuoteToken)

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

// RegisterNewPair registers a new pair on the trading engine
// In the final version, the token pair will be written to the database
func (s *Server) RegisterNewPair(w http.ResponseWriter, r *http.Request) {
	p := &TokenPair{}
	d := json.NewDecoder(r.Body)

	err := d.Decode(p)
	if err != nil {
		log.Printf("Error decoding pair struct: %v", err)
		http.Error(w, "Error decoding pair struct", http.StatusInternalServerError)
	}

	if err := s.engine.RegisterNewPair(*p, nil); err != nil {
		log.Printf("Error registering new pair: %v", err)
		http.Error(w, "Error registering pair", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "application/json")
}

// RegisterNewQuoteToken registers a new quote token on the trading engine
// In the final version, the quote token will be written to the database
func (s *Server) RegisterNewQuoteToken(w http.ResponseWriter, r *http.Request) {
	quote := &Token{}
	d := json.NewDecoder(r.Body)

	err := d.Decode(quote)
	if err != nil {
		log.Printf("Error decoding pair struct: %v", err)
		http.Error(w, "Error decoding pair struct", http.StatusInternalServerError)
	}

	if err := s.engine.RegisterNewQuoteToken(*quote); err != nil {
		log.Printf("Error registering new quote token: %v", err)
		http.Error(w, "Error registering quote token", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "application/json")
}
