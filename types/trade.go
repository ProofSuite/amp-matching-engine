package types

import (
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/sha3"

	"gopkg.in/mgo.v2/bson"
)

// Trade struct holds arguments corresponding to a "Taker Order"
// To be valid an accept by the matching engine (and ultimately the exchange smart-contract),
// the trade signature must be made from the trader Maker account
type Trade struct {
	ID           bson.ObjectId    `json:"id,omitempty" bson:"_id"`
	OrderHash    string           `json:"orderHash" bson:"orderHash"`
	Amount       int64            `json:"amount" bson:"amount"`
	Price        int64            `json:"price" bson:"price"`
	Side         OrderSide        `json:"side" bson:"side"`
	TradeNonce   int64            `json:"tradeNonce" bson:"tradeNonce"`
	Taker        string           `json:"taker" bson:"taker"`
	Maker        string           `json:"maker" bson:"maker"`
	TakerOrderID bson.ObjectId    `json:"takerOrderId" bson:"takerOrderId"`
	MakerOrderID bson.ObjectId    `json:"makerOrderId" bson:"makerOrderId"`
	Signature    *Signature       `json:"signature" bson:"signature"`
	Hash         string           `json:"hash" bson:"hash"`
	PairName     string           `json:"pairName" bson:"pairName"`
	BaseToken    string           `json:"baseToken" bson:"baseToken"`
	QuoteToken   string           `json:"quoteToken" bson:"quoteToken"`
	Tx           *eth.Transaction `json:"tx" bson:"tx"`
	CreatedAt    time.Time        `json:"createdAt" bson:"createdAt" redis:"createdAt"`
	UpdatedAt    time.Time        `json:"updatedAt" bson:"updatedAt" redis:"updatedAt"`
}

// NewTrade returns a new unsigned trade corresponding to an Order, amount and taker address
func NewTrade(o *Order, amount int64, price int64, taker string) *Trade {
	t := &Trade{
		OrderHash:  o.Hash,
		PairName:   o.PairName,
		Amount:     amount,
		Price:      price,
		TradeNonce: 0,
		Taker:      taker,
		Signature:  &Signature{},
	}
	if o.Side == SELL {
		t.Side = BUY
	} else {
		t.Side = SELL
	}
	return t
}

// ComputeHash returns hashes the trade
//
// The OrderHash, Amount, Taker and TradeNonce attributes must be
// set before attempting to compute the trade hash
func (t *Trade) ComputeHash() string {
	sha := sha3.NewKeccak256()

	sha.Write([]byte(t.OrderHash))
	sha.Write([]byte(fmt.Sprintf("%d", t.Amount)))
	sha.Write([]byte(t.Taker))
	sha.Write([]byte(fmt.Sprintf("%d", t.TradeNonce)))
	return common.BytesToHash(sha.Sum(nil)).Hex()
}

// VerifySignature verifies that the trade is correct and corresponds
// to the trade Taker address
func (t *Trade) VerifySignature() (bool, error) {
	address, err := t.Signature.Verify(common.HexToHash(t.Hash))
	if err != nil {
		return false, err
	}

	if address != common.HexToAddress(t.Taker) {
		return false, errors.New("Recovered address is incorrect")
	}

	return true, nil
}

// Tick is the format in which mongo aggregate pipeline returns data when queried for OHLCV data
type Tick struct {
	ID    TickID `json:"_id,omitempty" bson:"_id"`
	C     int64  `json:"c" bson:"c"`
	Count int64  `json:"count" bson:"count"`
	H     int64  `json:"h" bson:"h"`
	L     int64  `json:"l" bson:"l"`
	O     int64  `json:"o" bson:"o"`
	Ts    int64  `json:"ts" bson:"ts"`
	V     int64  `json:"v" bson:"v"`
}

// TickID is the subdocument for aggregate grouping for OHLCV data
type TickID struct {
	Pair       string `json:"pair" bson:"pair"`
	BaseToken  string `json:"baseToken" bson:"baseToken"`
	QuoteToken string `json:"quoteToken" bson:"quoteToken"`
}

type TickRequest struct {
	Pair     []PairSubDoc `json:"pair"`
	From     int64        `json:"from"`
	To       int64        `json:"to"`
	Duration int64        `json:"duration"`
	Units    string       `json:"units"`
}
