package types

import (
	"errors"
	"time"

	. "github.com/ethereum/go-ethereum/common"
	"gopkg.in/mgo.v2/bson"
)

// Balance holds both the address and the private key of an ethereum account
type Balance struct {
	ID        bson.ObjectId           `json:"id" bson:"_id"`
	Address   string                  `json:"address" bson:"address"`
	Tokens    map[string]TokenBalance `json:"tokens" bson:"tokens"`
	CreatedAt time.Time               `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time               `json:"updatedAt" bson:"updatedAt"`
}

type TokenBalance struct {
	TokenID      bson.ObjectId `json:"tokenId" bson:"tokenId"`
	Amount       int64         `json:"amount" bson:"amount"`
	LockedAmount int64         `json:"lockedAmount" bson:"lockedAmount"`
}

// NewBalance returns a new wallet object corresponding to a random private key
func NewBalance(address string) (w *Balance, err error) {
	if !IsHexAddress(address) {
		return nil, errors.New("Invalid Address")
	}
	return
}
