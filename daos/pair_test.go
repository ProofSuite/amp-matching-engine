package daos

import (
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	temp, _ := ioutil.TempDir("", "test")
	server.SetPath(temp)

	session := server.Session()
	db = &Database{session}
}

func ComparePair(t *testing.T, a, b *types.Pair) {
	assert.Equal(t, a.ID, b.ID)
	assert.Equal(t, a.Name, b.Name)
	assert.Equal(t, a.BaseTokenID, b.BaseTokenID)
	assert.Equal(t, a.BaseTokenSymbol, b.BaseTokenSymbol)
	assert.Equal(t, a.BaseTokenAddress, b.BaseTokenAddress)
	assert.Equal(t, a.QuoteTokenID, b.QuoteTokenID)
	assert.Equal(t, a.QuoteTokenSymbol, b.QuoteTokenSymbol)
	assert.Equal(t, a.QuoteTokenAddress, b.QuoteTokenAddress)
	assert.Equal(t, a.Active, b.Active)
	assert.Equal(t, a.MakeFee, b.MakeFee)
	assert.Equal(t, a.TakeFee, b.TakeFee)
}

func TestPairDao(t *testing.T) {
	dao := NewPairDao()

	pair := &types.Pair{
		ID:                bson.NewObjectId(),
		Name:              "REQ",
		BaseTokenID:       bson.NewObjectId(),
		BaseTokenSymbol:   "REQ",
		BaseTokenAddress:  common.HexToAddress("0xcf7389dc6c63637598402907d5431160ec8972a5"),
		QuoteTokenID:      bson.NewObjectId(),
		QuoteTokenSymbol:  "WETH",
		QuoteTokenAddress: common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		Active:            true,
		MakeFee:           big.NewInt(10000),
		TakeFee:           big.NewInt(10000),
	}

	err := dao.Create(pair)
	if err != nil {
		t.Errorf("Could not create pair object: %+v", err)
	}

	all, err := dao.GetAll()
	if err != nil {
		t.Errorf("Could not get pairs: %+v", err)
	}

	ComparePair(t, pair, &all[0])

	byID, err := dao.GetByID(pair.ID)
	if err != nil {
		t.Errorf("Could not get pair by ID: %v", err)
	}

	ComparePair(t, pair, byID)

	byAddress, err := dao.GetByTokenAddress(pair.BaseTokenAddress, pair.QuoteTokenAddress)
	if err != nil {
		t.Errorf("Could not get pair by address: %v", err)
	}

	ComparePair(t, pair, byAddress)
}
