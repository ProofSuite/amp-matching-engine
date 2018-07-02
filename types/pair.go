package types

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"labix.org/v2/mgo/bson"
)

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

func (t Pair) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Name, validation.Required),
		validation.Field(&t.BuyToken, validation.Required),
		validation.Field(&t.BuyTokenSymbol, validation.Required),
		validation.Field(&t.SellToken, validation.Required),
		validation.Field(&t.SellTokenSymbol, validation.Required),
	)
}
func (p *Pair) GetOrderBookKeys() (sell, buy string) {
	return p.BuyTokenSymbol + "::" + p.SellTokenSymbol + "::sell", p.BuyTokenSymbol + "::" + p.SellTokenSymbol + "::buy"
}
