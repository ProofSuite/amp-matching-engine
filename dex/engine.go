package dex

import (
	"errors"
	"fmt"
	"math/big"

	. "github.com/ethereum/go-ethereum/common"
)

// TradingEngine manages all the different orderbooks through a TokenPair->OrderBook
// It also holds the list of current quote tokens and token pairs
// quoteTokens is a mapping of token addresses to Token structs (address + symbol)
// pairs is a mapping of token pair ID to corresponding TokenPair structs
type TradingEngine struct {
	orderbooks  map[TokenPair]*OrderBook
	quoteTokens map[Address]Token
	pairs       map[Hash]TokenPair
	operator    *Operator
}

// NewTradingEngine returns an empty TradingEngine struct
func NewTradingEngine() *TradingEngine {
	e := new(TradingEngine)
	e.orderbooks = make(map[TokenPair]*OrderBook)
	e.quoteTokens = make(map[Address]Token)
	e.pairs = make(map[Hash]TokenPair)
	return e
}

// PrintLogs prints the logs for each token pair registered in the orderbook
func (e *TradingEngine) PrintLogs() {
	for symbol, orderbook := range e.orderbooks {
		fmt.Printf("Logs for %v:\n", symbol)
		fmt.Printf("%v\n", orderbook.GetLogs())
	}
}

// RegisterNewQuoteToken registers a new quote token on the engine
// Only an authorized user should be allowed to register a new token on the engine (TODO)
func (e *TradingEngine) RegisterNewQuoteToken(t Token) error {
	if _, ok := e.quoteTokens[t.Address]; ok {
		return errors.New("Quote token already exists")
	}

	e.quoteTokens[t.Address] = t
	return nil
}

// RegisterOperator registers a new operator on the engine. The operator is the account
// that sends the blockchain transaction associated to an order/trade to the exchange
// smart contract
func (e *TradingEngine) RegisterOperator(config *OperatorConfig) error {
	op, err := NewOperator(config)
	if err != nil {
		return err
	}

	e.operator = op
	return nil
}

// RegisterNewPair registers a new token pair on the engine
// To be valid, the quote token of the new token pair needs to be registered in the quote tokens of the engine
func (e *TradingEngine) RegisterNewPair(p TokenPair, done chan<- bool) error {
	if _, ok := e.quoteTokens[p.QuoteToken.Address]; !ok {
		return errors.New("Quote token is not registered")
	}

	id := p.ComputeID()
	if _, ok := e.pairs[id]; ok {
		return errors.New("Token Pair is already registered")
	}

	p.ID = id
	actions := make(chan *Action)
	logger := make([]*Action, 0)

	ob := new(OrderBook)
	ob.bid = 0
	ob.ask = MAX_PRICE
	ob.actions = actions
	ob.logger = logger
	ob.orderIndex = make(map[Hash]*Order)

	go func() {
		for {
			action := <-actions
			ob.logger = append(ob.logger, action)
			if action.actionType == AT_DONE {
				done <- true
				return
			}
		}
	}()

	for i := range ob.prices {
		ob.prices[i] = new(PricePoint)
	}

	e.pairs[id] = p
	e.orderbooks[p] = ob

	// fmt.Printf("\nRegistered new pair %v\n\n", p.String())
	return nil
}

// ComputeOrderPrice calculates the (Amount, Price) tuple corresponding
// to the (AmountBuy, TokenBuy, AmountSell, TokenSell) quadruplet
func (e *TradingEngine) ComputeOrderPrice(o *Order) error {
	tokenPair := e.pairs[o.PairID]

	if o.TokenBuy == tokenPair.BaseToken.Address {
		o.OrderType = BUY
		o.Amount = o.AmountBuy.Uint64()
		o.Price = o.AmountSell.Uint64() * (1e3) / o.AmountBuy.Uint64()
	} else if o.TokenBuy == tokenPair.QuoteToken.Address {
		o.OrderType = SELL
		o.Amount = o.AmountSell.Uint64()
		o.Price = o.AmountBuy.Uint64() * (1e3) / o.AmountSell.Uint64()
	} else {
		return errors.New("Token Buy Address should be equal to base token or quote token address")
	}

	return nil
}

// AddOrder computes the order price point
func (e *TradingEngine) AddOrder(o *Order) error {
	tokenPair, ok := e.pairs[o.PairID]
	if !ok {
		return errors.New("Token pair does not exist")
	}

	ob, ok := e.orderbooks[tokenPair]
	if !ok {
		return errors.New("Orderbook does not exist")
	}

	err := e.ComputeOrderPrice(o)
	if err != nil {
		return err
	}

	ok, err = o.VerifySignature()
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("Signature is not valid")
	}

	e.ComputeOrderPrice(o)
	ob.AddOrder(o)
	return nil
}

// CancelOrder cancels an order that was previously sent to the orderbook. To be valid,
// the orderCancel struct, need to correspond an existing token pair and order. The
// orderCancel's signature needs to be a valid signature from the order Maker.
func (e *TradingEngine) CancelOrder(oc *OrderCancel) error {
	tokenPair, ok := e.pairs[oc.PairID]
	if !ok {
		return errors.New("Token Pair does not exist")
	}

	ob, ok := e.orderbooks[tokenPair]
	if !ok {
		return errors.New("Orderbook does not exist")
	}

	order, ok := ob.orderIndex[oc.OrderHash]
	if !ok {
		return errors.New("Order does not exist")
	}

	ok, err := oc.VerifySignature(order)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("Signature is not valid")
	}

	ob.CancelOrder(oc.OrderHash)
	return nil
}

// Execute Order adds the order to the blockchain transaction execution queue
func (e *TradingEngine) ExecuteOrder(t *Trade) error {
	fmt.Printf("Now in the trading engine trying to trade: %v", t)

	tokenPair, ok := e.pairs[t.PairID]
	if !ok {
		return errors.New("Token Pair does not exist")
	}

	ob, ok := e.orderbooks[tokenPair]
	if !ok {
		return errors.New("Orderbook does not exist")
	}

	o, ok := ob.orderIndex[t.OrderHash]
	if !ok {
		return errors.New("Order does not exist")
	}

	err := e.operator.AddTradeToExecutionList(o, t)
	if err != nil {
		return err
	}

	return nil
}

func (e *TradingEngine) CancelTrade(t *Trade) error {
	tokenPair, ok := e.pairs[t.PairID]
	if !ok {
		return errors.New("Token Pair does not exist")
	}

	ob, ok := e.orderbooks[tokenPair]
	if !ok {
		return errors.New("Orderbook does not exist")
	}

	ob.CancelTrade(t)
	return nil
}

func (e *TradingEngine) TokenBalance(owner Address, token Address) (*big.Int, error) {
	ex := e.operator.Exchange
	balance, err := ex.TokenBalance(owner, token)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

// CloseOrderBook closes the orderbook associated to a pair ID
func (e *TradingEngine) CloseOrderBook(pairID Hash) (bool, error) {
	tokenPair := e.pairs[pairID]
	if ob, ok := e.orderbooks[tokenPair]; !ok {
		return false, errors.New("Orderbook does not exist")
	} else {
		ob.Done()
		return true, nil
	}
}
