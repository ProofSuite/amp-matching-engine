package daos

import (
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils"
	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	temp, _ := ioutil.TempDir("", "test")
	server.SetPath(temp)

	session := server.Session()
	db = &Database{session}
}

func TestPairDao(t *testing.T) {
	dao := NewPairDao()

	pair := &types.Pair{
		ID:                bson.NewObjectId(),
		BaseTokenSymbol:   "REQ",
		BaseTokenAddress:  common.HexToAddress("0xcf7389dc6c63637598402907d5431160ec8972a5"),
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

	testutils.ComparePair(t, pair, &all[0])

	byID, err := dao.GetByID(pair.ID)
	if err != nil {
		t.Errorf("Could not get pair by ID: %v", err)
	}

	testutils.ComparePair(t, pair, byID)

	byAddress, err := dao.GetByTokenAddress(pair.BaseTokenAddress, pair.QuoteTokenAddress)
	if err != nil {
		t.Errorf("Could not get pair by address: %v", err)
	}

	testutils.ComparePair(t, pair, byAddress)
}
