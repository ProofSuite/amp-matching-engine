package testutils

import (
	"log"
	"math/big"
	"math/rand"
	"time"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/math"
	"github.com/ethereum/go-ethereum/common"
)

// Orderfactory simplifies creating orders, trades and cancelOrders objects
// Pair is the token pair for which the order is created
// Exchange is the Ethereum address of the exchange smart contract
// CurrentOrderID increments for each new order
type OrderFactory struct {
	Wallet         *types.Wallet
	Pair           *types.Pair
	Params         *OrderParams
	TradeNonce     uint64
	OrderNonce     uint64
	NonceGenerator *rand.Rand
	// Client         *ethclient.Client
}

// OrderParams groups FeeMake, FeeTake, Nonce, Exipres
// FeeMake and FeeTake are the default fees imposed on makers and takers
// Nonce is the ethereum account nonce that tracks the numbers of transactions
// for the order factory account
// Expires adds a timeout after which an order can no longer be matched
type OrderParams struct {
	ExchangeAddress common.Address
	MakeFee         *big.Int
	TakeFee         *big.Int
	Nonce           *big.Int
	Expires         *big.Int
}

// NewOrderFactory returns an order factory from a given token pair and a given wallet
// TODO: Refactor this function to send back an error
func NewOrderFactory(p *types.Pair, w *types.Wallet, exchangeAddress common.Address) (*OrderFactory, error) {
	// rpcClient, err := rpc.DialWebsocket(context.Background(), "ws://127.0.0.1:8546", "")
	// if err != nil {
	// 	log.Printf("Could not create order factory")
	// 	return nil, err
	// }

	// client := ethclient.NewClient(rpcClient)

	params := &OrderParams{
		MakeFee:         big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Nonce:           big.NewInt(0),
		Expires:         big.NewInt(1e18),
		ExchangeAddress: exchangeAddress,
	}

	source := rand.NewSource(time.Now().UnixNano())
	ng := rand.New(source)

	return &OrderFactory{
		Pair:           p,
		Wallet:         w,
		Params:         params,
		NonceGenerator: ng,
		// Client:         client,
	}, nil
}

// GetWallet returns the order factory wallet
func (f *OrderFactory) GetWallet() *types.Wallet {
	return f.Wallet
}

// GetAddress returns the order factory address
func (f *OrderFactory) GetAddress() common.Address {
	return f.Wallet.Address
}

func (f *OrderFactory) GetExchangeAddress() common.Address {
	return f.Params.ExchangeAddress
}

// SetExchangeAddress changes the default exchange address for orders created by this factory
func (f *OrderFactory) SetExchangeAddress(addr common.Address) error {
	f.Params.ExchangeAddress = addr
	return nil
}

// NewOrderMessage creates an order with the given params and returns a new PLACE_ORDER message
func (f *OrderFactory) NewOrderMessage(buyToken common.Address, buyAmount int64, sellToken common.Address, sellAmount int64) (*types.WebsocketMessage, *types.Order, error) {
	o, err := f.NewOrder(buyToken, buyAmount, sellToken, sellAmount)
	if err != nil {
		return nil, nil, err
	}

	m := types.NewOrderWebsocketMessage(o)

	return m, o, nil
}

func (f *OrderFactory) NewCancelOrderMessage(o *types.Order) (*types.WebsocketMessage, *types.OrderCancel, error) {
	oc, err := f.NewCancelOrder(o)
	if err != nil {
		log.Print(err)
		return nil, nil, err
	}

	m := types.NewOrderCancelWebsocketMessage(oc)
	return m, oc, nil
}

// NewOrder returns a new order with the given params. The order is signed by the factory wallet.
// Currently the nonce is chosen randomly which will be changed in the future
func (f *OrderFactory) NewOrder(buyToken common.Address, buyAmount int64, sellToken common.Address, sellAmount int64) (*types.Order, error) {
	o := &types.Order{}

	o.UserAddress = f.Wallet.Address
	o.ExchangeAddress = f.Params.ExchangeAddress
	o.BuyToken = buyToken
	o.SellToken = sellToken
	o.BuyAmount = big.NewInt(buyAmount)
	o.SellAmount = big.NewInt(sellAmount)
	o.Status = "OPEN"
	o.Expires = f.Params.Expires
	o.MakeFee = f.Params.MakeFee
	o.TakeFee = f.Params.TakeFee
	o.Nonce = big.NewInt(int64(f.NonceGenerator.Intn(1e18)))
	o.Sign(f.Wallet)

	return o, nil
}

func (f *OrderFactory) NewLargeOrder(buyToken common.Address, buyAmount *big.Int, sellToken common.Address, sellAmount *big.Int) (*types.Order, error) {
	o := &types.Order{}

	o.UserAddress = f.Wallet.Address
	o.ExchangeAddress = f.Params.ExchangeAddress
	o.BuyToken = buyToken
	o.SellToken = sellToken
	o.BuyAmount = buyAmount
	o.SellAmount = sellAmount
	o.Status = "OPEN"
	o.Expires = f.Params.Expires
	o.MakeFee = f.Params.MakeFee
	o.TakeFee = f.Params.TakeFee
	o.Nonce = big.NewInt(int64(f.NonceGenerator.Intn(1e18)))
	o.Sign(f.Wallet)

	return o, nil
}

func (f *OrderFactory) NewBuyOrderMessage(price int64, amount float64) (*types.WebsocketMessage, *types.Order, error) {
	o, err := f.NewBuyOrder(price, amount)
	if err != nil {
		return nil, nil, err
	}

	m := types.NewOrderWebsocketMessage(&o)

	return m, &o, nil
}

func (f *OrderFactory) NewSellOrderMessage(price int64, amount float64) (*types.WebsocketMessage, *types.Order, error) {
	o, err := f.NewSellOrder(price, amount)
	if err != nil {
		return nil, nil, err
	}

	m := types.NewOrderWebsocketMessage(&o)

	return m, &o, nil
}

func (f *OrderFactory) NewCancelOrder(o *types.Order) (*types.OrderCancel, error) {
	oc := &types.OrderCancel{}

	oc.OrderHash = o.Hash
	oc.Sign(f.Wallet)
	return oc, nil
}

// NewBuyOrder creates a new buy order from the order factory
func (f *OrderFactory) NewBuyOrder(pricepoint int64, value float64, filled ...float64) (types.Order, error) {
	o := types.Order{}

	// Transform decimal value into rounded points value (ex: 0.01 ETH => 1 ETH)
	amountPoints := big.NewInt(int64(value * 100))
	etherPoints := big.NewInt(1e18)

	o.Amount = math.Div(math.Mul(etherPoints, amountPoints), big.NewInt(100))
	o.BuyAmount = o.Amount
	o.SellAmount = math.Mul(o.Amount, big.NewInt(pricepoint))

	o.UserAddress = f.Wallet.Address
	o.ExchangeAddress = f.Params.ExchangeAddress
	o.BuyToken = f.Pair.BaseTokenAddress
	o.SellToken = f.Pair.QuoteTokenAddress
	o.BaseToken = f.Pair.BaseTokenAddress
	o.QuoteToken = f.Pair.QuoteTokenAddress
	o.Expires = f.Params.Expires
	o.MakeFee = f.Params.MakeFee
	o.TakeFee = f.Params.TakeFee
	o.Nonce = big.NewInt(int64(f.NonceGenerator.Intn(1e8)))
	o.Side = "BUY"

	if filled == nil {
		o.FilledAmount = big.NewInt(0)
		o.Status = "OPEN"
	} else if value == filled[0] {
		o.FilledAmount = o.Amount
		o.Status = "FILLED"
	} else {
		filledPoints := big.NewInt(int64(filled[0] * 100))
		o.FilledAmount = math.Div(math.Mul(etherPoints, filledPoints), big.NewInt(100))
		o.Status = "PARTIAL_FILLED"
	}

	o.PairName = f.Pair.Name()
	o.PricePoint = big.NewInt(pricepoint)
	o.FilledAmount = big.NewInt(0)
	o.CreatedAt = time.Now()

	o.Sign(f.Wallet)
	return o, nil
}

// NewBuyOrder returns a new order with the given params. The order is signed by the factory wallet
// NewBuyOrder computes the AmountBuy and AmountSell parameters from the given amount and price.
// Currently, the amount, price and order type are also kept. This could be amended in the future
// (meaning we would let the engine compute OrderBuy, Amount and Price. Ultimately this does not really
// matter except maybe for convenience/readability purposes)
func (f *OrderFactory) NewSellOrder(pricepoint int64, value float64, filled ...float64) (types.Order, error) {
	o := types.Order{}

	amountPoints := big.NewInt(int64(value * 100))
	etherPoints := big.NewInt(1e18)

	o.Amount = math.Div(math.Mul(etherPoints, amountPoints), big.NewInt(100))
	o.SellAmount = o.Amount
	o.BuyAmount = math.Mul(o.Amount, big.NewInt(pricepoint))

	o.UserAddress = f.Wallet.Address
	o.ExchangeAddress = f.Params.ExchangeAddress
	o.BuyToken = f.Pair.QuoteTokenAddress
	o.SellToken = f.Pair.BaseTokenAddress
	o.BaseToken = f.Pair.BaseTokenAddress
	o.QuoteToken = f.Pair.QuoteTokenAddress
	o.Expires = f.Params.Expires
	o.MakeFee = f.Params.MakeFee
	o.TakeFee = f.Params.TakeFee
	o.Nonce = big.NewInt(int64(f.NonceGenerator.Intn(1e8)))
	o.Side = "SELL"

	o.PricePoint = big.NewInt(pricepoint)
	o.CreatedAt = time.Now()
	o.PairName = f.Pair.Name()

	if filled == nil {
		o.FilledAmount = big.NewInt(0)
		o.Status = "OPEN"
	} else if value == filled[0] {
		o.FilledAmount = o.Amount
		o.Status = "FILLED"
	} else {
		filledPoints := big.NewInt(int64(filled[0] * 100))
		o.FilledAmount = &big.Int{}
		o.FilledAmount.Mul(etherPoints, filledPoints)
		o.Status = "PARTIAL_FILLED"
	}

	o.Sign(f.Wallet)
	return o, nil
}

// NewTrade returns a new trade with the given params. The trade is signed by the factory wallet.
// Currently the nonce is chosen randomly which will be changed in the future
func (f *OrderFactory) NewTrade(o *types.Order, amount int64) (types.Trade, error) {
	t := types.Trade{}

	t.Maker = o.UserAddress
	t.Taker = f.Wallet.Address
	t.BaseToken = o.BaseToken
	t.QuoteToken = o.QuoteToken
	t.TradeNonce = big.NewInt(int64(f.NonceGenerator.Intn(1e8)))
	t.OrderHash = o.Hash
	t.Amount = big.NewInt(amount)

	t.Sign(f.Wallet)
	return t, nil
}

func (f *OrderFactory) NewPendingTrade(o *types.Order, amount int64) *types.PendingTradeMessage {
	t := &types.Trade{}

	t.Maker = o.UserAddress
	t.Taker = f.Wallet.Address
	t.BaseToken = o.BaseToken
	t.QuoteToken = o.QuoteToken
	t.TradeNonce = big.NewInt(int64(f.NonceGenerator.Intn(1e8)))
	t.OrderHash = o.Hash
	t.Amount = big.NewInt(amount)

	t.Sign(f.Wallet)

	return &types.PendingTradeMessage{o, t}
}

// NewOrderCancel creates a new OrderCancel object from a given order
// func (f *OrderFactory) NewOrderCancel(o *types.Order) (*OrderCancel, error) {
// 	oc := &OrderCancel{}

// 	oc.OrderId = o.Id
// 	oc.PairID = f.Pair.ID
// 	oc.OrderHash = o.Hash
// 	oc.Sign(f.Wallet)
// 	return oc, nil
// }

// // NewCancelOrderMessage creates a new OrderCancelMessage from a given order
// func (f *OrderFactory) NewCancelOrderMessage(o *types.Order) (*types.WebSocketMessage, error) {
// 	oc, err := f.NewOrderCancel(o)
// 	if err != nil {
// 		return nil, err
// 	}

// 	p := &OrderCancelPayload{OrderCancel: oc}
// 	return &Message{MessageType: CANCEL_ORDER, Payload: p}, nil
// }

// // NewBuyOrder returns a new order with the given params. The order is signed by the factory wallet
// // NewBuyOrder computes the AmountBuy and AmountSell parameters from the given amount and price.
// // Currently, the amount, price and order type are also kept. This could be amended in the future
// // (meaning we would let the engine compute OrderBuy, Amount and Price. Ultimately this does not really
// // matter except maybe for convenience/readability purposes)
