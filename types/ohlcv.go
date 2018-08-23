package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Proofsuite/amp-matching-engine/utils/math"
	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/mgo.v2/bson"
	"math/big"
)

// Tick is the format in which mongo aggregate pipeline returns data when queried for OHLCV data
type Tick struct {
	ID    TickID   `json:"_id,omitempty" bson:"_id"`
	C     *big.Int `json:"c" bson:"c"`
	Count *big.Int `json:"count" bson:"count"`
	H     *big.Int `json:"h" bson:"h"`
	L     *big.Int `json:"l" bson:"l"`
	O     *big.Int `json:"o" bson:"o"`
	Ts    int64    `json:"ts" bson:"ts"`
	V     *big.Int `json:"v" bson:"v"`
}

// TickID is the subdocument for aggregate grouping for OHLCV data
type TickID struct {
	Pair       string         `json:"pair" bson:"pair"`
	BaseToken  common.Address `json:"baseToken" bson:"baseToken"`
	QuoteToken common.Address `json:"quoteToken" bson:"quoteToken"`
}

type TickRequest struct {
	Pair     []PairSubDoc `json:"pair"`
	From     int64        `json:"from"`
	To       int64        `json:"to"`
	Duration int64        `json:"duration"`
	Units    string       `json:"units"`
}

// MarshalJSON returns the json encoded byte array representing the trade struct
func (t *Tick) MarshalJSON() ([]byte, error) {

	tick := map[string]interface{}{
		"_id": map[string]interface{}{
			"pair":       t.ID.Pair,
			"baseToken":  t.ID.BaseToken.Hex(),
			"quoteToken": t.ID.BaseToken.Hex(),
		},
		"ts": t.Ts,
		"o":  t.O.String(),
		"h":  t.H.String(),
		"l":  t.L.String(),
		"c":  t.C.String(),
		//"v":          t.V.String(),
		"count": t.Count.String(),
	}
	tab, err := json.Marshal(tick)
	return tab, err
}

// UnmarshalJSON creates a trade object from a json byte string
func (t *Tick) UnmarshalJSON(b []byte) error {
	tick := map[string]interface{}{}
	err := json.Unmarshal(b, &tick)

	if err != nil {
		return err
	}
	t.ID = TickID{}
	if tick["_id"] != nil {
		id := tick["_id"].(map[string]interface{})
		if id["quoteToken"] == nil {
			return errors.New("Quote token is not set")
		} else {
			t.ID.QuoteToken = common.HexToAddress(id["quoteToken"].(string))
		}

		if id["baseToken"] == nil {
			return errors.New("Base token is not set")
		} else {
			t.ID.BaseToken = common.HexToAddress(id["baseToken"].(string))
		}

		if id["pair"] == nil {
			return errors.New("Pair is not set")
		} else {
			t.ID.Pair = id["pair"].(string)
		}
	}

	if tick["ts"] == nil {
		return errors.New("ts is not set")
	} else {
		t.Ts = int64(tick["ts"].(float64))
	}
	t.O = new(big.Int)
	t.H = new(big.Int)
	t.L = new(big.Int)
	t.C = new(big.Int)
	t.V = new(big.Int)
	t.Count = new(big.Int)

	if tick["o"] != nil {
		t.O.UnmarshalJSON([]byte(fmt.Sprintf("%v", tick["o"])))
	}
	if tick["h"] != nil {
		t.H.UnmarshalJSON([]byte(fmt.Sprintf("%v", tick["h"])))
	}
	if tick["l"] != nil {
		t.L.UnmarshalJSON([]byte(fmt.Sprintf("%v", tick["l"])))
	}
	if tick["c"] != nil {
		t.C.UnmarshalJSON([]byte(fmt.Sprintf("%v", tick["c"])))
	}
	if tick["v"] != nil {
		t.V.UnmarshalJSON([]byte(fmt.Sprintf("%v", tick["v"])))
	}
	if tick["count"] != nil {
		t.Count.UnmarshalJSON([]byte(fmt.Sprintf("%v", tick["count"])))
	}
	return nil
}

func (t *Tick) GetBSON() (interface{}, error) {
	type TID struct {
		Pair       string `json:"pair" bson:"pair"`
		BaseToken  string `json:"baseToken" bson:"baseToken"`
		QuoteToken string `json:"quoteToken" bson:"quoteToken"`
	}
	return struct {
		ID    TID    `json:"_id,omitempty" bson:"_id"`
		C     string `json:"c" bson:"c"`
		Count string `json:"count" bson:"count"`
		H     string `json:"h" bson:"h"`
		L     string `json:"l" bson:"l"`
		O     string `json:"o" bson:"o"`
		Ts    int64  `json:"ts" bson:"ts"`
		V     string `json:"v" bson:"v"`
	}{
		ID: TID{
			t.ID.Pair,
			t.ID.BaseToken.Hex(),
			t.ID.QuoteToken.Hex(),
		},
		C:     t.C.String(),
		Count: t.Count.String(),
		H:     t.H.String(),
		L:     t.L.String(),
		O:     t.O.String(),
		V:     t.V.String(),
		Ts:    t.Ts,
	}, nil
}

func (t *Tick) SetBSON(raw bson.Raw) error {
	type TID struct {
		Pair       string `json:"pair" bson:"pair"`
		BaseToken  string `json:"baseToken" bson:"baseToken"`
		QuoteToken string `json:"quoteToken" bson:"quoteToken"`
	}
	decoded := new(struct {
		ID    TID    `json:"_id,omitempty" bson:"_id"`
		C     string `json:"c" bson:"c"`
		Count string `json:"count" bson:"count"`
		H     string `json:"h" bson:"h"`
		L     string `json:"l" bson:"l"`
		O     string `json:"o" bson:"o"`
		Ts    int64  `json:"ts" bson:"ts"`
		V     string `json:"v" bson:"v"`
	})

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}
	t.ID = TickID{
		Pair:       decoded.ID.Pair,
		BaseToken:  common.HexToAddress(decoded.ID.BaseToken),
		QuoteToken: common.HexToAddress(decoded.ID.QuoteToken),
	}
	t.Count = new(big.Int)
	t.C = new(big.Int)
	t.H = new(big.Int)
	t.L = new(big.Int)
	t.O = new(big.Int)
	t.V = new(big.Int)
	t.Count = math.ToBigInt(decoded.Count)
	t.C = math.ToBigInt(decoded.C)
	t.H = math.ToBigInt(decoded.H)
	t.L = math.ToBigInt(decoded.L)
	t.O = math.ToBigInt(decoded.O)
	t.V = math.ToBigInt(decoded.V)

	t.Ts = decoded.Ts
	return nil
}
