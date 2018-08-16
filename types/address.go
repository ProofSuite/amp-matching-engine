package types

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	validation "github.com/go-ozzo/ozzo-validation"
	"gopkg.in/mgo.v2/bson"
)

// UserAddress holds both the address and the private key of an ethereum account
type UserAddress struct {
	ID        bson.ObjectId `json:"-" bson:"_id"`
	Address   string        `json:"address" bson:"address"`
	Nonce     int64         `json:"nonce" bson:"nonce"`
	IsBlocked bool          `json:"isBlocked" bson:"isBlocked"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
}

// Validate function is used to verify if the struct instance is
// valid or not based on user defined rules
func (ua UserAddress) Validate() error {
	return validation.ValidateStruct(&ua,
		validation.Field(&ua.Address, validation.Required, validation.NewStringRule(common.IsHexAddress, "Invalid Address")),
	)
}
