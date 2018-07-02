package types

import (
	"time"

	. "github.com/ethereum/go-ethereum/common"
	validation "github.com/go-ozzo/ozzo-validation"
	"labix.org/v2/mgo/bson"
)

// Address holds both the address and the private key of an ethereum account
type UserAddress struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Address   string        `json:"address" bson:"address"`
	IsBlocked bool          `json:"isBlocked" bson:"isBlocked"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
}

func (ua UserAddress) Validate() error {
	return validation.ValidateStruct(&ua,
		validation.Field(&ua.Address, validation.Required, validation.NewStringRule(IsHexAddress, "Invalid Address")),
	)
}
