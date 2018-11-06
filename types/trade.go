package types

import (
	"encoding/json"
	"errors"
	"fmt"
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
	ID             bson.ObjectId  `json:"id,omitempty" bson:"_id"`
	Taker          common.Address `json:"taker" bson:"taker"`
	Maker          common.Address `json:"maker" bson:"maker"`
	BaseToken      common.Address `json:"baseToken" bson:"baseToken"`
	QuoteToken     common.Address `json:"quoteToken" bson:"quoteToken"`
	MakerOrderHash common.Hash    `json:"makerOrderHash" bson:"makerOrderHash"`
	TakerOrderHash common.Hash    `json:"takerOrderHash" bson:"takerOrderHash"`
	Hash           common.Hash    `json:"hash" bson:"hash"`
	TxHash         common.Hash    `json:"txHash" bson:"txHash"`
	PairName       string         `json:"pairName" bson:"pairName"`
	CreatedAt      time.Time      `json:"createdAt" bson:"createdAt" redis:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt" bson:"updatedAt" redis:"updatedAt"`
	PricePoint     *big.Int       `json:"pricepoint" bson:"pricepoint"`
	Status         string         `json:"status" bson:"status"`
	Amount         *big.Int       `json:"amount" bson:"amount"`
}

type TradeRecord struct {
	ID             bson.ObjectId `json:"id" bson:"_id"`
	Taker          string        `json:"taker" bson:"taker"`
	Maker          string        `json:"maker" bson:"maker"`
	BaseToken      string        `json:"baseToken" bson:"baseToken"`
	QuoteToken     string        `json:"quoteToken" bson:"quoteToken"`
	MakerOrderHash string        `json:"makerOrderHash" bson:"makerOrderHash"`
	TakerOrderHash string        `json:"takerOrderHash" bson:"takerOrderHash"`
	Hash           string        `json:"hash" bson:"hash"`
	TxHash         string        `json:"txHash" bson:"txHash"`
	PairName       string        `json:"pairName" bson:"pairName"`
	CreatedAt      time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt" bson:"updatedAt"`
	PricePoint     string        `json:"pricepoint" bson:"pricepoint"`
	Amount         string        `json:"amount" bson:"amount"`
	Status         string        `json:"status" bson:"status"`
}

// NewTrade returns a new unsigned trade corresponding to an Order, amount and taker address
func NewTrade(mo *Order, to *Order, amount *big.Int, pricepoint *big.Int) *Trade {
	t := &Trade{
		Maker:          mo.UserAddress,
		Taker:          to.UserAddress,
		BaseToken:      mo.BaseToken,
		QuoteToken:     mo.QuoteToken,
		MakerOrderHash: mo.Hash,
		TakerOrderHash: to.Hash,
		PairName:       mo.PairName,
		Amount:         amount,
		PricePoint:     pricepoint,
		Status:         "PENDING",
	}

	t.Hash = t.ComputeHash()

	return t
}

func (t *Trade) Validate() error {
	if (t.Taker == common.Address{}) {
		return errors.New("Trade 'taker' parameter is required'")
	}

	if (t.Maker == common.Address{}) {
		return errors.New("Trade 'maker' parameter is required")
	}

	if (t.TakerOrderHash == common.Hash{}) {
		return errors.New("Trade 'takerOrderHash' parameter is required")
	}

	if (t.MakerOrderHash == common.Hash{}) {
		return errors.New("Trade 'makerOrderHash' parameter is required")
	}

	if (t.BaseToken == common.Address{}) {
		return errors.New("Trade 'baseToken' parameter is required")
	}

	if (t.QuoteToken == common.Address{}) {
		return errors.New("Trade 'quoteToken' parameter is required")
	}

	if t.Amount == nil {
		return errors.New("Trade 'amount' parameter is required")
	}

	if t.PricePoint == nil {
		return errors.New("Trade 'pricepoint' paramter is required")
	}

	if math.IsEqualOrSmallerThan(t.PricePoint, big.NewInt(0)) {
		return errors.New("Trade 'pricepoint' parameter should be positive")
	}

	if math.IsEqualOrSmallerThan(t.Amount, big.NewInt(0)) {
		return errors.New("Trade 'amount' parameter should be positive")
	}

	//TODO add validations for hashes and addresses
	return nil
}

// MarshalJSON returns the json encoded byte array representing the trade struct
func (t *Trade) MarshalJSON() ([]byte, error) {
	trade := map[string]interface{}{
		"taker":      t.Taker,
		"maker":      t.Maker,
		"status":     t.Status,
		"hash":       t.Hash,
		"pairName":   t.PairName,
		"pricepoint": t.PricePoint.String(),
		"amount":     t.Amount.String(),
		"createdAt":  t.CreatedAt.Format(time.RFC3339Nano),
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

	if (t.TakerOrderHash != common.Hash{}) {
		trade["takerOrderHash"] = t.TakerOrderHash.Hex()
	}

	if (t.MakerOrderHash != common.Hash{}) {
		trade["makerOrderHash"] = t.MakerOrderHash.Hex()
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

	if trade["makerOrderHash"] == nil {
		return errors.New("Order Hash is not set")
	} else {
		t.MakerOrderHash = common.HexToHash(trade["makerOrderHash"].(string))
	}

	if trade["takerOrderHash"] != nil {
		t.TakerOrderHash = common.HexToHash(trade["takerOrderHash"].(string))
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

	if trade["status"] != nil {
		t.Status = trade["status"].(string)
	}

	if trade["pricepoint"] != nil {
		t.PricePoint = math.ToBigInt(fmt.Sprintf("%v", trade["pricepoint"]))
	}

	if trade["amount"] != nil {
		t.Amount = new(big.Int)
		t.Amount.UnmarshalJSON([]byte(fmt.Sprintf("%v", trade["amount"])))
	}

	if trade["createdAt"] != nil {
		tm, _ := time.Parse(time.RFC3339Nano, trade["createdAt"].(string))
		t.CreatedAt = tm
	}

	return nil
}

func (t *Trade) GetBSON() (interface{}, error) {
	tr := TradeRecord{
		ID:             t.ID,
		PairName:       t.PairName,
		Maker:          t.Maker.Hex(),
		Taker:          t.Taker.Hex(),
		BaseToken:      t.BaseToken.Hex(),
		QuoteToken:     t.QuoteToken.Hex(),
		MakerOrderHash: t.MakerOrderHash.Hex(),
		Hash:           t.Hash.Hex(),
		TxHash:         t.TxHash.Hex(),
		TakerOrderHash: t.TakerOrderHash.Hex(),
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
		PricePoint:     t.PricePoint.String(),
		Status:         t.Status,
		Amount:         t.Amount.String(),
	}

	return tr, nil
}

func (t *Trade) SetBSON(raw bson.Raw) error {
	decoded := new(struct {
		ID             bson.ObjectId `json:"id,omitempty" bson:"_id"`
		PairName       string        `json:"pairName" bson:"pairName"`
		Taker          string        `json:"taker" bson:"taker"`
		Maker          string        `json:"maker" bson:"maker"`
		BaseToken      string        `json:"baseToken" bson:"baseToken"`
		QuoteToken     string        `json:"quoteToken" bson:"quoteToken"`
		MakerOrderHash string        `json:"makerOrderHash" bson:"makerOrderHash"`
		TakerOrderHash string        `json:"takerOrderHash" bson:"takerOrderHash"`
		Hash           string        `json:"hash" bson:"hash"`
		TxHash         string        `json:"txHash" bson:"txHash"`
		CreatedAt      time.Time     `json:"createdAt" bson:"createdAt" redis:"createdAt"`
		UpdatedAt      time.Time     `json:"updatedAt" bson:"updatedAt" redis:"updatedAt"`
		PricePoint     string        `json:"pricepoint" bson:"pricepoint"`
		Status         string        `json:"status" bson:"status"`
		Amount         string        `json:"amount" bson:"amount"`
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
	t.MakerOrderHash = common.HexToHash(decoded.MakerOrderHash)
	t.TakerOrderHash = common.HexToHash(decoded.TakerOrderHash)
	t.Hash = common.HexToHash(decoded.Hash)
	t.TxHash = common.HexToHash(decoded.TxHash)
	t.Status = decoded.Status
	t.Amount = math.ToBigInt(decoded.Amount)
	t.PricePoint = math.ToBigInt(decoded.PricePoint)

	t.CreatedAt = decoded.CreatedAt
	t.UpdatedAt = decoded.UpdatedAt
	return nil
}

// ComputeHash returns hashes the trade
// The OrderHash, Amount, Taker and TradeNonce attributes must be
// set before attempting to compute the trade hash
func (t *Trade) ComputeHash() common.Hash {
	sha := sha3.NewKeccak256()
	sha.Write(t.MakerOrderHash.Bytes())
	sha.Write(t.TakerOrderHash.Bytes())
	return common.BytesToHash(sha.Sum(nil))
}

func (t *Trade) Pair() (*Pair, error) {
	if (t.BaseToken == common.Address{}) {
		return nil, errors.New("Base token is not set")
	}

	if (t.QuoteToken == common.Address{}) {
		return nil, errors.New("Quote token is set")
	}

	return &Pair{
		BaseTokenAddress:  t.BaseToken,
		QuoteTokenAddress: t.QuoteToken,
	}, nil
}

type TradeBSONUpdate struct {
	*Trade
}

func (t TradeBSONUpdate) GetBSON() (interface{}, error) {
	now := time.Now()

	set := bson.M{
		"taker":          t.Taker.Hex(),
		"maker":          t.Maker.Hex(),
		"baseToken":      t.BaseToken.Hex(),
		"quoteToken":     t.QuoteToken.Hex(),
		"makerOrderHash": t.MakerOrderHash.Hex(),
		"takerOrderHash": t.TakerOrderHash.Hex(),
		"txHash":         t.TxHash.Hex(),
		"pairName":       t.PairName,
		"status":         t.Status,
	}

	if t.PricePoint != nil {
		set["pricepoint"] = t.PricePoint.Int64()
	}

	if t.Amount != nil {
		set["amount"] = t.Amount.String()
	}

	setOnInsert := bson.M{
		"_id":       bson.NewObjectId(),
		"hash":      t.Hash.Hex(),
		"createdAt": now,
	}

	update := bson.M{
		"$set":         set,
		"$setOnInsert": setOnInsert,
	}

	return update, nil
}
