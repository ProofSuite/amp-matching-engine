package dex

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	. "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

type OrderStatus int

const (
	NEW OrderStatus = iota
	OPEN
	PARTIAL_FILLED
	FILLED
	CANCELLED
	PENDING
)

func (orderStatus *OrderStatus) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	value, ok := map[string]OrderStatus{
		"NEW":            NEW,
		"OPEN":           OPEN,
		"PARTIAL_FILLED": PARTIAL_FILLED,
		"FILLED":         FILLED,
		"CANCELLED":      CANCELLED,
		"PENDING":        PENDING}[s]
	if !ok {
		return errors.New("Invalid Enum Status Value")
	}
	*orderStatus = value
	return nil
}

func (orderStatus *OrderStatus) MarshalJSON() ([]byte, error) {
	value, ok := map[OrderStatus]string{
		NEW:            "NEW",
		OPEN:           "OPEN",
		PARTIAL_FILLED: "PARTIAL_FILLED",
		FILLED:         "FILLED",
		CANCELLED:      "CANCELLED",
		PENDING:        "PENDING"}[*orderStatus]
	if !ok {
		return nil, errors.New("Invalid Enum Type")
	}
	return json.Marshal(value)
}

type OrderType int

const (
	BUY OrderType = iota
	SELL
)

func (orderType *OrderType) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	value, ok := map[string]OrderType{"BUY": BUY, "SELL": SELL}[s]
	if !ok {
		return errors.New("Invalid Enum Type Value")
	}
	*orderType = value
	return nil
}

func (orderType *OrderType) MarshalJSON() ([]byte, error) {
	value, ok := map[OrderType]string{BUY: "BUY", SELL: "SELL"}[*orderType]
	if !ok {
		return nil, errors.New("Invalid Enum Type")
	}
	return json.Marshal(value)
}

type Order struct {
	Id              uint64      `json:"id,omitempty"`
	OrderType       OrderType   `json:"orderType,omitempty"`
	status          OrderStatus `json:"orderStatus,omitempty"`
	ExchangeAddress Address     `json:"exchangeAddress,omitempty"`
	Maker           Address     `json:"maker,omitempty"`
	TokenBuy        Address     `json:"tokenBuy,omitempty"`
	TokenSell       Address     `json:"tokenSell,omitempty"`
	SymbolBuy       string      `json:"symbolBuy,omitempty"`
	SymbolSell      string      `json:"symbolSell,omitempty"`
	AmountBuy       *big.Int    `json:"amountBuy,omitempty"`
	AmountSell      *big.Int    `json:"amountSell,omitempty"`
	Expires         *big.Int    `json:"expires,omitempty"`
	Nonce           *big.Int    `json:"nonce,omitempty"`
	FeeMake         *big.Int    `json:"feeMake,omitempty"`
	FeeTake         *big.Int    `json:"feeTake,omitempty"`
	Signature       *Signature  `json:"signature,omitempty"`
	PairID          Hash        `json:"pairID,omitempty"`
	Hash            Hash        `json:"hash,omitempty"`
	Price           uint64      `json:"price,omitempty"`
	Amount          uint64      `json:"amount,omitempty"`
	next            *Order
	events          chan *Event
}

func (o *Order) MarshalJSON() ([]byte, error) {
	order := map[string]interface{}{
		"id":              o.Id,
		"exchangeAddress": o.ExchangeAddress,
		"maker":           o.Maker,
		"tokenBuy":        o.TokenBuy,
		"tokenSell":       o.TokenSell,
		"symbolBuy":       o.SymbolBuy,
		"symbolSell":      o.SymbolSell,
		"amountBuy":       o.AmountBuy.String(),
		"amountSell":      o.AmountSell.String(),
		"expires":         o.Expires.String(),
		"nonce":           o.Nonce.String(),
		"feeMake":         o.FeeMake.String(),
		"feeTake":         o.FeeTake.String(),
		"signature": map[string]interface{}{
			"V": o.Signature.V,
			"R": o.Signature.R,
			"S": o.Signature.S,
		},
		"pairID": o.PairID,
		"hash":   o.Hash,
		"price":  o.Price,
		"amount": o.Amount,
	}
	return json.Marshal(order)
}

func (o *Order) UnmarshalJSON(b []byte) error {
	order := map[string]interface{}{}

	err := json.Unmarshal(b, &order)
	if err != nil {
		return err
	}

	if order["id"] == nil {
		return errors.New("Order ID not set")
	}

	o.Id, _ = strconv.ParseUint(order["id"].(string), 10, 64)

	if order["price"] != nil {
		o.Price, _ = strconv.ParseUint(order["price"].(string), 10, 64)
	}

	if order["amount"] != nil {
		o.Amount, _ = strconv.ParseUint(order["amount"].(string), 10, 64)
	}

	o.ExchangeAddress = HexToAddress(order["exchangeAddress"].(string))
	o.Maker = HexToAddress(order["maker"].(string))
	o.TokenBuy = HexToAddress(order["tokenBuy"].(string))
	o.TokenSell = HexToAddress(order["tokenSell"].(string))
	o.SymbolBuy = order["symbolBuy"].(string)
	o.SymbolSell = order["symbolSell"].(string)

	o.AmountBuy = new(big.Int)
	o.AmountSell = new(big.Int)
	o.Expires = new(big.Int)
	o.Nonce = new(big.Int)
	o.FeeMake = new(big.Int)
	o.FeeTake = new(big.Int)

	o.AmountBuy.UnmarshalJSON([]byte(order["amountBuy"].(string)))
	o.AmountSell.UnmarshalJSON([]byte(order["amountSell"].(string)))
	o.Expires.UnmarshalJSON([]byte(order["expires"].(string)))
	o.Nonce.UnmarshalJSON([]byte(order["nonce"].(string)))
	o.FeeMake.UnmarshalJSON([]byte(order["feeMake"].(string)))
	o.FeeTake.UnmarshalJSON([]byte(order["feeTake"].(string)))

	o.PairID = HexToHash(order["pairID"].(string))
	o.Hash = HexToHash(order["hash"].(string))

	signature := order["signature"].(map[string]interface{})
	o.Signature = &Signature{
		V: byte(signature["V"].(float64)),
		R: HexToHash(signature["R"].(string)),
		S: HexToHash(signature["S"].(string)),
	}

	return nil
}

func (o *Order) Decode(order map[string]interface{}) error {
	if order["id"] == nil {
		return errors.New("Order ID not set")
	}
	o.Id = uint64(order["id"].(float64))

	if order["price"] != nil {
		o.Price = uint64(order["price"].(float64))
	}

	if order["amount"] != nil {
		o.Amount = uint64(order["amount"].(float64))
	}

	o.ExchangeAddress = HexToAddress(order["exchangeAddress"].(string))
	o.Maker = HexToAddress(order["maker"].(string))
	o.TokenBuy = HexToAddress(order["tokenBuy"].(string))
	o.TokenSell = HexToAddress(order["tokenSell"].(string))
	o.SymbolBuy = order["symbolBuy"].(string)
	o.SymbolSell = order["symbolSell"].(string)

	o.AmountBuy = new(big.Int)
	o.AmountSell = new(big.Int)
	o.Expires = new(big.Int)
	o.Nonce = new(big.Int)
	o.FeeMake = new(big.Int)
	o.FeeTake = new(big.Int)

	o.AmountBuy.UnmarshalJSON([]byte(order["amountBuy"].(string)))
	o.AmountSell.UnmarshalJSON([]byte(order["amountSell"].(string)))
	o.Expires.UnmarshalJSON([]byte(order["expires"].(string)))
	o.Nonce.UnmarshalJSON([]byte(order["nonce"].(string)))
	o.FeeMake.UnmarshalJSON([]byte(order["feeMake"].(string)))
	o.FeeTake.UnmarshalJSON([]byte(order["feeTake"].(string)))

	o.PairID = HexToHash(order["pairID"].(string))
	o.Hash = HexToHash(order["hash"].(string))

	signature := order["signature"].(map[string]interface{})
	o.Signature = &Signature{
		V: byte(signature["V"].(float64)),
		R: HexToHash(signature["R"].(string)),
		S: HexToHash(signature["S"].(string)),
	}

	return nil
}

// Stringer method for order
func (o *Order) String() string {
	return fmt.Sprintf(
		"Order:\n"+
			"Id: %d\nOrderType: %v\nExchangeAddress: %x\nMaker: %x\nTokenBuy: %x\nTokenSell: %x\n"+
			"AmountBuy: %v\nAmountSell: %v\nSymbolBuy: %v\nSymbolSell: %v\nExpires: %v\nNonce: %v\n"+
			"FeeMake: %v\nFeeTake: %v\nSignature.V: %x\nSignature.R: %x\nSignature.S: %x\nPairID: %x\nHash: %x\nPrice: %v\nAmount: %v\n\n",
		o.Id, o.OrderType, o.ExchangeAddress, o.Maker, o.TokenBuy, o.TokenSell, o.AmountBuy, o.AmountSell, o.SymbolBuy, o.SymbolSell, o.Expires,
		o.Nonce, o.FeeMake, o.FeeTake, o.Signature.V, o.Signature.R, o.Signature.S, o.PairID, o.Hash, o.Price, o.Amount,
	)
}

// PriceInfo prints the following information:
// -ID
// -BuyTokenAmount
// -SellTokenAmount
// -Price
// -Amount
// -Type
func (o *Order) PriceInfo() string {
	return fmt.Sprintf("\nOrder Price Info:\nid: %d\nbuyTokenAmount: %v\nsellTokenAmount: %v\nprice: %v\namount: %v\ntype: %v\n\n", o.Id, o.AmountBuy, o.AmountSell, o.Price, o.Amount, o.OrderType)
}

// TokenInfo prints the following information:
// -BuyToken (address)
// -SellToken (address)
// -BuyToken Symbol
// -SellToken Symbol
func (o *Order) TokenInfo() string {
	return fmt.Sprintf("Order Token Info:\nbuyToken: %x\nsellToken: %x\nbuyTokenSymbol: %v\n, sellTokenSymbol: %v\n", o.TokenBuy, o.TokenSell, o.SymbolBuy, o.SymbolSell)
}

// ComputeHash calculates the order hash
func (o *Order) ComputeHash() Hash {
	sha := sha3.NewKeccak256()
	sha.Write(o.ExchangeAddress.Bytes())
	sha.Write(o.TokenBuy.Bytes())
	sha.Write(BigToHash(o.AmountBuy).Bytes())
	sha.Write(o.TokenSell.Bytes())
	sha.Write(BigToHash(o.AmountSell).Bytes())
	sha.Write(BigToHash(o.Expires).Bytes())
	sha.Write(BigToHash(o.Nonce).Bytes())
	sha.Write(o.Maker.Bytes())
	return BytesToHash(sha.Sum(nil))
}

// VerifySignature checks that the order signature corresponds to the address in the maker field
func (o *Order) VerifySignature() (bool, error) {

	message := crypto.Keccak256(
		[]byte("\x19Ethereum Signed Message:\n32"),
		o.Hash.Bytes(),
	)

	address, err := o.Signature.Verify(BytesToHash(message))
	if err != nil {
		return false, err
	}

	if address != o.Maker {
		return false, errors.New("Recovered address is incorrect")
	}

	return true, nil
}

// ValidateOrder checks the following elements:
//Order Type needs to be equal to BUY or SELL
//Exchange Address needs to be correct
//AmountBuy and AmountSell need to be positive
//OrderHash needs to be correct
func (o *Order) ValidateOrder() (bool, error) {
	return true, nil
}

// NewOrderPlacedEvent is called when an order is first placed in
// the orderbook.
func (o *Order) NewOrderPlacedEvent() *Event {
	return &Event{eventType: ORDER_PLACED, payload: o}
}

// NewOrderMatchedEvent is called when an order is matched (as taker)
// in the orderbook. This does not mean that the order is executed on the
// blockchain as of yet.
func (o *Order) NewOrderMatchedEvent() *Event {
	return &Event{eventType: ORDER_MATCHED, payload: o}
}

// NewOrderPartiallyFilledEvent is called when an order is mached (as taker)
// partially. This does not mean that the order is executed on the
// blockchain as of yet.
func (o *Order) NewOrderPartiallyFilledEvent() *Event {
	return &Event{eventType: ORDER_PARTIALLY_FILLED, payload: o}
}

// NewOrderPartiallyFilledEvent is called when an order is mached (as taker)
//This does not mean that the order is executed on the
// blockchain as of yet.
func (o *Order) NewOrderFilledEvent(t *Trade) *Event {
	payload := &TradePayload{Order: o, Trade: t}
	return &Event{eventType: ORDER_FILLED, payload: payload}
}

// NewOrderCanceled is called when an order is called
func (o *Order) NewOrderCanceledEvent() *Event {
	return &Event{eventType: ORDER_CANCELED, payload: o}
}

// NewOrderExecuted is called when an order is executed meaning that
// the operator has performed a blockchain transaction and is currently
// waiting for the transaction to resolve.
func (o *Order) NewOrderExecutedEvent(tx *types.Transaction) *Event {
	payload := &OrderExecutedPayload{Order: o, Tx: tx.Hash()}
	return &Event{eventType: ORDER_EXECUTED, payload: payload}
}

// NewOrderTransactionSuccessful is called when the operator receives confirmation
// that a trade was carried out successfully.
func (o *Order) NewOrderTxSuccess(t *Trade, tx *types.Transaction) *Event {
	p := &TxSuccessPayload{Order: o, Trade: t, Tx: tx.Hash()}
	return &Event{eventType: ORDER_TX_SUCCESS, payload: p}
}

// NewOrderTransactionError is called when the operator receives an error event
// from the exchange smart contract.
func (o *Order) NewOrderTxError(t *Trade, errorId uint8) *Event {
	p := &TxErrorPayload{Order: o, Trade: t, ErrorId: errorId}
	return &Event{eventType: ORDER_TX_ERROR, payload: p}
}

// NewDoneMessage is used to close certain channels
func NewDoneMessage() *Event {
	return &Event{eventType: DONE}
}

// Sign first calculates the order hash, then computes a signature of this hash
// with the given wallet
func (o *Order) Sign(w *Wallet) error {
	hash := o.ComputeHash()
	sig, err := w.SignHash(hash)
	if err != nil {
		return err
	}

	o.Hash = hash
	o.Signature = sig
	return nil
}
