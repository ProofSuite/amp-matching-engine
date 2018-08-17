package types

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	validation "github.com/go-ozzo/ozzo-validation"
	"gopkg.in/mgo.v2/bson"
)

// Token struct is used to model the token data in the system and DB
type Token struct {
	ID              bson.ObjectId  `json:"id" bson:"_id"`
	Name            string         `json:"name" bson:"name"`
	Symbol          string         `json:"symbol" bson:"symbol"`
	Image           Image          `json:"image" bson:"image"`
	ContractAddress common.Address `json:"contractAddress" bson:"contractAddress"`
	Decimal         int            `json:"decimal" bson:"decimal"`
	Active          bool           `json:"active" bson:"active"`
	Quote           bool           `json:"quote" bson:"quote"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// Image is a sub document used to store data related to images
type Image struct {
	URL  string                 `json:"url" bson:"url"`
	Meta map[string]interface{} `json:"meta" bson:"meta"`
}

// Validate function is used to verify if an instance of
// struct satisfies all the conditions for a valid instance
func (t Token) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Symbol, validation.Required),
		validation.Field(&t.ContractAddress, validation.Required),
		validation.Field(&t.Decimal, validation.Required),
	)
}
