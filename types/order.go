package types

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Proofsuite/amp-matching-engine/utils"
	"gopkg.in/mgo.v2/bson"
)

// OrderStatus is used to represent the current status of order.
// It is an enum
type OrderStatus string

// This block declares an enum of type OrderStatus
// containing all possible status of an order.
const (
	NEW           OrderStatus = "NEW"
	OPEN                      = "OPEN"
	MATCHED                   = "MATCHED"
	SUBMITTED                 = "SUBMITTED"
	PARTIALFILLED             = "PARTIAL_FILLED"
	FILLED                    = "FILLED"
	CANCELLED                 = "CANCELLED"
	PENDING                   = "PENDING"
	INVALIDORDER              = "INVALID_ORDER"
	ERROR                     = "ERROR"
)

// UnmarshalJSON unmarshals []byte to type orderStatus
func (orderStatus *OrderStatus) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	value, ok := map[string]OrderStatus{
		"NEW":            NEW,
		"OPEN":           OPEN,
		"MATCHED":        MATCHED,
		"SUBMITTED":      SUBMITTED,
		"PARTIAL_FILLED": PARTIALFILLED,
		"FILLED":         FILLED,
		"CANCELLED":      CANCELLED,
		"PENDING":        PENDING,
		"INVALID_ORDER":  INVALIDORDER,
		"ERROR":          ERROR,
	}[s]
	if !ok {
		return errors.New("Invalid Enum Status Value")
	}

	*orderStatus = value
	return nil
}

// MarshalJSON marshals type orderStatus to []byte.
func (orderStatus *OrderStatus) MarshalJSON() ([]byte, error) {

	value, ok := map[OrderStatus]string{
		NEW:           "NEW",
		OPEN:          "OPEN",
		MATCHED:       "MATCHED",
		SUBMITTED:     "SUBMITTED",
		PARTIALFILLED: "PARTIAL_FILLED",
		FILLED:        "FILLED",
		CANCELLED:     "CANCELLED",
		PENDING:       "PENDING",
		INVALIDORDER:  "INVALID_ORDER",
		ERROR:         "ERROR",
	}[*orderStatus]
	if !ok {
		return nil, errors.New("Invalid Enum Type")
	}
	return json.Marshal(value)
}

// OrderSide is an enum of various buy/sell type of orders
type OrderSide string

// This block declares various members of enum OrderType.
const (
	BUY  OrderSide = "BUY"
	SELL OrderSide = "SELL"
)

// UnmarshalJSON unmarshals []byte to type OrderType
func (orderType *OrderSide) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	value, ok := map[string]OrderSide{"BUY": BUY, "SELL": SELL}[s]
	if !ok {
		return errors.New("Invalid Enum Type Value")
	}
	*orderType = value
	return nil
}

// MarshalJSON marshals type OrderType to []byte
func (orderType *OrderSide) MarshalJSON() ([]byte, error) {
	value, ok := map[OrderSide]string{BUY: "BUY", SELL: "SELL"}[*orderType]
	if !ok {
		return nil, errors.New("Invalid Enum Type")
	}
	return json.Marshal(value)
}

// Order contains the data related to an order sent by the user
type Order struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	BuyToken    string        `json:"buyToken" bson:"buyToken"`
	SellToken   string        `json:"sellToken" bson:"sellToken"`
	BaseToken   string        `json:"baseToken" bson:"baseToken"`
	QuoteToken  string        `json:"quoteToken" bson:"quoteToken"`
	BuyAmount   int64         `json:"buyAmount" bson:"buyAmount"`
	SellAmount  int64         `json:"sellAmount" bson:"sellAmount"`
	Nonce       int64         `json:"nonce" bson:"nonce"`
	UserAddress string        `json:"userAddress" bson:"userAddress"`
	Hash        string        `json:"hash" bson:"hash"`
	Signature   *Signature    `json:"signature,omitempty" bson:"signature"`

	Side         OrderSide    `json:"side" bson:"side"`
	Amount       int64        `json:"amount" bson:"amount"`
	Price        int64        `json:"price" bson:"price"`
	FilledAmount int64        `json:"filledAmount" bson:"filledAmount"`
	Status       OrderStatus  `json:"status" bson:"status"`
	OrderBook    *OrderSubDoc `json:"orderBook" bson:"orderBook"`

	Fee     int64 `json:"fee" bson:"fee"`
	MakeFee int64 `json:"makeFee" bson:"makeFee"`
	TakeFee int64 `json:"takeFee" bson:"takeFee"`

	PairID          bson.ObjectId `json:"pairID" bson:"pairID"`
	PairName        string        `json:"pairName" bson:"pairName"`
	ExchangeAddress string        `json:"exchangeAddress" bson:"exchangeAddress"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// OrderSubDoc is a sub document, it is used to store the order in order book
// It contains the amount that was kept in orderbook alongwith the signature of maker
// It is particularly used in case of partially filled orders.
type OrderSubDoc struct {
	Amount    int64      `json:"amount" bson:"amount"`
	Signature *Signature `json:"signature,omitempty" bson:"signature" redis:"signature"`
}

// ComputeHash calculates the order hash
// func (o *Order) ComputeHash() (ch string) {
// 	sha := sha3.NewKeccak256()
// sha.Write(o.ExchangeAddress.Bytes())
// sha.Write([]byte(o.BaseToken))
// sha.Write([]byte(o.QuoteToken))
// sha.Write(strconv.ParseUint(o.Price))
// sha.Write(BigToHash(o.Amount).Bytes())
// sha.Write(BigToHash(o.Expires).Bytes())
// sha.Write(BigToHash(o.Nonce).Bytes())
// sha.Write(o.Maker.Bytes())
// return BytesToHash(sha.Sum(nil))
// 	return
// }

// GetKVPrefix returns the key value store(redis) prefix to be used
// by matching engine correspondind to a particular order.
func (o *Order) GetKVPrefix() string {
	return o.BaseToken + "::" + o.QuoteToken
}

// GetOBKeys returns the keys corresponding to an order
// orderbook price point key
// orderbook list key corresponding to order price.
func (o *Order) GetOBKeys() (ss, list string) {
	var k string
	if o.Side == BUY {
		k = "buy"
	} else if o.Side == SELL {
		k = "sell"
	}
	ss = o.GetKVPrefix() + "::" + k
	list = o.GetKVPrefix() + "::" + k + "::" + utils.UintToPaddedString(o.Price)
	return
}

// GetOBMatchKey returns the orderbook price point key
// aginst which the order needs to be matched
func (o *Order) GetOBMatchKey() (ss string) {
	var k string
	if o.Side == BUY {
		k = "sell"
	} else if o.Side == SELL {
		k = "buy"
	}

	ss = o.GetKVPrefix() + "::" + k
	return
}
