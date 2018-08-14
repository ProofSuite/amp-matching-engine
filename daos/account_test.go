package daos

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func CompareAccount(t *testing.T, a, b *types.Account) {
	assert.Equal(t, a.ID, b.ID)
	assert.Equal(t, a.Address, b.Address)
	assert.Equal(t, a.TokenBalances, b.TokenBalances)
	assert.Equal(t, a.IsBlocked, b.IsBlocked)
}

func TestAccountDao(t *testing.T) {
	address := common.HexToAddress("0xe8e84ee367bc63ddb38d3d01bccef106c194dc47")
	tokenAddress1 := common.HexToAddress("0xcf7389dc6c63637598402907d5431160ec8972a5")
	tokenAddress2 := common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa")

	tokenBalance1 := &types.TokenBalance{
		TokenID:       bson.NewObjectId(),
		TokenAddress:  tokenAddress1,
		TokenSymbol:   "EOS",
		Balance:       big.NewInt(10000),
		Allowance:     big.NewInt(10000),
		LockedBalance: big.NewInt(5000),
	}

	tokenBalance2 := &types.TokenBalance{
		TokenID:       bson.NewObjectId(),
		TokenAddress:  tokenAddress2,
		TokenSymbol:   "ZRX",
		Balance:       big.NewInt(10000),
		Allowance:     big.NewInt(10000),
		LockedBalance: big.NewInt(5000),
	}

	account := &types.Account{
		ID:      bson.NewObjectId(),
		Address: address,
		TokenBalances: map[common.Address]*types.TokenBalance{
			tokenAddress1: tokenBalance1,
			tokenAddress2: tokenBalance2,
		},
		IsBlocked: false,
	}

	dao := NewAccountDao()

	err := dao.Create(account)
	if err != nil {
		t.Errorf("Could not create order object")
	}

	a1, err := dao.GetByAddress(account.Address)
	if err != nil {
		t.Errorf("Could not get order by hash: %v", err)
	}

	CompareAccount(t, account, a1)
}

func TestAccountGetAllTokenBalances(t *testing.T) {
	address := common.HexToAddress("0xe8e84ee367bc63ddb38d3d01bccef106c194dc47")
	tokenAddress1 := common.HexToAddress("0xcf7389dc6c63637598402907d5431160ec8972a5")
	tokenAddress2 := common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa")

	tokenBalance1 := &types.TokenBalance{
		TokenID:       bson.NewObjectId(),
		TokenAddress:  tokenAddress1,
		TokenSymbol:   "EOS",
		Balance:       big.NewInt(10000),
		Allowance:     big.NewInt(10000),
		LockedBalance: big.NewInt(5000),
	}

	tokenBalance2 := &types.TokenBalance{
		TokenID:       bson.NewObjectId(),
		TokenAddress:  tokenAddress2,
		TokenSymbol:   "ZRX",
		Balance:       big.NewInt(10000),
		Allowance:     big.NewInt(10000),
		LockedBalance: big.NewInt(5000),
	}

	account := &types.Account{
		ID:      bson.NewObjectId(),
		Address: address,
		TokenBalances: map[common.Address]*types.TokenBalance{
			tokenAddress1: tokenBalance1,
			tokenAddress2: tokenBalance2,
		},
		IsBlocked: false,
	}

	dao := NewAccountDao()

	err := dao.Create(account)
	if err != nil {
		t.Errorf("Could not create account object")
	}

	balances, err := dao.GetAllTokenBalances(account.Address)

	if err != nil {
		t.Errorf("Could not retrieve token balances: %v", balances)
	}

	assert.Equal(t, balances[tokenAddress1], tokenBalance1)
	assert.Equal(t, balances[tokenAddress2], tokenBalance2)
}

func TestGetTokenBalance(t *testing.T) {
	address := common.HexToAddress("0xe8e84ee367bc63ddb38d3d01bccef106c194dc47")
	tokenAddress1 := common.HexToAddress("0xcf7389dc6c63637598402907d5431160ec8972a5")
	tokenAddress2 := common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498")

	tokenBalance1 := &types.TokenBalance{
		TokenID:       bson.NewObjectId(),
		TokenAddress:  tokenAddress1,
		TokenSymbol:   "EOS",
		Balance:       big.NewInt(10000),
		Allowance:     big.NewInt(10000),
		LockedBalance: big.NewInt(5000),
	}

	tokenBalance2 := &types.TokenBalance{
		TokenID:       bson.NewObjectId(),
		TokenAddress:  tokenAddress2,
		TokenSymbol:   "ZRX",
		Balance:       big.NewInt(10000),
		Allowance:     big.NewInt(10000),
		LockedBalance: big.NewInt(5000),
	}

	account := &types.Account{
		ID:      bson.NewObjectId(),
		Address: address,
		TokenBalances: map[common.Address]*types.TokenBalance{
			tokenAddress1: tokenBalance1,
			tokenAddress2: tokenBalance2,
		},
		IsBlocked: false,
	}

	dao := NewAccountDao()

	err := dao.Create(account)
	if err != nil {
		t.Errorf("Could not create account: %v", err)
	}

	balance, err := dao.GetTokenBalance(address, tokenAddress2)
	if err != nil {
		t.Errorf("Could not get token balance: %v", err)
	}

	assert.Equal(t, balance, tokenBalance2)
}

func TestAddress(t *testing.T) {
	address := common.HexToAddress("0xE8e84ee367bc63ddb38d3d01bccef106c194dc47")

	fmt.Printf("%v", address)

	res := address.Hex()
	fmt.Printf("%v", res)
}

func TestUpdateAccountBalance(t *testing.T) {
	address := common.HexToAddress("0xe8e84ee367bc63ddb38d3d01bccef106c194dc47")
	tokenAddress1 := common.HexToAddress("0xcf7389dc6c63637598402907d5431160ec8972a5")
	tokenAddress2 := common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa")

	tokenBalance1 := &types.TokenBalance{
		TokenID:       bson.NewObjectId(),
		TokenAddress:  tokenAddress1,
		TokenSymbol:   "EOS",
		Balance:       big.NewInt(10000),
		Allowance:     big.NewInt(10000),
		LockedBalance: big.NewInt(5000),
	}

	tokenBalance2 := &types.TokenBalance{
		TokenID:       bson.NewObjectId(),
		TokenAddress:  tokenAddress2,
		TokenSymbol:   "ZRX",
		Balance:       big.NewInt(10000),
		Allowance:     big.NewInt(10000),
		LockedBalance: big.NewInt(5000),
	}

	account := &types.Account{
		ID:      bson.NewObjectId(),
		Address: address,
		TokenBalances: map[common.Address]*types.TokenBalance{
			tokenAddress1: tokenBalance1,
			tokenAddress2: tokenBalance2,
		},
		IsBlocked: false,
	}

	dao := NewAccountDao()

	err := dao.Create(account)
	if err != nil {
		t.Errorf("Could not create account object")
	}

	err = dao.UpdateBalance(address, tokenAddress1, big.NewInt(20000))
	if err != nil {
		t.Errorf("Could not update balance")
	}

	balance, err := dao.GetTokenBalance(address, tokenAddress1)
	if err != nil {
		t.Errorf("Could not get token balance: %v", err)
	}

	assert.Equal(t, balance.Balance, big.NewInt(20000))
}

// func TestUpdateAccountBalance(t *testing.T) {
// 	address := common.HexToAddress("0xe8e84ee367bc63ddb38d3d01bccef106c194dc47")
// 	tokenAddress1 := common.HexToAddress("0xcf7389dc6c63637598402907d5431160ec8972a5")
// 	tokenAddress2 := common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa")

// 	tokenBalance1 := &types.TokenBalance{
// 		TokenID:       bson.NewObjectId(),
// 		TokenAddress:  tokenAddress1,
// 		TokenSymbol:   "EOS",
// 		Balance:       big.NewInt(10000),
// 		Allowance:     big.NewInt(10000),
// 		LockedBalance: big.NewInt(5000),
// 	}

// 	tokenBalance2 := &types.TokenBalance{
// 		TokenID:       bson.NewObjectId(),
// 		TokenAddress:  tokenAddress2,
// 		TokenSymbol:   "ZRX",
// 		Balance:       big.NewInt(10000),
// 		Allowance:     big.NewInt(10000),
// 		LockedBalance: big.NewInt(5000),
// 	}

// 	account := &types.Account{
// 		ID:      bson.NewObjectId(),
// 		Address: address,
// 		TokenBalances: map[common.Address]*types.TokenBalance{
// 			tokenAddress1: tokenBalance1,
// 			tokenAddress2: tokenBalance2,
// 		},
// 		IsBlocked: false,
// 	}

// 	dao := NewAccountDao()

// 	err := dao.Create(account)
// 	if err != nil {
// 		t.Errorf("Could not create account object")
// 	}

// 	err := dao.UpdateBalance(address, token common.Address, balance *big.Int)

// 			t.Errorf("Could not retrieve token balance: %v", err)
// 		balance, err := dao.GetTokenBalance(address, tokenAddress2)
// 		if err != nil {
// 	}

// 	assert.Equal(t, tokenBalance1, balance)

// }
