package types

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	validation "github.com/go-ozzo/ozzo-validation"
	"gopkg.in/mgo.v2/bson"
)

// Pair struct is used to model the pair data in the system and DB
type Pair struct {
	ID               bson.ObjectId `json:"id" bson:"_id"`
	Name             string        `json:"name" bson:"name"`
	BuyToken         bson.ObjectId `json:"buyToken" bson:"buyToken"`
	BuyTokenSymbol   string        `json:"buyTokenSymbol" bson:"buyTokenSymbol"`
	BuyTokenAddress  string        `json:"buyTokenAddress" bson:"buyTokenAddress"`
	SellToken        bson.ObjectId `json:"sellToken" bson:"sellToken"`
	SellTokenAddress string        `json:"sellTokenAddress" bson:"sellTokenAddress"`
	SellTokenSymbol  string        `json:"sellTokenSymbol" bson:"sellTokenSymbol"`

	MakerFee float64 `json:"makerFee" bson:"makerFee"`
	TakerFee float64 `json:"takerFee" bson:"takerFee"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// Validate function is used to verify if an instance of
// struct satisfies all the conditions for a valid instance
func (p Pair) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.BuyTokenAddress, validation.Required, validation.NewStringRule(common.IsHexAddress, "BuyTokenAddress must be of type HexAddress")),
		validation.Field(&p.SellTokenAddress, validation.Required, validation.NewStringRule(common.IsHexAddress, "SellTokenAddress must be of type HexAddress")),
	)
}

// GetOrderBookKeys returns the orderbook price point keys for corresponding pair
// It is used to fetch the orderbook from redis of a pair
func (p *Pair) GetOrderBookKeys() (sell, buy string) {
	return p.BuyTokenSymbol + "::" + p.SellTokenSymbol + "::sell", p.BuyTokenSymbol + "::" + p.SellTokenSymbol + "::buy"
}
