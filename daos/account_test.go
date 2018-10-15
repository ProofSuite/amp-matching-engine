package daos

import (
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

// var db *Database

// func TestMain(m *testing.M) {
// 	db := &Database{}
// 	dropTestServer := testutils.InitDBTestServer(db)
// 	defer dropTestServer()
// 	m.Run()
// }

func init() {
	server := testutils.NewDBTestServer()
	temp, _ := ioutil.TempDir("", "test")
	server.SetPath(temp)

	session := server.Session()
	db = &Database{Session: session}
}

func TestAccountDao(t *testing.T) {
	dao := NewAccountDao()
	dao.Drop()

	address := common.HexToAddress("0xe8e84ee367bc63ddb38d3d01bccef106c194dc47")
	tokenAddress1 := common.HexToAddress("0xcf7389dc6c63637598402907d5431160ec8972a5")
	tokenAddress2 := common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa")

	tokenBalance1 := &types.TokenBalance{
		Address:       tokenAddress1,
		Symbol:        "EOS",
		Balance:       big.NewInt(10000),
		Allowance:     big.NewInt(10000),
		LockedBalance: big.NewInt(5000),
	}

	tokenBalance2 := &types.TokenBalance{
		Address:       tokenAddress2,
		Symbol:        "ZRX",
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

	err := dao.Create(account)
	if err != nil {
		t.Errorf("Could not create order object")
	}

	a1, err := dao.GetByAddress(account.Address)
	if err != nil {
		t.Errorf("Could not get order by hash: %v", err)
	}

	testutils.CompareAccount(t, account, a1)
}

func TestAccountGetAllTokenBalances(t *testing.T) {
	dao := NewAccountDao()
	dao.Drop()

	address := common.HexToAddress("0xe8e84ee367bc63ddb38d3d01bccef106c194dc47")
	tokenAddress1 := common.HexToAddress("0xcf7389dc6c63637598402907d5431160ec8972a5")
	tokenAddress2 := common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa")

	tokenBalance1 := &types.TokenBalance{
		Address:       tokenAddress1,
		Symbol:        "EOS",
		Balance:       big.NewInt(10000),
		Allowance:     big.NewInt(10000),
		LockedBalance: big.NewInt(5000),
	}

	tokenBalance2 := &types.TokenBalance{
		Address:       tokenAddress2,
		Symbol:        "ZRX",
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

	err := dao.Create(account)
	if err != nil {
		t.Errorf("Could not create account object")
	}

	balances, err := dao.GetTokenBalances(account.Address)
	if err != nil {
		t.Errorf("Could not retrieve token balances: %v", balances)
	}

	assert.Equal(t, balances[tokenAddress1], tokenBalance1)
	assert.Equal(t, balances[tokenAddress2], tokenBalance2)
}

func TestGetTokenBalance(t *testing.T) {
	dao := NewAccountDao()
	dao.Drop()

	address := common.HexToAddress("0xe8e84ee367bc63ddb38d3d01bccef106c194dc47")
	tokenAddress1 := common.HexToAddress("0xcf7389dc6c63637598402907d5431160ec8972a5")
	tokenAddress2 := common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498")

	tokenBalance1 := &types.TokenBalance{
		Address:       tokenAddress1,
		Symbol:        "EOS",
		Balance:       big.NewInt(10000),
		Allowance:     big.NewInt(10000),
		LockedBalance: big.NewInt(5000),
	}

	tokenBalance2 := &types.TokenBalance{
		Address:       tokenAddress2,
		Symbol:        "ZRX",
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

func TestUpdateAccountBalance(t *testing.T) {
	dao := NewAccountDao()
	dao.Drop()

	address := common.HexToAddress("0xe8e84ee367bc63ddb38d3d01bccef106c194dc47")
	tokenAddress1 := common.HexToAddress("0xcf7389dc6c63637598402907d5431160ec8972a5")
	tokenAddress2 := common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa")

	tokenBalance1 := &types.TokenBalance{
		Address:       tokenAddress1,
		Symbol:        "EOS",
		Balance:       big.NewInt(10000),
		Allowance:     big.NewInt(10000),
		LockedBalance: big.NewInt(5000),
	}

	tokenBalance2 := &types.TokenBalance{
		Address:       tokenAddress2,
		Symbol:        "ZRX",
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
