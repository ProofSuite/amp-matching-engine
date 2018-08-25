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
		"v":          t.V.String(),
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
	fmt.Print(tick)
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
	fmt.Println("CAME HERE.. hopefully can be done")
	type TID struct {
		Pair       string `json:"pair" bson:"pair"`
		BaseToken  string `json:"baseToken" bson:"baseToken"`
		QuoteToken string `json:"quoteToken" bson:"quoteToken"`
	}

	count, err := bson.ParseDecimal128(t.Count.String())
	if err != nil {
		return nil, err
	}
	o, err := bson.ParseDecimal128(t.O.String())
	if err != nil {
		return nil, err
	}
	h, err := bson.ParseDecimal128(t.H.String())
	if err != nil {
		return nil, err
	}
	l, err := bson.ParseDecimal128(t.L.String())
	if err != nil {
		return nil, err
	}
	c, err := bson.ParseDecimal128(t.C.String())
	if err != nil {
		return nil, err
	}
	v, err := bson.ParseDecimal128(t.V.String())
	if err != nil {
		return nil, err
	}

	return struct {
		ID    TID             `json:"_id,omitempty" bson:"_id"`
		Count bson.Decimal128 `json:"count" bson:"count"`
		O     bson.Decimal128 `json:"o" bson:"o"`
		H     bson.Decimal128 `json:"h" bson:"h"`
		L     bson.Decimal128 `json:"l" bson:"l"`
		C     bson.Decimal128 `json:"c" bson:"c"`
		V     bson.Decimal128 `json:"v" bson:"v"`
		Ts    int64           `json:"ts" bson:"ts"`
	}{
		ID: TID{
			t.ID.Pair,
			t.ID.BaseToken.Hex(),
			t.ID.QuoteToken.Hex(),
		},

		O:     o,
		H:     h,
		L:     l,
		C:     c,
		V:     v,
		Count: count,
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
		ID    TID             `json:"_id,omitempty" bson:"_id"`
		Count bson.Decimal128 `json:"count" bson:"count"`
		O     bson.Decimal128 `json:"o" bson:"o"`
		H     bson.Decimal128 `json:"h" bson:"h"`
		L     bson.Decimal128 `json:"l" bson:"l"`
		C     bson.Decimal128 `json:"c" bson:"c"`
		V     bson.Decimal128 `json:"v" bson:"v"`
		Ts    int64           `json:"ts" bson:"ts"`
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

	count := decoded.Count.String()
	o := decoded.O.String()
	h := decoded.H.String()
	l := decoded.L.String()
	c := decoded.C.String()
	v := decoded.V.String()
	t.Count = math.ToBigInt(count)
	t.C = math.ToBigInt(c)
	t.H = math.ToBigInt(h)
	t.L = math.ToBigInt(l)
	t.O = math.ToBigInt(o)
	t.V = math.ToBigInt(v)

	t.Ts = decoded.Ts
	return nil
}
