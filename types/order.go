package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/Proofsuite/amp-matching-engine/utils/math"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	validation "github.com/go-ozzo/ozzo-validation"
	"gopkg.in/mgo.v2/bson"
)

// Order contains the data related to an order sent by the user
type Order struct {
	ID              bson.ObjectId  `json:"id" bson:"_id"`
	UserAddress     common.Address `json:"userAddress" bson:"userAddress"`
	ExchangeAddress common.Address `json:"exchangeAddress" bson:"exchangeAddress"`
	BuyToken        common.Address `json:"buyToken" bson:"buyToken"`
	SellToken       common.Address `json:"sellToken" bson:"sellToken"`
	BaseToken       common.Address `json:"baseToken" bson:"baseToken"`
	QuoteToken      common.Address `json:"quoteToken" bson:"quoteToken"`
	BuyAmount       *big.Int       `json:"buyAmount" bson:"buyAmount"`
	SellAmount      *big.Int       `json:"sellAmount" bson:"sellAmount"`
	Status          string         `json:"status" bson:"status"`
	Side            string         `json:"side" bson:"side"`
	Hash            common.Hash    `json:"hash" bson:"hash"`
	Signature       *Signature     `json:"signature,omitempty" bson:"signature"`
	Price           *big.Int       `json:"price" bson:"price"`
	PricePoint      *big.Int       `json:"pricepoint" bson:"pricepoint"`
	Amount          *big.Int       `json:"amount" bson:"amount"`
	FilledAmount    *big.Int       `json:"filledAmount" bson:"filledAmount"`
	Nonce           *big.Int       `json:"nonce" bson:"nonce"`
	Expires         *big.Int       `json:"expires" bson:"expires"`
	MakeFee         *big.Int       `json:"makeFee" bson:"makeFee"`
	TakeFee         *big.Int       `json:"takeFee" bson:"takeFee"`
	OrderBook       *OrderSubDoc   `json:"orderBook" bson:"orderBook"`
	PairName        string         `json:"pairName" bson:"pairName"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// OrderSubDoc is a sub document, it is used to store the order in order book
// It contains the amount that was kept in orderbook alongwith the signature of maker
// It is particularly used in case of partially filled orders.
// type OrderSubDoc struct {
// 	Amount    int64      `json:"amount" bson:"amount"`
// 	Signature *Signature `json:"signature,omitempty" bson:"signature" redis:"signature"`
// }

type OrderSubDoc struct {
	Amount    *big.Int   `json:"amount" bson:"amount"`
	Signature *Signature `json:"signature,omitempty" bson:"signature" redis:"signature"`
}

func (o Order) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ExchangeAddress, validation.Required),
		validation.Field(&o.UserAddress, validation.Required),
		validation.Field(&o.SellToken, validation.Required),
		validation.Field(&o.BuyToken, validation.Required),
		validation.Field(&o.MakeFee, validation.Required),
		validation.Field(&o.TakeFee, validation.Required),
		validation.Field(&o.Nonce, validation.Required),
		//validation.Field(&o.Expires, validation.Required),
		validation.Field(&o.SellAmount, validation.Required),
		validation.Field(&o.UserAddress, validation.Required),
		//validation.Field(&o.Signature, validation.Required),
		// validation.Field(&m.PairName, validation.Required),
	)
}

// ComputeHash calculates the orderRequest hash
func (o *Order) ComputeHash() common.Hash {
	sha := sha3.NewKeccak256()
	sha.Write(o.UserAddress.Bytes())
	sha.Write(o.ExchangeAddress.Bytes())
	sha.Write(o.BuyToken.Bytes())
	sha.Write(o.SellToken.Bytes())
	sha.Write(common.BigToHash(o.BuyAmount).Bytes())
	sha.Write(common.BigToHash(o.SellAmount).Bytes())
	sha.Write(common.BigToHash(o.Expires).Bytes())
	sha.Write(common.BigToHash(o.Nonce).Bytes())
	return common.BytesToHash(sha.Sum(nil))
}

// VerifySignature checks that the orderRequest signature corresponds to the address in the userAddress field
func (o *Order) VerifySignature() (bool, error) {
	o.Hash = o.ComputeHash()
	message := crypto.Keccak256(
		[]byte("\x19Ethereum Signed Message:\n32"),
		o.Hash.Bytes(),
	)
	return true, nil
	address, err := o.Signature.Verify(common.BytesToHash(message))
	if err != nil {
		return false, err
	}

	if address != o.UserAddress {
		return false, errors.New("Recovered address is incorrect")
	}

	return true, nil
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

func (o *Order) Process(p *Pair) error {
	if o.BuyToken == p.BaseTokenAddress {
		o.Side = "BUY"
		o.Amount = o.BuyAmount
		o.Price = math.Div(o.SellAmount, o.BuyAmount)
		o.PricePoint = math.Div(math.Mul(o.SellAmount, big.NewInt(1e8)), o.BuyAmount)
	} else if o.BuyToken == p.QuoteTokenAddress {
		o.Side = "SELL"
		o.Amount = o.SellAmount
		o.Price = math.Div(o.BuyAmount, o.SellAmount)
		o.PricePoint = math.Div(math.Mul(o.BuyAmount, big.NewInt(1e8)), o.SellAmount)
	} else {
		return errors.New("Could not determine o side")
	}

	o.BaseToken = p.BaseTokenAddress
	o.QuoteToken = p.QuoteTokenAddress
	o.PairName = p.GetPairName()
	return nil
}

// temp := big.NewInt(0)
// temp.Mul(o.SellAmount, big.NewInt(1e8))
// o.Price = o.Price.Div(temp, o.BuyAmount)

// o.Price = o.Price.Div(o.BuyAmount, o.SellAmount)
// o.Amount = o.BuyAmount.Int64()
// log.Print(o.SellAmount.Int64() * 1e8)
// o.Price = o.SellAmount.Int64() * 1e8 / o.BuyAmount.Int64()
// log.Println(o.SellAmount)
// log.Println(1e8)
// log.Println(o.SellAmount.Int64() * 1e8)

// Process computes the pricepoint and the amount corresponding to a token pair
// func (o *Order) Process(p *Pair) error {

// 	utils.PrintJSON(o)

// 	btPow := int64(math.Pow10(p.BaseTokenDecimal - app.Config.Decimal))
// 	qtPow := int64(math.Pow10(p.QuoteTokenDecimal - app.Config.Decimal))
// 	sa := new(big.Int)
// 	ba := new(big.Int)

// 	if o.BuyToken == p.BaseTokenAddress {
// 		o.Side = "BUY"
// 		sa.Div(o.SellAmount, o.BuyAmount)
// 		sa.Div(o.SellAmount, big.NewInt(qtPow))
// 		o.Amount = ba.Div(ba, big.NewInt(btPow)).Int64()
// 		o.Price = sa.Int64()
// 		log.Println(o.Price)

// 	} else if o.BuyToken == p.QuoteTokenAddress {
// 		o.Side = "SELL"
// 		ba.Div(o.BuyAmount, o.SellAmount)
// 		ba.Div(ba, big.NewInt(btPow))
// 		o.Amount = sa.Div(sa, big.NewInt(qtPow)).Int64()
// 		o.Price = ba.Int64()
// 		log.Println(o.Price)

// 	} else {
// 		return errors.New("Could not determine order side")
// 	}

// 	o.BaseToken = p.BaseTokenAddress
// 	o.QuoteToken = p.QuoteTokenAddress
// 	o.PairName = p.Name

// 	// utils.PrintJSON(o)

// 	return nil
// }

// GetKVPrefix returns the key value store(redis) prefix to be used
// by matching engine correspondind to a particular order.
func (o *Order) GetKVPrefix() string {
	return o.BaseToken.Hex() + "::" + o.QuoteToken.Hex()
}

// GetOBKeys returns the keys corresponding to an order
// orderbook price point key
// orderbook list key corresponding to order price.
func (o *Order) GetOBKeys() (ss, list string) {
	var k string
	if o.Side == "BUY" {
		k = "BUY"
	} else if o.Side == "SELL" {
		k = "SELL"
	}

	ss = o.GetKVPrefix() + "::" + k
	list = o.GetKVPrefix() + "::" + k + "::" + utils.UintToPaddedString(o.PricePoint.Int64())
	return
}

// GetOBMatchKey returns the orderbook price point key
// aginst which the order needs to be matched
func (o *Order) GetOBMatchKey() (ss string) {
	var k string
	if o.Side == "BUY" {
		k = "SELL"
	} else if o.Side == "SELL" {
		k = "BUY"
	}

	ss = o.GetKVPrefix() + "::" + k
	return
}

// JSON Marshal/Unmarshal interface

// MarshalJSON implements the json.Marshal interface
func (o *Order) MarshalJSON() ([]byte, error) {
	order := map[string]interface{}{
		"exchangeAddress": o.ExchangeAddress,
		"userAddress":     o.UserAddress,
		"buyToken":        o.BuyToken,
		"sellToken":       o.SellToken,
		"baseToken":       o.BaseToken,
		"quoteToken":      o.QuoteToken,
		"side":            o.Side,
		"status":          o.Status,
		"pairName":        o.PairName,
		"buyAmount":       o.BuyAmount.String(),
		"sellAmount":      o.SellAmount.String(),
		"makeFee":         o.MakeFee.String(),
		"takeFee":         o.TakeFee.String(),
		"expires":         o.Expires.String(),
		"nonce":           o.Nonce.String(),
		"price":           o.Price.String(),
		"pricepoint":      o.PricePoint.String(),
		"filledAmount":    o.FilledAmount.String(),
		"amount":          o.Amount.String(),
		"hash":            o.Hash.String(),
		// NOTE: Currently removing this to simplify public API, might reinclude
		// later. An alternative would be to create additional simplified type
		// "createdAt":       o.CreatedAt.Format(time.RFC3339Nano),
		// "updatedAt":       o.UpdatedAt.Format(time.RFC3339Nano),
	}

	// NOTE: Currently removing this to simplify public API, will reinclude
	// if needed. An alternative would be to create additional simplified type
	// if o.ID != bson.ObjectId("") {
	// 	order["id"] = o.ID
	// }

	if o.Signature != nil {
		order["signature"] = map[string]interface{}{
			"V": o.Signature.V,
			"R": o.Signature.R,
			"S": o.Signature.S,
		}
	}

	return json.Marshal(order)
}

func (o *Order) UnmarshalJSON(b []byte) error {
	order := map[string]interface{}{}

	err := json.Unmarshal(b, &order)
	if err != nil {
		return err
	}

	if order["id"] != nil && bson.IsObjectIdHex(order["id"].(string)) {
		o.ID = bson.ObjectIdHex(order["id"].(string))
	}

	if order["pairName"] != nil {
		o.PairName = order["pairName"].(string)
	}

	if order["exchangeAddress"] != nil {
		o.ExchangeAddress = common.HexToAddress(order["exchangeAddress"].(string))
	}

	if order["userAddress"] != nil {
		o.UserAddress = common.HexToAddress(order["userAddress"].(string))
	}

	if order["buyToken"] != nil {
		o.BuyToken = common.HexToAddress(order["buyToken"].(string))
	}

	if order["sellToken"] != nil {
		o.SellToken = common.HexToAddress(order["sellToken"].(string))
	}

	if order["baseToken"] != nil {
		o.BaseToken = common.HexToAddress(order["baseToken"].(string))
	}

	if order["quoteToken"] != nil {
		o.QuoteToken = common.HexToAddress(order["quoteToken"].(string))
	}

	if order["price"] != nil {
		o.Price = math.ToBigInt(order["price"].(string))
	}

	if order["pricepoint"] != nil {
		o.PricePoint = math.ToBigInt(order["pricepoint"].(string))
	}

	if order["amount"] != nil {
		o.Amount = math.ToBigInt(order["amount"].(string))
	}

	if order["filledAmount"] != nil {
		o.FilledAmount = math.ToBigInt(order["filledAmount"].(string))
	}

	if order["buyAmount"] != nil {
		o.BuyAmount = math.ToBigInt(order["buyAmount"].(string))
	}

	if order["sellAmount"] != nil {
		o.SellAmount = math.ToBigInt(order["sellAmount"].(string))
	}

	if order["expires"] != nil {
		o.Expires = math.ToBigInt(order["expires"].(string))
	}

	if order["nonce"] != nil {
		o.Nonce = math.ToBigInt(order["nonce"].(string))
	}

	if order["makeFee"] != nil {
		o.MakeFee = math.ToBigInt(order["makeFee"].(string))
	}

	if order["takeFee"] != nil {
		o.TakeFee = math.ToBigInt(order["takeFee"].(string))
	}

	if order["hash"] != nil {
		o.Hash = common.HexToHash(order["hash"].(string))
	}

	if order["side"] != nil {
		o.Side = order["side"].(string)
	}

	if order["status"] != nil {
		o.Status = order["status"].(string)
	}

	if order["signature"] != nil {
		signature := order["signature"].(map[string]interface{})
		o.Signature = &Signature{
			V: byte(signature["V"].(float64)),
			R: common.HexToHash(signature["R"].(string)),
			S: common.HexToHash(signature["S"].(string)),
		}
	}

	if order["orderBook"] != nil {
		subdoc := order["orderBook"].(map[string]interface{})
		sudocsig := subdoc["signature"].(map[string]interface{})
		o.OrderBook = &OrderSubDoc{
			Amount: math.ToBigInt(subdoc["amount"].(string)),
			Signature: &Signature{
				V: byte(sudocsig["V"].(float64)),
				R: common.HexToHash(sudocsig["R"].(string)),
				S: common.HexToHash(sudocsig["S"].(string)),
			},
		}
	}

	if order["createdAt"] != nil {
		t, _ := time.Parse(time.RFC3339Nano, order["createdAt"].(string))
		o.CreatedAt = t
	}

	if order["updatedAt"] != nil {
		t, _ := time.Parse(time.RFC3339Nano, order["updatedAt"].(string))
		o.UpdatedAt = t
	}

	return nil
}

func (o *OrderSubDoc) MarshalJSON() ([]byte, error) {
	order := map[string]interface{}{}
	if o.Amount != nil {
		order["amount"] = o.Amount.String()
	}

	if o.Signature != nil {
		order["signature"] = map[string]interface{}{
			"V": o.Signature.V,
			"R": o.Signature.R,
			"S": o.Signature.S,
		}
	}

	return json.Marshal(order)
}

func (o *OrderSubDoc) UnmarshalJSON(b []byte) error {
	order := map[string]interface{}{}

	if order["amount"] != nil {
		o.Amount = math.ToBigInt(order["amount"].(string))
	}

	if order["signature"] != nil {
		signature := order["signature"].(map[string]interface{})
		o.Signature = &Signature{
			V: byte(signature["V"].(float64)),
			R: common.HexToHash(signature["R"].(string)),
			S: common.HexToHash(signature["S"].(string)),
		}
	}

	return nil
}

// OrderRecord is the object that will be saved in the database
type OrderRecord struct {
	ID              bson.ObjectId      `json:"id" bson:"_id"`
	UserAddress     string             `json:"userAddress" bson:"userAddress"`
	ExchangeAddress string             `json:"exchangeAddress" bson:"exchangeAddress"`
	BuyToken        string             `json:"buyToken" bson:"buyToken"`
	SellToken       string             `json:"sellToken" bson:"sellToken"`
	BaseToken       string             `json:"baseToken" bson:"baseToken"`
	QuoteToken      string             `json:"quoteToken" bson:"quoteToken"`
	BuyAmount       string             `json:"buyAmount" bson:"buyAmount"`
	SellAmount      string             `json:"sellAmount" bson:"sellAmount"`
	Status          string             `json:"status" bson:"status"`
	Side            string             `json:"side" bson:"side"`
	Hash            string             `json:"hash" bson:"hash"`
	Price           string             `json:"price" bson:"price"`
	PricePoint      string             `json:"pricepoint" bson:"pricepoint"`
	Amount          string             `json:"amount" bson:"amount"`
	FilledAmount    string             `json:"filledAmount" bson:"filledAmount"`
	Nonce           string             `json:"nonce" bson:"nonce"`
	Expires         string             `json:"expires" bson:"expires"`
	MakeFee         string             `json:"makeFee" bson:"makeFee"`
	TakeFee         string             `json:"takeFee" bson:"takeFee"`
	Signature       *SignatureRecord   `json:"signature,omitempty" bson:"signature"`
	OrderBook       *OrderSubDocRecord `json:"orderBook" bson:"orderBook"`

	PairName  string    `json:"pairName" bson:"pairName"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type OrderSubDocRecord struct {
	Amount    string           `json:"amount" bson:"amount"`
	Signature *SignatureRecord `json:"signature" bson:"signature"`
}

func (o *Order) GetBSON() (interface{}, error) {
	or := OrderRecord{
		ID:              o.ID,
		PairName:        o.PairName,
		ExchangeAddress: o.ExchangeAddress.Hex(),
		UserAddress:     o.UserAddress.Hex(),
		BuyToken:        o.BuyToken.Hex(),
		SellToken:       o.SellToken.Hex(),
		BaseToken:       o.BaseToken.Hex(),
		QuoteToken:      o.QuoteToken.Hex(),
		BuyAmount:       o.BuyAmount.String(),
		SellAmount:      o.SellAmount.String(),
		Status:          o.Status,
		Side:            o.Side,
		Hash:            o.Hash.Hex(),
		Price:           o.Price.String(),
		PricePoint:      o.PricePoint.String(),
		Amount:          o.Amount.String(),
		FilledAmount:    o.FilledAmount.String(),
		Nonce:           o.Nonce.String(),
		Expires:         o.Expires.String(),
		MakeFee:         o.MakeFee.String(),
		TakeFee:         o.TakeFee.String(),
		CreatedAt:       o.CreatedAt,
		UpdatedAt:       o.UpdatedAt,
	}

	if o.Signature != nil {
		or.Signature = &SignatureRecord{
			V: o.Signature.V,
			R: o.Signature.R.Hex(),
			S: o.Signature.S.Hex(),
		}
	}

	if o.OrderBook != nil {
		or.OrderBook = &OrderSubDocRecord{
			Amount: o.OrderBook.Amount.String(),
			Signature: &SignatureRecord{
				V: o.OrderBook.Signature.V,
				R: o.OrderBook.Signature.R.Hex(),
				S: o.OrderBook.Signature.S.Hex(),
			},
		}
	}

	return or, nil
}

func (o *Order) SetBSON(raw bson.Raw) error {
	decoded := new(struct {
		ID              bson.ObjectId      `json:"id,omitempty" bson:"_id"`
		PairName        string             `json:"pairName" bson:"pairName"`
		ExchangeAddress string             `json:"exchangeAddress" bson:"exchangeAddress"`
		UserAddress     string             `json:"userAddress" bson:"userAddress"`
		BuyToken        string             `json:"buyToken" bson:"buyToken"`
		SellToken       string             `json:"sellToken" bson:"sellToken"`
		BaseToken       string             `json:"baseToken" bson:"baseToken"`
		QuoteToken      string             `json:"quoteToken" bson:"quoteToken"`
		BuyAmount       string             `json:"buyAmount" bson:"buyAmount"`
		SellAmount      string             `json:"sellAmount" bson:"sellAmount"`
		Status          string             `json:"status" bson:"status"`
		Side            string             `json:"side" bson:"side"`
		Hash            string             `json:"hash" bson:"hash"`
		Price           string             `json:"price" bson:"price"`
		PricePoint      string             `json:"pricepoint" bson:"pricepoint"`
		Amount          string             `json:"amount" bson:"amount"`
		FilledAmount    string             `json:"filledAmount" bson:"filledAmount"`
		Nonce           string             `json:"nonce" bson:"nonce"`
		Expires         string             `json:"expires" bson:"expires"`
		MakeFee         string             `json:"makeFee" bson:"makeFee"`
		TakeFee         string             `json:"takeFee" bson:"takeFee"`
		Signature       *SignatureRecord   `json:"signature" bson:"signature"`
		OrderBook       *OrderSubDocRecord `json:"orderBook" bson:"orderBook"`
		CreatedAt       time.Time          `json:"createdAt" bson:"createdAt"`
		UpdatedAt       time.Time          `json:"updatedAt" bson:"updatedAt"`
	})

	err := raw.Unmarshal(decoded)
	if err != nil {
		log.Print(err)
		return err
	}

	o.ID = decoded.ID
	o.PairName = decoded.PairName
	o.ExchangeAddress = common.HexToAddress(decoded.ExchangeAddress)
	o.UserAddress = common.HexToAddress(decoded.UserAddress)
	o.BuyToken = common.HexToAddress(decoded.BuyToken)
	o.SellToken = common.HexToAddress(decoded.SellToken)
	o.BaseToken = common.HexToAddress(decoded.BaseToken)
	o.QuoteToken = common.HexToAddress(decoded.QuoteToken)

	o.BuyAmount = math.ToBigInt(decoded.BuyAmount)
	o.SellAmount = math.ToBigInt(decoded.SellAmount)
	o.FilledAmount = math.ToBigInt(decoded.FilledAmount)
	o.Amount = math.ToBigInt(decoded.Amount)
	o.Nonce = math.ToBigInt(decoded.Nonce)
	o.Expires = math.ToBigInt(decoded.Expires)
	o.MakeFee = math.ToBigInt(decoded.MakeFee)
	o.TakeFee = math.ToBigInt(decoded.TakeFee)
	o.PricePoint = math.ToBigInt(decoded.PricePoint)
	o.Price = math.ToBigInt(decoded.Price)

	o.Status = decoded.Status
	o.Side = decoded.Side
	o.Hash = common.HexToHash(decoded.Hash)

	if decoded.Signature != nil {
		o.Signature = &Signature{
			V: byte(decoded.Signature.V),
			R: common.HexToHash(decoded.Signature.R),
			S: common.HexToHash(decoded.Signature.S),
		}
	}

	if decoded.OrderBook != nil {
		o.OrderBook = &OrderSubDoc{
			Amount: math.ToBigInt(decoded.OrderBook.Amount),
			Signature: &Signature{
				V: byte(decoded.OrderBook.Signature.V),
				R: common.HexToHash(decoded.OrderBook.Signature.R),
				S: common.HexToHash(decoded.OrderBook.Signature.S),
			},
		}
	}

	o.CreatedAt = decoded.CreatedAt
	o.UpdatedAt = decoded.UpdatedAt

	return nil
}

func (o *Order) Print() {
	b, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Print("\n", string(b))
}

// type OrderData struct {
// 	Order          *Order
// 	Trades         []*Trade
// 	RemainingOrder *Order
// 	MatchingOrders []*FillOrder
// 	FillStatus     FillStatus
// }

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

// UnmarshalJSON unmarshals []byte to type orderStatus
// func (orderStatus *OrderStatus) UnmarshalJSON(data []byte) error {
// 	var s string
// 	err := json.Unmarshal(data, &s)
// 	if err != nil {
// 		return err
// 	}

// 	value, ok := map[string]OrderStatus{
// 		"NEW":            NEW,
// 		"OPEN":           OPEN,
// 		"MATCHED":        MATCHED,
// 		"SUBMITTED":      SUBMITTED,
// 		"PARTIAL_FILLED": PARTIALFILLED,
// 		"FILLED":         FILLED,
// 		"CANCELLED":      CANCELLED,
// 		"PENDING":        PENDING,
// 		"INVALID_ORDER":  INVALIDORDER,
// 		"ERROR":          ERROR,
// 	}[s]
// 	if !ok {
// 		return errors.New("Invalid Enum Status Value")
// 	}

// 	*orderStatus = value
// 	return nil
// }

// // MarshalJSON marshals type orderStatus to []byte.
// func (orderStatus *OrderStatus) MarshalJSON() ([]byte, error) {

// 	value, ok := map[OrderStatus]string{
// 		NEW:           "NEW",
// 		OPEN:          "OPEN",
// 		MATCHED:       "MATCHED",
// 		SUBMITTED:     "SUBMITTED",
// 		PARTIALFILLED: "PARTIAL_FILLED",
// 		FILLED:        "FILLED",
// 		CANCELLED:     "CANCELLED",
// 		PENDING:       "PENDING",
// 		INVALIDORDER:  "INVALID_ORDER",
// 		ERROR:         "ERROR",
// 	}[*orderStatus]
// 	if !ok {
// 		return nil, errors.New("Invalid Enum Type")
// 	}
// 	return json.Marshal(value)
// }

// // OrderSide is an enum of various buy/sell type of orders
// type OrderSide string

// // This block declares various members of enum OrderType.
// const (
// 	BUY  OrderSide = "BUY"
// 	SELL OrderSide = "SELL"
// )

// // UnmarshalJSON unmarshals []byte to type OrderType
// func (orderType *OrderSide) UnmarshalJSON(data []byte) error {
// 	var s string
// 	err := json.Unmarshal(data, &s)
// 	if err != nil {
// 		return err
// 	}

// 	value, ok := map[string]OrderSide{"BUY": BUY, "SELL": SELL}[s]
// 	if !ok {
// 		return errors.New("Invalid Enum Type Value")
// 	}
// 	*orderType = value
// 	return nil
// }

// // MarshalJSON marshals type OrderType to []byte
// func (orderType *OrderSide) MarshalJSON() ([]byte, error) {
// 	value, ok := map[OrderSide]string{BUY: "BUY", SELL: "SELL"}[*orderType]
// 	if !ok {
// 		return nil, errors.New("Invalid Enum Type")
// 	}
// 	return json.Marshal(value)
// }

// OrderStatus is used to represent the current status of order.
// // It is an enum
// type OrderStatus string

// // This block declares an enum of type OrderStatus
// // containing all possible status of an order.
// const (
// 	NEW           OrderStatus = "NEW"
// 	OPEN                      = "OPEN"
// 	MATCHED                   = "MATCHED"
// 	SUBMITTED                 = "SUBMITTED"
// 	PARTIALFILLED             = "PARTIAL_FILLED"
// 	FILLED                    = "FILLED"
// 	CANCELLED                 = "CANCELLED"
// 	PENDING                   = "PENDING"
// 	INVALIDORDER              = "INVALID_ORDER"
// 	ERROR                     = "ERROR"
// )
