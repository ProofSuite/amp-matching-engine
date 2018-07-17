package types

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"gopkg.in/mgo.v2/bson"
)

type OrderStatus int

const (
	NEW OrderStatus = iota
	OPEN
	MATCHED
	SUBMITTED
	PARTIAL_FILLED
	FILLED
	CANCELLED
	PENDING
	INVALID_ORDER
	ERROR
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
		"MATCHED":        MATCHED,
		"SUBMITTED":      SUBMITTED,
		"PARTIAL_FILLED": PARTIAL_FILLED,
		"FILLED":         FILLED,
		"CANCELLED":      CANCELLED,
		"PENDING":        PENDING,
		"INVALID_ORDER":  INVALID_ORDER,
		"ERROR":          ERROR,
	}[s]
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
		MATCHED:        "MATCHED",
		SUBMITTED:      "SUBMITTED",
		PARTIAL_FILLED: "PARTIAL_FILLED",
		FILLED:         "FILLED",
		CANCELLED:      "CANCELLED",
		PENDING:        "PENDING",
		INVALID_ORDER:  "INVALID_ORDER",
		ERROR:          "ERROR",
	}[*orderStatus]
	if !ok {
		return nil, errors.New("Invalid Enum Type")
	}
	return json.Marshal(value)
}

type OrderType int

const (
	_ OrderType = iota
	BUY
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
	ID               bson.ObjectId `json:"id" bson:"_id" redis:"_id"`
	BuyToken         string        `json:"buyToken" bson:"buyToken" redis:"buyToken"`
	SellToken        string        `json:"sellToken" bson:"sellToken" redis:"sellToken"`
	BuyTokenAddress  string        `json:"buyTokenAddress" bson:"buyTokenAddress" redis:"buyTokenAddress"`
	SellTokenAddress string        `json:"sellTokenAddress" bson:"sellTokenAddress" redis:"sellTokenAddress"`
	FilledAmount     int64         `json:"filledAmount" bson:"filledAmount" redis:"filledAmount"`
	Amount           int64         `json:"amount" bson:"amount" redis:"amount"`
	Price            int64         `json:"price" bson:"price" redis:"price"`
	Fee              int64         `json:"fee" bson:"fee" redis:"fee"`
	Type             OrderType     `json:"type" bson:"type" redis:"type"`
	AmountBuy        int64         `json:"amountBuy" bson:"amountBuy" redis:"amountBuy"`
	AmountSell       int64         `json:"amountSell" bson:"amountSell" redis:"amountSell"`
	ExchangeAddress  string        `json:"exchangeAddress" bson:"exchangeAddress" redis:"exchangeAddress"`
	Status           OrderStatus   `json:"status" bson:"status" redis:"status"`
	Signature        *Signature    `json:"signature,omitempty" bson:"signature" redis:"signature"`
	PairID           bson.ObjectId `json:"pairID" bson:"pairID" redis:"pairID"`
	PairName         string        `json:"pairName" bson:"pairName" redis:"pairName"`
	Hash             string        `json:"hash" bson:"hash" redis:"hash"`
	UserAddress      string        `json:"userAddress" bson:"userAddress" redis:"userAddress"`
	OrderBook        *OrderSubDoc  `json:"orderBook" bson:"orderBook"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt" redis:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt" redis:"updatedAt"`
}

type OrderSubDoc struct {
	Amount    int64      `json:"amount" bson:"amount"`
	Signature *Signature `json:"signature,omitempty" bson:"signature" redis:"signature"`
}

// ComputeHash calculates the order hash
func (o *Order) ComputeHash() (ch string) {
	sha := sha3.NewKeccak256()
	// sha.Write(o.ExchangeAddress.Bytes())
	sha.Write([]byte(o.BuyToken))
	sha.Write([]byte(o.SellToken))
	// sha.Write(strconv.ParseUint(o.Price))
	// sha.Write(BigToHash(o.Amount).Bytes())
	// sha.Write(BigToHash(o.Expires).Bytes())
	// sha.Write(BigToHash(o.Nonce).Bytes())
	// sha.Write(o.Maker.Bytes())
	// return BytesToHash(sha.Sum(nil))
	return
}
func (o *Order) GetKVPrefix() string {
	return o.BuyToken + "::" + o.SellToken
}
func (o *Order) GetOBKeys() (ss, list string) {
	var k string
	if o.Type == BUY {
		k = "buy"
	} else if o.Type == SELL {
		k = "sell"
	}
	ss = o.GetKVPrefix() + "::" + k
	list = o.GetKVPrefix() + "::" + k + "::" + utils.UintToPaddedString(o.Price)
	return
}
func (o *Order) GetOBMatchKey() (ss string) {
	var k string
	if o.Type == BUY {
		k = "sell"
	} else if o.Type == SELL {
		k = "buy"
	}

	ss = o.GetKVPrefix() + "::" + k
	return
}
