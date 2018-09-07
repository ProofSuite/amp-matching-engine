package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/Proofsuite/amp-matching-engine/utils/math"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"

	"gopkg.in/mgo.v2/bson"
)

// Trade struct holds arguments corresponding to a "Taker Order"
// To be valid an accept by the matching engine (and ultimately the exchange smart-contract),
// the trade signature must be made from the trader Maker account
type Trade struct {
	ID         bson.ObjectId  `json:"id,omitempty" bson:"_id"`
	Taker      common.Address `json:"taker" bson:"taker"`
	Maker      common.Address `json:"maker" bson:"maker"`
	BaseToken  common.Address `json:"baseToken" bson:"baseToken"`
	QuoteToken common.Address `json:"quoteToken" bson:"quoteToken"`
	OrderHash  common.Hash    `json:"orderHash" bson:"orderHash"`
	Hash       common.Hash    `json:"hash" bson:"hash"`
	TxHash     common.Hash    `json:"txHash" bson:"txHash"`
	PairName   string         `json:"pairName" bson:"pairName"`
	TradeNonce *big.Int       `json:"tradeNonce" bson:"tradeNonce"`
	Signature  *Signature     `json:"signature" bson:"signature"`
	CreatedAt  time.Time      `json:"createdAt" bson:"createdAt" redis:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt" bson:"updatedAt" redis:"updatedAt"`
	Price      *big.Int       `json:"price" bson:"price"`
	PricePoint *big.Int       `json:"pricepoint" bson:"pricepoint"`
	Side       string         `json:"side" bson:"side"`
	Amount     *big.Int       `json:"amount" bson:"amount"`
}

type TradeRecord struct {
	ID         bson.ObjectId    `json:"id" bson:"_id"`
	Taker      string           `json:"taker" bson:"taker"`
	Maker      string           `json:"maker" bson:"maker"`
	BaseToken  string           `json:"baseToken" bson:"baseToken"`
	QuoteToken string           `json:"quoteToken" bson:"quoteToken"`
	OrderHash  string           `json:"orderHash" bson:"orderHash"`
	Hash       string           `json:"hash" bson:"hash"`
	TxHash     string           `json:"txHash" bson:"txHash"`
	PairName   string           `json:"pairName" bson:"pairName"`
	TradeNonce string           `json:"tradeNonce" bson:"tradeNonce"`
	Signature  *SignatureRecord `json:"signature" bson:"signature"`
	CreatedAt  time.Time        `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time        `json:"updatedAt" bson:"updatedAt"`
	Price      string           `json:"price" bson:"price"`
	PricePoint string           `json:"pricepoint" bson:"pricepoint"`
	Side       string           `json:"side" bson:"side"`
	Amount     string           `json:"amount" bson:"amount"`
}

// NewTrade returns a new unsigned trade corresponding to an Order, amount and taker address
func NewTrade(o *Order, amount *big.Int, price *big.Int, taker common.Address) *Trade {
	t := &Trade{
		OrderHash:  o.Hash,
		PairName:   o.PairName,
		Amount:     amount,
		Price:      price,
		TradeNonce: big.NewInt(0),
		Side:       o.Side,
		Taker:      taker,
		Signature:  &Signature{},
	}

	return t
}

// MarshalJSON returns the json encoded byte array representing the trade struct
func (t *Trade) MarshalJSON() ([]byte, error) {
	trade := map[string]interface{}{
		"taker":      t.Taker,
		"maker":      t.Maker,
		"baseToken":  t.BaseToken,
		"quoteToken": t.QuoteToken,
		"orderHash":  t.OrderHash,
		"side":       t.Side,
		"hash":       t.Hash,
		"txHash":     t.TxHash,
		"pairName":   t.PairName,
		"tradeNonce": t.TradeNonce.String(),
		// NOTE: I don't these are publicly needed but leaving this here until confirmation
		// "createdAt":    t.CreatedAt.String(),
		// "updatedAt":    t.UpdatedAt.String(),
		"price":      t.Price.String(),
		"pricepoint": t.PricePoint.String(),
		"amount":     t.Amount.String(),
	}

	if (t.BaseToken != common.Address{}) {
		trade["baseToken"] = t.BaseToken.Hex()
	}

	if (t.QuoteToken != common.Address{}) {
		trade["quoteToken"] = t.QuoteToken.Hex()
	}

	if (t.TxHash != common.Hash{}) {
		trade["txHash"] = t.TxHash.Hex()
	}

	// NOTE: Currently remove marshalling of IDs to simplify public API but will uncommnent
	// if needed.
	// if t.ID != bson.ObjectId("") {
	// 	trade["id"] = t.ID
	// }

	if t.Signature != nil {
		trade["signature"] = map[string]interface{}{
			"V": t.Signature.V,
			"R": t.Signature.R,
			"S": t.Signature.S,
		}
	}

	return json.Marshal(trade)
}

// UnmarshalJSON creates a trade object from a json byte string
func (t *Trade) UnmarshalJSON(b []byte) error {
	trade := map[string]interface{}{}

	err := json.Unmarshal(b, &trade)
	if err != nil {
		return err
	}

	if trade["orderHash"] == nil {
		return errors.New("Order Hash is not set")
	} else {
		t.OrderHash = common.HexToHash(trade["orderHash"].(string))
	}

	if trade["hash"] == nil {
		return errors.New("Hash is not set")
	} else {
		t.Hash = common.HexToHash(trade["hash"].(string))
	}

	if trade["quoteToken"] == nil {
		return errors.New("Quote token is not set")
	} else {
		t.QuoteToken = common.HexToAddress(trade["quoteToken"].(string))
	}

	if trade["baseToken"] == nil {
		return errors.New("Base token is not set")
	} else {
		t.BaseToken = common.HexToAddress(trade["baseToken"].(string))
	}

	if trade["maker"] == nil {
		return errors.New("Maker is not set")
	} else {
		t.Taker = common.HexToAddress(trade["taker"].(string))
	}

	if trade["taker"] == nil {
		return errors.New("Taker is not set")
	} else {
		t.Maker = common.HexToAddress(trade["maker"].(string))
	}

	if trade["id"] != nil && bson.IsObjectIdHex(trade["id"].(string)) {
		t.ID = bson.ObjectIdHex(trade["id"].(string))
	}

	if trade["txHash"] != nil {
		t.TxHash = common.HexToHash(trade["txHash"].(string))
	}

	if trade["pairName"] != nil {
		t.PairName = trade["pairName"].(string)
	}

	if trade["side"] != nil {
		t.Side = trade["side"].(string)
	}

	if trade["price"] != nil {
		t.Price = math.ToBigInt(fmt.Sprintf("%v", trade["price"]))
	}

	if trade["pricepoint"] != nil {
		t.PricePoint = math.ToBigInt(fmt.Sprintf("%v", trade["pricepoint"]))
	}

	if trade["amount"] != nil {
		t.Amount = new(big.Int)
		t.Amount.UnmarshalJSON([]byte(fmt.Sprintf("%v", trade["amount"])))
	}

	if trade["tradeNonce"] != nil {
		t.TradeNonce = new(big.Int)
		t.TradeNonce.UnmarshalJSON([]byte(fmt.Sprintf("%v", trade["tradeNonce"])))
	}

	if trade["signature"] != nil {
		signature := trade["signature"].(map[string]interface{})
		t.Signature = &Signature{
			V: byte(signature["V"].(float64)),
			R: common.HexToHash(signature["R"].(string)),
			S: common.HexToHash(signature["S"].(string)),
		}
	}

	return nil
}

func (t *Trade) GetBSON() (interface{}, error) {
	tr := TradeRecord{
		ID:         t.ID,
		PairName:   t.PairName,
		Maker:      t.Maker.Hex(),
		Taker:      t.Taker.Hex(),
		BaseToken:  t.BaseToken.Hex(),
		QuoteToken: t.QuoteToken.Hex(),
		OrderHash:  t.OrderHash.Hex(),
		TradeNonce: t.TradeNonce.String(),
		Hash:       t.Hash.Hex(),
		TxHash:     t.TxHash.Hex(),
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
		PricePoint: t.PricePoint.String(),
		Price:      t.Price.String(),
		Side:       t.Side,
		Amount:     t.Amount.String(),
	}

	if t.Signature != nil {
		tr.Signature = &SignatureRecord{
			V: t.Signature.V,
			R: t.Signature.R.Hex(),
			S: t.Signature.S.Hex(),
		}
	}

	return tr, nil
}

func (t *Trade) SetBSON(raw bson.Raw) error {
	decoded := new(struct {
		ID         bson.ObjectId    `json:"id,omitempty" bson:"_id"`
		PairName   string           `json:"pairName" bson:"pairName"`
		Taker      string           `json:"taker" bson:"taker"`
		Maker      string           `json:"maker" bson:"maker"`
		BaseToken  string           `json:"baseToken" bson:"baseToken"`
		QuoteToken string           `json:"quoteToken" bson:"quoteToken"`
		OrderHash  string           `json:"orderHash" bson:"orderHash"`
		Hash       string           `json:"hash" bson:"hash"`
		TxHash     string           `json:"txHash" bson:"txHash"`
		TradeNonce string           `json:"tradeNonce" bson:"tradeNonce"`
		Signature  *SignatureRecord `json:"signature" bson:"signature"`
		CreatedAt  time.Time        `json:"createdAt" bson:"createdAt" redis:"createdAt"`
		UpdatedAt  time.Time        `json:"updatedAt" bson:"updatedAt" redis:"updatedAt"`
		Price      string           `json:"price" bson:"price"`
		PricePoint string           `json:"pricepoint" bson:"pricepoint"`
		Side       string           `json:"side" bson:"side"`
		Amount     string           `json:"amount" bson:"amount"`
	})

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}

	t.ID = decoded.ID
	t.PairName = decoded.PairName
	t.Taker = common.HexToAddress(decoded.Taker)
	t.Maker = common.HexToAddress(decoded.Maker)
	t.BaseToken = common.HexToAddress(decoded.BaseToken)
	t.QuoteToken = common.HexToAddress(decoded.QuoteToken)
	t.OrderHash = common.HexToHash(decoded.OrderHash)
	t.Hash = common.HexToHash(decoded.Hash)
	t.TxHash = common.HexToHash(decoded.TxHash)

	t.TradeNonce = math.ToBigInt(decoded.TradeNonce)
	t.Amount = math.ToBigInt(decoded.Amount)
	t.Price = math.ToBigInt(decoded.Price)
	t.PricePoint = math.ToBigInt(decoded.PricePoint)
	t.Side = decoded.Side

	if decoded.Signature != nil {
		t.Signature = &Signature{
			V: byte(decoded.Signature.V),
			R: common.HexToHash(decoded.Signature.R),
			S: common.HexToHash(decoded.Signature.S),
		}
	}

	t.CreatedAt = decoded.CreatedAt
	t.UpdatedAt = decoded.UpdatedAt
	return nil
}

// ComputeHash returns hashes the trade
// The OrderHash, Amount, Taker and TradeNonce attributes must be
// set before attempting to compute the trade hash
func (t *Trade) ComputeHash() common.Hash {
	sha := sha3.NewKeccak256()

	sha.Write(t.OrderHash.Bytes())
	sha.Write(t.Taker.Bytes())
	sha.Write(common.BigToHash(t.Amount).Bytes())
	sha.Write(common.BigToHash(t.TradeNonce).Bytes())
	return common.BytesToHash(sha.Sum(nil))
}

// VerifySignature verifies that the trade is correct and corresponds
// to the trade Taker address
func (t *Trade) VerifySignature() (bool, error) {
	address, err := t.Signature.Verify(t.Hash)
	if err != nil {
		return false, err
	}

	if address != t.Taker {
		return false, errors.New("Recovered address is incorrect")
	}

	return true, nil
}

// Sign calculates ands sets the trade hash and signature with the
// given wallet
func (t *Trade) Sign(w *Wallet) error {
	hash := t.ComputeHash()
	signature, err := w.SignHash(hash)
	if err != nil {
		return err
	}

	t.Hash = hash
	t.Signature = signature
	return nil
}

func (t *Trade) Print() {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		log.Print(err)
	}

	fmt.Print(string(b))
}

// NewTrade returns a new trade with the given params. The trade is signed by the factory wallet.
// Currently the nonce is chosen randomly which will be changed in the future
func NewUnsignedTrade(o *Order, taker common.Address, amount *big.Int) (Trade, error) {
	t := Trade{}
	t.Maker = o.UserAddress
	t.BaseToken = o.BaseToken
	t.QuoteToken = o.QuoteToken
	t.Price = o.Price
	t.PricePoint = o.PricePoint
	t.OrderHash = o.Hash
	t.Taker = taker
	t.Amount = amount

	if o.Side == "BUY" {
		t.Side = "SELL"
	} else if o.Side == "SELL" {
		t.Side = "BUY"
	}

	return t, nil
}

//Replacement for function above
func NewUnsignedTrade1(maker *Order, taker *Order, amount *big.Int) (Trade, error) {
	t := Trade{}
	t.Maker = maker.UserAddress
	t.Taker = taker.UserAddress
	t.BaseToken = maker.BaseToken
	t.QuoteToken = maker.QuoteToken
	t.Price = taker.Price
	t.PricePoint = taker.PricePoint
	t.OrderHash = maker.Hash

	//TODO compute from taker amount and maker amount
	t.Amount = amount

	if maker.Side == "BUY" {
		t.Side = "SELL"
	} else if maker.Side == "SELL" {
		t.Side = "BUY"
	}

	return t, nil
}
