package types

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/Proofsuite/amp-matching-engine/utils/math"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/globalsign/mgo/bson"
)

// Token struct is used to model the token data in the system and DB
type Token struct {
	ID       bson.ObjectId  `json:"-" bson:"_id"`
	Symbol   string         `json:"symbol" bson:"symbol"`
	Address  common.Address `json:"address" bson:"address"`
	Decimals int            `json:"decimals" bson:"decimals"`
	Active   bool           `json:"active" bson:"active"`
	Listed   bool           `json:"listed" bson:"listed"`
	Quote    bool           `json:"quote" bson:"quote"`
	MakeFee  *big.Int       `json:"makeFee,omitempty" bson:"makeFee,omitempty"`
	TakeFee  *big.Int       `json:"takeFee,omitempty" bson:"makeFee,omitempty"`
	Rank     int            `json:"rank,omitempty" bson:"rank,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// TokenRecord is the struct which is stored in db
type TokenRecord struct {
	ID       bson.ObjectId `json:"-" bson:"_id"`
	Symbol   string        `json:"symbol" bson:"symbol"`
	Address  string        `json:"address" bson:"address"`
	Decimals int           `json:"decimals" bson:"decimals"`
	Active   bool          `json:"active" bson:"active"`
	Listed   bool          `json:"listed" bson:"listed"`
	Quote    bool          `json:"quote" bson:"quote"`
	MakeFee  string        `json:"makeFee,omitempty" bson:"makeFee,omitempty"`
	TakeFee  string        `json:"takeFee,omitempty" bson:"takeFee,omitempty"`
	Rank     int           `json:"rank,omitempty" bson:"rank,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// Validate function is used to verify if an instance of
// struct satisfies all the conditions for a valid instance
func (t Token) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Symbol, validation.Required),
		validation.Field(&t.Address, validation.Required),
		validation.Field(&t.Decimals, validation.Required),
	)
}

func (t *Token) MarshalJSON() ([]byte, error) {
	token := map[string]interface{}{
		"id":        t.ID,
		"symbol":    t.Symbol,
		"address":   t.Address.Hex(),
		"decimals":  t.Decimals,
		"active":    t.Active,
		"listed":    t.Listed,
		"quote":     t.Quote,
		"createdAt": t.CreatedAt.Format(time.RFC3339Nano),
		"updatedAt": t.UpdatedAt.Format(time.RFC3339Nano),
		"rank":      t.Rank,
	}

	if t.MakeFee != nil {
		token["makeFee"] = t.MakeFee.String()
	}

	if t.TakeFee != nil {
		token["takeFee"] = t.TakeFee.String()
	}

	return json.Marshal(token)
}

func (t *Token) UnmarshalJSON(b []byte) error {
	token := map[string]interface{}{}

	err := json.Unmarshal(b, &token)
	if err != nil {
		return err
	}

	if token["address"] != nil {
		t.Address = common.HexToAddress(token["address"].(string))
	}

	if token["listed"] != nil {
		t.Listed = token["listed"].(bool)
	}

	if token["quote"] != nil {
		t.Quote = token["quote"].(bool)
	}

	if token["active"] != nil {
		t.Active = token["active"].(bool)
	}

	if token["decimals"] != nil {
		t.Decimals = token["decimals"].(int)
	}

	if token["symbol"] != nil {
		t.Symbol = token["symbol"].(string)
	}

	if token["id"] != nil {
		t.ID = bson.ObjectIdHex(token["id"].(string))
	}

	if token["createdAt"] != nil {
		tm, _ := time.Parse(time.RFC3339Nano, token["createdAt"].(string))
		t.CreatedAt = tm
	}

	if token["updatedAt"] != nil {
		tm, _ := time.Parse(time.RFC3339Nano, token["updatedAt"].(string))
		t.UpdatedAt = tm
	}

	if token["makeFee"] != nil {
		t.MakeFee = math.ToBigInt(token["makeFee"].(string))
	}

	if token["takeFee"] != nil {
		t.TakeFee = math.ToBigInt(token["takeFee"].(string))
	}

	if token["rank"] != nil {
		t.Rank = token["rank"].(int)
	}

	return nil
}

// GetBSON implements bson.Getter
func (t *Token) GetBSON() (interface{}, error) {
	tr := TokenRecord{
		ID:        t.ID,
		Symbol:    t.Symbol,
		Address:   t.Address.Hex(),
		Decimals:  t.Decimals,
		Active:    t.Active,
		Listed:    t.Listed,
		Quote:     t.Quote,
		Rank:      t.Rank,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}

	if t.MakeFee != nil {
		tr.MakeFee = t.MakeFee.String()
	}

	if t.TakeFee != nil {
		tr.TakeFee = t.TakeFee.String()
	}

	return tr, nil
}

// SetBSON implemenets bson.Setter
func (t *Token) SetBSON(raw bson.Raw) error {
	decoded := &TokenRecord{}

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}

	t.ID = decoded.ID
	t.Symbol = decoded.Symbol
	if common.IsHexAddress(decoded.Address) {
		t.Address = common.HexToAddress(decoded.Address)
	}

	t.Decimals = decoded.Decimals
	t.Active = decoded.Active
	t.Listed = decoded.Listed
	t.Quote = decoded.Quote
	t.CreatedAt = decoded.CreatedAt
	t.UpdatedAt = decoded.UpdatedAt
	t.Rank = decoded.Rank

	if decoded.MakeFee != "" {
		t.MakeFee = math.ToBigInt(decoded.MakeFee)
	}

	if decoded.TakeFee != "" {
		t.TakeFee = math.ToBigInt(decoded.TakeFee)
	}

	return nil
}
