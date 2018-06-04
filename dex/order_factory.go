package dex

import (
	"context"
	"log"
	"math/big"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// Orderfactory simplifies creating orders, trades and cancelOrders objects
// Pair is the token pair for which the order is created
// Exchange is the Ethereum address of the exchange smart contract
// CurrentOrderID increments for each new order
type OrderFactory struct {
	Client         *ethclient.Client
	Pair           *TokenPair
	Wallet         *Wallet
	Exchange       common.Address
	Params         *OrderParams
	TradeNonce     uint64
	OrderNonce     uint64
	CurrentOrderID uint64
	NonceGenerator *rand.Rand
}

// OrderParams groups FeeMake, FeeTake, Nonce, Exipres
// FeeMake and FeeTake are the default fees imposed on makers and takers
// Nonce is the ethereum account nonce that tracks the numbers of transactions
// for the order factory account
// Expires adds a timeout after which an order can no longer be matched
type OrderParams struct {
	FeeMake *big.Int
	FeeTake *big.Int
	Nonce   *big.Int
	Expires *big.Int
}

// NewOrderFactory returns an order factory from a given token pair and a given wallet
// TODO: Refactor this function to send back an error
func NewOrderFactory(p *TokenPair, w *Wallet) *OrderFactory {

	rpcClient, err := rpc.DialWebsocket(context.Background(), "ws://127.0.0.1:8546", "")
	if err != nil {
		log.Printf("Could not create order factory")
		return nil
	}

	client := ethclient.NewClient(rpcClient)

	params := &OrderParams{
		FeeMake: big.NewInt(0),
		FeeTake: big.NewInt(0),
		Nonce:   big.NewInt(0),
		Expires: big.NewInt(1e18),
	}

	source := rand.NewSource(time.Now().UnixNano())
	ng := rand.New(source)

	return &OrderFactory{
		Pair:           p,
		Wallet:         w,
		Exchange:       config.Exchange,
		Params:         params,
		CurrentOrderID: 0,
		Client:         client,
		NonceGenerator: ng,
	}
}

// SetExchangeAddress changes the default exchange address for orders created by this factory
func (f *OrderFactory) SetExchangeAddress(exchange common.Address) error {
	f.Exchange = exchange
	return nil
}

// NewOrderMessage creates an order with the given params and returns a new PLACE_ORDER message
func (f *OrderFactory) NewOrderMessage(tokenBuy Token, amountBuy int64, tokenSell Token, amountSell int64) (*Message, *Order, error) {
	o, err := f.NewOrder(tokenBuy, amountBuy, tokenSell, amountBuy)
	if err != nil {
		return nil, nil, err
	}

	p := &OrderPayload{Order: o}
	return &Message{MessageType: PLACE_ORDER, Payload: p}, o, nil
}

// NewOrder returns a new order with the given params. The order is signed by the factory wallet.
// Currently the nonce is chosen randomly which will be changed in the future
func (f *OrderFactory) NewOrder(tokenBuy Token, amountBuy int64, tokenSell Token, amountSell int64) (*Order, error) {
	o := &Order{}

	o.Id = f.CurrentOrderID
	o.ExchangeAddress = f.Exchange
	o.TokenBuy = tokenBuy.Address
	o.SymbolBuy = tokenBuy.Symbol
	o.TokenSell = tokenSell.Address
	o.SymbolSell = tokenSell.Symbol
	o.AmountBuy = big.NewInt(amountBuy)
	o.AmountSell = big.NewInt(amountSell)
	o.Expires = f.Params.Expires
	o.FeeMake = f.Params.FeeMake
	o.FeeTake = f.Params.FeeTake
	o.Nonce = big.NewInt(int64(f.NonceGenerator.Intn(1e8)))
	o.Maker = f.Wallet.Address
	o.Price = 0
	o.Amount = 0
	o.PairID = f.Pair.ID
	o.Sign(f.Wallet)

	f.OrderNonce++
	f.CurrentOrderID++
	return o, nil
}

func (f *OrderFactory) NewOrderWithEvents(tokenBuy Token, amountBuy int64, tokenSell Token, amountSell int64) (*Order, error) {
	o, err := f.NewOrder(tokenBuy, amountBuy, tokenSell, amountSell)
	if err != nil {
		return nil, err
	}

	o.events = make(chan *Event)
	return o, nil
}

// NewBuyOrder returns a new order with the given params. The order is signed by the factory wallet
// NewBuyOrder computes the AmountBuy and AmountSell parameters from the given amount and price.
// Currently, the amount, price and order type are also kept. This could be amended in the future
// (meaning we would let the engine compute OrderBuy, Amount and Price. Ultimately this does not really
// matter except maybe for convenience/readability purposes)
func (f *OrderFactory) NewBuyOrder(price uint64, amount uint64) (*Order, error) {
	o := &Order{}

	o.Id = f.CurrentOrderID
	o.ExchangeAddress = f.Exchange
	o.TokenBuy = f.Pair.BaseToken.Address
	o.TokenSell = f.Pair.QuoteToken.Address
	o.SymbolBuy = f.Pair.BaseToken.Symbol
	o.SymbolSell = f.Pair.QuoteToken.Symbol
	o.Expires = f.Params.Expires
	o.FeeMake = f.Params.FeeMake
	o.FeeTake = f.Params.FeeTake
	o.Maker = f.Wallet.Address
	o.Price = price
	o.Amount = amount
	o.Nonce = big.NewInt(int64(f.NonceGenerator.Intn(1e8)))
	o.PairID = f.Pair.ID

	o.AmountBuy = big.NewInt(int64(o.Amount))

	amountSell := &big.Int{}
	amountSell.Mul(big.NewInt(int64(o.Amount)), big.NewInt(int64(o.Price)))
	o.AmountSell = amountSell
	o.OrderType = BUY

	o.Sign(f.Wallet)

	f.OrderNonce++
	f.CurrentOrderID++
	return o, nil
}

// NewSellOrderWithEvents adds an event channel to the trade returned by NewSellOrder
func (f *OrderFactory) NewBuyOrderWithEvent(price uint64, amount uint64) (*Order, error) {
	o, err := f.NewBuyOrder(price, amount)
	if err != nil {
		return nil, err
	}

	o.events = make(chan *Event)
	return o, nil
}

// NewBuyOrder returns a new order with the given params. The order is signed by the factory wallet
// NewBuyOrder computes the AmountBuy and AmountSell parameters from the given amount and price.
// Currently, the amount, price and order type are also kept. This could be amended in the future
// (meaning we would let the engine compute OrderBuy, Amount and Price. Ultimately this does not really
// matter except maybe for convenience/readability purposes)
func (f *OrderFactory) NewSellOrder(price uint64, amount uint64) (*Order, error) {
	o := &Order{}

	o.Id = f.CurrentOrderID
	o.ExchangeAddress = f.Exchange
	o.TokenBuy = f.Pair.QuoteToken.Address
	o.TokenSell = f.Pair.BaseToken.Address
	o.SymbolBuy = f.Pair.QuoteToken.Symbol
	o.SymbolSell = f.Pair.BaseToken.Symbol
	o.Expires = f.Params.Expires
	o.FeeMake = f.Params.FeeMake
	o.FeeTake = f.Params.FeeTake
	o.Maker = f.Wallet.Address
	o.Price = price
	o.Amount = amount
	o.Nonce = big.NewInt(int64(f.NonceGenerator.Intn(1e8)))
	o.PairID = f.Pair.ID

	o.AmountSell = big.NewInt(int64(o.Amount))

	amountBuy := &big.Int{}
	amountBuy.Mul(big.NewInt(int64(o.Amount)), big.NewInt(int64(o.Price)))

	o.AmountBuy = amountBuy
	o.OrderType = SELL

	o.Sign(f.Wallet)
	f.OrderNonce++
	f.CurrentOrderID++
	return o, nil
}

// NewSellOrderWithEvents adds an event channel to the trade returned by NewSellOrder
func (f *OrderFactory) NewSellOrderWithEvent(price uint64, amount uint64) (*Order, error) {
	o, err := f.NewSellOrder(price, amount)
	if err != nil {
		return nil, err
	}

	o.events = make(chan *Event)
	return o, nil
}

// NewTrade returns a new trade with the given params. The trade is signed by the factory wallet.
// Currently the nonce is chosen randomly which will be changed in the future
func (f *OrderFactory) NewTrade(o *Order, amount int64) (*Trade, error) {
	t := &Trade{}

	t.OrderHash = o.Hash
	t.PairID = f.Pair.ID
	t.Taker = f.Wallet.Address
	t.TradeNonce = big.NewInt(int64(f.NonceGenerator.Intn(1e8)))
	t.Amount = big.NewInt(amount)
	t.Sign(f.Wallet)

	f.TradeNonce++
	return t, nil
}

// NewTradeWithEvents adds an event channel to the trade returned by NewTrade
func (f *OrderFactory) NewTradeWithEvents(o *Order, amount int64) (*Trade, error) {
	t, err := f.NewTrade(o, amount)
	if err != nil {
		return nil, err
	}

	t.events = make(chan *Event)
	return t, nil
}

// NewOrderCancel creates a new OrderCancel object from a given order
func (f *OrderFactory) NewOrderCancel(o *Order) (*OrderCancel, error) {
	oc := &OrderCancel{}

	oc.OrderId = o.Id
	oc.PairID = f.Pair.ID
	oc.OrderHash = o.Hash
	oc.Sign(f.Wallet)
	return oc, nil
}

// NewCancelOrderMessage creates a new OrderCancelMessage from a given order
func (f *OrderFactory) NewCancelOrderMessage(o *Order) (*Message, error) {
	oc, err := f.NewOrderCancel(o)
	if err != nil {
		return nil, err
	}

	p := &OrderCancelPayload{OrderCancel: oc}
	return &Message{MessageType: CANCEL_ORDER, Payload: p}, nil
}
