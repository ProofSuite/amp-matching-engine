package dex

import (
	"errors"
	"fmt"

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
}

// NewTradingEngine returns an empty TradingEngine struct
func NewTradingEngine() *TradingEngine {
	engine := new(TradingEngine)
	engine.orderbooks = make(map[TokenPair]*OrderBook)
	engine.quoteTokens = make(map[Address]Token)
	engine.pairs = make(map[Hash]TokenPair)
	return engine
}

// PrintLogs prints the logs for each token pair registered in the orderbook
func (engine *TradingEngine) PrintLogs() {
	for symbol, orderbook := range engine.orderbooks {
		fmt.Printf("Logs for %v:\n", symbol)
		fmt.Printf("%v\n", orderbook.GetLogs())
	}
}

// RegisterNewQuoteToken registers a new quote token on the engine
// Only an authorized user should be allowed to register a new token on the engine (TODO)
func (engine *TradingEngine) RegisterNewQuoteToken(t Token) error {
	if _, ok := engine.quoteTokens[t.Address]; ok {
		return errors.New("Quote token already exists")
	}

	engine.quoteTokens[t.Address] = t
	return nil
}

// RegisterNewPair registers a new token pair on the engine
// To be valid, the quote token of the new token pair needs to be registered in the quote tokens of the engine
func (engine *TradingEngine) RegisterNewPair(p TokenPair, done chan<- bool) error {
	if _, ok := engine.quoteTokens[p.QuoteToken.Address]; !ok {
		return errors.New("Quote token is not registered")
	}

	id := p.ComputeID()
	if _, ok := engine.pairs[id]; ok {
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

	engine.pairs[id] = p
	engine.orderbooks[p] = ob

	// fmt.Printf("\nRegistered new pair %v\n\n", p.String())
	return nil
}

// ComputeOrderPrice calculates the (Amount, Price) tuple corresponding
// to the (AmountBuy, TokenBuy, AmountSell, TokenSell) quadruplet
func (engine *TradingEngine) ComputeOrderPrice(o *Order) error {
	tokenPair := engine.pairs[o.PairID]

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
func (engine *TradingEngine) AddOrder(o *Order) error {

	tokenPair, ok := engine.pairs[o.PairID]
	if !ok {
		return errors.New("Token pair does not exist")
	}

	orderbook, ok := engine.orderbooks[tokenPair]
	if !ok {
		return errors.New("Orderbook does not exist")
	}

	err := engine.ComputeOrderPrice(o)
	if err != nil {
		return err
	}

	if ok, _ := o.VerifySignature(); !ok {
		return errors.New("Signature is not valid")
	}

	engine.ComputeOrderPrice(o)
	orderbook.AddOrder(o)
	return nil
}

// CancelOrder cancels an order that was previously sent to the orderbook. To be valid,
// the orderCancel struct, need to correspond an existing token pair and order. The
// orderCancel's signature needs to be a valid signature from the order Maker.
func (engine *TradingEngine) CancelOrder(oc *OrderCancel) error {
	tokenPair, ok := engine.pairs[oc.PairID]
	if !ok {
		return errors.New("Token Pair does not exist")
	}

	orderbook, ok := engine.orderbooks[tokenPair]
	if !ok {
		return errors.New("Orderbook does not exist")
	}

	order, ok := orderbook.orderIndex[oc.OrderHash]
	if !ok {
		return errors.New("Order does not exist")
	}

	ok, _ = oc.VerifySignature(order)
	if !ok {
		return errors.New("Signature is not valid")
	}

	orderbook.CancelOrder(oc.OrderHash)
	return nil
}

// Execute Order adds the order to the blockchain transaction execution queue
func (engine *TradingEngine) ExecuteOrder(o *Order) error {
	fmt.Printf("Executing order: %v", o)
	return nil
}

// CloseOrderBook closes the orderbook associated to a pair ID
func (engine *TradingEngine) CloseOrderBook(pairID Hash) (bool, error) {
	tokenPair := engine.pairs[pairID]
	if orderbook, ok := engine.orderbooks[tokenPair]; !ok {
		return false, errors.New("Orderbook does not exist")
	} else {
		orderbook.Done()
		return true, nil
	}
}
