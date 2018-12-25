package types

import (
	"encoding/json"
	"math/big"

	"github.com/Proofsuite/amp-matching-engine/utils/math"
	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo/bson"
)

// Tick is the format in which mongo aggregate pipeline returns data when queried for OHLCV data
type Tick struct {
	Pair      PairID   `json:"id,omitempty" bson:"_id"`
	Close     *big.Int `json:"close,omitempty" bson:"close"`
	Count     *big.Int `json:"count,omitempty" bson:"count"`
	High      *big.Int `json:"high,omitempty" bson:"high"`
	Low       *big.Int `json:"low,omitempty" bson:"low"`
	Open      *big.Int `json:"open,omitempty" bson:"open"`
	Volume    *big.Int `json:"volume,omitempty" bson:"volume"`
	Timestamp int64    `json:"timestamp,omitempty" bson:"timestamp"`
}

// PairID is the subdocument for aggregate grouping for OHLCV data
type PairID struct {
	PairName   string         `json:"pair" bson:"pairName"`
	BaseToken  common.Address `json:"baseToken" bson:"baseToken"`
	QuoteToken common.Address `json:"quoteToken" bson:"quoteToken"`
}

type OHLCVParams struct {
	Pair     []PairAddresses `json:"pair"`
	From     int64           `json:"from"`
	To       int64           `json:"to"`
	Duration int64           `json:"duration"`
	Units    string          `json:"units"`
}

// MarshalJSON returns the json encoded byte array representing the trade struct
func (t *Tick) MarshalJSON() ([]byte, error) {
	tick := map[string]interface{}{
		"pair": map[string]interface{}{
			"pairName":   t.Pair.PairName,
			"baseToken":  t.Pair.BaseToken.Hex(),
			"quoteToken": t.Pair.QuoteToken.Hex(),
		},
		"timestamp": t.Timestamp,
	}

	if t.Open != nil {
		tick["open"] = t.Open.String()
	}

	if t.High != nil {
		tick["high"] = t.High.String()
	}

	if t.Low != nil {
		tick["low"] = t.Low.String()
	}

	if t.Volume != nil {
		tick["volume"] = t.Volume.String()
	}

	if t.Close != nil {
		tick["close"] = t.Close.String()
	}

	if t.Count != nil {
		tick["count"] = t.Count.String()
	}

	bytes, err := json.Marshal(tick)
	return bytes, err
}

// UnmarshalJSON creates a trade object from a json byte string
func (t *Tick) UnmarshalJSON(b []byte) error {
	tick := map[string]interface{}{}
	err := json.Unmarshal(b, &tick)

	if err != nil {
		return err
	}

	if tick["pair"] != nil {
		pair := tick["pair"].(map[string]interface{})
		t.Pair = PairID{
			PairName:   pair["pairName"].(string),
			BaseToken:  common.HexToAddress(pair["baseToken"].(string)),
			QuoteToken: common.HexToAddress(pair["quoteToken"].(string)),
		}
	}

	if tick["timestamp"] != nil {
		t.Timestamp = int64(tick["timestamp"].(float64))
	}

	if tick["open"] != nil {
		t.Open = math.ToBigInt(tick["open"].(string))
	}

	if tick["high"] != nil {
		t.High = math.ToBigInt(tick["high"].(string))
	}

	if tick["low"] != nil {
		t.Low = math.ToBigInt(tick["low"].(string))
	}

	if tick["close"] != nil {
		t.Close = math.ToBigInt(tick["close"].(string))
	}

	if tick["volume"] != nil {
		t.Volume = math.ToBigInt(tick["volume"].(string))
	}

	if tick["count"] != nil {
		t.Count = math.ToBigInt(tick["count"].(string))
	}

	return nil
}

func (t *Tick) GetBSON() (interface{}, error) {
	type PairID struct {
		PairName   string `json:"pairName" bson:"pairName"`
		BaseToken  string `json:"baseToken" bson:"baseToken"`
		QuoteToken string `json:"quoteToken" bson:"quoteToken"`
	}

	count, err := bson.ParseDecimal128(t.Count.String())
	if err != nil {
		return nil, err
	}

	o := t.Open.String()
	h := t.High.String()
	l := t.Low.String()
	c := t.Close.String()

	v, err := bson.ParseDecimal128(t.Volume.String())
	if err != nil {
		return nil, err
	}

	return struct {
		ID        PairID          `json:"id,omitempty" bson:"_id"`
		Count     bson.Decimal128 `json:"count" bson:"count"`
		Open      string          `json:"open" bson:"open"`
		High      string          `json:"high" bson:"high"`
		Low       string          `json:"low" bson:"low"`
		Close     string          `json:"close" bson:"close"`
		Volume    bson.Decimal128 `json:"volume" bson:"volume"`
		Timestamp int64           `json:"timestamp" bson:"timestamp"`
	}{
		ID: PairID{
			t.Pair.PairName,
			t.Pair.BaseToken.Hex(),
			t.Pair.QuoteToken.Hex(),
		},

		Open:      o,
		High:      h,
		Low:       l,
		Close:     c,
		Volume:    v,
		Count:     count,
		Timestamp: t.Timestamp,
	}, nil
}

func (t *Tick) SetBSON(raw bson.Raw) error {
	type PairIDRecord struct {
		PairName   string `json:"pairName" bson:"pairName"`
		BaseToken  string `json:"baseToken" bson:"baseToken"`
		QuoteToken string `json:"quoteToken" bson:"quoteToken"`
	}

	decoded := new(struct {
		Pair      PairIDRecord    `json:"pair,omitempty" bson:"_id"`
		Count     bson.Decimal128 `json:"count" bson:"count"`
		Open      string          `json:"open" bson:"open"`
		High      string          `json:"high" bson:"high"`
		Low       string          `json:"low" bson:"low"`
		Close     string          `json:"close" bson:"close"`
		Volume    bson.Decimal128 `json:"volume" bson:"volume"`
		Timestamp int64           `json:"timestamp" bson:"timestamp"`
	})

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}

	t.Pair = PairID{
		PairName:   decoded.Pair.PairName,
		BaseToken:  common.HexToAddress(decoded.Pair.BaseToken),
		QuoteToken: common.HexToAddress(decoded.Pair.QuoteToken),
	}

	count := decoded.Count.String()
	o := decoded.Open
	h := decoded.High
	l := decoded.Low
	c := decoded.Close
	v := decoded.Volume.String()

	t.Count = math.ToBigInt(count)
	t.Close = math.ToBigInt(c)
	t.High = math.ToBigInt(h)
	t.Low = math.ToBigInt(l)
	t.Open = math.ToBigInt(o)
	t.Volume = math.ToBigInt(v)

	t.Timestamp = decoded.Timestamp
	return nil
}

func (t *Tick) AddressCode() string {
	code := t.Pair.BaseToken.Hex() + "::" + t.Pair.QuoteToken.Hex()
	return code
}
