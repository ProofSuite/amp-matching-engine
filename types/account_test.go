package types

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestAccountBSON(t *testing.T) {
	assert := assert.New(t)

	address := NewAddressFromString("0xe8e84ee367bc63ddb38d3d01bccef106c194dc47")
	tokenAddress1 := NewAddressFromString("0xcf7389dc6c63637598402907d5431160ec8972a5")
	tokenAddress2 := NewAddressFromString("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa")

	tokenBalance1 := &TokenBalance{
		ID:            bson.NewObjectId(),
		Address:       tokenAddress1,
		Symbol:        "EOS",
		Balance:       NewBigInt("10000"),
		Allowance:     NewBigInt("10000"),
		LockedBalance: NewBigInt("5000"),
	}

	tokenBalance2 := &TokenBalance{
		ID:            bson.NewObjectId(),
		Address:       tokenAddress2,
		Symbol:        "ZRX",
		Balance:       NewBigInt("10000"),
		Allowance:     NewBigInt("10000"),
		LockedBalance: NewBigInt("5000"),
	}

	account := &Account{
		ID:      bson.NewObjectId(),
		Address: address,
		TokenBalances: map[common.Address]*TokenBalance{
			tokenAddress1: tokenBalance1,
			tokenAddress2: tokenBalance2,
		},
		IsBlocked: false,
	}

	data, err := bson.Marshal(account)
	if err != nil {
		t.Error(err)
	}

	decoded := &Account{}
	if err := bson.Unmarshal(data, decoded); err != nil {
		t.Error(err)
	}

	assert.Equal(decoded, account)
}
