package daos

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func init() {
	temp, _ := ioutil.TempDir("", "test")
	server.SetPath(temp)

	session := server.Session()
	db = &Database{session}
}

func Compare(t *testing.T, a, b *types.Token) {
	assert.Equal(t, a.Name, b.Name)
	assert.Equal(t, a.Symbol, b.Symbol)
	assert.Equal(t, a.ContractAddress, b.ContractAddress)
	assert.Equal(t, a.Decimal, b.Decimal)
	assert.Equal(t, a.Active, b.Active)
	assert.Equal(t, a.Quote, b.Quote)
	assert.Equal(t, a.ID, b.ID)
}

func TestTokenDao(t *testing.T) {
	dao := NewTokenDao()

	token := &types.Token{
		Name:            "PRFT",
		Symbol:          "PRFT",
		ContractAddress: common.HexToAddress("0x6e9a406696617ec5105f9382d33ba3360fcfabcc"),
		Decimal:         18,
		Active:          true,
		Quote:           true,
	}

	err := dao.Create(token)
	if err != nil {
		t.Errorf("Could not create token object: %+v", err)
	}

	all, err := dao.GetAll()
	if err != nil {
		t.Errorf("Could not get wallets: %+v", err)
	}

	Compare(t, token, &all[0])

	byId, err := dao.GetByID(token.ID)
	if err != nil {
		t.Errorf("Could not get token by ID: %+v", err)
	}

	Compare(t, token, byId)

	byAddress, err := dao.GetByAddress(common.HexToAddress("0x6e9a406696617ec5105f9382d33ba3360fcfabcc"))
	if err != nil {
		t.Errorf("Could not get token by address: %+v", err)
	}

	fmt.Printf(":Headsfasdfasdf%+v", byAddress)

	Compare(t, token, byAddress)
}
