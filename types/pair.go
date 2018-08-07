package types

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	validation "github.com/go-ozzo/ozzo-validation"
	"gopkg.in/mgo.v2/bson"
)

// Pair struct is used to model the pair data in the system and DB
type Pair struct {
	ID                bson.ObjectId `json:"id" bson:"_id"`
	Name              string        `json:"name" bson:"name"`
	BaseTokenID       bson.ObjectId `json:"baseTokenId" bson:"baseTokenId"`
	BaseTokenSymbol   string        `json:"baseTokenSymbol" bson:"baseTokenSymbol"`
	BaseTokenAddress  string        `json:"baseTokenAddress" bson:"baseTokenAddress"`
	QuoteTokenID      bson.ObjectId `json:"quoteTokenId" bson:"quoteTokenId"`
	QuoteTokenAddress string        `json:"quoteTokenAddress" bson:"quoteTokenAddress"`
	QuoteTokenSymbol  string        `json:"quoteTokenSymbol" bson:"quoteTokenSymbol"`

	Active   bool    `json:"active" bson:"active"`
	MakerFee float64 `json:"makerFee" bson:"makerFee"`
	TakerFee float64 `json:"takerFee" bson:"takerFee"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type PairSubDoc struct {
	Name       string `json:"name" bson:"name"`
	BaseToken  string `json:"baseToken" bson:"baseToken"`
	QuoteToken string `json:"quoteToken" bson:"quoteToken"`
}

// Validate function is used to verify if an instance of
// struct satisfies all the conditions for a valid instance
func (p Pair) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.BaseTokenAddress, validation.Required, validation.NewStringRule(common.IsHexAddress, "BaseTokenAddress must be of type HexAddress")),
		validation.Field(&p.QuoteTokenAddress, validation.Required, validation.NewStringRule(common.IsHexAddress, "QuoteTokenAddress must be of type HexAddress")),
	)
}

// GetOrderBookKeys returns the orderbook price point keys for corresponding pair
// It is used to fetch the orderbook from redis of a pair
func (p *Pair) GetOrderBookKeys() (sell, buy string) {
	return p.BaseTokenAddress + "::" + p.QuoteTokenAddress + "::sell", p.BaseTokenAddress + "::" + p.QuoteTokenAddress + "::buy"
}
